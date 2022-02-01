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

	"github.com/certikfoundation/shentu/v2/common"
	"github.com/certikfoundation/shentu/v2/simapp"
	vesting "github.com/certikfoundation/shentu/v2/x/auth/types"
	"github.com/certikfoundation/shentu/v2/x/bank/keeper"
	"github.com/certikfoundation/shentu/v2/x/bank/types"
)

var (
	denom    = common.MicroCTKDenom
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
	params      types.AccountKeeper
	keeper      keeper.Keeper
	queryClient bankTypes.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = suite.app.BankKeeper
	suite.params = suite.app.AccountKeeper

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	bankTypes.RegisterQueryServer(queryHelper, suite.app.BankKeeper)
	suite.queryClient = bankTypes.NewQueryClient(queryHelper)

	for _, acc := range []sdk.AccAddress{acc1, acc2, acc3, acc4} {
		err := sdksimapp.FundAccount(
			suite.app.BankKeeper,
			suite.ctx,
			acc,
			sdk.NewCoins(
				sdk.NewCoin(denom, sdk.NewInt(1000)),
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
		sendAmt    sdk.Coins
		accBalance sdk.Coins
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
		{"Sender(1) isValid: SendAmt < Total Amount",
			args{
				sendAmt: sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(200))},
				accBalance: sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(800))},
				fromAddr:   suite.address[0],
				toAddr:     suite.address[1],
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Sender(1) inValid Coins: SendAmt = 0",
			args{
				sendAmt: sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(0))},
				accBalance: sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(1000))},
				fromAddr:   suite.address[0],
				toAddr:     suite.address[1],
			},
			errArgs{
				shouldPass: false,
				contains:   "",
			},
		},
		{"Sender(1) inValid: SendAmt > Total Amount",
			args{
				sendAmt: sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(11000))},
				accBalance: sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(1000))},
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
			err := suite.keeper.SendCoins(suite.ctx, tc.args.fromAddr, tc.args.toAddr, tc.args.sendAmt)
			balance := suite.keeper.GetAllBalances(suite.ctx, tc.args.fromAddr)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().NotEqual(tc.args.sendAmt, balance)
				suite.Require().Equal(tc.args.accBalance, balance)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestInputOutputCoins() {
	type args struct {
		Addr1        sdk.AccAddress
		Addr2        sdk.AccAddress
		Addr3        sdk.AccAddress
		amount       sdk.Coins
		addr1Balance sdk.Coins
		addr2Balance sdk.Coins
		addr3Balance sdk.Coins
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
		{"Sender(1) isValid: Sufficient Amount",
			args{
				amount:       sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(200))},
				Addr1:        suite.address[0],
				Addr2:        suite.address[1],
				Addr3:        suite.address[2],
				addr1Balance: sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(600))},
				addr2Balance: sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(1200))},
				addr3Balance: sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(1200))},
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Sender(1) inValid: InSufficient Amount",
			args{
				amount:       sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(5000))},
				Addr1:        suite.address[0],
				Addr2:        suite.address[1],
				Addr3:        suite.address[2],
				addr1Balance: sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(1000))},
				addr2Balance: sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(1000))},
				addr3Balance: sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(1000))},
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
				{Address: tc.args.Addr1.String(), Coins:  tc.args.amount},
				{Address: tc.args.Addr1.String(), Coins:  tc.args.amount},
			}
			outputs := []bankTypes.Output{
				{Address: tc.args.Addr2.String(), Coins: tc.args.amount},
				{Address: tc.args.Addr3.String(), Coins: tc.args.amount},
			}
			err := suite.keeper.InputOutputCoins(suite.ctx, inputs, outputs)
			addr1Balance := suite.keeper.GetAllBalances(suite.ctx, tc.args.Addr1)
			addr2Balance := suite.keeper.GetAllBalances(suite.ctx, tc.args.Addr2)
			addr3Balance := suite.keeper.GetAllBalances(suite.ctx, tc.args.Addr3)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.args.addr1Balance, addr1Balance)
				suite.Require().Equal(tc.args.addr2Balance, addr2Balance)
				suite.Require().Equal(tc.args.addr3Balance, addr3Balance)
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
		lockedSendAmt   sdk.Coins
		accBalance      sdk.Coins
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
		{"Sender(1) isValid: lockedSendAmt < Total Amount",
			args{
				lockedSendAmt:   sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(200))},
				accBalance:      sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(1200))),
				fromAddr:        suite.address[0],
				toAddr:          suite.address[1],
				UnlockerAddress: suite.address[2],
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Sender(1) inValid: lockedSendAmt > Total Amount",
			args{
				lockedSendAmt:   sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(2000))},
				accBalance:      sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(3000))),
				fromAddr:        suite.address[0],
				toAddr:          suite.address[1],
				UnlockerAddress: suite.address[2],
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
			//primary checks
			from := suite.params.GetAccount(suite.ctx, tc.args.fromAddr)
			acc := suite.params.GetAccount(suite.ctx, tc.args.toAddr)
			suite.Require().NotNil(acc)
			baseAcc := authtypes.NewBaseAccount(tc.args.toAddr, acc.GetPubKey(), acc.GetAccountNumber(), acc.GetSequence())
			toAcc := vesting.NewManualVestingAccount(baseAcc, sdk.NewCoins(), sdk.NewCoins(), tc.args.UnlockerAddress)
			suite.params.SetAccount(suite.ctx, toAcc)
			//send coin
			sendCoins := tc.args.lockedSendAmt
			//send some coins to the vesting account
			err := suite.keeper.SendCoins(suite.ctx, tc.args.fromAddr, tc.args.toAddr, sendCoins)
			//require that the coin is spendable
			toAcc = suite.params.GetAccount(suite.ctx, tc.args.toAddr).(*vesting.ManualVestingAccount)
			balances := suite.keeper.GetAllBalances(suite.ctx, toAcc.GetAddress())
			if tc.errArgs.shouldPass {
				suite.Require().NotNil(from)
				suite.Require().NotEqual(tc.args.toAddr, tc.args.UnlockerAddress)
				suite.Require().NotNil(baseAcc)
				suite.Require().NoError(err)
				suite.Require().NotNil(toAcc)
				suite.Require().Equal(balances, tc.args.accBalance)
				//LockedCoinsFromVesting returns all the coins that are not spendable (i.e. locked)
				suite.Require().Equal(toAcc.LockedCoinsFromVesting(sendCoins), tc.args.lockedSendAmt)
				//LockedCoins returns the set of coins that are  spendable plus any have vested
				suite.Require().Equal(balances.Sub(toAcc.LockedCoins(suite.ctx.BlockTime())), tc.args.accBalance)
			} else {
				suite.Require().NotEqual(balances, tc.args.accBalance)
				suite.Require().Error(err)
			}
		})
	}
}
