package keeper_test

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// TestAddGrant tests the functionality of adding a grant
func (suite *KeeperTestSuite) TestAddGrant() {
	// First, create a theorem
	theoremID := uint64(1)
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	theorem := types.Theorem{
		Id:          theoremID,
		Title:       "Test Theorem",
		Description: "Test Description",
		Proposer:    suite.programAddr.String(),
		Status:      types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		TotalGrant:  sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(0))),
	}
	require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem))

	// Add a grant
	grantAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(100)))
	err = suite.keeper.AddGrant(suite.ctx, theoremID, suite.normalAddr, grantAmount)
	require.NoError(suite.T(), err)

	// Verify grant record
	grant, err := suite.keeper.Grants.Get(suite.ctx, collections.Join(theoremID, suite.normalAddr))
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), theoremID, grant.TheoremId)
	require.Equal(suite.T(), suite.normalAddr.String(), grant.Grantor)
	require.True(suite.T(), grantAmount.Equal(grant.Amount))

	// Verify theorem's total grant
	updatedTheorem, err := suite.keeper.Theorems.Get(suite.ctx, theoremID)
	require.NoError(suite.T(), err)
	require.True(suite.T(), grantAmount.Equal(updatedTheorem.TotalGrant))

	// Test adding a second grant from the same address
	secondGrantAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(50)))
	err = suite.keeper.AddGrant(suite.ctx, theoremID, suite.normalAddr, secondGrantAmount)
	require.NoError(suite.T(), err)

	// Verify updated grant record - amounts should be added together
	updatedGrant, err := suite.keeper.Grants.Get(suite.ctx, collections.Join(theoremID, suite.normalAddr))
	require.NoError(suite.T(), err)
	expectedTotalAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(150))) // 100 + 50
	require.True(suite.T(), expectedTotalAmount.Equal(updatedGrant.Amount), "expected: %v, got: %v", expectedTotalAmount, updatedGrant.Amount)

	// Verify theorem's total grant has been updated correctly
	updatedTheoremAgain, err := suite.keeper.Theorems.Get(suite.ctx, theoremID)
	require.NoError(suite.T(), err)
	require.True(suite.T(), expectedTotalAmount.Equal(updatedTheoremAgain.TotalGrant), "expected: %v, got: %v", expectedTotalAmount, updatedTheoremAgain.TotalGrant)

	// Test adding a grant from a different address
	differentAddress := suite.whiteHatAddr
	differentGrantAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(200)))
	err = suite.keeper.AddGrant(suite.ctx, theoremID, differentAddress, differentGrantAmount)
	require.NoError(suite.T(), err)

	// Check module account balance increased
	moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
	require.Equal(suite.T(), math.NewInt(350), moduleBalance.Amount) // 150 + 200

	// Verify grant record for the different address
	differentGrant, err := suite.keeper.Grants.Get(suite.ctx, collections.Join(theoremID, differentAddress))
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), theoremID, differentGrant.TheoremId)
	require.Equal(suite.T(), differentAddress.String(), differentGrant.Grantor)
	require.True(suite.T(), differentGrantAmount.Equal(differentGrant.Amount), "expected: %v, got: %v", differentGrantAmount, differentGrant.Amount)

	// Verify theorem's total grant has been updated with both addresses' contributions
	finalTheorem, err := suite.keeper.Theorems.Get(suite.ctx, theoremID)
	require.NoError(suite.T(), err)
	expectedFinalAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(350))) // 150 + 200
	require.True(suite.T(), expectedFinalAmount.Equal(finalTheorem.TotalGrant), "expected: %v, got: %v", expectedFinalAmount, finalTheorem.TotalGrant)

	// Test adding grant for non-existent theorem
	invalidTheoremID := uint64(999)
	err = suite.keeper.AddGrant(suite.ctx, invalidTheoremID, suite.normalAddr, grantAmount)
	require.Error(suite.T(), err)

	// Test adding grant for a closed theorem
	closedTheorem := theorem
	closedTheorem.Id = uint64(2)
	closedTheorem.Status = types.TheoremStatus_THEOREM_STATUS_CLOSED
	require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, closedTheorem.Id, closedTheorem))

	err = suite.keeper.AddGrant(suite.ctx, closedTheorem.Id, suite.normalAddr, grantAmount)
	require.Error(suite.T(), err)
}

// TestAddGrantEdgeCases tests edge cases and error handling for AddGrant
func (suite *KeeperTestSuite) TestAddGrantEdgeCases() {
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	// Create a valid theorem for testing
	theoremID := uint64(500)
	theorem := types.Theorem{
		Id:          theoremID,
		Title:       "Test Theorem Edge Cases",
		Description: "Test Description",
		Proposer:    suite.programAddr.String(),
		Status:      types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		TotalGrant:  sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(0))),
	}
	require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem))

	// Test 1: Grant with large amount (within account balance)
	largeAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1000000000))) // 1B stake, within initial balance of 10B
	err = suite.keeper.AddGrant(suite.ctx, theoremID, suite.normalAddr, largeAmount)
	require.NoError(suite.T(), err)

	// Verify large amount was recorded correctly
	grant, err := suite.keeper.Grants.Get(suite.ctx, collections.Join(theoremID, suite.normalAddr))
	require.NoError(suite.T(), err)
	require.True(suite.T(), largeAmount.Equal(grant.Amount))

	// Test 2: Grant with theorem in different invalid statuses
	testCases := []struct {
		name   string
		status types.TheoremStatus
	}{
		{"Closed status", types.TheoremStatus_THEOREM_STATUS_CLOSED},
	}

	for i, tc := range testCases {
		invalidTheorem := types.Theorem{
			Id:          uint64(501 + i),
			Title:       fmt.Sprintf("Invalid Theorem %d", i),
			Description: "Test Description",
			Proposer:    suite.programAddr.String(),
			Status:      tc.status,
			TotalGrant:  sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(0))),
		}
		require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, invalidTheorem.Id, invalidTheorem))

		grantAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(100)))
		err = suite.keeper.AddGrant(suite.ctx, invalidTheorem.Id, suite.normalAddr, grantAmount)
		require.Error(suite.T(), err, "should fail for status: %s", tc.name)
		require.ErrorIs(suite.T(), err, types.ErrTheoremProposal)
	}

	// Test 3: Grant from account with exact balance (edge case)
	// Create a new account with exact balance
	exactBalanceAddr := sdk.AccAddress([]byte("exact_balance_addr"))
	exactAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1000)))
	err = suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, exactAmount)
	require.NoError(suite.T(), err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, exactBalanceAddr, exactAmount)
	require.NoError(suite.T(), err)

	// Create another theorem for this test
	exactBalanceTheoremID := uint64(510)
	exactBalanceTheorem := types.Theorem{
		Id:          exactBalanceTheoremID,
		Title:       "Exact Balance Theorem",
		Description: "Test Description",
		Proposer:    suite.programAddr.String(),
		Status:      types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		TotalGrant:  sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(0))),
	}
	require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, exactBalanceTheorem.Id, exactBalanceTheorem))

	// Grant exact balance
	err = suite.keeper.AddGrant(suite.ctx, exactBalanceTheoremID, exactBalanceAddr, exactAmount)
	require.NoError(suite.T(), err)

	// Verify balance is now zero
	balance := suite.app.BankKeeper.GetBalance(suite.ctx, exactBalanceAddr, bondDenom)
	require.True(suite.T(), balance.Amount.IsZero())

	// Test 5: Try to grant more than available balance (should fail)
	insufficientAddr := sdk.AccAddress([]byte("insufficient_addr"))
	smallBalance := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(50)))
	err = suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, smallBalance)
	require.NoError(suite.T(), err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, insufficientAddr, smallBalance)
	require.NoError(suite.T(), err)

	largeGrantAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1000)))
	err = suite.keeper.AddGrant(suite.ctx, theoremID, insufficientAddr, largeGrantAmount)
	require.Error(suite.T(), err) // Should fail due to insufficient balance
}

// TestRefundAndDeleteGrants tests refunding and deleting all grants
func (suite *KeeperTestSuite) TestRefundAndDeleteGrants() {
	// Get bond denom
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	// Set initial balance
	initialBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.normalAddr, bondDenom)

	// Create a theorem
	theoremID := uint64(3)
	theorem := types.Theorem{
		Id:          theoremID,
		Title:       "Test Theorem",
		Description: "Test Description",
		Proposer:    suite.programAddr.String(),
		Status:      types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		TotalGrant:  sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(0))),
	}
	require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem))

	// Add a grant
	grantAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(100)))
	err = suite.keeper.AddGrant(suite.ctx, theoremID, suite.normalAddr, grantAmount)
	require.NoError(suite.T(), err)

	// Check module account balance increased
	moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
	require.Equal(suite.T(), math.NewInt(100), moduleBalance.Amount)

	// Check user balance decreased
	userBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.normalAddr, bondDenom)
	require.Equal(suite.T(), initialBalance.Amount.Sub(math.NewInt(100)), userBalance.Amount)

	// Add a second grantor
	err = suite.keeper.AddGrant(suite.ctx, theoremID, suite.whiteHatAddr, grantAmount)
	require.NoError(suite.T(), err)

	// Check module account balance increased again
	moduleBalance = suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
	require.Equal(suite.T(), math.NewInt(200), moduleBalance.Amount) // 100 + 100

	// Refund and delete all grants
	err = suite.keeper.RefundAndDeleteGrants(suite.ctx, theoremID)
	require.NoError(suite.T(), err)

	// Check module account balance is zero after refund
	moduleBalance = suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
	require.Equal(suite.T(), math.NewInt(0), moduleBalance.Amount)

	// Verify grant records are deleted
	_, err = suite.keeper.Grants.Get(suite.ctx, collections.Join(theoremID, suite.normalAddr))
	require.True(suite.T(), errors.IsOf(err, collections.ErrNotFound))

	_, err = suite.keeper.Grants.Get(suite.ctx, collections.Join(theoremID, suite.whiteHatAddr))
	require.True(suite.T(), errors.IsOf(err, collections.ErrNotFound))

	// Verify user balance restored
	finalUserBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.normalAddr, bondDenom)
	require.Equal(suite.T(), initialBalance.Amount, finalUserBalance.Amount)
}

// TestIterateGrants tests iterating over grangts
func (suite *KeeperTestSuite) TestIterateGrants() {
	// Get bond denom
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	// Create a theorem
	theoremID := uint64(4)
	theorem := types.Theorem{
		Id:          theoremID,
		Title:       "Test Theorem",
		Description: "Test Description",
		Proposer:    suite.programAddr.String(),
		Status:      types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		TotalGrant:  sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(0))),
	}
	require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem))

	// Add multiple grants
	grantors := []sdk.AccAddress{suite.normalAddr, suite.whiteHatAddr, suite.programAddr}
	amounts := []int64{100, 200, 300}
	total := int64(0)

	for i, grantor := range grantors {
		grantAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(amounts[i])))
		err := suite.keeper.AddGrant(suite.ctx, theoremID, grantor, grantAmount)
		require.NoError(suite.T(), err)
		total += amounts[i]
	}

	// Verify using IterateGrants
	count := 0
	totalAmount := int64(0)
	err = suite.keeper.IterateGrants(suite.ctx, theoremID, func(key collections.Pair[uint64, sdk.AccAddress], grant types.Grant) (bool, error) {
		count++
		coin := grant.Amount[0]
		totalAmount += coin.Amount.Int64()
		return false, nil
	})
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), len(grantors), count)
	require.Equal(suite.T(), total, totalAmount)
}

// TestDistributionGrants tests distributing grants as rewards
func (suite *KeeperTestSuite) TestDistributionGrants() {
	// Get bond denom
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	// Get parameters
	params, err := suite.keeper.Params.Get(suite.ctx)
	require.NoError(suite.T(), err)

	// Create a reference theorem (to be imported)
	refTheoremID := uint64(100)
	refTheorem := types.Theorem{
		Id:            refTheoremID,
		Title:         "Imported Theorem",
		Description:   "A theorem that will be imported",
		Proposer:      suite.normalAddr.String(),
		Status:        types.TheoremStatus_THEOREM_STATUS_CLOSED,
		Complexity:    50, // Reference theorem complexity
		ImportedCount: 0,  // Initially not cited
		TotalGrant:    sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(0))),
	}
	require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, refTheorem.Id, refTheorem))

	// Create a theorem with import
	theoremID := uint64(5)
	complexity := int64(10)
	theorem := types.Theorem{
		Id:          theoremID,
		Title:       "Test Theorem",
		Description: "Test Description",
		Proposer:    suite.programAddr.String(),
		Status:      types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		Complexity:  complexity,
		Imports:     []uint64{refTheoremID}, // Reference the first theorem
		TotalGrant:  sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(0))),
	}
	require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem))

	// Add a grant
	// Grant amount must be sufficient for:
	// - Checker rewards: complexity (10) * complexityFee (10000) = 100000
	// - Imported rewards: refComplexity (50) / (importedCount+1) * complexityFee = 50 * 10000 = 500000
	// - Prover rewards: remainder
	// Total needed: at least 600000, using 1000000 to leave room for prover
	grantAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1000000)))
	err = suite.keeper.AddGrant(suite.ctx, theoremID, suite.whiteHatAddr, grantAmount)
	require.NoError(suite.T(), err)

	// Update theorem with the grant
	theorem.TotalGrant = grantAmount

	// Check module account balance increased after grant
	moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
	expectedModuleBalance := sdk.NewCoin(bondDenom, math.NewInt(1000000))
	require.True(suite.T(), moduleBalance.Equal(expectedModuleBalance))

	// Distribute rewards
	checker := suite.whiteHatAddr
	prover := suite.programAddr
	err = suite.keeper.DistributionGrants(suite.ctx, theorem, checker, prover)
	require.NoError(suite.T(), err)

	// Calculate expected rewards based on actual implementation
	totalGrant := sdk.NewDecCoinsFromCoins(grantAmount...)
	complexityFeeAmount := math.LegacyNewDecFromInt(params.ComplexityFee.Amount)

	// 1. Checker rewards: complexity * complexityFee
	expectedCheckerRewardAmount := complexityFeeAmount.MulInt64(complexity)
	expectedCheckerReward := sdk.NewDecCoins(sdk.NewDecCoinFromDec(params.ComplexityFee.Denom, expectedCheckerRewardAmount))

	// 2. Imported rewards: (Complexity / (ImportedCount + 1)) * ComplexityFee
	complexityDec := math.LegacyNewDec(refTheorem.Complexity)
	importedCountDec := math.LegacyNewDec(refTheorem.ImportedCount + 1) // 0 + 1 = 1
	normalizedComplexity := complexityDec.Quo(importedCountDec)
	expectedImportedRewardAmount := complexityFeeAmount.Mul(normalizedComplexity)
	expectedImportedReward := sdk.NewDecCoins(sdk.NewDecCoinFromDec(params.ComplexityFee.Denom, expectedImportedRewardAmount))

	// 3. Prover rewards: remaining after checker and imported rewards
	expectedProverReward := totalGrant.Sub(expectedCheckerReward).Sub(expectedImportedReward)

	// Verify checker's reward
	checkerReward, err := suite.keeper.Rewards.Get(suite.ctx, checker)
	require.NoError(suite.T(), err)
	require.True(suite.T(), expectedCheckerReward.Equal(checkerReward.Reward),
		"expected checker reward: %v, got: %v", expectedCheckerReward, checkerReward.Reward)

	// Verify imported reward for reference theorem proposer
	importedReward, err := suite.keeper.ImportedRewards.Get(suite.ctx, suite.normalAddr)
	require.NoError(suite.T(), err)
	require.True(suite.T(), expectedImportedReward.Equal(importedReward.Reward),
		"expected imported reward: %v, got: %v", expectedImportedReward, importedReward.Reward)

	// Verify prover's reward
	proverReward, err := suite.keeper.Rewards.Get(suite.ctx, prover)
	require.NoError(suite.T(), err)
	require.True(suite.T(), expectedProverReward.Equal(proverReward.Reward),
		"expected prover reward: %v, got: %v", expectedProverReward, proverReward.Reward)

	// Verify reference theorem imported count was incremented
	updatedRefTheorem, err := suite.keeper.Theorems.Get(suite.ctx, refTheoremID)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), int64(1), updatedRefTheorem.ImportedCount)

	// Verify module account balance remains the same after distribution
	// (since funds aren't actually transferred, just recorded as rewards)
	moduleBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
	require.True(suite.T(), moduleBalance.Equal(moduleBalanceAfter))
}

// TestDistributionGrantsInsufficientFunds tests distribution with insufficient grant amounts
func (suite *KeeperTestSuite) TestDistributionGrantsInsufficientFunds() {
	// Get bond denom
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	// Get parameters
	params, err := suite.keeper.Params.Get(suite.ctx)
	require.NoError(suite.T(), err)

	// Test Case 1: Insufficient funds for checker rewards
	refTheoremID1 := uint64(200)
	refTheorem1 := types.Theorem{
		Id:            refTheoremID1,
		Title:         "Reference Theorem 1",
		Description:   "A reference theorem",
		Proposer:      suite.normalAddr.String(),
		Status:        types.TheoremStatus_THEOREM_STATUS_CLOSED,
		Complexity:    10,
		ImportedCount: 0,
		TotalGrant:    sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(0))),
	}
	require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, refTheorem1.Id, refTheorem1))

	theoremID1 := uint64(201)
	complexity1 := int64(100) // High complexity
	theorem1 := types.Theorem{
		Id:          theoremID1,
		Title:       "Test Theorem Insufficient Checker",
		Description: "Test insufficient funds for checker",
		Proposer:    suite.programAddr.String(),
		Status:      types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		Complexity:  complexity1,
		Imports:     []uint64{},
		// Checker needs: 100 * 10000 = 1,000,000
		// Grant only: 500,000
		TotalGrant: sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(500000))),
	}
	require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, theorem1.Id, theorem1))

	err = suite.keeper.DistributionGrants(suite.ctx, theorem1, suite.whiteHatAddr, suite.programAddr)
	require.Error(suite.T(), err)
	require.ErrorIs(suite.T(), err, types.ErrInsufficientGrantChecker)

	// Test Case 2: Insufficient funds for imported rewards
	refTheoremID2 := uint64(202)
	refTheorem2 := types.Theorem{
		Id:            refTheoremID2,
		Title:         "Reference Theorem 2",
		Description:   "A reference theorem with high complexity",
		Proposer:      suite.normalAddr.String(),
		Status:        types.TheoremStatus_THEOREM_STATUS_CLOSED,
		Complexity:    200, // High complexity reference
		ImportedCount: 0,
		TotalGrant:    sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(0))),
	}
	require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, refTheorem2.Id, refTheorem2))

	theoremID2 := uint64(203)
	complexity2 := int64(10) // Low complexity
	theorem2 := types.Theorem{
		Id:          theoremID2,
		Title:       "Test Theorem Insufficient Imported",
		Description: "Test insufficient funds for imported rewards",
		Proposer:    suite.programAddr.String(),
		Status:      types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		Complexity:  complexity2,
		Imports:     []uint64{refTheoremID2},
		// Checker needs: 10 * 10000 = 100,000
		// Imported needs: 200 * 10000 = 2,000,000
		// Total needed: 2,100,000
		// Grant only: 1,500,000
		TotalGrant: sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1500000))),
	}
	require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, theorem2.Id, theorem2))

	err = suite.keeper.DistributionGrants(suite.ctx, theorem2, suite.whiteHatAddr, suite.programAddr)
	require.Error(suite.T(), err)
	require.ErrorIs(suite.T(), err, types.ErrInsufficientGrantTotal)

	// Test Case 3: Exactly sufficient funds (boundary case)
	theoremID3 := uint64(204)
	complexity3 := int64(10)
	// Checker: 10 * 10000 = 100,000
	// No imports, so prover gets 0
	theorem3 := types.Theorem{
		Id:          theoremID3,
		Title:       "Test Theorem Exact Funds",
		Description: "Test exact sufficient funds",
		Proposer:    suite.programAddr.String(),
		Status:      types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		Complexity:  complexity3,
		Imports:     []uint64{},
		TotalGrant:  sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(100000))),
	}
	require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, theorem3.Id, theorem3))

	err = suite.keeper.DistributionGrants(suite.ctx, theorem3, suite.whiteHatAddr, suite.programAddr)
	require.NoError(suite.T(), err)

	// Verify checker got exact amount
	checkerReward, err := suite.keeper.Rewards.Get(suite.ctx, suite.whiteHatAddr)
	require.NoError(suite.T(), err)
	expectedCheckerAmount := math.LegacyNewDecFromInt(params.ComplexityFee.Amount).MulInt64(complexity3)
	expectedChecker := sdk.NewDecCoins(sdk.NewDecCoinFromDec(params.ComplexityFee.Denom, expectedCheckerAmount))
	require.True(suite.T(), expectedChecker.Equal(checkerReward.Reward))

	// Verify prover got zero
	proverReward, err := suite.keeper.Rewards.Get(suite.ctx, suite.programAddr)
	require.NoError(suite.T(), err)
	require.True(suite.T(), proverReward.Reward.IsZero())
}

// TestDistributionGrantsMultipleImports tests distribution with multiple imported theorems
func (suite *KeeperTestSuite) TestDistributionGrantsMultipleImports() {
	// Get bond denom
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	// Get parameters
	params, err := suite.keeper.Params.Get(suite.ctx)
	require.NoError(suite.T(), err)

	// Create 5 reference theorems with different complexities and proposers
	refTheoremIDs := []uint64{300, 301, 302, 303, 304}
	refComplexities := []int64{10, 20, 30, 40, 50}
	refProposers := []sdk.AccAddress{
		suite.normalAddr,
		suite.whiteHatAddr,
		suite.programAddr,
		suite.normalAddr, // Reuse address to test multiple rewards to same address
		suite.whiteHatAddr,
	}

	for i, refID := range refTheoremIDs {
		refTheorem := types.Theorem{
			Id:            refID,
			Title:         fmt.Sprintf("Reference Theorem %d", i),
			Description:   "A reference theorem",
			Proposer:      refProposers[i].String(),
			Status:        types.TheoremStatus_THEOREM_STATUS_CLOSED,
			Complexity:    refComplexities[i],
			ImportedCount: int64(i), // Different imported counts
			TotalGrant:    sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(0))),
		}
		require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, refTheorem.Id, refTheorem))
	}

	// Create main theorem that imports all reference theorems
	theoremID := uint64(305)
	complexity := int64(20)
	theorem := types.Theorem{
		Id:          theoremID,
		Title:       "Test Theorem Multiple Imports",
		Description: "Test with multiple imported theorems",
		Proposer:    suite.programAddr.String(),
		Status:      types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		Complexity:  complexity,
		Imports:     refTheoremIDs,
		// Grant should be large enough for all rewards
		TotalGrant: sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(10000000))),
	}
	require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem))

	// Add grant
	grantAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(10000000)))
	err = suite.keeper.AddGrant(suite.ctx, theoremID, suite.whiteHatAddr, grantAmount)
	require.NoError(suite.T(), err)

	// Distribute rewards
	checker := suite.whiteHatAddr
	prover := suite.normalAddr // Different from checker
	err = suite.keeper.DistributionGrants(suite.ctx, theorem, checker, prover)
	require.NoError(suite.T(), err)

	// Calculate expected rewards
	totalGrant := sdk.NewDecCoinsFromCoins(grantAmount...)
	complexityFeeAmount := math.LegacyNewDecFromInt(params.ComplexityFee.Amount)

	// 1. Expected checker reward
	expectedCheckerAmount := complexityFeeAmount.MulInt64(complexity)
	expectedCheckerReward := sdk.NewDecCoins(sdk.NewDecCoinFromDec(params.ComplexityFee.Denom, expectedCheckerAmount))

	// 2. Expected imported rewards for each reference
	expectedImportedRewards := make(map[string]sdk.DecCoins)
	totalImportedReward := sdk.NewDecCoins()
	for i, refID := range refTheoremIDs {
		// Calculate: Complexity / (ImportedCount + 1) * ComplexityFee
		complexityDec := math.LegacyNewDec(refComplexities[i])
		normalizedComplexity := complexityDec.QuoInt64(int64(i) + 1) // original ImportedCount + 1
		refRewardAmount := complexityFeeAmount.Mul(normalizedComplexity)
		refReward := sdk.NewDecCoinFromDec(params.ComplexityFee.Denom, refRewardAmount)

		proposerStr := refProposers[i].String()
		if existing, ok := expectedImportedRewards[proposerStr]; ok {
			expectedImportedRewards[proposerStr] = existing.Add(refReward)
		} else {
			expectedImportedRewards[proposerStr] = sdk.NewDecCoins(refReward)
		}
		totalImportedReward = totalImportedReward.Add(refReward)

		// Verify imported count was incremented
		updatedRefTheorem, err := suite.keeper.Theorems.Get(suite.ctx, refID)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), int64(i)+1, updatedRefTheorem.ImportedCount)
	}

	// 3. Expected prover reward
	expectedProverReward := totalGrant.Sub(expectedCheckerReward).Sub(totalImportedReward)

	// Verify checker reward
	checkerReward, err := suite.keeper.Rewards.Get(suite.ctx, checker)
	require.NoError(suite.T(), err)
	require.True(suite.T(), expectedCheckerReward.Equal(checkerReward.Reward),
		"expected checker: %v, got: %v", expectedCheckerReward, checkerReward.Reward)

	// Verify imported rewards for each proposer
	for proposerStr, expectedReward := range expectedImportedRewards {
		proposerAddr, err := sdk.AccAddressFromBech32(proposerStr)
		require.NoError(suite.T(), err)
		importedReward, err := suite.keeper.ImportedRewards.Get(suite.ctx, proposerAddr)
		require.NoError(suite.T(), err)
		require.True(suite.T(), expectedReward.Equal(importedReward.Reward),
			"proposer %s: expected imported: %v, got: %v", proposerStr, expectedReward, importedReward.Reward)
	}

	// Verify prover reward
	proverReward, err := suite.keeper.Rewards.Get(suite.ctx, prover)
	require.NoError(suite.T(), err)
	require.True(suite.T(), expectedProverReward.Equal(proverReward.Reward),
		"expected prover: %v, got: %v", expectedProverReward, proverReward.Reward)
}

// TestDistributionGrantsNoImports tests distribution without any imported theorems
func (suite *KeeperTestSuite) TestDistributionGrantsNoImports() {
	// Get bond denom
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	// Get parameters
	params, err := suite.keeper.Params.Get(suite.ctx)
	require.NoError(suite.T(), err)

	// Create theorem without imports
	theoremID := uint64(400)
	complexity := int64(50)
	grantAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(2000000)))
	theorem := types.Theorem{
		Id:          theoremID,
		Title:       "Test Theorem No Imports",
		Description: "Test without imported theorems",
		Proposer:    suite.programAddr.String(),
		Status:      types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD,
		Complexity:  complexity,
		Imports:     []uint64{}, // No imports
		TotalGrant:  grantAmount,
	}
	require.NoError(suite.T(), suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem))

	// Add grant
	err = suite.keeper.AddGrant(suite.ctx, theoremID, suite.whiteHatAddr, grantAmount)
	require.NoError(suite.T(), err)

	// Distribute rewards
	checker := suite.whiteHatAddr
	prover := suite.normalAddr
	err = suite.keeper.DistributionGrants(suite.ctx, theorem, checker, prover)
	require.NoError(suite.T(), err)

	// Calculate expected rewards
	totalGrant := sdk.NewDecCoinsFromCoins(grantAmount...)
	complexityFeeAmount := math.LegacyNewDecFromInt(params.ComplexityFee.Amount)

	// Checker reward
	expectedCheckerAmount := complexityFeeAmount.MulInt64(complexity)
	expectedCheckerReward := sdk.NewDecCoins(sdk.NewDecCoinFromDec(params.ComplexityFee.Denom, expectedCheckerAmount))

	// Prover gets everything else (no imported rewards)
	expectedProverReward := totalGrant.Sub(expectedCheckerReward)

	// Verify rewards
	checkerReward, err := suite.keeper.Rewards.Get(suite.ctx, checker)
	require.NoError(suite.T(), err)
	require.True(suite.T(), expectedCheckerReward.Equal(checkerReward.Reward))

	proverReward, err := suite.keeper.Rewards.Get(suite.ctx, prover)
	require.NoError(suite.T(), err)
	require.True(suite.T(), expectedProverReward.Equal(proverReward.Reward))
}

// TestAddDeposit tests adding a deposit for a proof
func (suite *KeeperTestSuite) TestAddDeposit() {
	// Get bond denom
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	// Create a proof
	proofID := "test-proof-id"
	proof := types.Proof{
		Id:        proofID,
		TheoremId: uint64(6),
		Status:    types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD,
		Prover:    suite.programAddr.String(),
	}
	require.NoError(suite.T(), suite.keeper.Proofs.Set(suite.ctx, proof.Id, proof))

	// Set initial balance
	initialBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.normalAddr, bondDenom)

	// Add deposit
	depositAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(100)))
	err = suite.keeper.AddDeposit(suite.ctx, proofID, suite.normalAddr, depositAmount)
	require.NoError(suite.T(), err)

	// Verify deposit record
	deposit, err := suite.keeper.Deposits.Get(suite.ctx, collections.Join(proofID, suite.normalAddr))
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), proofID, deposit.ProofId)
	require.Equal(suite.T(), suite.normalAddr.String(), deposit.Depositor)
	require.True(suite.T(), depositAmount.Equal(deposit.Amount))

	// Check module account balance increased
	moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
	expectedModuleBalance := sdk.NewCoin(bondDenom, math.NewInt(100))
	require.True(suite.T(), moduleBalance.Equal(expectedModuleBalance))

	// Check user balance decreased
	userBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.normalAddr, bondDenom)
	require.Equal(suite.T(), initialBalance.Amount.Sub(math.NewInt(100)), userBalance.Amount)

	// Test adding deposit for non-existent proof
	invalidProofID := "invalid-proof-id"
	err = suite.keeper.AddDeposit(suite.ctx, invalidProofID, suite.normalAddr, depositAmount)
	require.Error(suite.T(), err)

	// Test adding deposit for proof with invalid status
	invalidStatusProof := proof
	invalidStatusProof.Id = "invalid-status-proof"
	invalidStatusProof.Status = types.ProofStatus_PROOF_STATUS_HASH_DETAIL_PERIOD
	require.NoError(suite.T(), suite.keeper.Proofs.Set(suite.ctx, invalidStatusProof.Id, invalidStatusProof))

	err = suite.keeper.AddDeposit(suite.ctx, invalidStatusProof.Id, suite.normalAddr, depositAmount)
	require.Error(suite.T(), err)
}

// TestAddDepositEdgeCases tests edge cases and error handling for AddDeposit
func (suite *KeeperTestSuite) TestAddDepositEdgeCases() {
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	// Create a valid proof for testing
	proofID := "test-deposit-edge-cases-proof"
	proof := types.Proof{
		Id:        proofID,
		TheoremId: uint64(600),
		Status:    types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD,
		Prover:    suite.programAddr.String(),
	}
	require.NoError(suite.T(), suite.keeper.Proofs.Set(suite.ctx, proof.Id, proof))

	// Test 1: Deposit with large amount (within account balance)
	largeAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1000000000))) // 1B stake, within initial balance of 10B
	err = suite.keeper.AddDeposit(suite.ctx, proofID, suite.normalAddr, largeAmount)
	require.NoError(suite.T(), err)

	// Verify large amount was recorded correctly
	deposit, err := suite.keeper.Deposits.Get(suite.ctx, collections.Join(proofID, suite.normalAddr))
	require.NoError(suite.T(), err)
	require.True(suite.T(), largeAmount.Equal(deposit.Amount))

	// Test 2: Multiple deposits from same depositor
	secondDeposit := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(5000)))
	err = suite.keeper.AddDeposit(suite.ctx, proofID, suite.normalAddr, secondDeposit)
	require.NoError(suite.T(), err)

	// Verify amounts are added together
	updatedDeposit, err := suite.keeper.Deposits.Get(suite.ctx, collections.Join(proofID, suite.normalAddr))
	require.NoError(suite.T(), err)
	expectedTotal := largeAmount.Add(secondDeposit...)
	require.True(suite.T(), expectedTotal.Equal(updatedDeposit.Amount))

	// Test 3: Deposit with proof in different invalid statuses
	testCases := []struct {
		name   string
		status types.ProofStatus
	}{
		{"Hash detail period status", types.ProofStatus_PROOF_STATUS_HASH_DETAIL_PERIOD},
	}

	for i, tc := range testCases {
		invalidProof := types.Proof{
			Id:        fmt.Sprintf("invalid-proof-%d", i),
			TheoremId: uint64(601 + i),
			Status:    tc.status,
			Prover:    suite.programAddr.String(),
		}
		require.NoError(suite.T(), suite.keeper.Proofs.Set(suite.ctx, invalidProof.Id, invalidProof))

		depositAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(100)))
		err = suite.keeper.AddDeposit(suite.ctx, invalidProof.Id, suite.normalAddr, depositAmount)
		require.Error(suite.T(), err, "should fail for status: %s", tc.name)
		require.ErrorIs(suite.T(), err, types.ErrProofStatusInvalid)
	}

	// Test 4: Deposit from account with exact balance
	exactBalanceAddr := sdk.AccAddress([]byte("exact_deposit_addr"))
	exactAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1000)))
	err = suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, exactAmount)
	require.NoError(suite.T(), err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, exactBalanceAddr, exactAmount)
	require.NoError(suite.T(), err)

	// Create another proof for this test
	exactProofID := "exact-balance-deposit-proof"
	exactProof := types.Proof{
		Id:        exactProofID,
		TheoremId: uint64(610),
		Status:    types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD,
		Prover:    suite.programAddr.String(),
	}
	require.NoError(suite.T(), suite.keeper.Proofs.Set(suite.ctx, exactProof.Id, exactProof))

	// Deposit exact balance
	err = suite.keeper.AddDeposit(suite.ctx, exactProofID, exactBalanceAddr, exactAmount)
	require.NoError(suite.T(), err)

	// Verify balance is now zero
	balance := suite.app.BankKeeper.GetBalance(suite.ctx, exactBalanceAddr, bondDenom)
	require.True(suite.T(), balance.Amount.IsZero())

	// Test 5: Try to deposit more than available balance (should fail)
	insufficientAddr := sdk.AccAddress([]byte("insufficient_deposit_addr"))
	smallBalance := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(30)))
	err = suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, smallBalance)
	require.NoError(suite.T(), err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, insufficientAddr, smallBalance)
	require.NoError(suite.T(), err)

	largeDepositAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1000)))
	err = suite.keeper.AddDeposit(suite.ctx, proofID, insufficientAddr, largeDepositAmount)
	require.Error(suite.T(), err) // Should fail due to insufficient balance

	// Test 6: Multiple depositors for same proof
	depositors := []sdk.AccAddress{suite.whiteHatAddr, suite.programAddr}
	for _, depositor := range depositors {
		amount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(500)))
		err = suite.keeper.AddDeposit(suite.ctx, proofID, depositor, amount)
		require.NoError(suite.T(), err)
	}

	// Verify all deposits exist
	for _, depositor := range depositors {
		deposit, err := suite.keeper.Deposits.Get(suite.ctx, collections.Join(proofID, depositor))
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), proofID, deposit.ProofId)
		require.Equal(suite.T(), depositor.String(), deposit.Depositor)
	}
}

// TestRefundAndDeleteDeposit tests refunding and deleting a deposit
func (suite *KeeperTestSuite) TestRefundAndDeleteDeposit() {
	// Get bond denom
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	// Create a proof
	proofID := "test-refund-proof-id"
	proof := types.Proof{
		Id:        proofID,
		TheoremId: uint64(6),
		Status:    types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD,
		Prover:    suite.programAddr.String(),
	}
	require.NoError(suite.T(), suite.keeper.Proofs.Set(suite.ctx, proof.Id, proof))

	// Set initial balance
	initialBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.normalAddr, bondDenom)

	// Add deposit
	depositAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(100)))
	err = suite.keeper.AddDeposit(suite.ctx, proofID, suite.normalAddr, depositAmount)
	require.NoError(suite.T(), err)

	// Check module account balance increased
	moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
	expectedModuleBalance := sdk.NewCoin(bondDenom, math.NewInt(100))
	require.True(suite.T(), moduleBalance.Equal(expectedModuleBalance))

	// Refund and delete the deposit
	err = suite.keeper.RefundAndDeleteDeposit(suite.ctx, proofID, suite.normalAddr)
	require.NoError(suite.T(), err)

	// Verify deposit record is deleted
	_, err = suite.keeper.Deposits.Get(suite.ctx, collections.Join(proofID, suite.normalAddr))
	require.True(suite.T(), errors.IsOf(err, collections.ErrNotFound))

	// Check module account balance is zero after refund
	moduleBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
	expectedModuleBalanceAfter := sdk.NewCoin(bondDenom, math.NewInt(0))
	require.True(suite.T(), moduleBalanceAfter.Equal(expectedModuleBalanceAfter))

	// Verify user balance restored
	finalUserBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.normalAddr, bondDenom)
	require.Equal(suite.T(), initialBalance.Amount, finalUserBalance.Amount)
}

// TestRefundAndDeleteDeposits tests refunding and deleting all deposits for a proof
func (suite *KeeperTestSuite) TestRefundAndDeleteDeposits() {
	// Get bond denom
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	// Create a proof
	proofID := "test-refund-deposits-proof-id"
	proof := types.Proof{
		Id:        proofID,
		TheoremId: uint64(7),
		Status:    types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD,
		Prover:    suite.programAddr.String(),
	}
	require.NoError(suite.T(), suite.keeper.Proofs.Set(suite.ctx, proof.Id, proof))

	// Add multiple deposits
	depositAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(100)))
	depositors := []sdk.AccAddress{suite.normalAddr, suite.whiteHatAddr}

	for _, depositor := range depositors {
		err := suite.keeper.AddDeposit(suite.ctx, proofID, depositor, depositAmount)
		require.NoError(suite.T(), err)
	}

	// Check module account balance increased after deposits
	moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
	expectedModuleBalance := sdk.NewCoin(bondDenom, math.NewInt(200)) // 100 * 2
	require.True(suite.T(), moduleBalance.Equal(expectedModuleBalance))

	// Refund and delete all deposits
	err = suite.keeper.RefundAndDeleteDeposits(suite.ctx, proofID)
	require.NoError(suite.T(), err)

	// Verify all deposit records are deleted
	for _, depositor := range depositors {
		_, err = suite.keeper.Deposits.Get(suite.ctx, collections.Join(proofID, depositor))
		require.True(suite.T(), errors.IsOf(err, collections.ErrNotFound))
	}

	// Check module account balance is zero after refund
	moduleBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
	expectedModuleBalanceAfter := sdk.NewCoin(bondDenom, math.NewInt(0))
	require.True(suite.T(), moduleBalanceAfter.Equal(expectedModuleBalanceAfter))
}

// TestIterateDeposits tests iterating over deposits for a proof
func (suite *KeeperTestSuite) TestIterateDeposits() {
	// Get bond denom
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	// Create a proof
	proofID := "test-iterate-deposits-proof-id"
	proof := types.Proof{
		Id:        proofID,
		TheoremId: uint64(8),
		Status:    types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD,
		Prover:    suite.programAddr.String(),
	}
	require.NoError(suite.T(), suite.keeper.Proofs.Set(suite.ctx, proof.Id, proof))

	// Add multiple deposits
	depositors := []sdk.AccAddress{suite.normalAddr, suite.whiteHatAddr, suite.programAddr}
	amounts := []int64{100, 200, 300}
	total := int64(0)

	for i, depositor := range depositors {
		depositAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(amounts[i])))
		err := suite.keeper.AddDeposit(suite.ctx, proofID, depositor, depositAmount)
		require.NoError(suite.T(), err)
		total += amounts[i]
	}

	// Verify using IterateDeposits
	count := 0
	totalAmount := int64(0)
	err = suite.keeper.IterateDeposits(suite.ctx, proofID, func(key collections.Pair[string, sdk.AccAddress], deposit types.Deposit) (bool, error) {
		count++
		coin := deposit.Amount[0]
		totalAmount += coin.Amount.Int64()
		return false, nil
	})
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), len(depositors), count)
	require.Equal(suite.T(), total, totalAmount)

	// Test iterating over empty deposit list
	emptyProofID := "empty-proof-id"
	emptyProof := types.Proof{
		Id:        emptyProofID,
		TheoremId: uint64(9),
		Status:    types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD,
		Prover:    suite.programAddr.String(),
	}
	require.NoError(suite.T(), suite.keeper.Proofs.Set(suite.ctx, emptyProof.Id, emptyProof))

	emptyCount := 0
	err = suite.keeper.IterateDeposits(suite.ctx, emptyProofID, func(key collections.Pair[string, sdk.AccAddress], deposit types.Deposit) (bool, error) {
		emptyCount++
		return false, nil
	})
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 0, emptyCount)

	// Test early termination by returning true
	stopCount := 0
	err = suite.keeper.IterateDeposits(suite.ctx, proofID, func(key collections.Pair[string, sdk.AccAddress], deposit types.Deposit) (bool, error) {
		stopCount++
		return stopCount >= 2, nil // Stop after 2 iterations
	})
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 2, stopCount)

	// Test callback returning error
	errorCount := 0
	expectedErr := fmt.Errorf("callback error")
	err = suite.keeper.IterateDeposits(suite.ctx, proofID, func(key collections.Pair[string, sdk.AccAddress], deposit types.Deposit) (bool, error) {
		errorCount++
		if errorCount == 2 {
			return false, expectedErr
		}
		return false, nil
	})
	require.Error(suite.T(), err)
	require.Equal(suite.T(), 2, errorCount)
}

// TestValidateFunds tests fund validation functionality
func (suite *KeeperTestSuite) TestValidateFunds() {
	// Get bond denom
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	// Test valid grant amount
	validGrantAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(100)))
	_, err = suite.keeper.ValidateFunds(suite.ctx, validGrantAmount, types.FundTypeGrant)
	require.NoError(suite.T(), err)

	// Test invalid grant amount (too small)
	invalidGrantAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(10)))
	_, err = suite.keeper.ValidateFunds(suite.ctx, invalidGrantAmount, types.FundTypeGrant)
	require.Error(suite.T(), err)
	require.ErrorIs(suite.T(), err, types.ErrMinGrantTooSmall)

	// Test valid deposit amount
	validDepositAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(50)))
	_, err = suite.keeper.ValidateFunds(suite.ctx, validDepositAmount, types.FundTypeDeposit)
	require.NoError(suite.T(), err)

	// Test invalid deposit amount (too small)
	invalidDepositAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(10)))
	_, err = suite.keeper.ValidateFunds(suite.ctx, invalidDepositAmount, types.FundTypeDeposit)
	require.Error(suite.T(), err)
	require.ErrorIs(suite.T(), err, types.ErrMinDepositTooSmall)

	// Test invalid denomination
	invalidDenom := sdk.NewCoins(sdk.NewCoin("invalid_denom", math.NewInt(100)))
	_, err = suite.keeper.ValidateFunds(suite.ctx, invalidDenom, types.FundTypeGrant)
	require.Error(suite.T(), err)
	require.ErrorIs(suite.T(), err, types.ErrInvalidDepositDenom)

	// Test negative amount
	negativeAmount := sdk.Coins{sdk.Coin{Denom: bondDenom, Amount: math.NewInt(-100)}}
	_, err = suite.keeper.ValidateFunds(suite.ctx, negativeAmount, types.FundTypeGrant)
	require.Error(suite.T(), err)

	// Test invalid funds type
	_, err = suite.keeper.ValidateFunds(suite.ctx, validGrantAmount, "invalid_type")
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "invalid funds type")
}

// TestValidateFundsExtended tests extended validation scenarios
func (suite *KeeperTestSuite) TestValidateFundsExtended() {
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	// Test 1: Empty coins (nil)
	var nilCoins sdk.Coins
	_, err = suite.keeper.ValidateFunds(suite.ctx, nilCoins, types.FundTypeGrant)
	require.Error(suite.T(), err)

	// Test 2: Empty coins (initialized but empty)
	emptyCoins := sdk.NewCoins()
	_, err = suite.keeper.ValidateFunds(suite.ctx, emptyCoins, types.FundTypeGrant)
	require.Error(suite.T(), err)

	// Test 3: Amount exactly equal to minimum grant
	params, err := suite.keeper.Params.Get(suite.ctx)
	require.NoError(suite.T(), err)

	exactMinGrant := params.MinGrant
	_, err = suite.keeper.ValidateFunds(suite.ctx, exactMinGrant, types.FundTypeGrant)
	require.NoError(suite.T(), err)

	// Test 4: Amount exactly equal to minimum deposit
	exactMinDeposit := params.MinDeposit
	_, err = suite.keeper.ValidateFunds(suite.ctx, exactMinDeposit, types.FundTypeDeposit)
	require.NoError(suite.T(), err)

	// Test 5: Amount slightly below minimum grant
	if len(params.MinGrant) > 0 {
		belowMinGrant := sdk.NewCoins(sdk.NewCoin(params.MinGrant[0].Denom, params.MinGrant[0].Amount.SubRaw(1)))
		_, err = suite.keeper.ValidateFunds(suite.ctx, belowMinGrant, types.FundTypeGrant)
		require.Error(suite.T(), err)
		require.ErrorIs(suite.T(), err, types.ErrMinGrantTooSmall)
	}

	// Test 6: Amount slightly below minimum deposit
	if len(params.MinDeposit) > 0 {
		belowMinDeposit := sdk.NewCoins(sdk.NewCoin(params.MinDeposit[0].Denom, params.MinDeposit[0].Amount.SubRaw(1)))
		_, err = suite.keeper.ValidateFunds(suite.ctx, belowMinDeposit, types.FundTypeDeposit)
		require.Error(suite.T(), err)
		require.ErrorIs(suite.T(), err, types.ErrMinDepositTooSmall)
	}

	// Test 7: Multiple valid denominations (if params allow)
	// Assume params allow bondDenom
	multiDenomValid := sdk.NewCoins(
		sdk.NewCoin(bondDenom, math.NewInt(100)),
	)
	_, err = suite.keeper.ValidateFunds(suite.ctx, multiDenomValid, types.FundTypeGrant)
	require.NoError(suite.T(), err)

	// Test 8: Mixed valid and invalid denominations
	mixedDenoms := sdk.NewCoins(
		sdk.NewCoin(bondDenom, math.NewInt(100)),
		sdk.NewCoin("invalid_denom", math.NewInt(50)),
	)
	_, err = suite.keeper.ValidateFunds(suite.ctx, mixedDenoms, types.FundTypeGrant)
	require.Error(suite.T(), err)
	require.ErrorIs(suite.T(), err, types.ErrInvalidDepositDenom)

	// Test 9: First denomination invalid
	firstInvalidDenoms := sdk.NewCoins(
		sdk.NewCoin("invalid_first", math.NewInt(100)),
		sdk.NewCoin(bondDenom, math.NewInt(100)),
	)
	_, err = suite.keeper.ValidateFunds(suite.ctx, firstInvalidDenoms, types.FundTypeGrant)
	require.Error(suite.T(), err)
	require.ErrorIs(suite.T(), err, types.ErrInvalidDepositDenom)

	// Test 10: Last denomination invalid
	lastInvalidDenoms := sdk.NewCoins(
		sdk.NewCoin(bondDenom, math.NewInt(100)),
		sdk.NewCoin("invalid_last", math.NewInt(50)),
	)
	_, err = suite.keeper.ValidateFunds(suite.ctx, lastInvalidDenoms, types.FundTypeGrant)
	require.Error(suite.T(), err)
	require.ErrorIs(suite.T(), err, types.ErrInvalidDepositDenom)

	// Test 11: Zero amount (should be invalid)
	zeroAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(0)))
	_, err = suite.keeper.ValidateFunds(suite.ctx, zeroAmount, types.FundTypeGrant)
	require.Error(suite.T(), err)

	// Test 12: Very large amount
	veryLargeAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1000000000000000)))
	_, err = suite.keeper.ValidateFunds(suite.ctx, veryLargeAmount, types.FundTypeGrant)
	require.NoError(suite.T(), err)

	// Test 13: Consolidated coin amount (testing that a single large coin works)
	// Note: sdk.NewCoins panics on duplicate denominations, so we test with already consolidated amount
	consolidatedAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(110))) // 50 + 60
	_, err = suite.keeper.ValidateFunds(suite.ctx, consolidatedAmount, types.FundTypeGrant)
	require.NoError(suite.T(), err) // Should succeed with valid consolidated amount

	// Test 14: Invalid coins format
	invalidCoins := sdk.Coins{sdk.Coin{Denom: bondDenom, Amount: math.NewInt(-50)}}
	_, err = suite.keeper.ValidateFunds(suite.ctx, invalidCoins, types.FundTypeGrant)
	require.Error(suite.T(), err)

	// Test 15: Deposit type with grant minimum
	// This may succeed or fail depending on whether MinGrant >= MinDeposit
	// If MinDeposit < MinGrant, this should succeed
	// We'll just verify it doesn't panic
	_, _ = suite.keeper.ValidateFunds(suite.ctx, params.MinGrant, types.FundTypeDeposit)

	// Test 16: Grant type with deposit minimum
	// This may succeed or fail depending on whether MinDeposit >= MinGrant
	// We'll just verify it doesn't panic
	_, _ = suite.keeper.ValidateFunds(suite.ctx, params.MinDeposit, types.FundTypeGrant)
}
