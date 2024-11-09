package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cometbft/cometbft/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
)

// shared setup
type KeeperTestSuite struct {
	suite.Suite
	app *shentuapp.ShentuApp
	ctx sdk.Context
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = shentuapp.Setup(suite.T(), false)
	suite.ctx = suite.app.BaseApp.NewContext(false)
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	suite.Require().NoError(err)
	coins := sdk.Coins{sdk.NewInt64Coin(bondDenom, 1e10)}
	suite.Require().NoError(testutil.FundModuleAccount(suite.ctx, suite.app.BankKeeper, minttypes.ModuleName, coins))
}

func (suite *KeeperTestSuite) TestKeeper_SendToCommunityPool() {
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	suite.Require().NoError(err)

	tests := []struct {
		name  string
		coins sdk.Coins
		err   bool
	}{
		{
			name:  "Funding Community Pool",
			coins: sdk.NewCoins(sdk.NewInt64Coin(bondDenom, 1e9)),
			err:   false,
		},
		{
			name:  "Funding With 0",
			coins: sdk.NewCoins(sdk.NewInt64Coin(bondDenom, 0)),
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
			initalMintBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAcct, bondDenom)
			distributionBalance1 := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(crypto.AddressHash([]byte(distrtypes.ModuleName))), bondDenom)

			err := suite.app.MintKeeper.SendToCommunityPool(suite.ctx, tc.coins)
			suite.Require().NoError(err)
			deductedMintBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAcct, bondDenom)
			distributionBalance2 := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(crypto.AddressHash([]byte(distrtypes.ModuleName))), bondDenom)
			suite.Require().Equal(initalMintBalance.Sub(deductedMintBalance), distributionBalance2.Sub(distributionBalance1))
		}
	}
}
