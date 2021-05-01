package keeper_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/oracle/keeper"
	"github.com/certikfoundation/shentu/x/oracle/types"
)

// ----------------- TO-DO ----------------- //
//
// ...
// ----------------------------------------- //

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
		{"One Operator: Create",
			args{
				params:       suite.keeper.GetLockedPoolParams(suite.ctx),
				collateral:   sdk.Coins{sdk.NewInt64Coin("uctk", suite.keeper.GetLockedPoolParams(suite.ctx).MinimumCollateral)},
				senderAddr:   suite.address[0],
				proposerAddr: suite.address[1],
				operatorName: "Operator",
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

func (suite *KeeperTestSuite) TestGetOperators() {
	type args struct {
		params        types.LockedPoolParams
		collateral    sdk.Coins
		senderAddr1   sdk.AccAddress
		proposerAddr1 sdk.AccAddress
		operatorName1 string
		senderAddr2   sdk.AccAddress
		proposerAddr2 sdk.AccAddress
		operatorName2 string
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
		{"Two Operators: Get One then Get All",
			args{
				params:        suite.keeper.GetLockedPoolParams(suite.ctx),
				collateral:    sdk.Coins{sdk.NewInt64Coin("uctk", suite.keeper.GetLockedPoolParams(suite.ctx).MinimumCollateral)},
				senderAddr1:   suite.address[0],
				proposerAddr1: suite.address[1],
				operatorName1: "Operator1",
				senderAddr2:   suite.address[2],
				proposerAddr2: suite.address[3],
				operatorName2: "Operator2",
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
			err := suite.keeper.CreateOperator(suite.ctx, tc.args.senderAddr1, tc.args.collateral, tc.args.proposerAddr1, tc.args.operatorName1)
			suite.Require().NoError(err, tc.name)
			err = suite.keeper.CreateOperator(suite.ctx, tc.args.senderAddr2, tc.args.collateral, tc.args.proposerAddr2, tc.args.operatorName2)
			suite.Require().NoError(err, tc.name)
			operator1, err := suite.keeper.GetOperator(suite.ctx, tc.args.senderAddr1)
			allOperators := suite.keeper.GetAllOperators(suite.ctx)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				suite.Equal(tc.args.senderAddr1.String(), operator1.Address)
				suite.Equal(tc.args.collateral, operator1.Collateral)
				suite.Equal(tc.args.proposerAddr1.String(), operator1.Proposer)
				suite.Len(allOperators, 2)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestRemoveOperator() {
	type args struct {
		params        types.LockedPoolParams
		collateral    sdk.Coins
		senderAddr1   sdk.AccAddress
		proposerAddr1 sdk.AccAddress
		operatorName1 string
		senderAddr2   sdk.AccAddress
		proposerAddr2 sdk.AccAddress
		operatorName2 string
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
		{"Two Operators: Remove One",
			args{
				params:        suite.keeper.GetLockedPoolParams(suite.ctx),
				collateral:    sdk.Coins{sdk.NewInt64Coin("uctk", suite.keeper.GetLockedPoolParams(suite.ctx).MinimumCollateral)},
				senderAddr1:   suite.address[0],
				proposerAddr1: suite.address[1],
				operatorName1: "Operator1",
				senderAddr2:   suite.address[2],
				proposerAddr2: suite.address[3],
				operatorName2: "Operator2",
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
			err := suite.keeper.CreateOperator(suite.ctx, tc.args.senderAddr1, tc.args.collateral, tc.args.proposerAddr1, tc.args.operatorName1)
			suite.Require().NoError(err, tc.name)
			err = suite.keeper.CreateOperator(suite.ctx, tc.args.senderAddr2, tc.args.collateral, tc.args.proposerAddr2, tc.args.operatorName2)
			suite.Require().NoError(err, tc.name)
			operator1, err := suite.keeper.GetOperator(suite.ctx, tc.args.senderAddr1)
			suite.Require().NoError(err, tc.name)
			// convert operator1.Address (string) back to sdk.AccAddress
			operator1Addr, _ := sdk.AccAddressFromBech32(operator1.Address)
			operator2, err := suite.keeper.GetOperator(suite.ctx, tc.args.senderAddr2)
			suite.Require().NoError(err, tc.name)
			// remove operator1
			err = suite.keeper.RemoveOperator(suite.ctx, operator1Addr)
			allOperators := suite.keeper.GetAllOperators(suite.ctx)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				suite.Len(allOperators, 1)
				suite.Equal(operator2, allOperators[0])
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}
