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
	pk = genAccs(10)

	baseVAcc = vestingtypes.NewBaseVestingAccount(
		authtypes.NewBaseAccountWithAddress(toAddr(pk[0])), sdk.NewCoins(), math.MaxInt64)
	baseMVA = types.ManualVestingAccount{
		BaseVestingAccount: baseVAcc,
		VestedCoins:        sdk.NewCoins(),
		Unlocker:           toAddr(pk[1]).String(),
	}

	// 2000000uctk vested out of 5000000uctk
	baseVAcc2 = vestingtypes.NewBaseVestingAccount(
		authtypes.NewBaseAccountWithAddress(toAddr(pk[1])), sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 5000000)}, math.MaxInt64)
	baseMVA2 = types.ManualVestingAccount{
		BaseVestingAccount: baseVAcc2,
		VestedCoins:        sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 2000000)},
		Unlocker:           toAddr(pk[2]).String(),
	}

	// fully vested
	baseVAcc3 = vestingtypes.NewBaseVestingAccount(
		authtypes.NewBaseAccountWithAddress(toAddr(pk[2])), sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 5000000)}, math.MaxInt64)
	baseMVA3 = types.ManualVestingAccount{
		BaseVestingAccount: baseVAcc3,
		VestedCoins:        sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 5000000)},
		Unlocker:           toAddr(pk[3]).String(),
	}

	// fully vesting (locked)
	baseVAcc4 = vestingtypes.NewBaseVestingAccount(
		authtypes.NewBaseAccountWithAddress(toAddr(pk[3])), sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 5000000)}, math.MaxInt64)
	baseMVA4 = types.ManualVestingAccount{
		BaseVestingAccount: baseVAcc4,
		VestedCoins:        sdk.NewCoins(),
		Unlocker:           toAddr(pk[4]).String(),
	}
)

// shared setup
type ForkTestSuite struct {
	suite.Suite

	// cdc    *codec.LegacyAminogenAccs
	app      *simapp.SimApp
	ctx      sdk.Context
	ak       sdkauthkeeper.AccountKeeper
	bk       bankkeeper.Keeper
	sk       stakingkeeper.Keeper
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

func genAccs(n int) []ed25519.PrivKey {
	var pks []ed25519.PrivKey
	for i := 0; i < n; i++ {
		pks = append(pks, ed25519.GenPrivKey())
	}
	return pks
}

func toAddr(key ed25519.PrivKey) sdk.AccAddress {
	return key.PubKey().Address().Bytes()
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

	for _, pk := range pk {
		acc := toAddr(pk)
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
		{
			"manual vesting account with some delegated vesting coins", args{
				&baseMVA2,
				[]sdk.Int{sdk.NewInt(2000000)},
				[]sdk.Int{sdk.NewInt(1000000)},
			},
			errArgs{
				true,
				&baseMVA2,
			},
		},
		{
			"manual vesting account with some delegated vesting and delegated free coins", args{
				&baseMVA2,
				[]sdk.Int{sdk.NewInt(3500000)},
				[]sdk.Int{},
			},
			errArgs{
				true,
				&baseMVA2,
			},
		},
		{
			"fully vested manual vesting account", args{
				&baseMVA3,
				[]sdk.Int{sdk.NewInt(3500000)},
				[]sdk.Int{sdk.NewInt(1500000)},
			},
			errArgs{
				true,
				&baseMVA3,
			},
		},
		{
			"fully vested manual vesting account", args{
				&baseMVA3,
				[]sdk.Int{sdk.NewInt(3500000)},
				[]sdk.Int{sdk.NewInt(1500000)},
			},
			errArgs{
				true,
				&baseMVA3,
			},
		},
		{
			"fully vesting (locked) manual vesting account", args{
				&baseMVA4,
				[]sdk.Int{sdk.NewInt(3500000)},
				[]sdk.Int{sdk.NewInt(1500000)},
			},
			errArgs{
				true,
				&baseMVA4,
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

			// reset account delegation tracking to test the migration
			var mva types.ManualVestingAccount
			acc := tc.args.acc.(*types.ManualVestingAccount)
			mva = *acc

			mva.BaseVestingAccount.DelegatedFree = sdk.NewCoins()
			mva.BaseVestingAccount.DelegatedVesting = sdk.NewCoins()

			res, err := MigrateAccount(suite.ctx, &mva, suite.bk, &suite.sk)
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
