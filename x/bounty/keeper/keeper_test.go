package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
	suite.ctx = suite.app.BaseApp.NewContext(false)
	suite.keeper = suite.app.BountyKeeper
	suite.address = shentuapp.AddTestAddrs(suite.app, suite.ctx, 5, math.NewInt(1e10))
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServer(suite.app.BountyKeeper))
	suite.queryClient = types.NewQueryClient(queryHelper)
	suite.msgServer = keeper.NewMsgServerImpl(suite.keeper)

	err := suite.app.CertKeeper.SetCertifier(suite.ctx, certTypes.NewCertifier(suite.address[2], "", suite.address[2], ""))
	suite.Require().NoError(err)
	certificate, err := certTypes.NewCertificate("bountyadmin", suite.address[3].String(), "", "", "", suite.address[2])
	if err != nil {
		panic(err)
	}
	id, _ := suite.app.CertKeeper.GetNextCertificateID(suite.ctx)
	certificate.CertificateId = id
	// set the cert and its ID in the store
	err = suite.app.CertKeeper.SetNextCertificateID(suite.ctx, id+1)
	suite.Require().NoError(err)
	err = suite.app.CertKeeper.SetCertificate(suite.ctx, certificate)
	suite.Require().NoError(err)
	suite.programAddr = suite.address[0]
	suite.whiteHatAddr = suite.address[1]
	// suite.address[2] is certifier addr
	suite.bountyAdminAddr = suite.address[3]
	suite.normalAddr = suite.address[4]

	// Set module parameters
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	suite.Require().NoError(err)

	minGrant := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(50)))
	minDeposit := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(30)))
	theoremMaxProofPeriod := 14 * 24 * time.Hour
	proofMaxLockPeriod := 10 * time.Minute
	checkerRate := math.LegacyMustNewDecFromStr("0.2") // 0.2

	params := types.Params{
		MinGrant:              minGrant,
		MinDeposit:            minDeposit,
		TheoremMaxProofPeriod: &theoremMaxProofPeriod,
		ProofMaxLockPeriod:    &proofMaxLockPeriod,
		CheckerRate:           checkerRate,
	}
	err = suite.keeper.Params.Set(suite.ctx, params)
	suite.Require().NoError(err)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
