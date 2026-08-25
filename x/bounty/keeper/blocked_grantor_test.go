package keeper_test

import (
	"time"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/shentufoundation/shentu/v2/x/bounty"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// blockedGrantorFixture stages an expired theorem carrying a single grant from
// grantor, with the module account funded so the refund cannot fail for lack of
// balance.
func (suite *KeeperTestSuite) blockedGrantorFixture(grantor sdk.AccAddress) (uint64, sdk.Coins) {
	t := suite.T()

	// The fixture ctx carries no block time, and the queue key needs a real one.
	suite.ctx = suite.ctx.WithBlockTime(time.Date(2026, 8, 25, 0, 0, 0, 0, time.UTC))

	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(t, err)
	amount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(100)))

	endTime := suite.ctx.BlockTime().Add(-time.Hour) // already expired
	theorem := types.Theorem{
		Id:         42,
		Title:      "blocked grantor",
		Proposer:   suite.normalAddr.String(),
		Status:     types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		EndTime:    &endTime,
		TotalGrant: amount,
	}
	require.NoError(t, suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem))
	require.NoError(t, suite.keeper.ActiveTheoremsQueue.Set(
		suite.ctx, collections.Join(endTime, theorem.Id), theorem.Id))
	require.NoError(t, suite.keeper.Grants.Set(
		suite.ctx, collections.Join(theorem.Id, grantor),
		types.Grant{TheoremId: theorem.Id, Grantor: grantor.String(), Amount: amount}))
	require.NoError(t, suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, amount))

	return theorem.Id, amount
}

// A grant whose grantor sits in the bank blocked set used to abort EndBlocker.
// That error reaches FinalizeBlock, which has no recover, so every node fails at
// the same height and the chain halts. The refund must now land in the community
// pool instead and EndBlocker must succeed.
//
// A blocked address can hold a grant even though AddGrant rejects one today:
// MsgGrant declares grantor as its signer, and a governance proposal can make
// the gov module account that signer. Records predating the AddGrant guard are
// reachable the same way.
func (suite *KeeperTestSuite) TestEndBlocker_BlockedGrantorRefundGoesToCommunityPool() {
	t := suite.T()

	govAddr := authtypes.NewModuleAddress("gov")
	require.True(t, suite.app.BankKeeper.BlockedAddr(govAddr),
		"precondition: the gov module address must be in the blocked set")

	theoremID, amount := suite.blockedGrantorFixture(govAddr)

	poolBefore, err := suite.app.DistrKeeper.FeePool.Get(suite.ctx)
	require.NoError(t, err)

	require.NoError(t, bounty.EndBlocker(suite.ctx, &suite.keeper),
		"EndBlocker must not fail on a blocked grantor")

	// The grant is settled, not left dangling.
	_, err = suite.keeper.Grants.Get(suite.ctx, collections.Join(theoremID, govAddr))
	require.Error(t, err, "grant record should be removed after the refund")

	// The funds went to the community pool.
	poolAfter, err := suite.app.DistrKeeper.FeePool.Get(suite.ctx)
	require.NoError(t, err)
	delta := poolAfter.CommunityPool.Sub(poolBefore.CommunityPool)
	require.True(t, delta.Equal(sdk.NewDecCoinsFromCoins(amount...)),
		"community pool should grow by the refunded amount: got %s want %s",
		delta, sdk.NewDecCoinsFromCoins(amount...))

	// The redirect is observable by indexers.
	var found bool
	for _, e := range suite.ctx.EventManager().Events() {
		if e.Type == types.EventTypeGrantRefundRedirected {
			found = true
		}
	}
	require.True(t, found, "expected a %s event", types.EventTypeGrantRefundRedirected)
}

// The ordinary path must keep paying refunds straight back to the grantor.
func (suite *KeeperTestSuite) TestEndBlocker_UnblockedGrantorRefundGoesBack() {
	t := suite.T()

	grantor := suite.normalAddr
	require.False(t, suite.app.BankKeeper.BlockedAddr(grantor))

	theoremID, amount := suite.blockedGrantorFixture(grantor)
	balBefore := suite.app.BankKeeper.GetAllBalances(suite.ctx, grantor)

	require.NoError(t, bounty.EndBlocker(suite.ctx, &suite.keeper))

	_, err := suite.keeper.Grants.Get(suite.ctx, collections.Join(theoremID, grantor))
	require.Error(t, err, "grant record should be removed after the refund")

	balAfter := suite.app.BankKeeper.GetAllBalances(suite.ctx, grantor)
	require.True(t, balAfter.Sub(balBefore...).Equal(amount),
		"grantor should be refunded: got %s want %s", balAfter.Sub(balBefore...), amount)
}

// AddGrant rejects a blocked grantor up front, so no new unrefundable grant is
// created.
func (suite *KeeperTestSuite) TestAddGrant_RejectsBlockedGrantor() {
	t := suite.T()

	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(t, err)

	theorem := types.Theorem{
		Id: 7, Title: "t", Proposer: suite.normalAddr.String(),
		Status:     types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		TotalGrant: sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(0))),
	}
	require.NoError(t, suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem))

	govAddr := authtypes.NewModuleAddress("gov")
	amount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(100)))

	err = suite.keeper.AddGrant(suite.ctx, theorem.Id, govAddr, amount)
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	_, err = suite.keeper.Grants.Get(suite.ctx, collections.Join(theorem.Id, govAddr))
	require.Error(t, err, "no grant record should have been created")
}

// An underfunded module account means the books are already wrong: every grant
// was funded into that account when it was created. This case is deliberately
// left fatal rather than skipped-and-retried — the account is shared by every
// theorem's grants and every proof's deposits, so a retry could settle this
// refund out of unrelated funds and leave the newer obligation unbacked.
func (suite *KeeperTestSuite) TestEndBlocker_UnderfundedRefundIsFatal() {
	t := suite.T()

	suite.ctx = suite.ctx.WithBlockTime(time.Date(2026, 8, 25, 0, 0, 0, 0, time.UTC))
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(t, err)

	grantor := suite.normalAddr
	// More than the module account holds, so the refund cannot go through.
	amount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1_000_000_000_000)))

	endTime := suite.ctx.BlockTime().Add(-time.Hour)
	theorem := types.Theorem{
		Id: 99, Title: "underfunded", Proposer: grantor.String(),
		Status:  types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		EndTime: &endTime, TotalGrant: amount,
	}
	require.NoError(t, suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem))
	require.NoError(t, suite.keeper.ActiveTheoremsQueue.Set(
		suite.ctx, collections.Join(endTime, theorem.Id), theorem.Id))
	require.NoError(t, suite.keeper.Grants.Set(
		suite.ctx, collections.Join(theorem.Id, grantor),
		types.Grant{TheoremId: theorem.Id, Grantor: grantor.String(), Amount: amount}))

	require.Error(t, bounty.EndBlocker(suite.ctx, &suite.keeper),
		"a corrupt ledger must surface, not be retried against pooled funds")
}
