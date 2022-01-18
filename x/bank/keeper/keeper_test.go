package keeper_test

import (
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/certikfoundation/shentu/v2/simapp"
	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// shared setup
type KeeperTestSuite struct {
	suite.Suite
	app *simapp.SimApp
	ctx sdk.Context
}

var (
	acc1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc2 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc3 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
)

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	coins := sdk.Coins{sdk.NewInt64Coin("uctk", 80000*1e6)}
	suite.Require().NoError(sdksimapp.FundAccount(suite.app.BankKeeper, suite.ctx, acc1, coins))
}

func (suite *KeeperTestSuite) TestKeeper_SendCoins() {
	tests := []struct {
		name  string
		coins sdk.Coins
		err   bool
	}{
		{
			name:  "Transferring some coins",
			coins: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), 100*1e6)),
			err:   false,
		},
		{
			name:  "Transferring 0 coins",
			coins: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), 0)),
			err:   false,
		},
	}

	for _, tc := range tests {
		suite.T().Log(tc.name)
		if tc.err {
			err := suite.app.BankKeeper.SendCoins(suite.ctx, acc1, acc2, tc.coins)
			suite.Require().Nil(err)
		} else {
			initalBalance := suite.app.BankKeeper.GetBalance(suite.ctx, acc2, "uctk")
			err := suite.app.BankKeeper.SendCoins(suite.ctx, acc1, acc2, tc.coins)
			suite.Require().NoError(err)
			finalBalance := suite.app.BankKeeper.GetBalance(suite.ctx, acc2, "uctk")
			suite.Require().Equal(sdk.NewCoins(finalBalance.Sub(initalBalance)), tc.coins)
		}
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
