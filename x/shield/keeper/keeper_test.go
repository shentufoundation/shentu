package keeper_test

/*
import (
	"encoding/hex"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/common/tests"
	"github.com/certikfoundation/shentu/simapp"

	"github.com/certikfoundation/shentu/x/gov/testgov"
	"github.com/certikfoundation/shentu/x/shield/testshield"
	"github.com/certikfoundation/shentu/x/staking/teststaking"
)

// nextBlock calls staking, shield, and gov endblockers and updates
// their test helpers' contexts.
func nextBlock(ctx sdk.Context, tstaking *teststaking.Helper, tshield *testshield.Helper, tgov *testgov.Helper) sdk.Context {
	newTime := ctx.BlockTime().Add(time.Second * time.Duration(int64(common.SecondsPerBlock)))
	ctx = ctx.WithBlockTime(newTime).WithBlockHeight(ctx.BlockHeight() + 1)

	tstaking.TurnBlock(ctx)
	tshield.TurnBlock(ctx)
	tgov.TurnBlock(ctx)

	return ctx
}

func skipBlocks(ctx sdk.Context, numBlocks int64, tstaking *teststaking.Helper, tshield *testshield.Helper, tgov *testgov.Helper) sdk.Context {
	newTime := ctx.BlockTime().Add(time.Second * time.Duration(int64(common.SecondsPerBlock)*numBlocks))
	ctx = ctx.WithBlockTime(newTime).WithBlockHeight(ctx.BlockHeight() + 1)

	tstaking.TurnBlock(ctx)
	tshield.TurnBlock(ctx)
	tgov.TurnBlock(ctx)

	return ctx
}

// TestWithdraw tests withdraws triggered by staking undelegation.
func TestWithdrawsByUndelegate(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()})

	// create and add addresses
	addresses := simapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(2e8))
	delAddr, delAddr2, accAddr, accAddr2 := addresses[0], addresses[1], addresses[2], addresses[3]

	// validator addresses
	valAddr, valAddr2 := sdk.ValAddress(accAddr), sdk.ValAddress(accAddr2)
	pubKey := tests.MakeTestPubKey()
	pubKey2 := tests.MakeTestPubKey()

	// set up testing helpers
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)
	tshield := testshield.NewHelper(t, ctx, app.ShieldKeeper, tstaking.Denom)
	tgov := testgov.NewHelper(t, ctx, app.GovKeeper, tstaking.Denom)

	// set up two validators
	tstaking.CreateValidatorWithValPower(valAddr, pubKey, 100, true)
	ctx = nextBlock(ctx, tstaking, tshield, tgov)
	tstaking.CheckValidator(valAddr, sdk.Bonded, false)

	tstaking.CreateValidatorWithValPower(valAddr2, pubKey2, 100, true)
	ctx = nextBlock(ctx, tstaking, tshield, tgov)
	tstaking.CheckValidator(valAddr2, sdk.Bonded, false)

	// attempt depositing collateral
	tshield.DepositCollateral(delAddr, 75, false)

	// both delegators delegate 50 to each validator
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

	// both delegators deposit collateral of amount 75
	tshield.DepositCollateral(delAddr, 75, true)
	tshield.DepositCollateral(delAddr2, 75, true)

	// undelegate total 50 to trigger total withdrawal of 25
	tstaking.Undelegate(delAddr, valAddr, 30, true)
	tstaking.Undelegate(delAddr2, valAddr2, 10, true)
	tstaking.Undelegate(delAddr, valAddr2, 20, true)
	tstaking.Undelegate(delAddr2, valAddr2, 40, true)

	ctx = nextBlock(ctx, tstaking, tshield, tgov)

	withdraws := app.ShieldKeeper.GetAllWithdraws(ctx)
	numWithdraws := len(withdraws)
	require.True(t, withdraws[0].Amount.Equal(sdk.NewInt(5)))
	require.True(t, withdraws[0].Address.Equals(delAddr))
	require.True(t, withdraws[1].Amount.Equal(sdk.NewInt(20)))
	require.True(t, withdraws[1].Address.Equals(delAddr))
	require.True(t, withdraws[2].Amount.Equal(sdk.NewInt(25)))
	require.True(t, withdraws[2].Address.Equals(delAddr2))

	// undelegate 5 and trigger another withdrawal of 5.
	tstaking.Undelegate(delAddr, valAddr, 5, true)

	numWithdraws++
	withdraws = app.ShieldKeeper.GetAllWithdraws(ctx)
	require.True(t, len(withdraws) == numWithdraws)
	require.True(t, withdraws[numWithdraws-1].Amount.Equal(sdk.NewInt(5)))

	// must fail deposit of 10
	tshield.DepositCollateral(delAddr, 10, false)

	// delegate 25
	tstaking.Delegate(delAddr, valAddr, 25)
	ctx = nextBlock(ctx, tstaking, tshield, tgov)

	// withdraw 5
	tshield.WithdrawCollateral(delAddr, 5, true)
	numWithdraws++
	withdraws = app.ShieldKeeper.GetAllWithdraws(ctx)
	require.True(t, len(withdraws) == numWithdraws)

	// undelegate 25 without triggering withdrawal
	tstaking.Undelegate(delAddr, valAddr, 25, true)
	ctx = nextBlock(ctx, tstaking, tshield, tgov)
	withdraws = app.ShieldKeeper.GetAllWithdraws(ctx)
	require.True(t, len(withdraws) == numWithdraws)
}

// TestWithdraw tests withdraws triggered by staking redelegation.
func TestWithdrawsByRedelegate(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()})

	// create and add addresses
	addresses := simapp.AddTestAddrs(app, ctx, 3, sdk.NewInt(2e8))
	delAddr, accAddr, accAddr2 := addresses[0], addresses[1], addresses[2]

	// validator addresses
	valAddr, valAddr2 := sdk.ValAddress(accAddr), sdk.ValAddress(accAddr2)
	pubKey := tests.MakeTestPubKey()
	pubKey2 := tests.MakeTestPubKey()

	// set up testing helpers
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)
	tshield := testshield.NewHelper(t, ctx, app.ShieldKeeper, tstaking.Denom)
	tgov := testgov.NewHelper(t, ctx, app.GovKeeper, tstaking.Denom)

	// set up two validators
	tstaking.CreateValidatorWithValPower(valAddr, pubKey, 100, true)
	ctx = nextBlock(ctx, tstaking, tshield, tgov)
	tstaking.CheckValidator(valAddr, sdk.Bonded, false)

	tstaking.CreateValidatorWithValPower(valAddr2, pubKey2, 100, true)
	ctx = nextBlock(ctx, tstaking, tshield, tgov)
	tstaking.CheckValidator(valAddr2, sdk.Bonded, false)

	// must fail at depositing collateral
	tshield.DepositCollateral(delAddr, 75, false)

	// delegate 100 to the validator
	tstaking.CheckDelegator(delAddr, valAddr, false)
	tstaking.Delegate(delAddr, valAddr, 100)
	tstaking.CheckDelegator(delAddr, valAddr, true)

	// deposit collateral of amount 75
	tshield.DepositCollateral(delAddr, 75, true)

	// redelegate 50 to trigger withdrawal of 25
	// remaining staking: 100, remaining deposit: 50
	tstaking.Redelegate(delAddr, valAddr, valAddr2, 50, true)
	ctx = nextBlock(ctx, tstaking, tshield, tgov)
	withdraws := app.ShieldKeeper.GetAllWithdraws(ctx)
	require.True(t, withdraws[0].Amount.Equal(sdk.NewInt(25)))
	numWithdraws := len(withdraws)

	// Redelegation hopping must fail
	tstaking.Redelegate(delAddr, valAddr2, valAddr, 10, false)

	// Redelegate 30 but do not trigger withdrawal
	// Remaining staking: 100, remaining deposit: 50
	tstaking.Redelegate(delAddr, valAddr, valAddr2, 30, true)
	ctx = nextBlock(ctx, tstaking, tshield, tgov)
	withdraws = app.ShieldKeeper.GetAllWithdraws(ctx)
	require.True(t, len(withdraws) == numWithdraws)

	// must fail deposit of 60
	tshield.DepositCollateral(delAddr, 60, false)

	// must succeed deposit of 50
	tshield.DepositCollateral(delAddr, 50, true)
}

// TestClaimProposal tests a claim proposal process that involves
// withdrawal and unbonding delays.
func TestClaimProposal(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()})

	// create and add addresses
	shieldAdmin := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(250e9))[0]
	app.ShieldKeeper.SetAdmin(ctx, shieldAdmin)
	sponsorAddr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1))[0]
	purchaser := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10e9))[0]
	delAddr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(125e9))[0]

	var adminDeposit int64 = 200e9
	var delegatorDeposit int64 = 125e9
	totalDeposit := adminDeposit + delegatorDeposit

	// validator addresses
	valAddr := sdk.ValAddress(simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(100e6))[0])
	pubKey := tests.MakeTestPubKey()

	// set up testing helpers
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)
	bondDenom := tstaking.Denom
	tshield := testshield.NewHelper(t, ctx, app.ShieldKeeper, bondDenom)
	tgov := testgov.NewHelper(t, ctx, app.GovKeeper, bondDenom)

	// set up a validator
	tstaking.CreateValidatorWithValPower(valAddr, pubKey, 100, true)
	ctx = nextBlock(ctx, tstaking, tshield, tgov)
	tstaking.CheckValidator(valAddr, sdk.Bonded, false)

	// shield admin deposit and create pool
	// $BondDenom pool with shield = 100,000 $BondDenom, limit = 500,000 $BondDenom, serviceFees = 200 $BondDenom
	tstaking.Delegate(shieldAdmin, valAddr, adminDeposit)
	tshield.DepositCollateral(shieldAdmin, adminDeposit, true)
	tshield.CreatePool(shieldAdmin, sponsorAddr, 200e6, 100e9, 500e9, "CertiK", "fake_description")

	pools := app.ShieldKeeper.GetAllPools(ctx)
	require.True(t, len(pools) == 1)
	require.True(t, pools[0].SponsorAddress.Equals(sponsorAddr))
	poolID := pools[0].ID

	// delegator deposits
	tstaking.CheckDelegator(delAddr, valAddr, false)
	tstaking.Delegate(delAddr, valAddr, delegatorDeposit)
	tstaking.CheckDelegator(delAddr, valAddr, true)
	tshield.DepositCollateral(delAddr, delegatorDeposit, true)

	// purchaser purhcases a shield
	var shield int64 = 50e9
	tshield.PurchaseShield(purchaser, shield, poolID, true)

	// delegator undelegates all delegations, triggering a withdrawal
	tstaking.Undelegate(delAddr, valAddr, 25e9, true)
	withdraw1End := ctx.BlockTime().Add(app.ShieldKeeper.GetPoolParams(ctx).WithdrawPeriod)
	ctx = nextBlock(ctx, tstaking, tshield, tgov)
	tstaking.Undelegate(delAddr, valAddr, 90e9, true)
	ctx = nextBlock(ctx, tstaking, tshield, tgov)
	tstaking.Undelegate(delAddr, valAddr, 10e9, true)

	withdraws := app.ShieldKeeper.GetAllWithdraws(ctx)
	require.True(t, len(withdraws) == 3)
	require.True(t, withdraws[0].Amount.Equal(sdk.NewInt(25e9)))
	require.True(t, withdraws[1].Amount.Equal(sdk.NewInt(90e9)))
	require.True(t, withdraws[2].Amount.Equal(sdk.NewInt(10e9)))

	delUBD := app.StakingKeeper.GetAllUnbondingDelegations(ctx, delAddr)[0]
	require.True(t, delUBD.Entries[0].Balance.Equal(sdk.NewInt(25e9)))
	require.True(t, delUBD.Entries[1].Balance.Equal(sdk.NewInt(90e9)))
	require.True(t, delUBD.Entries[2].Balance.Equal(sdk.NewInt(10e9)))

	// 20 days later (345,600 blocks)
	ctx = skipBlocks(ctx, 345600, tstaking, tshield, tgov)

	// the purchaser submits a claim proposal
	loss := shield
	tgov.ShieldClaimProposal(purchaser, loss, poolID, 2, true)
	var proposalID uint64 = 1 // TODO: unmarshal sdk.Result to obtain proposal ID

	// verify that the withdrawal and unbonding have been delayed
	// about 19e9 must be secured (two of three withdraws & ubds are delayed)
	withdraws = app.ShieldKeeper.GetAllWithdraws(ctx)
	delayedWithdrawEnd := ctx.BlockTime().Add(app.GovKeeper.GetVotingParams(ctx).VotingPeriod * 2)
	require.True(t, withdraws[0].Amount.Equal(sdk.NewInt(25e9)))
	require.True(t, withdraws[0].CompletionTime.Equal(withdraw1End)) //25e9 not delayed
	require.True(t, withdraws[1].Amount.Equal(sdk.NewInt(10e9)))
	require.True(t, withdraws[1].CompletionTime.Equal(delayedWithdrawEnd)) // 10e9 delayed
	require.True(t, withdraws[2].Amount.Equal(sdk.NewInt(90e9)))
	require.True(t, withdraws[2].CompletionTime.Equal(delayedWithdrawEnd)) // 90e9 delayed

	delUBD = app.StakingKeeper.GetAllUnbondingDelegations(ctx, delAddr)[0]
	require.True(t, delUBD.Entries[0].Balance.Equal(sdk.NewInt(25e9)))
	require.True(t, delUBD.Entries[0].CompletionTime.Equal(withdraw1End)) //25e9 not delayed
	require.True(t, delUBD.Entries[1].Balance.Equal(sdk.NewInt(90e9)))
	require.True(t, delUBD.Entries[1].CompletionTime.Equal(delayedWithdrawEnd)) // 90e9 delayed
	require.True(t, delUBD.Entries[2].Balance.Equal(sdk.NewInt(10e9)))
	require.True(t, delUBD.Entries[2].CompletionTime.Equal(delayedWithdrawEnd)) // 10e9 delayed

	// create reimbursement
	lossCoins := sdk.NewCoins(sdk.NewInt64Coin(bondDenom, loss))
	err := app.ShieldKeeper.CreateReimbursement(ctx, proposalID, lossCoins, purchaser)
	require.NoError(t, err)
	reimbursement, err := app.ShieldKeeper.GetReimbursement(ctx, proposalID)
	require.NoError(t, err)
	require.True(t, reimbursement.Amount.IsEqual(lossCoins))

	// confirm admin delegation reduction
	lossRatio := float64(loss) / float64(totalDeposit)
	expected := adminDeposit - int64(math.Round(float64(adminDeposit)*lossRatio))
	if hex.EncodeToString(shieldAdmin) < hex.EncodeToString(delAddr) {
		expected -= 1 // adjust for discrepancy due to sorting
	}

	adminDels := app.StakingKeeper.GetAllDelegatorDelegations(ctx, shieldAdmin)
	validator, _ := app.StakingKeeper.GetValidator(ctx, valAddr)
	require.True(t, validator.TokensFromShares(adminDels[0].Shares).Equal(sdk.NewDec(expected)))

	// confirm delegator unbonding reduction
	expected = 25e9 + 10e9 + 90e9 - int64(math.Round(float64(125e9)*lossRatio))
	if hex.EncodeToString(shieldAdmin) < hex.EncodeToString(delAddr) {
		expected += 1 // adjust for discrepancy due to sorting
	}
	withdraws = app.ShieldKeeper.GetAllWithdraws(ctx)
	require.True(t, withdraws[0].Amount.Add(withdraws[1].Amount.Add(withdraws[2].Amount)).Equal(sdk.NewInt(expected)))
	delUBD = app.StakingKeeper.GetAllUnbondingDelegations(ctx, delAddr)[0]
	require.True(t, delUBD.Entries[0].Balance.Add(delUBD.Entries[1].Balance.Add(delUBD.Entries[2].Balance)).Equal(sdk.NewInt(expected)))

	// test withdraw reimbursement
	// 56 days later (967,680 blocks)
	ctx = skipBlocks(ctx, 967680, tstaking, tshield, tgov)

	beforeInt := app.BankKeeper.GetCoins(ctx, purchaser).AmountOf(bondDenom)
	tshield.WithdrawReimbursement(purchaser, proposalID, true)
	afterInt := app.BankKeeper.GetCoins(ctx, purchaser).AmountOf(bondDenom)
	require.True(t, beforeInt.Add(sdk.NewInt(loss)).Equal(afterInt))
}
*/
