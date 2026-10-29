package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
)

// setupTallyFixture creates one bonded validator (1000 self-delegation) plus a
// delegator (500 delegation to that validator), and a proposal to vote on.
//
// After bonding: validator DelegatorShares = 3_000_000, BondedTokens = 3_000_000
// (2M self + 1M delegated, 1:1 ratio); the delegator's 1M shares = 1M power.
// Amounts are millions so that consensus power (tokens / 10^6) is non-zero
// and the validator enters the bonded set.
func setupTallyFixture(t *testing.T) (*shentuapp.ShentuApp, sdk.Context, govtypesv1.Proposal, sdk.AccAddress, sdk.AccAddress) {
	t.Helper()
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	addrs := shentuapp.AddTestAddrsIncremental(app, ctx, 3, math.NewInt(10_000_000))
	valAcc, delegator := addrs[0], addrs[1]

	stakingSrv := stakingkeeper.NewMsgServerImpl(app.StakingKeeper)

	// create validator with 1000 self-delegation
	pubKey := ed25519.GenPrivKey().PubKey()
	msgCreate := &stakingtypes.MsgCreateValidator{
		Description:       stakingtypes.Description{Moniker: "test-val"},
		Commission:        stakingtypes.CommissionRates{Rate: math.LegacyZeroDec(), MaxRate: math.LegacyOneDec(), MaxChangeRate: math.LegacyZeroDec()},
		MinSelfDelegation: math.OneInt(),
		DelegatorAddress:  valAcc.String(),
		ValidatorAddress:  sdk.ValAddress(valAcc).String(),
		Pubkey:            codectypes.UnsafePackAny(pubKey),
		Value:             sdk.NewCoin("stake", math.NewInt(2_000_000)),
	}
	_, err := stakingSrv.CreateValidator(ctx, msgCreate)
	require.NoError(t, err)

	// apply validator set updates so the validator enters the bonded set
	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)

	// delegator delegates 500
	_, err = stakingSrv.Delegate(ctx, stakingtypes.NewMsgDelegate(
		delegator.String(), sdk.ValAddress(valAcc).String(), sdk.NewCoin("stake", math.NewInt(1_000_000)),
	))
	require.NoError(t, err)

	val, err := app.StakingKeeper.Validator(ctx, sdk.ValAddress(valAcc))
	require.NoError(t, err)
	t.Logf("DEBUG status=%v bonded=%v shares=%v", val.GetStatus(), val.GetBondedTokens(), val.GetDelegatorShares())
	require.True(t, val.IsBonded(), "validator must be bonded for tally math")
	require.Equal(t, int64(3_000_000), val.GetBondedTokens().Int64())
	require.Equal(t, "3000000.000000000000000000", val.GetDelegatorShares().String())

	// sanity: delegator voting power for 1M shares must be 1M
	shares := math.LegacyNewDec(1_000_000)
	power := shares.MulInt(val.GetBondedTokens()).Quo(val.GetDelegatorShares())
	require.Equal(t, "1000000.000000000000000000", power.String())

	// submit a proposal (any msg signed by gov; tally only uses the id)
	govAddr := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()
	msgSend := banktypes.NewMsgSend(govAddr, addrs[2], sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1))))
	proposer, err := app.AccountKeeper.AddressCodec().StringToBytes(valAcc.String())
	require.NoError(t, err)
	proposal, err := app.GovKeeper.SubmitProposal(ctx, []sdk.Msg{msgSend}, "", "tally test", "summary", proposer, false)
	require.NoError(t, err)

	return app, ctx, proposal, valAcc, delegator
}

// setBallot writes a single-option weight-1 ballot directly to the Votes
// collection.
func setBallot(t *testing.T, app *shentuapp.ShentuApp, ctx sdk.Context, proposalID uint64, voter sdk.AccAddress, opt govtypesv1.VoteOption) {
	t.Helper()
	vote := govtypesv1.NewVote(proposalID, voter, govtypesv1.WeightedVoteOptions{
		{Option: opt, Weight: "1"},
	}, "")
	require.NoError(t, app.GovKeeper.Votes.Set(ctx, collections.Join(proposalID, voter), vote))
}

// TestTallyDelegatorVoteSingleCount guards against the delegator
// double-counting regression: a delegator's voting power must be tallied
// exactly once, and the validator's own power must be deducted exactly once.
func TestTallyDelegatorVoteSingleCount(t *testing.T) {
	app, ctx, proposal, _, delegator := setupTallyFixture(t)

	// delegator votes YES; validator does not vote
	setBallot(t, app, ctx, proposal.Id, delegator, govtypesv1.OptionYes)

	_, _, tally, err := app.GovKeeper.Tally(ctx, proposal)
	require.NoError(t, err)

	// before the fix the delegator's 1M power was counted twice (2M)
	require.Equal(t, "1000000", tally.YesCount, "delegator voting power must be counted exactly once")
	require.Equal(t, "0", tally.NoCount)
}

// TestTallyValidatorVoteNotDoubleDeducted verifies that when both the
// validator and its delegator vote (with different options), each power is
// counted once: YES = 500 (delegator), NO = 1000 (validator self share).
func TestTallyValidatorVoteNotDoubleDeducted(t *testing.T) {
	app, ctx, proposal, valAcc, delegator := setupTallyFixture(t)

	setBallot(t, app, ctx, proposal.Id, delegator, govtypesv1.OptionYes)
	setBallot(t, app, ctx, proposal.Id, valAcc, govtypesv1.OptionNo)

	_, _, tally, err := app.GovKeeper.Tally(ctx, proposal)
	require.NoError(t, err)

	// before the fix the double deduction left the validator with negative
	// shares-after-deductions and inflated the delegator side
	require.Equal(t, "1000000", tally.YesCount)
	require.Equal(t, "2000000", tally.NoCount)
}
