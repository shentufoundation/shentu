package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/common"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

// shared setup
type KeeperTestSuite struct {
	suite.Suite
	app *shentuapp.ShentuApp
	ctx sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = shentuapp.Setup(suite.T(), false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	coins := sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 1e10)}
	suite.Require().NoError(testutil.FundModuleAccount(suite.app.BankKeeper, suite.ctx, minttypes.ModuleName, coins))
}

func (suite *KeeperTestSuite) TestKeeper_SendToCommunityPool() {
	tests := []struct {
		name  string
		coins sdk.Coins
		err   bool
	}{
		{
			name:  "Funding Community Pool",
			coins: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), 1e9)),
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
		moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(minttypes.ModuleName)))
		if tc.err {
			err := suite.app.MintKeeper.SendToCommunityPool(suite.ctx, tc.coins)
			suite.Require().Nil(err)
		} else {
			initalMintBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAcct, common.MicroCTKDenom)
			distributionBalance1 := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(crypto.AddressHash([]byte(distrtypes.ModuleName))), common.MicroCTKDenom)

			err := suite.app.MintKeeper.SendToCommunityPool(suite.ctx, tc.coins)
			suite.Require().NoError(err)
			deductedMintBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAcct, common.MicroCTKDenom)
			distributionBalance2 := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(crypto.AddressHash([]byte(distrtypes.ModuleName))), common.MicroCTKDenom)
			suite.Require().Equal(initalMintBalance.Sub(deductedMintBalance), distributionBalance2.Sub(distributionBalance1))
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
			coins: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), 1e9)),
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
		moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(minttypes.ModuleName)))
		if tc.err {
			err := suite.app.MintKeeper.SendToShieldRewards(suite.ctx, tc.coins)
			suite.Require().Nil(err)
		} else {
			initalMintBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAcct, common.MicroCTKDenom)
			shieldBalance1 := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(crypto.AddressHash([]byte(shieldtypes.ModuleName))), common.MicroCTKDenom)
			err := suite.app.MintKeeper.SendToShieldRewards(suite.ctx, tc.coins)
			suite.Require().NoError(err)
			deductedMintBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAcct, common.MicroCTKDenom)
			shieldBalance2 := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(crypto.AddressHash([]byte(shieldtypes.ModuleName))), common.MicroCTKDenom)
			suite.Require().Equal(initalMintBalance.Sub(deductedMintBalance), shieldBalance2.Sub(shieldBalance1))
		}
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
