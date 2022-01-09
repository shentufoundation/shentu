package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/simapp"
	"github.com/tendermint/tendermint/crypto"
)

// shared setup
type KeeperTestSuite struct {
	suite.Suite
	app *simapp.SimApp
	ctx sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	coins := sdk.Coins{sdk.NewInt64Coin("uctk", 80000*1e6)}
	suite.Require().NoError(sdksimapp.FundModuleAccount(suite.app.BankKeeper, suite.ctx, "mint", coins))
}

func (suite *KeeperTestSuite) TestKeeper_SendToCommunityPool() {
	tests := []struct {
		name  string
		coins sdk.Coins
		err   bool
	}{
		{
			name:  "Funding Community Pool",
			coins: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), 100*1e6)),
			err:   false,
		},
		{
			name:  "Funding With 0",
			coins: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), 0)),
			err:   true,
		},
	}

	for _, tc := range tests {
		suite.T().Log(tc.name)
		moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte("mint")))
		if tc.err {
			err := suite.app.MintKeeper.SendToCommunityPool(suite.ctx, tc.coins)
			suite.Require().Nil(err)
		} else {
			initalMintBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAcct, "uctk")
			err := suite.app.MintKeeper.SendToCommunityPool(suite.ctx, tc.coins)
			suite.Require().NoError(err)
			deductedMintBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAcct, "uctk")
			distributionBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(crypto.AddressHash([]byte("distribution"))), "uctk")
			suite.Require().Equal(initalMintBalance.Sub(deductedMintBalance), distributionBalance)
		}
	}
}

func (suite *KeeperTestSuite) TestKeeper_SendToShieldRewards() {
	tests := []struct {
		name  string
		coins sdk.Coins
		err   bool
	}{
		{
			name:  "Funding Shield Rewards",
			coins: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), 100*1e6)),
			err:   false,
		},
		{
			name:  "Funding With 0",
			coins: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), 0)),
			err:   true,
		},
	}

	for _, tc := range tests {
		suite.T().Log(tc.name)
		moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte("mint")))
		if tc.err {
			err := suite.app.MintKeeper.SendToShieldRewards(suite.ctx, tc.coins)
			suite.Require().Nil(err)
		} else {
			initalMintBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAcct, "uctk")
			err := suite.app.MintKeeper.SendToShieldRewards(suite.ctx, tc.coins)
			suite.Require().NoError(err)
			deductedMintBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAcct, "uctk")
			shieldBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(crypto.AddressHash([]byte("shield"))), "uctk")
			suite.Require().Equal(initalMintBalance.Sub(deductedMintBalance), shieldBalance)
		}
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
