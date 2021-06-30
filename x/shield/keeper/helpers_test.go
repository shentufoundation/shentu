package keeper_test

import (
	"testing"
	"time"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/gov/testgov"
	"github.com/certikfoundation/shentu/x/shield/keeper"
	"github.com/certikfoundation/shentu/x/shield/testshield"
	"github.com/certikfoundation/shentu/x/shield/types"
	"github.com/certikfoundation/shentu/x/staking/teststaking"
)

var (
	pks  = simapp.CreateTestPubKeys(5)
	acc1 = sdk.AccAddress(pks[0].Address().Bytes())
	acc2 = sdk.AccAddress(pks[1].Address().Bytes())
	acc3 = sdk.AccAddress(pks[2].Address().Bytes())
	acc4 = sdk.AccAddress(pks[3].Address().Bytes())
	acc5 = sdk.AccAddress(pks[4].Address().Bytes())

	val1 = sdk.ValAddress(pks[0].Address())
	val2 = sdk.ValAddress(pks[1].Address())
	val3 = sdk.ValAddress(pks[2].Address())
	val4 = sdk.ValAddress(pks[3].Address())
	val5 = sdk.ValAddress(pks[4].Address())

	basePurchase = types.Purchase{
		PurchaseId:        1,
		ProtectionEndTime: time.Time{},
		DeletionTime:      time.Time{},
		Description:       "",
		Shield:            sdk.OneInt(),
		ServiceFees:       OneMixedDecCoins(common.MicroCTKDenom),
	}
)

type TestSuite struct {
	*testing.T

	app         *simapp.SimApp
	ctx         sdk.Context
	keeper      keeper.Keeper
	tshield     *testshield.Helper
	tstaking    *teststaking.Helper
	tgov        *testgov.Helper
	accounts    []sdk.AccAddress
	vals        []stakingtypes.Validator
	queryClient types.QueryClient
}

func setup(t *testing.T) TestSuite {
	var ts TestSuite
	ts.T = t
	ts.app = simapp.Setup(false)
	ts.ctx = ts.app.BaseApp.NewContext(false, tmproto.Header{})
	ts.keeper = ts.app.ShieldKeeper
	ts.accounts = []sdk.AccAddress{acc1, acc2, acc3, acc4, acc5}
	for _, acc := range []sdk.AccAddress{acc1, acc2, acc3, acc4, acc5} {
		err := ts.app.BankKeeper.AddCoins(
			ts.ctx,
			acc,
			sdk.NewCoins(
				sdk.NewCoin(ts.app.StakingKeeper.BondDenom(ts.ctx), sdk.NewInt(100000000000)), // 10,000 stake
			),
		)
		if err != nil {
			panic(err)
		}
	}
	ts.tstaking = teststaking.NewHelper(ts.T, ts.ctx, ts.app.StakingKeeper)
	ts.tgov = testgov.NewHelper(ts.T, ts.ctx, ts.app.GovKeeper, ts.tstaking.Denom)
	ts.tshield = testshield.NewHelper(ts.T, ts.ctx, ts.keeper, ts.tstaking.Denom)

	queryHelper := baseapp.NewQueryServerTestHelper(ts.ctx, ts.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, ts.app.ShieldKeeper)
	ts.queryClient = types.NewQueryClient(queryHelper)

	ts.setupProviders()
	return ts
}

func OneMixedCoins(nativeDenom string) types.MixedCoins {
	native := sdk.NewCoins(sdk.NewCoin(nativeDenom, sdk.NewInt(1)))
	foreign := sdk.NewCoins(sdk.NewCoin("dummy", sdk.NewInt(1)))
	return types.MixedCoins{
		Native:  native,
		Foreign: foreign,
	}
}

func OneMixedDecCoins(nativeDenom string) types.MixedDecCoins {
	native := sdk.NewDecCoins(sdk.NewDecCoin(nativeDenom, sdk.NewInt(1)))
	foreign := sdk.NewDecCoins(sdk.NewDecCoin("dummy", sdk.NewInt(1)))
	return types.MixedDecCoins{
		Native:  native,
		Foreign: foreign,
	}
}

func DummyPool(id uint64) types.Pool {
	return types.Pool{
		Id:          id,
		Description: "w",
		Sponsor:     acc1.String(),
		SponsorAddr: acc2.String(),
		ShieldLimit: sdk.NewInt(1),
		Active:      true,
		Shield:      sdk.NewInt(1),
	}
}

type poolpurchase struct {
	poolID    uint64
	purchases []types.Purchase
}

func (suite TestSuite) setupProviders() {
	simapp.AddTestAddrsFromPubKeys(suite.app, suite.ctx, pks, sdk.NewInt(2e8))
	for _, pk := range pks {
		suite.tstaking.CreateValidatorWithValPower(sdk.ValAddress(pk.Address()), pk, 10000, true)
		val := suite.tstaking.CheckValidator(sdk.ValAddress(pk.Address()), -1, false)
		suite.vals = append(suite.vals, val)
		suite.tshield.DepositCollateral(acc1, 500000000, true)
		suite.tstaking.TurnBlock(suite.ctx.BlockTime().Add(time.Second))
		suite.tshield.TurnBlock(suite.ctx.BlockTime().Add(time.Second))
	}
}

func (suite TestSuite) setupUndelegate() {
	suite.tstaking.Undelegate(acc5, val5, sdk.NewInt(1000000000), true)
	suite.tstaking.TurnBlock(suite.ctx.BlockTime().Add(time.Second))
	suite.tshield.TurnBlock(suite.ctx.BlockTime().Add(time.Second))
}
