package keeper_test

import (
	certTypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *shentuapp.ShentuApp
	ctx         sdk.Context
	keeper      keeper.Keeper
	address     []sdk.AccAddress
	msgServer   types.MsgServer
	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = shentuapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = suite.app.BountyKeeper
	suite.address = shentuapp.AddTestAddrs(suite.app, suite.ctx, 4, sdk.NewInt(1e10))
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.BountyKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
	suite.msgServer = keeper.NewMsgServerImpl(suite.keeper)

	suite.app.CertKeeper.SetCertifier(suite.ctx, certTypes.NewCertifier(suite.address[3], "", suite.address[3], ""))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
