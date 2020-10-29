package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/certikfoundation/shentu/common/tests"
	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/shield/types"
	"github.com/certikfoundation/shentu/x/staking/teststaking"
)

// TestWithdraw tests withdraws triggered by staking undelegation.
func TestWithdrawsByUndelegate(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()})

	// create and add addresses
	delAddr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(2e8))[0]

	delAddr2 := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(2e8))[0]

	accAddr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(2e8))[0]
	valAddr := sdk.ValAddress(accAddr)
	pubKey := tests.MakeTestPubKey()

	accAddr2 := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(2e8))[0]
	valAddr2 := sdk.ValAddress(accAddr2)
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

	// Both delegators delegate 50 to each validator
	tstaking.CheckDelegator(delAddr, valAddr, false)
	tstaking.Delegate(delAddr, valAddr, 50)
	tstaking.CheckDelegator(delAddr, valAddr, true)

	tstaking.CheckDelegator(delAddr, valAddr2, false)
	tstaking.Delegate(delAddr, valAddr2, 50)
	tstaking.CheckDelegator(delAddr, valAddr2, true)

	tstaking.CheckDelegator(delAddr2, valAddr, false)
	tstaking.Delegate(delAddr2, valAddr, 50)
	tstaking.CheckDelegator(delAddr2, valAddr, true)

	tstaking.CheckDelegator(delAddr2, valAddr2, false)
	tstaking.Delegate(delAddr2, valAddr2, 50)
	tstaking.CheckDelegator(delAddr2, valAddr2, true)

	// Both delegators deposit collateral of amount 75
	err = app.ShieldKeeper.DepositCollateral(ctx, delAddr, sdk.NewInt(75))
	require.Nil(t, err)

	err = app.ShieldKeeper.DepositCollateral(ctx, delAddr2, sdk.NewInt(75))
	require.Nil(t, err)

	// Undelegate total 50 to trigger total withdrawal of 25
	tstaking.Undelegate(delAddr, valAddr, sdk.NewInt(30), true)
	tstaking.Undelegate(delAddr2, valAddr2, sdk.NewInt(10), true)
	tstaking.Undelegate(delAddr, valAddr2, sdk.NewInt(20), true)
	tstaking.Undelegate(delAddr2, valAddr2, sdk.NewInt(40), true)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	withdraws := app.ShieldKeeper.GetAllWithdraws(ctx)
	numWithdraws := len(withdraws)

	require.True(t, withdraws[0].Amount.Equal(sdk.NewInt(5)))
	require.True(t, withdraws[0].Address.Equals(delAddr))
	require.True(t, withdraws[1].Amount.Equal(sdk.NewInt(20)))
	require.True(t, withdraws[1].Address.Equals(delAddr))
	require.True(t, withdraws[2].Amount.Equal(sdk.NewInt(25)))
	require.True(t, withdraws[2].Address.Equals(delAddr2))

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
	err = app.ShieldKeeper.WithdrawCollateral(ctx, delAddr, sdk.NewInt(5), nil)
	require.Nil(t, err)
	numWithdraws++
	withdraws = app.ShieldKeeper.GetAllWithdraws(ctx) // GetAllWithdraws NOT WORKING?
	require.True(t, len(withdraws) == numWithdraws)

	// Undelegate 25. Shouldn't trigger withdrawal
	tstaking.Undelegate(delAddr, valAddr, sdk.NewInt(25), true)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	withdraws = app.ShieldKeeper.GetAllWithdraws(ctx) // GetAllWithdraws NOT WORKING?
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
	pubKey := tests.MakeTestPubKey()

	accAddr2 := simapp.AddTestAddrs(app, ctx, 1, sdk.TokensFromConsensusPower(200))[0]
	valAddr2 := sdk.ValAddress(accAddr2)
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

	// Redelegation hopping must fail
	tstaking.Redelegate(delAddr, valAddr2, valAddr, 10, false)

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
}

func TestSecureCollaterals(t *testing.T) {
	// Test setup
	app := simapp.Setup(false)
	curTime := time.Now().UTC()
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: curTime})
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)
	bondDenom := tstaking.Denom
	shieldKeeper := app.ShieldKeeper
	govKeeper := app.GovKeeper
	stakingKeeper := app.StakingKeeper

	// Create and add addresses
	shieldAddr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(250000000000))[0] // 250,000ctk
	shieldKeeper.SetAdmin(ctx, shieldAddr)
	sponsorAddr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1))[0]
	purchaser := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(2000000000))[0]
	delAddr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(125000000000))[0] // 125,000ctk

	// Set up a validator
	accAddr := simapp.AddTestAddrs(app, ctx, 1, sdk.TokensFromConsensusPower(200))[0]
	valAddr := sdk.ValAddress(accAddr)
	pubKey := tests.MakeTestPubKey()

	tstaking.CreateValidatorWithValPower(valAddr, pubKey, 100, true)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	tstaking.CheckValidator(valAddr, sdk.Bonded, false)

	// Shield admin and depositor delegate
	tstaking.Delegate(shieldAddr, valAddr, 200000000000)
	tstaking.Delegate(delAddr, valAddr, 125000000000)

	// Shield admin deposits collateral. 200,000 CTK
	err := shieldKeeper.DepositCollateral(ctx, shieldAddr, sdk.NewInt(200000000000))
	require.Nil(t, err)

	// ShieldAdmin creates CTK Pool with Shield = 100,000 CTK, limit = 500,000 CTK, serviceFees = 200 CTK
	shield := sdk.NewCoins(sdk.NewInt64Coin(bondDenom, 100000000000))
	deposit := types.MixedCoins{Native: sdk.NewCoins(sdk.NewInt64Coin(bondDenom, 200000000))}
	shieldLimit := sdk.NewInt(500000000000)
	poolID, err := shieldKeeper.CreatePool(ctx, shieldAddr, shield, deposit, "CertiK", sponsorAddr, "fake_description", shieldLimit)
	require.Nil(t, err)
	_, found := shieldKeeper.GetPool(ctx, poolID)
	require.True(t, found)

	// Deposit collateral of amount 125,000 ctk
	var collateral int64 = 125000000000
	err = shieldKeeper.DepositCollateral(ctx, delAddr, sdk.NewInt(collateral))
	require.Nil(t, err)

	// Purchase shield
	var shieldAmt int64 = 50000000000
	purchaseShield := sdk.NewCoins(sdk.NewInt64Coin(bondDenom, shieldAmt))
	purchase, err := shieldKeeper.PurchaseShield(ctx, poolID, purchaseShield, "fake_purchase_description", purchaser)
	require.Nil(t, err)

	// Delegator undelegates all and triggers complete withdrawal
	var partial int64 = 115000000000 // fail with 116000000000
	tstaking.Undelegate(delAddr, valAddr, sdk.NewInt(partial), true)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// another 
	curTime = curTime.Add(time.Hour * 24)
	ctx = ctx.WithBlockTime(curTime)
	tstaking.TurnBlock(curTime)

	tstaking.Undelegate(delAddr, valAddr, sdk.NewInt(collateral - partial), true)
	staking.EndBlocker(ctx, app.StakingKeeper.Keeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	withdraws := shieldKeeper.GetAllWithdraws(ctx)
	require.True(t, len(withdraws) == 2)
	require.True(t, withdraws[0].Amount.Equal(sdk.NewInt(partial)))
	require.True(t, withdraws[1].Amount.Equal(sdk.NewInt(collateral - partial)))

	ubdTime1 := withdraws[0].LinkedUnbonding.CompletionTime
	timeSlice := stakingKeeper.GetUBDQueueTimeSlice(ctx, ubdTime1)
	require.True(t, len(timeSlice) == 1)
	ubdTime2 := withdraws[1].LinkedUnbonding.CompletionTime
	timeSlice = stakingKeeper.GetUBDQueueTimeSlice(ctx, ubdTime2)
	require.True(t, len(timeSlice) == 1)

	// 19 days later...
	curTime = curTime.Add(time.Hour * 24 * 19)
	ctx = ctx.WithBlockTime(curTime)

	// secure collaterals
	claimDuration := govKeeper.GetVotingParams(ctx).VotingPeriod * 2
	lossAmt := shieldAmt / 2
	loss := sdk.NewCoins(sdk.NewInt64Coin(bondDenom, lossAmt))
	err = shieldKeeper.SecureCollaterals(ctx, poolID, purchaser, purchase.PurchaseID, loss, claimDuration)
	require.Nil(t, err)

	// confirm withdraw & unbonding extensions (no splitting)
	withdraws = shieldKeeper.GetAllWithdraws(ctx)
	delayedTime := curTime.Add(claimDuration)
	
	w := withdraws[0] // not delayed
	require.True(t, w.CompletionTime.Equal(ubdTime1)) 
	require.True(t, w.LinkedUnbonding.CompletionTime.Equal(ubdTime1))
	unbonding, found := stakingKeeper.GetUnbondingDelegation(ctx, w.Address, w.LinkedUnbonding.ValidatorAddress)
	require.True(t, found)
	require.True(t, unbonding.Entries[0].CompletionTime.Equal(ubdTime1))

	w = withdraws[1] // delayed
	require.True(t, w.CompletionTime.Equal(delayedTime)) 
	require.True(t, w.LinkedUnbonding.CompletionTime.Equal(delayedTime))
	require.True(t, unbonding.Entries[1].CompletionTime.Equal(delayedTime))

	// check UBD queue
	timeSlice = stakingKeeper.GetUBDQueueTimeSlice(ctx, ubdTime1)
	require.True(t, len(timeSlice) == 1)
	timeSlice = stakingKeeper.GetUBDQueueTimeSlice(ctx, delayedTime)
	require.True(t, len(timeSlice) == 1)

	// check the purchase
	purchaseList, found := shieldKeeper.GetPurchaseList(ctx, poolID, purchaser)
	require.True(t, found)
	var index int
	for i, entry := range purchaseList.Entries {
		if entry.PurchaseID == purchase.PurchaseID {
			index = i
			break
		}
	}
	purchase = purchaseList.Entries[index]
	require.True(t, purchase.Shield.Equal(sdk.NewInt(25000000000)))
}
