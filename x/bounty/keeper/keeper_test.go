package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
	certTypes "github.com/shentufoundation/shentu/v2/x/cert/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *shentuapp.ShentuApp
	ctx         sdk.Context
	keeper      keeper.Keeper
	address     []sdk.AccAddress
	msgServer   types.MsgServer
	queryClient types.QueryClient

	// addr type
	programAddr, whiteHatAddr, bountyAdminAddr, normalAddr sdk.AccAddress
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = shentuapp.Setup(suite.T(), false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = suite.app.BountyKeeper
	suite.address = shentuapp.AddTestAddrs(suite.app, suite.ctx, 5, sdk.NewInt(1e10))
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.BountyKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
	suite.msgServer = keeper.NewMsgServerImpl(suite.keeper)

	suite.app.CertKeeper.SetCertifier(suite.ctx, certTypes.NewCertifier(suite.address[2], "", suite.address[2], ""))
	certificate, err := certTypes.NewCertificate("bountyadmin", suite.address[3].String(), "", "", "", suite.address[2])
	if err != nil {
		panic(err)
	}
	id := suite.app.CertKeeper.GetNextCertificateID(suite.ctx)
	certificate.CertificateId = id
	// set the cert and its ID in the store
	suite.app.CertKeeper.SetNextCertificateID(suite.ctx, id+1)
	suite.app.CertKeeper.SetCertificate(suite.ctx, certificate)

	suite.programAddr = suite.address[0]
	suite.whiteHatAddr = suite.address[1]
	// suite.address[2] is certifier addr
	suite.bountyAdminAddr = suite.address[3]
	suite.normalAddr = suite.address[4]
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
