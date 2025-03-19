package keeper_test

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (suite *KeeperTestSuite) TestAddGrant() {
	// Setup a theorem for testing
	// Create a theorem first
	submitTime := suite.ctx.BlockTime()
	params, err := suite.keeper.Params.Get(suite.ctx)
	suite.Require().NoError(err)
	endTime := submitTime.Add(*params.TheoremMaxProofPeriod)
	theorem, err := types.NewTheorem(1, suite.normalAddr, "Test Theorem", "Description", "Code", submitTime, endTime)
	suite.Require().NoError(err)
	theorem.Status = types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD
	err = suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem)
	suite.Require().NoError(err)

	// Test cases
	testCases := []struct {
		name          string
		theoremID     uint64
		grantor       sdk.AccAddress
		grantAmount   sdk.Coins
		expectedError bool
		errorMessage  string
	}{
		{
			name:          "Valid grant",
			theoremID:     1,
			grantor:       suite.normalAddr,
			grantAmount:   params.MinGrant,
			expectedError: false,
		},
		{
			name:          "Theorem does not exist",
			theoremID:     999,
			grantor:       suite.normalAddr,
			grantAmount:   params.MinGrant,
			expectedError: true,
			errorMessage:  "theorem 999 doesn't exist",
		},
		{
			name:          "Grant too small",
			theoremID:     1,
			grantor:       suite.normalAddr,
			grantAmount:   sdk.NewCoins(sdk.NewCoin(params.MinGrant[0].Denom, math.NewInt(1))),
			expectedError: true,
			errorMessage:  "min grant too small",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Check initial balances
			initialBalance := suite.app.BankKeeper.GetAllBalances(suite.ctx, tc.grantor)

			// Execute the function
			err := suite.keeper.AddGrant(suite.ctx, tc.theoremID, tc.grantor, tc.grantAmount)

			if tc.expectedError {
				suite.Require().Error(err)
				if tc.errorMessage != "" {
					suite.Require().Contains(err.Error(), tc.errorMessage)
				}
			} else {
				suite.Require().NoError(err)

				// Verify theorem was updated
				updatedTheorem, err := suite.keeper.Theorems.Get(suite.ctx, tc.theoremID)
				suite.Require().NoError(err)
				suite.Require().True(sdk.NewCoins(updatedTheorem.TotalGrant...).Equal(tc.grantAmount))

				// Verify grant was created
				grant, err := suite.keeper.Grants.Get(suite.ctx, collections.Join(tc.theoremID, tc.grantor))
				suite.Require().NoError(err)
				suite.Require().Equal(tc.theoremID, grant.TheoremId)
				suite.Require().Equal(tc.grantor.String(), grant.Grantor)
				suite.Require().True(sdk.NewCoins(grant.Amount...).Equal(tc.grantAmount))

				// Verify funds were transferred
				finalBalance := suite.app.BankKeeper.GetAllBalances(suite.ctx, tc.grantor)
				expectedBalance := initialBalance.Sub(tc.grantAmount...)
				suite.Require().True(finalBalance.Equal(expectedBalance))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestAddGrantMultiple() {
	// Setup a theorem for testing
	submitTime := suite.ctx.BlockTime()
	params, err := suite.keeper.Params.Get(suite.ctx)
	suite.Require().NoError(err)
	endTime := submitTime.Add(*params.TheoremMaxProofPeriod)
	theorem, err := types.NewTheorem(2, suite.normalAddr, "Test Theorem", "Description", "Code", submitTime, endTime)
	suite.Require().NoError(err)
	theorem.Status = types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD
	err = suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem)
	suite.Require().NoError(err)

	// Add first grant
	err = suite.keeper.AddGrant(suite.ctx, theorem.Id, suite.normalAddr, params.MinGrant)
	suite.Require().NoError(err)

	// Verify first grant
	firstGrant, err := suite.keeper.Grants.Get(suite.ctx, collections.Join(theorem.Id, suite.normalAddr))
	suite.Require().NoError(err)
	suite.Require().True(sdk.NewCoins(firstGrant.Amount...).Equal(params.MinGrant))

	// Add second grant from same grantor
	err = suite.keeper.AddGrant(suite.ctx, theorem.Id, suite.normalAddr, params.MinGrant)
	suite.Require().NoError(err)

	// Verify grant was updated
	updatedGrant, err := suite.keeper.Grants.Get(suite.ctx, collections.Join(theorem.Id, suite.normalAddr))
	suite.Require().NoError(err)
	expectedAmount := sdk.NewCoins(params.MinGrant...).Add(params.MinGrant...)
	suite.Require().True(sdk.NewCoins(updatedGrant.Amount...).Equal(expectedAmount))

	// Verify theorem total was updated
	updatedTheorem, err := suite.keeper.Theorems.Get(suite.ctx, theorem.Id)
	suite.Require().NoError(err)
	suite.Require().True(sdk.NewCoins(updatedTheorem.TotalGrant...).Equal(expectedAmount))
}

func (suite *KeeperTestSuite) TestRefundAndDeleteGrants() {
	// Setup a theorem with grants
	submitTime := suite.ctx.BlockTime()
	params, err := suite.keeper.Params.Get(suite.ctx)
	suite.Require().NoError(err)
	endTime := submitTime.Add(*params.TheoremMaxProofPeriod)
	theorem, err := types.NewTheorem(3, suite.normalAddr, "Test Theorem", "Description", "Code", submitTime, endTime)
	suite.Require().NoError(err)
	theorem.Status = types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD
	err = suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem)
	suite.Require().NoError(err)

	// Add grants from multiple addresses
	grantors := []sdk.AccAddress{suite.normalAddr, suite.whiteHatAddr, suite.programAddr}
	for _, grantor := range grantors {
		initialBalance := suite.app.BankKeeper.GetAllBalances(suite.ctx, grantor)
		suite.Require().NoError(suite.keeper.AddGrant(suite.ctx, theorem.Id, grantor, params.MinGrant))

		// Verify funds were transferred
		afterGrantBalance := suite.app.BankKeeper.GetAllBalances(suite.ctx, grantor)
		expectedBalance := initialBalance.Sub(params.MinGrant...)
		suite.Require().True(afterGrantBalance.Equal(expectedBalance))
	}

	// Record balances before refund
	balancesBeforeRefund := make(map[string]sdk.Coins)
	for _, grantor := range grantors {
		balancesBeforeRefund[grantor.String()] = suite.app.BankKeeper.GetAllBalances(suite.ctx, grantor)
	}

	// Execute refund function
	err = suite.keeper.RefundAndDeleteGrants(suite.ctx, theorem.Id)
	suite.Require().NoError(err)

	// Verify grants were deleted and funds were returned
	for _, grantor := range grantors {
		// Check grant no longer exists
		_, err := suite.keeper.Grants.Get(suite.ctx, collections.Join(theorem.Id, grantor))
		suite.Require().Error(err)
		suite.Require().True(errors.IsOf(err, collections.ErrNotFound))

		// Check funds returned
		afterRefundBalance := suite.app.BankKeeper.GetAllBalances(suite.ctx, grantor)
		expectedBalance := balancesBeforeRefund[grantor.String()].Add(params.MinGrant...)
		suite.Require().True(afterRefundBalance.Equal(expectedBalance))
	}
}

func (suite *KeeperTestSuite) TestIterateGrants() {
	// Setup a theorem with grants
	submitTime := suite.ctx.BlockTime()
	params, err := suite.keeper.Params.Get(suite.ctx)
	suite.Require().NoError(err)
	endTime := submitTime.Add(*params.TheoremMaxProofPeriod)
	theorem, err := types.NewTheorem(4, suite.normalAddr, "Test Theorem", "Description", "Code", submitTime, endTime)
	suite.Require().NoError(err)
	theorem.Status = types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD
	err = suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem)
	suite.Require().NoError(err)

	// Add grants from multiple addresses
	grantors := []sdk.AccAddress{suite.normalAddr, suite.whiteHatAddr, suite.programAddr}
	for _, grantor := range grantors {
		suite.Require().NoError(suite.keeper.AddGrant(suite.ctx, theorem.Id, grantor, params.MinGrant))
	}

	// Count the grants using iterator
	var count int
	err = suite.keeper.IterateGrants(suite.ctx, theorem.Id, func(_ collections.Pair[uint64, sdk.AccAddress], _ types.Grant) (bool, error) {
		count++
		return false, nil
	})
	suite.Require().NoError(err)
	suite.Require().Equal(len(grantors), count)

	// Verify the grants using iterator
	grantorsFound := make(map[string]bool)
	err = suite.keeper.IterateGrants(suite.ctx, theorem.Id, func(key collections.Pair[uint64, sdk.AccAddress], grant types.Grant) (bool, error) {
		suite.Require().Equal(theorem.Id, grant.TheoremId)
		suite.Require().Equal(theorem.Id, key.K1())
		suite.Require().Equal(grant.Grantor, key.K2().String())
		suite.Require().True(sdk.NewCoins(grant.Amount...).Equal(params.MinGrant))
		grantorsFound[grant.Grantor] = true
		return false, nil
	})
	suite.Require().NoError(err)

	// All grantors should be found
	for _, grantor := range grantors {
		suite.Require().True(grantorsFound[grantor.String()])
	}
}

func (suite *KeeperTestSuite) TestDistributionGrants() {
	// Setup a theorem with grants
	submitTime := suite.ctx.BlockTime()
	params, err := suite.keeper.Params.Get(suite.ctx)
	suite.Require().NoError(err)
	endTime := submitTime.Add(*params.TheoremMaxProofPeriod)
	theorem, err := types.NewTheorem(5, suite.normalAddr, "Test Theorem", "Description", "Code", submitTime, endTime)
	suite.Require().NoError(err)
	theorem.Status = types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD
	err = suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem)
	suite.Require().NoError(err)

	// Add a grant
	grantAmount := sdk.NewCoins(sdk.NewCoin(params.MinGrant[0].Denom, math.NewInt(1000000)))
	suite.Require().NoError(suite.keeper.AddGrant(suite.ctx, theorem.Id, suite.normalAddr, grantAmount))

	// Update theorem with the grant
	theorem.TotalGrant = grantAmount
	err = suite.keeper.Theorems.Set(suite.ctx, theorem.Id, theorem)
	suite.Require().NoError(err)

	// Define checker and prover
	checker := suite.whiteHatAddr
	prover := suite.programAddr

	// Verify no rewards before distribution
	_, err = suite.keeper.Rewards.Get(suite.ctx, checker)
	suite.Require().Error(err)
	suite.Require().True(errors.IsOf(err, collections.ErrNotFound))

	_, err = suite.keeper.Rewards.Get(suite.ctx, prover)
	suite.Require().Error(err)
	suite.Require().True(errors.IsOf(err, collections.ErrNotFound))

	// Distribute grants
	err = suite.keeper.DistributionGrants(suite.ctx, theorem.Id, checker, prover)
	suite.Require().NoError(err)

	// Verify rewards after distribution
	checkerReward, err := suite.keeper.Rewards.Get(suite.ctx, checker)
	suite.Require().NoError(err)

	proverReward, err := suite.keeper.Rewards.Get(suite.ctx, prover)
	suite.Require().NoError(err)

	// Calculate expected rewards
	totalGrantDec := sdk.NewDecCoinsFromCoins(grantAmount...)
	expectedCheckerReward := totalGrantDec.MulDec(params.CheckerRate)
	expectedProverReward := totalGrantDec.Sub(expectedCheckerReward)

	// Verify rewards match expected values
	suite.Require().True(checkerReward.Reward.Equal(expectedCheckerReward))
	suite.Require().True(proverReward.Reward.Equal(expectedProverReward))
}

func (suite *KeeperTestSuite) TestSetGrant() {
	// Create a new grant
	theoremID := uint64(6)
	grantor := suite.normalAddr

	// Get current params
	params, err := suite.keeper.Params.Get(suite.ctx)
	suite.Require().NoError(err)

	grantAmount := params.MinGrant
	grant := types.NewGrant(theoremID, grantor, grantAmount)

	// Set the grant
	err = suite.keeper.SetGrant(suite.ctx, grant)
	suite.Require().NoError(err)

	// Verify the grant was set correctly
	storedGrant, err := suite.keeper.Grants.Get(suite.ctx, collections.Join(theoremID, grantor))
	suite.Require().NoError(err)

	suite.Require().Equal(grant.TheoremId, storedGrant.TheoremId)
	suite.Require().Equal(grant.Grantor, storedGrant.Grantor)
	suite.Require().True(sdk.NewCoins(grant.Amount...).Equal(sdk.NewCoins(storedGrant.Amount...)))
}

func (suite *KeeperTestSuite) TestMinGrantValidation() {
	// Get current params
	params, err := suite.keeper.Params.Get(suite.ctx)
	suite.Require().NoError(err)

	// Test with valid grant
	validGrant := params.MinGrant
	err = suite.keeper.AddGrant(suite.ctx, 0, suite.normalAddr, validGrant)
	// We expect error about theorem not existing, not about min grant
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "doesn't exist")

	// Test with too small grant
	smallGrant := sdk.NewCoins(sdk.NewCoin(params.MinGrant[0].Denom, math.NewInt(1)))
	err = suite.keeper.AddGrant(suite.ctx, 0, suite.normalAddr, smallGrant)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "min grant too small")

	// Test with invalid denom
	invalidDenom := sdk.NewCoins(sdk.NewCoin("invalid", math.NewInt(1000000)))
	err = suite.keeper.AddGrant(suite.ctx, 0, suite.normalAddr, invalidDenom)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid deposit denom")
}
