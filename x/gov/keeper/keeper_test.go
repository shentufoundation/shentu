package keeper_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/v2/simapp"
	. "github.com/certikfoundation/shentu/v2/x/gov/keeper"
	"github.com/certikfoundation/shentu/v2/x/gov/types"
)

func TestKeeper_ProposeAndVote(t *testing.T) {
	t.Log("Test keeper AddVote")
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))

	tp := govtypes.NewTextProposal("title0", "desc0")
	t.Run("Test submitting a proposal and adding a vote with yes", func(t *testing.T) {
		pp, err := app.GovKeeper.SubmitProposal(ctx, tp, addrs[0])
		if err != nil {
			panic(err)
		}
		options := govtypes.NewNonSplitVoteOption(govtypes.OptionYes)
		vote := govtypes.NewVote(pp.ProposalId, addrs[0], options)
		coins700 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 700*1e6))
		require.NoError(t, simapp.FundAccount(app.BankKeeper, ctx, addrs[1], coins700))

		votingPeriodActivated, err := app.GovKeeper.AddDeposit(ctx, pp.ProposalId, addrs[1], coins700)
		require.Equal(t, nil, err)
		require.Equal(t, true, votingPeriodActivated)

		voter, err := sdk.AccAddressFromBech32(vote.Voter)
		if err != nil {
			panic(err)
		}
		err = app.GovKeeper.AddVote(ctx, pp.ProposalId, voter, options)
		require.Equal(t, nil, err)

		// the vote does not count since addr[0] is not a validator
		results := map[govtypes.VoteOption]sdk.Dec{
			govtypes.OptionYes:        sdk.ZeroDec(),
			govtypes.OptionAbstain:    sdk.ZeroDec(),
			govtypes.OptionNo:         sdk.ZeroDec(),
			govtypes.OptionNoWithVeto: sdk.ZeroDec(),
		}

		pass, veto, res := Tally(ctx, app.GovKeeper, pp)
		require.Equal(t, false, pass)
		require.Equal(t, false, veto)
		require.Equal(t, govtypes.NewTallyResultFromMap(results), res)
	})

	// TODO: more tests. validator cases
}

func TestKeeper_GetVotes(t *testing.T) {
	t.Log("Test keeper GetVotes")
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(80000*1e6))

	tp := govtypes.TextProposal{Title: "title0", Description: "desc0"}
	t.Run("Test adding a lot of votes and retrieving them", func(t *testing.T) {
		pp, err := app.GovKeeper.SubmitProposal(ctx, &tp, addrs[0])
		require.Equal(t, nil, err)
		coins700 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 700*1e6))
		votingPeriodActivated, err := app.GovKeeper.AddDeposit(ctx, pp.ProposalId, addrs[0], coins700)
		require.Equal(t, nil, err)
		require.Equal(t, true, votingPeriodActivated)

		var addr sdk.AccAddress
		for i := 0; i < 880; i++ {
			addr = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
			options := govtypes.NewNonSplitVoteOption(govtypes.OptionYes)
			vote := govtypes.NewVote(pp.ProposalId, addr, options)
			voter, err := sdk.AccAddressFromBech32(vote.Voter)
			if err != nil {
				panic(err)
			}
			err = app.GovKeeper.AddVote(ctx, vote.ProposalId, voter, options)
			require.Equal(t, nil, err)
		}

		retrievedVotes := app.GovKeeper.GetVotesPaginated(ctx, pp.ProposalId, 1, 2000)
		require.Equal(t, 880, len(retrievedVotes))
		retrievedVotes = app.GovKeeper.GetVotesPaginated(ctx, pp.ProposalId, 2, 200)
		require.Equal(t, 200, len(retrievedVotes))
		retrievedVotes = app.GovKeeper.GetVotesPaginated(ctx, pp.ProposalId, 5, 200)
		require.Equal(t, 80, len(retrievedVotes))

		retrievedVotesNoPage := app.GovKeeper.GetVotes(ctx, pp.ProposalId)
		require.Equal(t, 880, len(retrievedVotesNoPage))

		for i := range retrievedVotes[:10] {
			require.True(t, reflect.DeepEqual(retrievedVotes[i], retrievedVotesNoPage[i+800]))
		}
	})
}

func TestKeeper_AddDeposit(t *testing.T) {
	t.Log("Test keeper AddDeposit")
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(10000))

	coins := sdk.Coins{sdk.NewInt64Coin("uctk", 80000*1e6)}
	require.NoError(t, simapp.FundAccount(app.BankKeeper, ctx, addrs[1], coins))

	tp := govtypes.TextProposal{Title: "title0", Description: "desc0"}

	t.Run("adding deposit and proposal doesn't exist", func(t *testing.T) {
		pp, err := app.GovKeeper.SubmitProposal(ctx, &tp, addrs[0])
		require.Equal(t, nil, err)
		coins100 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 100*1e6))

		votingPeriodActivated, err := app.GovKeeper.AddDeposit(ctx, pp.ProposalId+1, addrs[1], coins100)
		errString := fmt.Sprintf("%d: unknown proposal", pp.ProposalId+1)
		require.EqualError(t, err, errString)
		require.Equal(t, false, votingPeriodActivated)
	})

	t.Run("adding deposit not enough balance", func(t *testing.T) {
		pp, err := app.GovKeeper.SubmitProposal(ctx, &tp, addrs[0])
		require.Equal(t, nil, err)
		coins15000 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 15000*1e6))

		votingPeriodActivated, err := app.GovKeeper.AddDeposit(ctx, pp.ProposalId, addrs[0], coins15000)
		errString := "10000uctk is smaller than 15000000000uctk: insufficient funds"
		require.EqualError(t, err, errString)
		require.Equal(t, false, votingPeriodActivated)
	})

	t.Run("adding deposit and waiting for more deposits", func(t *testing.T) {
		pp, err := app.GovKeeper.SubmitProposal(ctx, &tp, addrs[0])
		require.Equal(t, nil, err)
		coins100 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 100*1e6))

		votingPeriodActivated, err := app.GovKeeper.AddDeposit(ctx, pp.ProposalId, addrs[1], coins100)
		require.Equal(t, nil, err)
		require.Equal(t, false, votingPeriodActivated)
	})

	t.Run("adding more deposit and still waiting for more", func(t *testing.T) {
		pp, err := app.GovKeeper.SubmitProposal(ctx, &tp, addrs[0])
		require.Equal(t, nil, err)
		coins100 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 100*1e6))
		coins200 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 200*1e6))

		votingPeriodActivated, err := app.GovKeeper.AddDeposit(ctx, pp.ProposalId, addrs[1], coins100)
		require.Equal(t, nil, err)
		require.Equal(t, false, votingPeriodActivated)

		votingPeriodActivated, err = app.GovKeeper.AddDeposit(ctx, pp.ProposalId, addrs[1], coins200)
		require.Equal(t, nil, err)
		require.Equal(t, false, votingPeriodActivated)
	})

	t.Run("adding deposit and entering votingPeriod", func(t *testing.T) {
		pp, err := app.GovKeeper.SubmitProposal(ctx, &tp, addrs[0])
		require.Equal(t, nil, err)
		coins700 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 700*1e6))

		votingPeriodActivated, err := app.GovKeeper.AddDeposit(ctx, pp.ProposalId, addrs[1], coins700)
		require.Equal(t, nil, err)
		require.Equal(t, true, votingPeriodActivated)
	})

	t.Run("entering votingPeriod and trying to add more deposit", func(t *testing.T) {
		pp, err := app.GovKeeper.SubmitProposal(ctx, &tp, addrs[0])
		require.Equal(t, nil, err)
		coins700 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 700*1e6))
		coinsAfterAvtivated := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 1))

		votingPeriodActivated, err := app.GovKeeper.AddDeposit(ctx, pp.ProposalId, addrs[1], coins700)
		require.Equal(t, nil, err)
		require.Equal(t, true, votingPeriodActivated)

		votingPeriodActivated, err = app.GovKeeper.AddDeposit(ctx, pp.ProposalId, addrs[1], coinsAfterAvtivated)
		errString := fmt.Sprintf("%d: proposal already active", pp.ProposalId)
		require.EqualError(t, err, errString)
		require.Equal(t, false, votingPeriodActivated)
	})
}

func TestKeeper_DepositOperation(t *testing.T) {
	t.Log("Test keeper DepositOperation")
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := simapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(80000*1e6))

	tp := govtypes.TextProposal{Title: "title0", Description: "desc0"}

	t.Run("refund all deposits in a specific proposal", func(t *testing.T) {
		pp, err := app.GovKeeper.SubmitProposal(ctx, &tp, addrs[0])
		require.Equal(t, nil, err)
		coins100 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 100*1e6))
		coins50 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 50*1e6))
		coins20 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 20*1e6))

		_, _ = app.GovKeeper.AddDeposit(ctx, pp.ProposalId, addrs[1], coins100)
		_, _ = app.GovKeeper.AddDeposit(ctx, pp.ProposalId, addrs[2], coins50)
		votingPeriodActivated, err := app.GovKeeper.AddDeposit(ctx, pp.ProposalId, addrs[3], coins20)
		require.Equal(t, nil, err)
		require.Equal(t, false, votingPeriodActivated)

		addr1Amount := app.BankKeeper.GetAllBalances(ctx, addrs[1])
		addr2Amount := app.BankKeeper.GetAllBalances(ctx, addrs[2])
		addr3Amount := app.BankKeeper.GetAllBalances(ctx, addrs[3])
		require.Equal(t, sdk.NewInt(79900*1e6).Int64(), addr1Amount.AmountOf(app.StakingKeeper.BondDenom(ctx)).Int64())
		require.Equal(t, sdk.NewInt(79950*1e6).Int64(), addr2Amount.AmountOf(app.StakingKeeper.BondDenom(ctx)).Int64())
		require.Equal(t, sdk.NewInt(79980*1e6).Int64(), addr3Amount.AmountOf(app.StakingKeeper.BondDenom(ctx)).Int64())

		app.GovKeeper.RefundDepositsByProposalID(ctx, pp.ProposalId)
		depositsRemaining := app.GovKeeper.GetAllDeposits(ctx)
		require.Equal(t, types.Deposits(nil), depositsRemaining)
		addr1Amount = app.BankKeeper.GetAllBalances(ctx, addrs[1])
		addr2Amount = app.BankKeeper.GetAllBalances(ctx, addrs[2])
		addr3Amount = app.BankKeeper.GetAllBalances(ctx, addrs[3])
		require.Equal(t, sdk.NewInt(80000*1e6).Int64(), addr1Amount.AmountOf(app.StakingKeeper.BondDenom(ctx)).Int64())
		require.Equal(t, sdk.NewInt(80000*1e6).Int64(), addr2Amount.AmountOf(app.StakingKeeper.BondDenom(ctx)).Int64())
		require.Equal(t, sdk.NewInt(80000*1e6).Int64(), addr3Amount.AmountOf(app.StakingKeeper.BondDenom(ctx)).Int64())
	})
	t.Run("delete all deposits in a specific proposal", func(t *testing.T) {
		pp, err := app.GovKeeper.SubmitProposal(ctx, &tp, addrs[0])
		require.Equal(t, nil, err)
		coins10 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 10*1e6))
		coins50 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 50*1e6))
		coins20 := sdk.NewCoins(sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 20*1e6))

		_, _ = app.GovKeeper.AddDeposit(ctx, pp.ProposalId, addrs[1], coins10)
		_, _ = app.GovKeeper.AddDeposit(ctx, pp.ProposalId, addrs[2], coins20)
		votingPeriodActivated, err := app.GovKeeper.AddDeposit(ctx, pp.ProposalId, addrs[3], coins50)
		require.Equal(t, nil, err)
		require.Equal(t, false, votingPeriodActivated)

		addr1Amount := app.BankKeeper.GetAllBalances(ctx, addrs[1])
		addr2Amount := app.BankKeeper.GetAllBalances(ctx, addrs[2])
		addr3Amount := app.BankKeeper.GetAllBalances(ctx, addrs[3])
		require.Equal(t, sdk.NewInt(79990*1e6).Int64(), addr1Amount.AmountOf(app.StakingKeeper.BondDenom(ctx)).Int64())
		require.Equal(t, sdk.NewInt(79980*1e6).Int64(), addr2Amount.AmountOf(app.StakingKeeper.BondDenom(ctx)).Int64())
		require.Equal(t, sdk.NewInt(79950*1e6).Int64(), addr3Amount.AmountOf(app.StakingKeeper.BondDenom(ctx)).Int64())

		app.GovKeeper.DeleteDepositsByProposalID(ctx, pp.ProposalId)
		depositsRemaining := app.GovKeeper.GetAllDeposits(ctx)
		require.Equal(t, types.Deposits(nil), depositsRemaining)

		addr1Amount = app.BankKeeper.GetAllBalances(ctx, addrs[1])
		addr2Amount = app.BankKeeper.GetAllBalances(ctx, addrs[2])
		addr3Amount = app.BankKeeper.GetAllBalances(ctx, addrs[3])
		require.Equal(t, sdk.NewInt(79990*1e6).Int64(), addr1Amount.AmountOf(app.StakingKeeper.BondDenom(ctx)).Int64())
		require.Equal(t, sdk.NewInt(79980*1e6).Int64(), addr2Amount.AmountOf(app.StakingKeeper.BondDenom(ctx)).Int64())
		require.Equal(t, sdk.NewInt(79950*1e6).Int64(), addr3Amount.AmountOf(app.StakingKeeper.BondDenom(ctx)).Int64())
	})
}
