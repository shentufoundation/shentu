package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/certikfoundation/shentu/common/tests"
	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/staking/teststaking"
)

// TestWithdraw tests withdraws triggered by staking undelegation.
func TestWithdrawsByUndelegate(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()})

	//p := app.StakingKeeper.GetParams(ctx)
	//p.MaxValidators = 5
	//app.StakingKeeper.SetParams(ctx, p)

	// create and add addresses
	delAddr := simapp.AddTestAddrs(app, ctx, 1, sdk.TokensFromConsensusPower(200))[0]

	accAddr := simapp.AddTestAddrs(app, ctx, 1, sdk.TokensFromConsensusPower(200))[0]
	valAddr := sdk.ValAddress(accAddr)

	pubKey := tests.MakeTestPubKey()

	// get testing helpers - no need?
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

	// Set up a validator
	tstaking.CreateValidatorWithValPower(valAddr, pubKey, 100, true)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	tstaking.CheckValidator(valAddr, sdk.Bonded, false)

	// Attempt depositing collateral
	// TODO: Create shield test helper
	err := app.ShieldKeeper.DepositCollateral(ctx, delAddr, sdk.NewInt(75))
	require.Error(t, err)

	// Delegate 100 to the validator
	tstaking.CheckDelegator(delAddr, valAddr, false)
	tstaking.Delegate(delAddr, valAddr, 100)
	tstaking.CheckDelegator(delAddr, valAddr, true)

	// Deposit collateral of amount 75
	err = app.ShieldKeeper.DepositCollateral(ctx, delAddr, sdk.NewInt(75))
	require.Nil(t, err)
	
	// Undelegate 50 to trigger withdrawal of 25
	//tstaking.Undelegate(sdk.AccAddress(valAddr), valAddr, sdk.TokensFromConsensusPower(1), true)
	tstaking.Undelegate(delAddr, valAddr, sdk.NewInt(50), true)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	withdraws := app.ShieldKeeper.GetAllWithdraws(ctx)
	require.True(t, withdraws[0].Amount.Equal(sdk.NewInt(25)))
	numWithdraws := len(withdraws)

	// Undelegate 5 and trigger another withdrawal of 5.
	tstaking.Undelegate(delAddr, valAddr, sdk.NewInt(5), true)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	numWithdraws++
	withdraws = app.ShieldKeeper.GetAllWithdraws(ctx)
	require.True(t, len(withdraws) == numWithdraws)
	require.True(t, withdraws[numWithdraws-1].Amount.Equal(sdk.NewInt(5)))

	// Must fail deposit of 10
	err = app.ShieldKeeper.DepositCollateral(ctx, delAddr, sdk.NewInt(10))
	require.Error(t, err)

	// Delegate 25
	tstaking.Delegate(delAddr, valAddr, 25)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// Withdraw 5
	err = app.ShieldKeeper.WithdrawCollateral(ctx, delAddr, sdk.NewInt(5))
	require.Nil(t, err)
	numWithdraws++
	withdraws = app.ShieldKeeper.GetAllWithdraws(ctx) // GetAllWithdraws NOT WORKING?
	require.True(t, len(withdraws) == numWithdraws)

	// Undelegate 25. Shouldn't trigger withdrawal
	tstaking.Undelegate(delAddr, valAddr, sdk.NewInt(25), true)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	withdraws = app.ShieldKeeper.GetAllWithdraws(ctx) // GetAllWithdraws NOT WORKING?
	fmt.Printf("\n WITHDRAWS: %v\n", withdraws)
	require.True(t, len(withdraws) == numWithdraws)
}


// TestWithdraw tests withdraws triggered by staking redelegation.
func TestWithdrawsByRedelegate(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()})

	// create and add addresses
	delAddr := simapp.AddTestAddrs(app, ctx, 1, sdk.TokensFromConsensusPower(200))[0]

	accAddr := simapp.AddTestAddrs(app, ctx, 1, sdk.TokensFromConsensusPower(200))[0]
	valAddr := sdk.ValAddress(accAddr)
	accAddr2 := simapp.AddTestAddrs(app, ctx, 1, sdk.TokensFromConsensusPower(200))[0]
	valAddr2 := sdk.ValAddress(accAddr2)

	pubKey := tests.MakeTestPubKey()
	pubKey2 := tests.MakeTestPubKey()

	// get testing helpers
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

	// Set up validators
	tstaking.CreateValidatorWithValPower(valAddr, pubKey, 100, true)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	tstaking.CheckValidator(valAddr, sdk.Bonded, false)

	tstaking.CreateValidatorWithValPower(valAddr2, pubKey2, 100, true)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	tstaking.CheckValidator(valAddr2, sdk.Bonded, false)

	// Attempt depositing collateral
	err := app.ShieldKeeper.DepositCollateral(ctx, delAddr, sdk.NewInt(75))
	require.Error(t, err)

	// Delegate 100 to the validator
	tstaking.CheckDelegator(delAddr, valAddr, false)
	tstaking.Delegate(delAddr, valAddr, 100)
	tstaking.CheckDelegator(delAddr, valAddr, true)

	// Deposit collateral of amount 75
	err = app.ShieldKeeper.DepositCollateral(ctx, delAddr, sdk.NewInt(75))
	require.Nil(t, err)
	
	// Redelegate 50 to trigger withdrawal of 25
	// Remaining staking: 100, remaining deposit: 50
	tstaking.Redelegate(delAddr, valAddr, valAddr2, 50, true)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	withdraws := app.ShieldKeeper.GetAllWithdraws(ctx)
	require.True(t, withdraws[0].Amount.Equal(sdk.NewInt(25)))
	numWithdraws := len(withdraws)

	// Redelegate 30 but do not trigger withdrawal
	// Remaining staking: 100, remaining deposit: 50
	tstaking.Redelegate(delAddr, valAddr, valAddr2, 30, true)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	withdraws = app.ShieldKeeper.GetAllWithdraws(ctx)
	require.True(t, len(withdraws) == numWithdraws)

	// Must fail deposit of 60
	err = app.ShieldKeeper.DepositCollateral(ctx, delAddr, sdk.NewInt(60))
	require.Error(t, err)

	// Must succeed deposit of 50
	err = app.ShieldKeeper.DepositCollateral(ctx, delAddr, sdk.NewInt(50))
	require.Nil(t, err)

	// Redelegation hopping must fail
	tstaking.Redelegate(delAddr, valAddr2, valAddr, 10, false)
}
