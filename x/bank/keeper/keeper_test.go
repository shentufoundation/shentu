package keeper_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
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
	queryClient types.QueryClient
	params      types.AccountKeeper
	keeper      keeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.params = suite.app.AccountKeeper
	suite.keeper = suite.app.BankKeeper

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.BankKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

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
		{"Operator(1) Create: first send",
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
		{"Operator(1) Create: second send if balance is less",
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
			err := suite.app.BankKeeper.SendCoins(suite.ctx, tc.args.fromAddr, tc.args.toAddr, sdk.Coins{sdk.NewInt64Coin("ctk", tc.args.amount)})
			balance := suite.app.BankKeeper.GetAllBalances(suite.ctx, tc.args.fromAddr)
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

func (suite *KeeperTestSuite) TestLockedSend() {
	//testing for locksend
}
