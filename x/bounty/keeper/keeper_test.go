package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app       *shentuapp.ShentuApp
	ctx       sdk.Context
	addrs     []sdk.AccAddress
	msgServer types.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	suite.app = app
	suite.ctx = ctx
	suite.addrs = shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(1e10))
	suite.msgServer = keeper.NewMsgServerImpl(suite.app.BountyKeeper)

}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
