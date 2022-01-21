package mva

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkauthkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/certikfoundation/shentu/v2/common"
	"github.com/certikfoundation/shentu/v2/simapp"
	"github.com/certikfoundation/shentu/v2/x/auth/types"
	bankkeeper "github.com/certikfoundation/shentu/v2/x/bank/keeper"
	stakingkeeper "github.com/certikfoundation/shentu/v2/x/staking/keeper"
	"github.com/certikfoundation/shentu/v2/x/staking/teststaking"
)

var (
	acc1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc2 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc3 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc4 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	baseVAcc = vestingtypes.NewBaseVestingAccount(
		authtypes.NewBaseAccountWithAddress(acc1), sdk.NewCoins(), math.MaxInt64)
	baseMVA = types.ManualVestingAccount{
		BaseVestingAccount: baseVAcc,
		VestedCoins:        sdk.NewCoins(),
		Unlocker:           acc2.String(),
	}
)

// shared setup
type ForkTestSuite struct {
	suite.Suite

	// cdc    *codec.LegacyAmino
	app      *simapp.SimApp
	ctx      sdk.Context
	ak       sdkauthkeeper.AccountKeeper
	bk       bankkeeper.Keeper
	sk       stakingkeeper.Keeper
	address  []sdk.AccAddress
	tstaking *teststaking.Helper
}

func TestForkTestSuite(t *testing.T) {
	suite.Run(t, new(ForkTestSuite))
}

func nextBlock(ctx sdk.Context, tstaking *teststaking.Helper) sdk.Context {
	newTime := ctx.BlockTime().Add(time.Second * time.Duration(int64(common.SecondsPerBlock)))
	ctx = ctx.WithBlockTime(newTime).WithBlockHeight(ctx.BlockHeight() + 1)

	tstaking.TurnBlock(ctx)

	return ctx
}

func (suite *ForkTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.ak = suite.app.AccountKeeper
	suite.bk = suite.app.BankKeeper
	suite.sk = suite.app.StakingKeeper

	pks := simapp.CreateTestPubKeys(4)
	simapp.AddTestAddrsFromPubKeys(suite.app, suite.ctx, pks, sdk.NewInt(2e8))
	val1pk, val2pk := pks[2], pks[3]
	val1addr, val2addr := sdk.ValAddress(val1pk.Address()), sdk.ValAddress(val2pk.Address())

	// set up testing helpers
	tstaking := teststaking.NewHelper(suite.T(), suite.ctx, suite.app.StakingKeeper)

	for _, acc := range []sdk.AccAddress{acc1, acc2, acc3, acc4} {
		err := sdksimapp.FundAccount(
			suite.app.BankKeeper,
			suite.ctx,
			acc,
			sdk.NewCoins(
				sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(10000000000)), // 1,000 CTK
			),
		)
		if err != nil {
			panic(err)
		}
	}

	// set up two validators
	tstaking.CreateValidatorWithValPower(val1addr, val1pk, 100, true)
	suite.ctx = nextBlock(suite.ctx, tstaking)
	tstaking.CheckValidator(val1addr, stakingtypes.Bonded, false)

	tstaking.CreateValidatorWithValPower(val2addr, val2pk, 100, true)
	suite.ctx = nextBlock(suite.ctx, tstaking)
	tstaking.CheckValidator(val2addr, stakingtypes.Bonded, false)

	suite.tstaking = tstaking
	suite.address = []sdk.AccAddress{acc1, acc2, acc3, acc4}
}

func (suite *ForkTestSuite) TestFork() {
	type args struct {
		acc        authtypes.AccountI
		stakings   []sdk.Int
		unbondings []sdk.Int
	}
	type errArgs struct {
		shouldPass bool
		expected   *types.ManualVestingAccount
	}
	tests := []struct {
		name     string
		args     args
		expected errArgs
	}{
		{
			"empty acc", args{
				&baseMVA,
				[]sdk.Int{},
				[]sdk.Int{},
			},
			errArgs{
				true,
				&baseMVA,
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			for i, s := range tc.args.stakings {
				operAddr, err := sdk.ValAddressFromBech32(suite.sk.GetAllValidators(suite.ctx)[i].OperatorAddress)
				if err != nil {
					panic(err)
				}
				suite.tstaking.Delegate(tc.args.acc.GetAddress(), operAddr, s.Int64())
			}
			for i, u := range tc.args.unbondings {
				operAddr, err := sdk.ValAddressFromBech32(suite.sk.GetAllValidators(suite.ctx)[i].OperatorAddress)
				if err != nil {
					panic(err)
				}
				suite.tstaking.Undelegate(tc.args.acc.GetAddress(), operAddr, u.Int64(), true)
			}
			res, err := MigrateAccount(suite.ctx, tc.args.acc, suite.bk, &suite.sk)
			if tc.expected.shouldPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().Equal(res, tc.expected.expected)
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}
