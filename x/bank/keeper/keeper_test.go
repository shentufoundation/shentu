package keeper_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/certikfoundation/shentu/v2/simapp"
	vesting "github.com/certikfoundation/shentu/v2/x/auth/types"
	"github.com/certikfoundation/shentu/v2/x/bank/keeper"
	"github.com/certikfoundation/shentu/v2/x/bank/types"
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

	address     []sdk.AccAddress
	app         *simapp.SimApp
	ctx         sdk.Context
	queryClient types.MsgServer
	params      types.AccountKeeper
	keeper      keeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = suite.app.BankKeeper
	suite.params = suite.app.AccountKeeper

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterMsgServer(queryHelper, suite.queryClient)
	suite.queryClient = &types.UnimplementedMsgServer{}

	for _, acc := range []sdk.AccAddress{acc1, acc2, acc3, acc4} {
		err := sdksimapp.FundAccount(
			suite.app.BankKeeper,
			suite.ctx,
			acc,
			sdk.NewCoins(
				sdk.NewCoin("ctk", sdk.NewInt(1000)), // 1,000 CTK
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

func (suite *KeeperTestSuite) TestSendCoins() {
	type args struct {
		fromAddr   sdk.AccAddress
		toAddr     sdk.AccAddress
		amount     int64
		accBalance int64
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
		{"Operator(1) Create: first send test case if coins is not greater than total amount",
			args{
				amount:     200,
				accBalance: 800,
				fromAddr:   suite.address[0],
				toAddr:     suite.address[1],
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Operator(1) Create: second send test case if coins is greater than  total amount",
			args{
				amount:     11000,
				accBalance: 1000,
				fromAddr:   suite.address[0],
				toAddr:     suite.address[1],
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
			err := suite.keeper.SendCoins(suite.ctx, tc.args.fromAddr, tc.args.toAddr, sdk.Coins{sdk.NewInt64Coin("ctk", tc.args.amount)})
			balance := suite.keeper.GetAllBalances(suite.ctx, tc.args.fromAddr)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().NotEqual(tc.args.amount, balance)
				suite.Require().Equal(sdk.Coins{sdk.NewInt64Coin("ctk", tc.args.accBalance)}, balance)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().NotEqual(tc.args.amount, balance)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestInputOutputCoins() {
	type args struct {
		Addr1  sdk.AccAddress
		Addr2  sdk.AccAddress
		Addr3  sdk.AccAddress
		amount int64
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
		{"Operator(1) Create: first send test case",
			args{
				amount: 200,
				Addr1:  suite.address[0],
				Addr2:  suite.address[1],
				Addr3:  suite.address[2],
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Operator(1) Create: second send if amount is insufficient",
			args{
				amount: 5000,
				Addr1:  suite.address[0],
				Addr2:  suite.address[1],
				Addr3:  suite.address[2],
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
			inputs := []bankTypes.Input{
				{Address: tc.args.Addr1.String(), Coins: sdk.Coins{sdk.NewInt64Coin("ctk", tc.args.amount)}},
				{Address: tc.args.Addr1.String(), Coins: sdk.Coins{sdk.NewInt64Coin("ctk", tc.args.amount)}},
			}
			outputs := []bankTypes.Output{
				{Address: tc.args.Addr2.String(), Coins: sdk.Coins{sdk.NewInt64Coin("ctk", tc.args.amount)}},
				{Address: tc.args.Addr3.String(), Coins: sdk.Coins{sdk.NewInt64Coin("ctk", tc.args.amount)}},
			}
			err := suite.keeper.InputOutputCoins(suite.ctx, inputs, outputs)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestLockedSend() {
	type args struct {
		fromAddr        sdk.AccAddress
		toAddr          sdk.AccAddress
		UnlockerAddress sdk.AccAddress
		amount          int64
		accBalance      int64
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
		{"Operator(1) Create: first send test case",
			args{
				amount:          200,
				accBalance:      1200,
				fromAddr:        suite.address[0],
				toAddr:          suite.address[1],
				UnlockerAddress: suite.address[2],
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Operator(1) Create: second send test case",
			args{
				amount:          100,
				accBalance:      1100,
				fromAddr:        suite.address[0],
				toAddr:          suite.address[1],
				UnlockerAddress: suite.address[2],
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
			//primary checks
			from := suite.params.GetAccount(suite.ctx, tc.args.fromAddr)
			suite.Require().NotNil(from)
			suite.Require().NotEqual(tc.args.toAddr, tc.args.UnlockerAddress)
			acc := suite.params.GetAccount(suite.ctx, tc.args.toAddr)
			suite.Require().NotNil(acc)
			baseAcc := authtypes.NewBaseAccount(tc.args.toAddr, acc.GetPubKey(), acc.GetAccountNumber(), acc.GetSequence())
			suite.Require().NotNil(baseAcc)
			toAcc := vesting.NewManualVestingAccount(baseAcc, sdk.NewCoins(), sdk.NewCoins(), tc.args.UnlockerAddress)
			suite.Require().NotNil(toAcc)
			suite.params.SetAccount(suite.ctx, toAcc)
			//send coin
			sendCoins := sdk.Coins{sdk.NewInt64Coin("ctk", tc.args.amount)}
			// send some coins to the vesting account
			err := suite.keeper.SendCoins(suite.ctx, tc.args.fromAddr, tc.args.toAddr, sendCoins)
			suite.Require().NoError(err)
			//require that the coin is spendable
			toAcc = suite.params.GetAccount(suite.ctx, tc.args.toAddr).(*vesting.ManualVestingAccount)
			balances := suite.keeper.GetAllBalances(suite.ctx, toAcc.GetAddress())
			if tc.errArgs.shouldPass {
				suite.Require().Equal(balances, sdk.Coins{sdk.NewInt64Coin("ctk", tc.args.accBalance)})
				//LockedCoinsFromVesting returns all the coins that are not spendable (i.e. locked)
				suite.Require().Equal(toAcc.LockedCoinsFromVesting(sendCoins), sdk.Coins{sdk.NewInt64Coin("ctk", tc.args.amount)})
				//LockedCoins returns the set of coins that are  spendable plus any have vested
				suite.Require().Equal(balances.Sub(toAcc.LockedCoins(suite.ctx.BlockTime())), sdk.Coins{sdk.NewInt64Coin("ctk", tc.args.accBalance)})
			} else {
				suite.Require().Equal(balances, sdk.Coins{sdk.NewInt64Coin("ctk", tc.args.accBalance)})
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}
