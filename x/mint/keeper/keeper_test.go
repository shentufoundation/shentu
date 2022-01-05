package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/simapp"
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
		err := suite.app.MintKeeper.SendToCommunityPool(suite.ctx, tc.coins)
		if tc.err {
			suite.Require().Nil(err)
		}
		suite.Require().NoError(err)
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
		err := suite.app.MintKeeper.SendToShieldRewards(suite.ctx, tc.coins)
		if tc.err {
			suite.Require().Nil(err)
		}
		suite.Require().NoError(err)
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
