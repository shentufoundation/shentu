package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	shentuapp "github.com/certikfoundation/shentu/v2/app"
	"github.com/certikfoundation/shentu/v2/common"
	"github.com/certikfoundation/shentu/v2/x/shield/keeper"
	"github.com/certikfoundation/shentu/v2/x/shield/types"
	"github.com/certikfoundation/shentu/v2/x/staking/teststaking"
)

var (
	acc1       = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc2       = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc3       = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc4       = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	validator  = sdk.ValAddress{}
	PKS        = shentuapp.CreateTestPubKeys(5)
	valConsPk1 = PKS[0]
)

// shared setup
type KeeperTestSuite struct {
	suite.Suite
	app       *shentuapp.ShentuApp
	ctx       sdk.Context
	msgServer types.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = shentuapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	suite.msgServer = keeper.NewMsgServerImpl(suite.app.ShieldKeeper)
	coins := sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 80000*1e6)}
	suite.Require().NoError(sdksimapp.FundModuleAccount(suite.app.BankKeeper, suite.ctx, types.ModuleName, coins))
	suite.Require().NoError(sdksimapp.FundAccount(suite.app.BankKeeper, suite.ctx, acc2, coins))
	// validator set up
	tstaking := teststaking.NewHelper(suite.T(), suite.ctx, suite.app.StakingKeeper)
	addr := shentuapp.AddTestAddrs(suite.app, suite.ctx, 2, sdk.NewInt(1000000000))
	valAddrs := sdksimapp.ConvertAddrsToValAddrs(addr)
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[0], valConsPk1, 100, true)
	suite.app.StakingKeeper.Delegation(suite.ctx, sdk.AccAddress(valAddrs[0]), valAddrs[0])
	acc1 = sdk.AccAddress(valAddrs[0])
	validator = valAddrs[0]
	suite.app.ShieldKeeper.SetAdmin(suite.ctx, acc1)
}

func (suite *KeeperTestSuite) Test_DepositCollateral() {
	tests := []struct {
		name    string
		message *types.MsgDepositCollateral
		err     bool
	}{
		{
			name:    "Should fail for insufficient delegation",
			message: types.NewMsgDepositCollateral(acc1, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 80000*1e6)}),
			err:     true,
		},
		{
			name:    "Should fail for non validator",
			message: types.NewMsgDepositCollateral(acc2, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 1)}),
			err:     true,
		},
		{
			name:    "Should pass for sufficient delegation",
			message: types.NewMsgDepositCollateral(acc1, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 1)}),
			err:     false,
		},
	}

	for _, tc := range tests {
		suite.T().Log(tc.name)
		initCollateral := suite.app.ShieldKeeper.GetTotalCollateral(suite.ctx)
		_, err := suite.msgServer.DepositCollateral(sdk.WrapSDKContext(suite.ctx), tc.message)
		finalCollateral := suite.app.ShieldKeeper.GetTotalCollateral(suite.ctx)
		if tc.err {
			suite.Require().Error(err)
		} else {
			suite.Require().NoError(err)
			suite.Require().Truef(finalCollateral.Sub(initCollateral).Equal(tc.message.Collateral.AmountOf(common.MicroCTKDenom)), "unexpected err")
		}
	}
}

func (suite *KeeperTestSuite) Test_WithdrawCollateral() {
	_, err := suite.msgServer.DepositCollateral(sdk.WrapSDKContext(suite.ctx), types.NewMsgDepositCollateral(acc1, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 10)}))
	suite.Require().NoError(err)
	tests := []struct {
		name    string
		message *types.MsgWithdrawCollateral
		err     bool
	}{
		{
			name:    "Should fail for insufficient collateral deposit",
			message: types.NewMsgWithdrawCollateral(acc1, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 80000*1e6)}),
			err:     true,
		},
		{
			name:    "Should fail for non validator",
			message: types.NewMsgWithdrawCollateral(acc2, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 1)}),
			err:     true,
		},
		{
			name:    "Should pass for sufficient collateral",
			message: types.NewMsgWithdrawCollateral(acc1, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 1)}),
			err:     false,
		},
	}

	for _, tc := range tests {
		suite.T().Log(tc.name)
		initCollateral := suite.app.ShieldKeeper.GetTotalWithdrawing(suite.ctx)
		_, err := suite.msgServer.WithdrawCollateral(sdk.WrapSDKContext(suite.ctx), tc.message)
		finalCollateral := suite.app.ShieldKeeper.GetTotalWithdrawing(suite.ctx)
		if tc.err {
			suite.Require().Error(err)
		} else {
			suite.Require().NoError(err)
			suite.Require().Truef(finalCollateral.Sub(initCollateral).Equal(tc.message.Collateral.AmountOf(common.MicroCTKDenom)), "unexpected err")
		}
	}
}

func (suite *KeeperTestSuite) Test_Donate() {
	tests := []struct {
		name    string
		message *types.MsgDonate
		err     bool
	}{
		{
			name:    "Should fail for insufficient balance",
			message: types.NewMsgDonate(acc1, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 80000*1e6)}),
			err:     true,
		},
		{
			name:    "Should pass for non validator",
			message: types.NewMsgDonate(acc2, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 1)}),
			err:     false,
		},
		{
			name:    "Should pass for sufficient collateral",
			message: types.NewMsgDonate(acc1, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 1)}),
			err:     false,
		},
	}

	for _, tc := range tests {
		suite.T().Log(tc.name)
		initCollateral := suite.app.ShieldKeeper.GetReserve(suite.ctx).Amount
		_, err := suite.msgServer.Donate(sdk.WrapSDKContext(suite.ctx), tc.message)
		finalCollateral := suite.app.ShieldKeeper.GetReserve(suite.ctx).Amount
		if tc.err {
			suite.Require().Error(err)
		} else {
			suite.Require().NoError(err)
			suite.Require().Truef(finalCollateral.Sub(initCollateral).Equal(tc.message.Amount.AmountOf(common.MicroCTKDenom)), "unexpected err")
		}
	}
}

func (suite *KeeperTestSuite) Test_CreatePool() {
	tests := []struct {
		name    string
		message *types.MsgCreatePool
		err     bool
	}{
		{
			name:    "Should create a pool",
			message: types.NewMsgCreatePool(acc1, acc1, "pool creation", sdk.NewDec(1), sdk.NewInt(1e14)),
			err:     false,
		},
		{
			name:    "Should pass for non validator",
			message: types.NewMsgCreatePool(acc2, acc2, "pool creation", sdk.NewDec(2), sdk.NewInt(1e14)),
			err:     true,
		},
		{
			name:    "Should fail for duplicate pool",
			message: types.NewMsgCreatePool(acc1, acc1, "pool creation", sdk.NewDec(1), sdk.NewInt(1e14)),
			err:     true,
		},
		{
			name:    "Should create pool for others",
			message: types.NewMsgCreatePool(acc1, acc2, "pool creation", sdk.NewDec(1), sdk.NewInt(1e14)),
			err:     false,
		},
	}

	for _, tc := range tests {
		suite.T().Log(tc.name)
		_, err := suite.msgServer.CreatePool(sdk.WrapSDKContext(suite.ctx), tc.message)
		_, check := suite.app.ShieldKeeper.GetPool(suite.ctx, 1)
		if tc.err {
			suite.Require().Error(err)
		} else {
			suite.Require().NoError(err)
			suite.Require().Truef(check, "should contain the pool")
		}
	}
}

func (suite *KeeperTestSuite) Test_UpdatePool() {
	_, err := suite.msgServer.CreatePool(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreatePool(acc1, acc1, "pool creation", sdk.NewDec(1), sdk.NewInt(1e14)))
	suite.Require().NoError(err)
	tests := []struct {
		name    string
		message *types.MsgUpdatePool
		err     bool
	}{
		{
			name:    "Should update a pool",
			message: types.NewMsgUpdatePool(acc1, 1, "pool updated", true, sdk.NewDec(1), sdk.NewInt(1e14)),
			err:     false,
		},
		{
			name:    "Should fail for different sponsor",
			message: types.NewMsgUpdatePool(acc2, 1, "pool updated", true, sdk.NewDec(1), sdk.NewInt(1e14)),
			err:     true,
		},
		{
			name:    "Should fail for non existent pool",
			message: types.NewMsgUpdatePool(acc2, 5, "pool updated", true, sdk.NewDec(1), sdk.NewInt(1e14)),
			err:     true,
		},
	}

	for _, tc := range tests {
		suite.T().Log(tc.name)
		_, err := suite.msgServer.UpdatePool(sdk.WrapSDKContext(suite.ctx), tc.message)
		pool, check := suite.app.ShieldKeeper.GetPool(suite.ctx, 1)
		if tc.err {
			suite.Require().Error(err)
		} else {
			suite.Require().NoError(err)
			suite.Require().Truef(pool.Description == "pool updated" && check, "should contain and update the pool")
		}
	}
}

func (suite *KeeperTestSuite) Test_PausePool() {
	_, err := suite.msgServer.CreatePool(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreatePool(acc1, acc1, "pool creation", sdk.NewDec(1), sdk.NewInt(1e14)))
	suite.Require().NoError(err)
	tests := []struct {
		name    string
		message *types.MsgPausePool
		err     bool
	}{
		{
			name:    "Should pause a pool",
			message: types.NewMsgPausePool(acc1, 1),
			err:     false,
		},
		{
			name:    "Should fail for different sponsor",
			message: types.NewMsgPausePool(acc2, 1),
			err:     true,
		},
		{
			name:    "Should fail for non existent pool",
			message: types.NewMsgPausePool(acc2, 5),
			err:     true,
		},
	}

	for _, tc := range tests {
		suite.T().Log(tc.name)
		_, err := suite.msgServer.PausePool(sdk.WrapSDKContext(suite.ctx), tc.message)
		pool, check := suite.app.ShieldKeeper.GetPool(suite.ctx, 1)
		if tc.err {
			suite.Require().Error(err)
		} else {
			suite.Require().NoError(err)
			suite.Require().Truef(!pool.Active && check, "should contain and pause the pool")
		}
	}
}

func (suite *KeeperTestSuite) Test_ResumePool() {
	_, err := suite.msgServer.CreatePool(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreatePool(acc1, acc1, "pool creation", sdk.NewDec(1), sdk.NewInt(1e14)))
	suite.Require().NoError(err)
	_, err = suite.msgServer.PausePool(sdk.WrapSDKContext(suite.ctx), types.NewMsgPausePool(acc1, 1))
	suite.Require().NoError(err)
	tests := []struct {
		name    string
		message *types.MsgResumePool
		err     bool
	}{
		{
			name:    "Should pause a pool",
			message: types.NewMsgResumePool(acc1, 1),
			err:     false,
		},
		{
			name:    "Should fail for different sponsor",
			message: types.NewMsgResumePool(acc2, 1),
			err:     true,
		},
		{
			name:    "Should fail for non existent pool",
			message: types.NewMsgResumePool(acc2, 5),
			err:     true,
		},
	}

	for _, tc := range tests {
		suite.T().Log(tc.name)
		_, err := suite.msgServer.ResumePool(sdk.WrapSDKContext(suite.ctx), tc.message)
		pool, check := suite.app.ShieldKeeper.GetPool(suite.ctx, 1)
		if tc.err {
			suite.Require().Error(err)
		} else {
			suite.Require().NoError(err)
			suite.Require().Truef(pool.Active && check, "should contain and resume the pool")
		}
	}
}

func (suite *KeeperTestSuite) Test_UpdateSponsor() {
	_, err := suite.msgServer.CreatePool(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreatePool(acc1, acc1, "pool creation", sdk.NewDec(1), sdk.NewInt(1e14)))
	suite.Require().NoError(err)
	tests := []struct {
		name    string
		message *types.MsgUpdateSponsor
		err     bool
	}{
		{
			name:    "Should fail for invalid from",
			message: types.NewMsgUpdateSponsor(1, "new sponsor", acc2, acc2),
			err:     true,
		},
		{
			name:    "Should fail for non existent pool",
			message: types.NewMsgUpdateSponsor(1, "new sponsor", acc1, acc1),
			err:     false,
		},
		{
			name:    "Should update a sponsor",
			message: types.NewMsgUpdateSponsor(1, "new sponsor", acc2, acc1),
			err:     false,
		},
	}

	for _, tc := range tests {
		suite.T().Log(tc.name)
		_, err := suite.msgServer.UpdateSponsor(sdk.WrapSDKContext(suite.ctx), tc.message)
		pool, check := suite.app.ShieldKeeper.GetPool(suite.ctx, 1)
		if tc.err {
			suite.Require().Error(err)
		} else {
			suite.Require().NoError(err)
			suite.Require().Truef(pool.SponsorAddr == tc.message.SponsorAddr && check, "should contain and update sponsor")
		}
	}
}

func (suite *KeeperTestSuite) Test_Purchase() {
	_, err := suite.msgServer.CreatePool(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreatePool(acc1, acc1, "pool creation", sdk.NewDec(1), sdk.NewInt(1e14)))
	suite.Require().NoError(err)
	tests := []struct {
		name    string
		message *types.MsgPurchase
		err     bool
	}{
		{
			name:    "Should pass for purchase",
			message: types.NewMsgPurchase(1, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 100000000)}, "purchasing", acc1),
			err:     false,
		},
		{
			name:    "Should fail for small purchase",
			message: types.NewMsgPurchase(1, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 10)}, "purchasing", acc1),
			err:     true,
		},
		{
			name:    "Should pass for other than spnosor purchase",
			message: types.NewMsgPurchase(1, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 100000000)}, "purchasing", acc2),
			err:     false,
		},
	}

	for _, tc := range tests {
		suite.T().Log(tc.name)
		_, err := suite.msgServer.Purchase(sdk.WrapSDKContext(suite.ctx), tc.message)
		purchase, check := suite.app.ShieldKeeper.GetPurchase(suite.ctx, 1, acc1)
		if tc.err {
			suite.Require().Error(err)
		} else {
			suite.Require().NoError(err)
			suite.Require().Truef(purchase.Amount.Equal(tc.message.Amount.AmountOf(common.MicroCTKDenom)) && check, "should contain and purchase")
		}
	}
}

func (suite *KeeperTestSuite) Test_Unstake() {
	_, err := suite.msgServer.CreatePool(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreatePool(acc1, acc1, "pool creation", sdk.NewDec(1), sdk.NewInt(1e14)))
	suite.Require().NoError(err)
	_, err = suite.msgServer.Purchase(sdk.WrapSDKContext(suite.ctx), types.NewMsgPurchase(1, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 100000000)}, "purchasing", acc1))
	suite.Require().NoError(err)

	tests := []struct {
		name    string
		message *types.MsgUnstake
		err     bool
	}{
		{
			name:    "Should pass for purchase",
			message: types.NewMsgUnstake(1, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 10)}, acc1),
			err:     false,
		},
		{
			name:    "Should fail for unstaking more",
			message: types.NewMsgUnstake(1, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 100000000)}, acc2),
			err:     true,
		},
		{
			name:    "Should fail for absent pool",
			message: types.NewMsgUnstake(3, sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 100000000)}, acc1),
			err:     true,
		},
	}

	for _, tc := range tests {
		suite.T().Log(tc.name)
		initialPurchase, _ := suite.app.ShieldKeeper.GetPurchase(suite.ctx, 1, acc1)
		_, err := suite.msgServer.Unstake(sdk.WrapSDKContext(suite.ctx), tc.message)
		finalPurchase, check := suite.app.ShieldKeeper.GetPurchase(suite.ctx, 1, acc1)
		if tc.err {
			suite.Require().Error(err)
		} else {
			suite.Require().NoError(err)
			suite.Require().Truef(initialPurchase.Amount.Sub(finalPurchase.Amount).Equal(tc.message.Amount.AmountOf(common.MicroCTKDenom)) && check, "unsuccessful unstaking")
		}
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
