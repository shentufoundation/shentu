package keeper_test

import (
	"testing"
	"time"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

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
	PKS  = simapp.CreateTestPubKeys(5)
	acc1 = sdk.AccAddress(PKS[0].Address().Bytes())
	acc2 = sdk.AccAddress(PKS[1].Address().Bytes())
	acc3 = sdk.AccAddress(PKS[2].Address().Bytes())
	acc4 = sdk.AccAddress(PKS[3].Address().Bytes())
	acc5 = sdk.AccAddress(PKS[4].Address().Bytes())

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
	simapp.AddTestAddrsFromPubKeys(suite.app, suite.ctx, PKS, sdk.NewInt(2e8))
	for _, pk := range PKS {
		suite.tstaking.CreateValidatorWithValPower(sdk.ValAddress(pk.Address()), pk, 10000, true)
		suite.tshield.DepositCollateral(acc1, 500000000, true)
		suite.tstaking.TurnBlock(suite.ctx.BlockTime().Add(time.Second))
		suite.tshield.TurnBlock(suite.ctx.BlockTime().Add(time.Second))
	}
}
