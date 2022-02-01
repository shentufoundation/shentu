package types_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/simapp"
	"github.com/certikfoundation/shentu/v2/x/bank/types"
)

var (
	acc1 = sdk.AccAddress([]byte("input1"))
	acc2 = sdk.AccAddress([]byte("input2"))
	acc3 = sdk.AccAddress([]byte("input3"))
	acc4 = sdk.AccAddress([]byte("input4"))
)

// shared setup
type TypesTestSuite struct {
	suite.Suite

	address []sdk.AccAddress
	app     *simapp.SimApp
	ctx     sdk.Context
	params  types.AccountKeeper
}

func (suite *TypesTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.params = suite.app.AccountKeeper

	suite.address = []sdk.AccAddress{acc1, acc2, acc3, acc4}
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

func (suite *TypesTestSuite) TestMsgSendRoute() {
	type args struct {
		expectedRoute   string
		expectedType    string
		fromAddr        sdk.AccAddress
		toAddr          sdk.AccAddress
		unlockerAddress sdk.AccAddress
		amount          int64
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
		{"Bank Module: Valid Route and Type",
			args{
				amount:          200,
				fromAddr:        suite.address[0],
				toAddr:          suite.address[1],
				unlockerAddress: suite.address[2],
				expectedRoute:   "bank",
				expectedType:    "locked_send",
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Bank Module: Invalid Route and Type",
			args{
				amount:          110,
				fromAddr:        suite.address[0],
				toAddr:          suite.address[1],
				unlockerAddress: suite.address[2],
				expectedRoute:   "auth",
				expectedType:    "locked_coins",
			},
			errArgs{
				shouldPass: false,
				contains:   "",
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			coins := sdk.NewCoins(sdk.NewInt64Coin("uctk", tc.args.amount))
			var msg = types.NewMsgLockedSend(tc.args.fromAddr, tc.args.toAddr, tc.args.unlockerAddress.String(), coins)
			if tc.errArgs.shouldPass {
				suite.Require().Equal(msg.Route(), tc.args.expectedRoute)
				suite.Require().Equal(msg.Type(), tc.args.expectedType)
			} else {

			}
		})
	}
}

func (suite *TypesTestSuite) TestMsgSendValidation() {
	type args struct {
		fromAddr        sdk.AccAddress
		toAddr          sdk.AccAddress
		unlockerAddress sdk.AccAddress
		amount          int64
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
		{"Valid recipient address",
			args{
				amount:          200,
				fromAddr:        suite.address[0],
				toAddr:          suite.address[1],
				unlockerAddress: suite.address[2],
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
			coins := sdk.NewCoins(sdk.NewInt64Coin("uctk", tc.args.amount))
			var msg = types.NewMsgLockedSend(tc.args.fromAddr, tc.args.toAddr, tc.args.unlockerAddress.String(), coins)
			err := msg.ValidateBasic()
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().Nil(err)
			} else {
				suite.Require().NotNil(err)
			}
		})
	}
}

func (suite *TypesTestSuite) TestMsgSendGetSignBytes() {
	type args struct {
		fromAddr        sdk.AccAddress
		toAddr          sdk.AccAddress
		unlockerAddress sdk.AccAddress
		amount          int64
		res             string
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
		{"Expected Response",
			args{
				amount:          200,
				fromAddr:        suite.address[0],
				toAddr:          suite.address[1],
				unlockerAddress: suite.address[2],
				res:             `{"type":"bank/MsgLockedSend","value":{"amount":[{"amount":"200","denom":"uctk"}],"from_address":"cosmos1d9h8qat5xyj6yfmj","to_address":"cosmos1d9h8qat5xgryzr24","unlocker_address":"cosmos1d9h8qat5xvvwq990"}}`,
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Unexpected Response",
			args{
				amount:          2000,
				fromAddr:        suite.address[0],
				toAddr:          suite.address[1],
				unlockerAddress: suite.address[2],
				res:             `{"type":"bank/MsgLockedSend","value":{"amount":[{"amount":"12000","denom":"uctk"}],"from_address":"cosmos1d9h8qat5xyj6yfmj","to_address":"cosmos1d9h8qat5xgryzr24","unlocker_address":"cosmos1d9h8qat5xvvwq990"}}`,
			},
			errArgs{
				shouldPass: false,
				contains:   "",
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			coins := sdk.NewCoins(sdk.NewInt64Coin("uctk", tc.args.amount))
			var msg = types.NewMsgLockedSend(tc.args.fromAddr, tc.args.toAddr, tc.args.unlockerAddress.String(), coins)
			res := msg.GetSignBytes()
			if tc.errArgs.shouldPass {
				suite.Require().Equal(tc.args.res, string(res))
			} else {
				suite.Require().NotEqual(tc.args.res, string(res))
			}
		})
	}
}

func (suite *TypesTestSuite) TestMsgSendGetSigners() {
	type args struct {
		fromAddr        sdk.AccAddress
		toAddr          sdk.AccAddress
		unlockerAddress sdk.AccAddress
		signature       string
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
		{"Expected Signature",
			args{
				fromAddr:        suite.address[0],
				toAddr:          suite.address[1],
				unlockerAddress: suite.address[2],
				signature:       "[696E70757431]",
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Unexpected Signature",
			args{
				fromAddr:        suite.address[0],
				toAddr:          suite.address[1],
				unlockerAddress: suite.address[2],
				signature:       "[696E70757431999]",
			},
			errArgs{
				shouldPass: false,
				contains:   "",
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			var msg = types.NewMsgLockedSend(tc.args.fromAddr, tc.args.toAddr, tc.args.unlockerAddress.String(), sdk.NewCoins())
			res := msg.GetSigners()
			if tc.errArgs.shouldPass {
				suite.Require().Equal(fmt.Sprintf("%v", res), tc.args.signature)
			} else {
				suite.Require().NotEqual(fmt.Sprintf("%v", res), tc.args.signature)
			}
		})
	}
}
