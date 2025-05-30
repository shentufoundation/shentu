package keeper_test

import (
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

	// Create a theorem
	theoremID := uint64(5)
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
	grantAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1001)))
	err = suite.keeper.AddGrant(suite.ctx, theoremID, suite.normalAddr, grantAmount)
	require.NoError(suite.T(), err)

	// Check module account balance increased after grant
	moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
	expectedModuleBalance := sdk.NewCoin(bondDenom, math.NewInt(1001))
	require.True(suite.T(), moduleBalance.Equal(expectedModuleBalance))

	// Distribute rewards
	checker := suite.whiteHatAddr
	prover := suite.programAddr
	err = suite.keeper.DistributionGrants(suite.ctx, theoremID, checker, prover)
	require.NoError(suite.T(), err)

	// Get parameters to check distribution ratio
	params, err := suite.keeper.Params.Get(suite.ctx)
	require.NoError(suite.T(), err)

	// Verify checker's reward
	checkerReward, err := suite.keeper.Rewards.Get(suite.ctx, checker)
	require.NoError(suite.T(), err)
	expectedCheckerReward := sdk.NewDecCoinsFromCoins(grantAmount...).MulDec(params.CheckerRate)
	require.True(suite.T(), expectedCheckerReward.Equal(checkerReward.Reward))

	// Verify prover's reward
	proverReward, err := suite.keeper.Rewards.Get(suite.ctx, prover)
	require.NoError(suite.T(), err)
	expectedProverReward := sdk.NewDecCoinsFromCoins(grantAmount...).Sub(expectedCheckerReward)
	require.True(suite.T(), expectedProverReward.Equal(proverReward.Reward))

	// Verify module account balance remains the same after distribution
	// (since funds aren't actually transferred, just recorded as rewards)
	moduleBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
	require.True(suite.T(), moduleBalance.Equal(moduleBalanceAfter))
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

// TestValidateFunds tests fund validation functionality
func (suite *KeeperTestSuite) TestValidateFunds() {
	// Get bond denom
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	require.NoError(suite.T(), err)

	// Test valid grant amount
	validGrantAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(100)))
	_, err = suite.keeper.ValidateFunds(suite.ctx, validGrantAmount, "grant")
	require.NoError(suite.T(), err)

	// Test invalid grant amount (too small)
	invalidGrantAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(10)))
	_, err = suite.keeper.ValidateFunds(suite.ctx, invalidGrantAmount, "grant")
	require.Error(suite.T(), err)

	// Test valid deposit amount
	validDepositAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(50)))
	_, err = suite.keeper.ValidateFunds(suite.ctx, validDepositAmount, "deposit")
	require.NoError(suite.T(), err)

	// Test invalid deposit amount (too small)
	invalidDepositAmount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(10)))
	_, err = suite.keeper.ValidateFunds(suite.ctx, invalidDepositAmount, "deposit")
	require.Error(suite.T(), err)
}
