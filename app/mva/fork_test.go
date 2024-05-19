package mva_test

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkauthkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/app/mva"
	"github.com/shentufoundation/shentu/v2/common"
	"github.com/shentufoundation/shentu/v2/x/auth/types"
	bankkeeper "github.com/shentufoundation/shentu/v2/x/bank/keeper"
	stakingkeeper "github.com/shentufoundation/shentu/v2/x/staking/keeper"
	"github.com/shentufoundation/shentu/v2/x/staking/teststaking"
)

var (
	pk = genAccs(10)

	baseVAcc = vestingtypes.NewBaseVestingAccount(
		authtypes.NewBaseAccountWithAddress(toAddr(pk[0])), sdk.NewCoins(), 0)
	baseMVA = types.ManualVestingAccount{
		BaseVestingAccount: baseVAcc,
		VestedCoins:        sdk.NewCoins(),
		Unlocker:           toAddr(pk[1]).String(),
	}

	// 2000000uctk vested out of 5000000uctk
	baseVAcc2 = vestingtypes.NewBaseVestingAccount(
		authtypes.NewBaseAccountWithAddress(toAddr(pk[1])), sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 5000000)}, 0)
	baseMVA2 = types.ManualVestingAccount{
		BaseVestingAccount: baseVAcc2,
		VestedCoins:        sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 2000000)},
		Unlocker:           toAddr(pk[2]).String(),
	}

	// fully vested
	baseVAcc3 = vestingtypes.NewBaseVestingAccount(
		authtypes.NewBaseAccountWithAddress(toAddr(pk[2])), sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 5000000)}, 0)
	baseMVA3 = types.ManualVestingAccount{
		BaseVestingAccount: baseVAcc3,
		VestedCoins:        sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 5000000)},
		Unlocker:           toAddr(pk[3]).String(),
	}

	// fully vesting (locked)
	baseVAcc4 = vestingtypes.NewBaseVestingAccount(
		authtypes.NewBaseAccountWithAddress(toAddr(pk[3])), sdk.Coins{sdk.NewInt64Coin(common.MicroCTKDenom, 5000000)}, 0)
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
	app      *shentuapp.ShentuApp
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

func copyMVA(mva types.ManualVestingAccount) *types.ManualVestingAccount {
	unlocker, err := sdk.AccAddressFromBech32(mva.Unlocker)
	if err != nil {
		panic(err)
	}
	return types.NewManualVestingAccount(mva.BaseAccount, mva.OriginalVesting, mva.VestedCoins, unlocker)
}

func (suite *ForkTestSuite) SetupTest() {
	suite.app = shentuapp.Setup(suite.T(), false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.ak = suite.app.AccountKeeper
	suite.bk = suite.app.BankKeeper
	suite.sk = suite.app.StakingKeeper

	pks := shentuapp.CreateTestPubKeys(4)
	shentuapp.AddTestAddrsFromPubKeys(suite.app, suite.ctx, pks, sdk.NewInt(2e8))
	val1pk, val2pk := pks[2], pks[3]
	val1addr, val2addr := sdk.ValAddress(val1pk.Address()), sdk.ValAddress(val2pk.Address())

	// set up testing helpers
	tstaking := teststaking.NewHelper(suite.T(), suite.ctx, suite.app.StakingKeeper)

	for _, pk := range pk {
		acc := toAddr(pk)
		err := testutil.FundAccount(
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
		dV         sdk.Coins
		dF         sdk.Coins
	}
	tests := []struct {
		name     string
		args     args
		expected errArgs
	}{
		{
			"empty acc", args{
				copyMVA(baseMVA),
				[]sdk.Int{},
				[]sdk.Int{},
			},
			errArgs{
				true,
				sdk.NewCoins(),
				sdk.NewCoins(),
			},
		},
		{
			"manual vesting account with some delegated vesting coins", args{
				copyMVA(baseMVA2),
				[]sdk.Int{sdk.NewInt(2000000)},
				[]sdk.Int{sdk.NewInt(1000000)},
			},
			errArgs{
				true,
				sdk.NewCoins(sdk.NewInt64Coin(common.MicroCTKDenom, 2000000)),
				sdk.NewCoins(),
			},
		},
		{
			"manual vesting account with some delegated vesting and delegated free coins", args{
				copyMVA(baseMVA2),
				[]sdk.Int{sdk.NewInt(3500000)},
				[]sdk.Int{},
			},
			errArgs{
				true,
				sdk.NewCoins(sdk.NewInt64Coin(common.MicroCTKDenom, 3000000)),
				sdk.NewCoins(sdk.NewInt64Coin(common.MicroCTKDenom, 500000)),
			},
		},
		{
			"fully vested manual vesting account", args{
				copyMVA(baseMVA3),
				[]sdk.Int{sdk.NewInt(3500000)},
				[]sdk.Int{sdk.NewInt(1500000)},
			},
			errArgs{
				true,
				sdk.NewCoins(),
				sdk.NewCoins(sdk.NewInt64Coin(common.MicroCTKDenom, 3500000)),
			},
		},
		{
			"fully vesting (locked) manual vesting account", args{
				copyMVA(baseMVA4),
				[]sdk.Int{sdk.NewInt(3500000)},
				[]sdk.Int{sdk.NewInt(1500000)},
			},
			errArgs{
				true,
				sdk.NewCoins(sdk.NewInt64Coin(common.MicroCTKDenom, 3500000)),
				sdk.NewCoins(),
			},
		},
		{
			"manual vesting account with some delegated vesting coins with multiple validators", args{
				copyMVA(baseMVA2),
				[]sdk.Int{sdk.NewInt(1000000), sdk.NewInt(1000000)},
				[]sdk.Int{sdk.NewInt(500000), sdk.NewInt(500000)},
			},
			errArgs{
				true,
				sdk.NewCoins(sdk.NewInt64Coin(common.MicroCTKDenom, 2000000)),
				sdk.NewCoins(),
			},
		},
		{
			"manual vesting account with some delegated vesting and delegated free coins with multiple validators", args{
				copyMVA(baseMVA2),
				[]sdk.Int{sdk.NewInt(2000000), sdk.NewInt(1500000)},
				[]sdk.Int{},
			},
			errArgs{
				true,
				sdk.NewCoins(sdk.NewInt64Coin(common.MicroCTKDenom, 3000000)),
				sdk.NewCoins(sdk.NewInt64Coin(common.MicroCTKDenom, 500000)),
			},
		},
		{
			"fully vested manual vesting account with multiple validators", args{
				copyMVA(baseMVA3),
				[]sdk.Int{sdk.NewInt(2000000), sdk.NewInt(1500000)},
				[]sdk.Int{sdk.NewInt(1000000), sdk.NewInt(500000)},
			},
			errArgs{
				true,
				sdk.NewCoins(),
				sdk.NewCoins(sdk.NewInt64Coin(common.MicroCTKDenom, 3500000)),
			},
		},
		{
			"fully vesting (locked) manual vesting account with multiple validators", args{
				copyMVA(baseMVA4),
				[]sdk.Int{sdk.NewInt(2000000), sdk.NewInt(1500000)},
				[]sdk.Int{sdk.NewInt(1000000), sdk.NewInt(500000)},
			},
			errArgs{
				true,
				sdk.NewCoins(sdk.NewInt64Coin(common.MicroCTKDenom, 3500000)),
				sdk.NewCoins(),
			},
		},
		{
			"sample failing test case", args{
				copyMVA(baseMVA4),
				[]sdk.Int{sdk.NewInt(2000000), sdk.NewInt(1500000)},
				[]sdk.Int{sdk.NewInt(1000000), sdk.NewInt(500000)},
			},
			errArgs{
				false,
				sdk.NewCoins(sdk.NewInt64Coin(common.MicroCTKDenom, 3000000)),
				sdk.NewCoins(),
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
				suite.tstaking.Delegate(tc.args.acc.GetAddress(), operAddr, math.NewInt(s.Int64()))
			}
			for i, u := range tc.args.unbondings {
				operAddr, err := sdk.ValAddressFromBech32(suite.sk.GetAllValidators(suite.ctx)[i].OperatorAddress)
				if err != nil {
					panic(err)
				}
				suite.tstaking.Undelegate(tc.args.acc.GetAddress(), operAddr, u.Int64(), true)
			}
			suite.tstaking.TurnBlock(suite.ctx)
			res := mva.MigrateAccount(suite.ctx, tc.args.acc, suite.bk, &suite.sk)

			resMVA := res.(*types.ManualVestingAccount)
			if tc.expected.shouldPass {
				suite.Require().Equal(resMVA.BaseVestingAccount.DelegatedVesting, tc.expected.dV)
				suite.Require().Equal(resMVA.BaseVestingAccount.DelegatedFree, tc.expected.dF)
			} else {
				suite.Require().True(!resMVA.BaseVestingAccount.DelegatedVesting.IsEqual(tc.expected.dV) ||
					resMVA.BaseVestingAccount.DelegatedFree.IsEqual(tc.expected.dF))
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
