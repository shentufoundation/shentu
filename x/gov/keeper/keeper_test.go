package keeper_test

import (
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/common"
	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
)

var address1 = "cosmos1ghekyjucln7y67ntx7cf27m9dpuxxemn4c8g4r"

// shared setup
type KeeperTestSuite struct {
	suite.Suite

	app               *shentuapp.ShentuApp
	ctx               sdk.Context
	queryClient       govtypesv1.QueryClient
	legacyQueryClient govtypesv1beta1.QueryClient
	addrs             []sdk.AccAddress
	msgSrvr           govtypesv1.MsgServer
	legacyMsgSrvr     govtypesv1beta1.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	app := shentuapp.Setup(suite.T(), false)
	ctx := app.BaseApp.NewContext(false)

	// Populate the gov account with some coins, as the TestProposal we have
	// is a MsgSend from the gov account.
	coins := sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, math.NewInt(1e10)))
	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	suite.NoError(err)
	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, govtypes.ModuleName, coins)
	suite.NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	//govtypesv1.RegisterQueryServer(queryHelper, app.GovKeeper.Keeper)
	legacyQueryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	//govtypesv1beta1.RegisterQueryServer(legacyQueryHelper, keeper.NewLegacyQueryServer(app.GovKeeper))
	queryClient := govtypesv1.NewQueryClient(queryHelper)
	legacyQueryClient := govtypesv1beta1.NewQueryClient(legacyQueryHelper)

	suite.app = app
	suite.ctx = ctx
	suite.queryClient = queryClient
	suite.legacyQueryClient = legacyQueryClient
	suite.msgSrvr = keeper.NewMsgServerImpl(suite.app.GovKeeper)
	govAcct := suite.app.GovKeeper.GetGovernanceAccount(suite.ctx).GetAddress()
	suite.legacyMsgSrvr = keeper.NewLegacyMsgServerImpl(govAcct.String(), suite.msgSrvr, app.GovKeeper)
	suite.addrs = shentuapp.AddTestAddrsIncremental(app, ctx, 2, math.NewInt(1e10))
}

func TestIncrementProposalNumber(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	addrBz, err := app.AccountKeeper.AddressCodec().StringToBytes(address1)
	require.NoError(t, err)

	tp := TestProposal
	_, err = app.GovKeeper.SubmitProposal(ctx, tp, "", "test", "summary", addrBz, false)
	require.NoError(t, err)
	_, err = app.GovKeeper.SubmitProposal(ctx, tp, "", "test", "summary", addrBz, false)
	require.NoError(t, err)
	_, err = app.GovKeeper.SubmitProposal(ctx, tp, "", "test", "summary", addrBz, false)
	require.NoError(t, err)
	_, err = app.GovKeeper.SubmitProposal(ctx, tp, "", "test", "summary", addrBz, false)
	require.NoError(t, err)
	_, err = app.GovKeeper.SubmitProposal(ctx, tp, "", "test", "summary", addrBz, false)
	require.NoError(t, err)
	proposal6, err := app.GovKeeper.SubmitProposal(ctx, tp, "", "test", "summary", addrBz, false)
	require.NoError(t, err)

	require.Equal(t, uint64(6), proposal6.Id)
}

func TestProposalQueues(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	addrBz, err := app.AccountKeeper.AddressCodec().StringToBytes(address1)
	require.NoError(t, err)

	// create test proposals
	tp := TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp, "", "test", "summary", addrBz, false)
	require.NoError(t, err)

	has, err := app.GovKeeper.InactiveProposalsQueue.Has(ctx, collections.Join(*proposal.DepositEndTime, proposal.Id))
	require.NoError(t, err)
	require.True(t, has)

	require.NoError(t, app.GovKeeper.ActivateVotingPeriod(ctx, proposal))

	proposal, err = app.GovKeeper.Proposals.Get(ctx, proposal.Id)
	require.Nil(t, err)

	has, err = app.GovKeeper.ActiveProposalsQueue.Has(ctx, collections.Join(*proposal.VotingEndTime, proposal.Id))
	require.NoError(t, err)
	require.True(t, has)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
