package keeper_test

import (
	"strings"
	"testing"

	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/oracle/keeper"
	"github.com/certikfoundation/shentu/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

var (
	acc1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc2 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc3 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc4 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
)

// shared setup
type KeeperTestSuite struct {
	suite.Suite

	// cdc    *codec.LegacyAmino
	ctx     sdk.Context
	app     *simapp.SimApp
	keeper  keeper.Keeper
	address []sdk.AccAddress
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = suite.app.OracleKeeper

	for _, acc := range []sdk.AccAddress{acc1, acc2, acc3, acc4} {
		err := suite.app.BankKeeper.AddCoins(
			suite.ctx,
			acc,
			sdk.NewCoins(
				sdk.NewCoin("uctk", sdk.NewInt(10000000)),
			),
		)
		if err != nil {
			panic(err)
		}
	}

	suite.address = []sdk.AccAddress{acc1, acc2, acc3, acc4}

}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestCreateOperator() {
	type args struct {
		params       types.LockedPoolParams
		collateral   sdk.Coins
		senderAddr   sdk.AccAddress
		proposerAddr sdk.AccAddress
		operatorName string
	}

	type errArgs struct {
		shouldPass bool
		contains   string
	}

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{"CreateOperator1",
			args{
				params:       suite.keeper.GetLockedPoolParams(suite.ctx),
				collateral:   sdk.Coins{sdk.NewInt64Coin("uctk", suite.keeper.GetLockedPoolParams(suite.ctx).MinimumCollateral)},
				senderAddr:   suite.address[0],
				proposerAddr: suite.address[1],
				operatorName: "Operator1",
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"CreateOperator2",
			args{
				params:       suite.keeper.GetLockedPoolParams(suite.ctx),
				collateral:   sdk.Coins{sdk.NewInt64Coin("uctk", suite.keeper.GetLockedPoolParams(suite.ctx).MinimumCollateral)},
				senderAddr:   suite.address[2],
				proposerAddr: suite.address[3],
				operatorName: "Operator2",
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			err := suite.keeper.CreateOperator(suite.ctx, tc.args.senderAddr, tc.args.collateral, tc.args.proposerAddr, tc.args.operatorName)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}

// TO-DO
// {"GetOperator"},
// {"GetAllOperators"},
// {"RemoveOperator"},
// ...
