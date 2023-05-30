package gov_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkauth "github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/gov"
	"github.com/shentufoundation/shentu/v2/x/gov/types"
)

func TestImportExportQueues(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, valTokens)

	SortAddresses(addrs)

	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx = app.BaseApp.NewContext(false, tmproto.Header{})

	// Create two proposals, put the second into the voting period
	proposal := TestProposal
	proposal1, err := app.GovKeeper.SubmitProposal(ctx, proposal)
	require.NoError(t, err)
	proposalID1 := proposal1.ProposalId

	proposal2, err := app.GovKeeper.SubmitProposal(ctx, proposal)
	require.NoError(t, err)
	proposalID2 := proposal2.ProposalId

	votingStarted, err := app.GovKeeper.AddDeposit(ctx, proposalID2, addrs[0], app.GovKeeper.GetDepositParams(ctx).MinDeposit)
	require.NoError(t, err)
	require.True(t, votingStarted)

	proposal1, ok := app.GovKeeper.GetProposal(ctx, proposalID1)
	require.True(t, ok)
	proposal2, ok = app.GovKeeper.GetProposal(ctx, proposalID2)
	require.True(t, ok)
	require.True(t, proposal1.Status == govtypes.StatusDepositPeriod)
	require.True(t, proposal2.Status == govtypes.StatusVotingPeriod)

	authGenState := sdkauth.ExportGenesis(ctx, app.AccountKeeper)
	bankGenState := app.BankKeeper.ExportGenesis(ctx)

	// export the state and import it into a new app
	govGenState := gov.ExportGenesis(ctx, app.GovKeeper)
	genesisState := shentuapp.NewDefaultGenesisState()

	genesisState[authtypes.ModuleName] = app.Codec().MustMarshalJSON(authGenState)
	genesisState[banktypes.ModuleName] = app.Codec().MustMarshalJSON(bankGenState)
	genesisState[govtypes.ModuleName] = app.Codec().MustMarshalJSON(govGenState)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	db := dbm.NewMemDB()
	app2 := shentuapp.NewShentuApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, shentuapp.DefaultNodeHome, 1, shentuapp.MakeEncodingConfig(), shentuapp.EmptyAppOptions{})

	app2.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	app2.Commit()
	app2.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app2.LastBlockHeight() + 1}})

	header = tmproto.Header{Height: app2.LastBlockHeight() + 1}
	app2.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx2 := app2.BaseApp.NewContext(false, tmproto.Header{})

	// Jump the time forward past the DepositPeriod and VotingPeriod
	ctx2 = ctx2.WithBlockTime(ctx2.BlockHeader().Time.Add(app2.GovKeeper.GetDepositParams(ctx2).MaxDepositPeriod).Add(app2.GovKeeper.GetVotingParams(ctx2).VotingPeriod))

	// Make sure that they are still in the DepositPeriod and VotingPeriod respectively
	proposal1, ok = app2.GovKeeper.GetProposal(ctx2, proposalID1)
	require.True(t, ok)
	proposal2, ok = app2.GovKeeper.GetProposal(ctx2, proposalID2)
	require.True(t, ok)
	require.True(t, proposal1.Status == govtypes.StatusDepositPeriod)
	require.True(t, proposal2.Status == govtypes.StatusVotingPeriod)

	macc := app2.GovKeeper.GetGovernanceAccount(ctx2)
	require.Equal(t, app2.GovKeeper.GetDepositParams(ctx2).MinDeposit, app2.BankKeeper.GetAllBalances(ctx2, macc.GetAddress()))

	// Run the endblocker. Check to make sure that proposal1 is removed from state, and proposal2 is finished VotingPeriod.
	gov.EndBlocker(ctx2, app2.GovKeeper)

	proposal1, ok = app2.GovKeeper.GetProposal(ctx2, proposalID1)
	require.False(t, ok)

	proposal2, ok = app2.GovKeeper.GetProposal(ctx2, proposalID2)
	require.True(t, ok)
	require.True(t, proposal2.Status == govtypes.StatusRejected)
}

func TestImportExportQueues_ErrorUnconsistentState(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	require.Panics(t, func() {
		gov.InitGenesis(ctx, app.GovKeeper, app.AccountKeeper, app.BankKeeper, types.GenesisState{
			Deposits: govtypes.Deposits{
				{
					ProposalId: 1234,
					Depositor:  "me",
					Amount: sdk.Coins{
						sdk.NewCoin(
							"stake",
							sdk.NewInt(1234),
						),
					},
				},
			},
		})
	})
}

func TestEqualProposals(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, valTokens)

	SortAddresses(addrs)

	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// Submit two proposals
	proposal := TestProposal
	proposal1, err := app.GovKeeper.SubmitProposal(ctx, proposal)
	require.NoError(t, err)

	proposal2, err := app.GovKeeper.SubmitProposal(ctx, proposal)
	require.NoError(t, err)

	// They are similar but their IDs should be different
	require.NotEqual(t, proposal1, proposal2)
	require.NotEqual(t, proposal1, proposal2)

	// Now create two genesis blocks
	state1 := types.GenesisState{Proposals: []govtypes.Proposal{proposal1}}
	state2 := types.GenesisState{Proposals: []govtypes.Proposal{proposal2}}
	require.NotEqual(t, state1, state2)

	// Now make proposals identical by setting both IDs to 55
	proposal1.ProposalId = 55
	proposal2.ProposalId = 55
	require.Equal(t, proposal1, proposal1)
	require.Equal(t, proposal1, proposal2)

	// Reassign proposals into state
	state1.Proposals[0] = proposal1
	state2.Proposals[0] = proposal2

	// State should be identical now.
	require.Equal(t, state1, state2)
}
