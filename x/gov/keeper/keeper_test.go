package keeper_test

//
//import (
//	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
//	"github.com/shentufoundation/shentu/v2/common"
//	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
//	"github.com/stretchr/testify/require"
//	"testing"
//
//	"github.com/stretchr/testify/suite"
//
//	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
//
//	"github.com/cosmos/cosmos-sdk/baseapp"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
//	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
//	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
//	shentuapp "github.com/shentufoundation/shentu/v2/app"
//)
//
//// shared setup
//type KeeperTestSuite struct {
//	suite.Suite
//
//	app               *shentuapp.ShentuApp
//	ctx               sdk.Context
//	queryClient       govtypesv1.QueryClient
//	legacyQueryClient govtypesv1beta1.QueryClient
//	addrs             []sdk.AccAddress
//	msgSrvr           govtypesv1.MsgServer
//	legacyMsgSrvr     govtypesv1beta1.MsgServer
//}
//
//func (suite *KeeperTestSuite) SetupTest() {
//	app := shentuapp.Setup(suite.T(), false)
//	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
//
//	// Populate the gov account with some coins, as the TestProposal we have
//	// is a MsgSend from the gov account.
//	coins := sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(1e10)))
//	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
//	suite.NoError(err)
//	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, govtypes.ModuleName, coins)
//	suite.NoError(err)
//
//	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
//	govtypesv1.RegisterQueryServer(queryHelper, app.GovKeeper.Keeper)
//	legacyQueryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
//	//govtypesv1beta1.RegisterQueryServer(legacyQueryHelper, keeper.NewLegacyQueryServer(app.GovKeeper))
//	queryClient := govtypesv1.NewQueryClient(queryHelper)
//	legacyQueryClient := govtypesv1beta1.NewQueryClient(legacyQueryHelper)
//
//	suite.app = app
//	suite.ctx = ctx
//	suite.queryClient = queryClient
//	suite.legacyQueryClient = legacyQueryClient
//	suite.msgSrvr = keeper.NewMsgServerImpl(suite.app.GovKeeper)
//	govAcct := suite.app.GovKeeper.GetGovernanceAccount(suite.ctx).GetAddress()
//	suite.legacyMsgSrvr = keeper.NewLegacyMsgServerImpl(govAcct.String(), suite.msgSrvr, app.GovKeeper)
//	suite.addrs = shentuapp.AddTestAddrsIncremental(app, ctx, 2, sdk.NewInt(1e10))
//}
//
//func TestIncrementProposalNumber(t *testing.T) {
//	app := shentuapp.Setup(t, false)
//	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
//
//	tp := TestProposal
//	_, err := app.GovKeeper.SubmitProposal(ctx, tp, "")
//	require.NoError(t, err)
//	_, err = app.GovKeeper.SubmitProposal(ctx, tp, "")
//	require.NoError(t, err)
//	_, err = app.GovKeeper.SubmitProposal(ctx, tp, "")
//	require.NoError(t, err)
//	_, err = app.GovKeeper.SubmitProposal(ctx, tp, "")
//	require.NoError(t, err)
//	_, err = app.GovKeeper.SubmitProposal(ctx, tp, "")
//	require.NoError(t, err)
//	proposal6, err := app.GovKeeper.SubmitProposal(ctx, tp, "")
//	require.NoError(t, err)
//
//	require.Equal(t, uint64(6), proposal6.Id)
//}
//
//func TestProposalQueues(t *testing.T) {
//	app := shentuapp.Setup(t, false)
//	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
//
//	// create test proposals
//	tp := TestProposal
//	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp, "")
//	require.NoError(t, err)
//
//	inactiveIterator := app.GovKeeper.InactiveProposalQueueIterator(ctx, *proposal.DepositEndTime)
//	require.True(t, inactiveIterator.Valid())
//
//	proposalID := govtypes.GetProposalIDFromBytes(inactiveIterator.Value())
//	require.Equal(t, proposalID, proposal.Id)
//	inactiveIterator.Close()
//
//	app.GovKeeper.ActivateVotingPeriod(ctx, proposal)
//
//	proposal, ok := app.GovKeeper.GetProposal(ctx, proposal.Id)
//	require.True(t, ok)
//
//	activeIterator := app.GovKeeper.ActiveProposalQueueIterator(ctx, *proposal.VotingEndTime)
//	require.True(t, activeIterator.Valid())
//
//	proposalID, _ = govtypes.SplitActiveProposalQueueKey(activeIterator.Key())
//	require.Equal(t, proposalID, proposal.Id)
//
//	activeIterator.Close()
//}
//
//func TestKeeperTestSuite(t *testing.T) {
//	suite.Run(t, new(KeeperTestSuite))
//}
