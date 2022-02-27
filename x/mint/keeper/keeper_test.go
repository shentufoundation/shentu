package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/crypto"

	shentuapp "github.com/certikfoundation/shentu/v2/app"
)

// shared setup
type KeeperTestSuite struct {
	suite.Suite
	app *shentuapp.ShentuApp
	ctx sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = shentuapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	coins := sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 80000*1e6)}
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
			initalMintBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAcct, common.MicroCTKDenom)
			err := suite.app.MintKeeper.SendToCommunityPool(suite.ctx, tc.coins)
			suite.Require().NoError(err)
			deductedMintBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAcct, common.MicroCTKDenom)
			distributionBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(crypto.AddressHash([]byte("distribution"))), common.MicroCTKDenom)
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
			initalMintBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAcct, common.MicroCTKDenom)
			err := suite.app.MintKeeper.SendToShieldRewards(suite.ctx, tc.coins)
			suite.Require().NoError(err)
			deductedMintBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAcct, common.MicroCTKDenom)
			shieldBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(crypto.AddressHash([]byte("shield"))), common.MicroCTKDenom)
			suite.Require().Equal(initalMintBalance.Sub(deductedMintBalance), shieldBalance)
		}
	}
}

func (suite *KeeperTestSuite) TestGetPoolMint() {
	tests := []struct {
		name       string
		ratio      sdk.Dec
		mintedCoin sdk.Coin
		expCoins   sdk.Coins
	}{
		{
			name:       "Get Pool Mint with Zero Ratio and Zero Minted Coin",
			ratio:      sdk.ZeroDec(), // 0%
			mintedCoin: sdk.NewCoin(common.MicroCTKDenom, sdk.ZeroInt()),
			expCoins:   sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, sdk.ZeroInt())},
		},
		{
			name:       "Get Pool Mint with Zero Ratio",
			ratio:      sdk.ZeroDec(), // 0%
			mintedCoin: sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(1e4)),
			expCoins:   sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, sdk.ZeroInt())},
		},
		{
			name:       "Get Pool Mint with Zero Minted Coin",
			ratio:      sdk.NewDecWithPrec(20, 2), // 20%
			mintedCoin: sdk.NewCoin(common.MicroCTKDenom, sdk.ZeroInt()),
			expCoins:   sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, sdk.ZeroInt())},
		},
		{
			name:       "Get Pool Mint with 100% Ratio",
			ratio:      sdk.NewDecWithPrec(100, 2), // 100%
			mintedCoin: sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(1e4)),
			expCoins:   sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(1e4))},
		},
		{
			name:       "Get Pool Mint with No Remainder",
			ratio:      sdk.NewDecWithPrec(20, 2), // 20%
			mintedCoin: sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(1e4)),
			expCoins:   sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(2e3))},
		},
		{
			name:       "Get Pool Mint with Remainder 1",
			ratio:      sdk.NewDecWithPrec(4928, 4), // 49.28%
			mintedCoin: sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(100)),
			expCoins:   sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(49))},
		},
		{
			name:       "Get Pool Mint with Remainder 2",
			ratio:      sdk.NewDecWithPrec(20, 2), // 20%
			mintedCoin: sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(9999)),
			expCoins:   sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(1999))},
		},
	}

	for _, tc := range tests {
		suite.T().Log(tc.name)
		poolMint := suite.app.MintKeeper.GetPoolMint(suite.ctx, tc.ratio, tc.mintedCoin)
		suite.Require().True(tc.expCoins.IsEqual(poolMint))
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
