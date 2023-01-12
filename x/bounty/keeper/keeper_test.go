package keeper_test

import (
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
	addrs       []sdk.AccAddress
	queryClient types.QueryClient
	msgServer   types.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.BountyKeeper)

	suite.app = app
	suite.ctx = ctx
	suite.addrs = shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(1e10))
	suite.queryClient = types.NewQueryClient(queryHelper)
	suite.msgServer = keeper.NewMsgServerImpl(suite.app.BountyKeeper)

}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
