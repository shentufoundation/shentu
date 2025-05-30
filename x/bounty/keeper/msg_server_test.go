package keeper_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/google/uuid"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (suite *KeeperTestSuite) TestCreateProgram() {
	type args struct {
		msgCreatePrograms []types.MsgCreateProgram
	}

	type errArgs struct {
		shouldPass bool
	}

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{"Program(1)  -> Set: Simple",
			args{
				msgCreatePrograms: []types.MsgCreateProgram{
					{
						Name:            "Name",
						Detail:          "detail",
						OperatorAddress: suite.programAddr.String(),
						ProgramId:       "1",
					},
				},
			},
			errArgs{
				shouldPass: true,
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			for _, program := range tc.args.msgCreatePrograms {
				ctx := sdk.WrapSDKContext(suite.ctx)

				_, err := suite.msgServer.CreateProgram(ctx, &program)
				suite.Require().NoError(err)

				// Directly retrieve Program using collections
				storedProgram, err := suite.keeper.Programs.Get(suite.ctx, program.ProgramId)

				if tc.errArgs.shouldPass {
					suite.Require().NoError(err)
					suite.Require().Equal(storedProgram.ProgramId, program.ProgramId)
				} else {
					suite.Require().Error(err)
				}
			}
		})
	}

	// Test for non-existent program
	suite.Run("Program -> Get: Non-existent", func() {
		nonExistentProgramID := "non-existent-id"

		// Try to retrieve a program that doesn't exist
		_, err := suite.keeper.Programs.Get(suite.ctx, nonExistentProgramID)

		// Verify that we get a "not found" error
		suite.Require().Error(err)
		suite.Require().True(errors.IsOf(err, collections.ErrNotFound), "Expected ErrNotFound, got: %v", err)
	})
}

// TestEditProgram tests the EditProgram message handler
func (suite *KeeperTestSuite) TestEditProgram() {
	// Create a test program first
	pid := uuid.NewString()
	suite.InitCreateProgram(pid)

	testCases := []struct {
		name    string
		req     *types.MsgEditProgram
		expPass bool
	}{
		{
			"empty request",
			&types.MsgEditProgram{},
			false,
		},
		{
			"invalid program ID",
			&types.MsgEditProgram{
				ProgramId:       "non-existent-id",
				OperatorAddress: suite.programAddr.String(),
				Name:            "Updated Name",
				Detail:          "Updated Detail",
			},
			false,
		},
		{
			"unauthorized operator",
			&types.MsgEditProgram{
				ProgramId:       pid,
				OperatorAddress: suite.normalAddr.String(), // Not program admin or bounty admin
				Name:            "Updated Name",
				Detail:          "Updated Detail",
			},
			false,
		},
		{
			"valid request - program admin",
			&types.MsgEditProgram{
				ProgramId:       pid,
				OperatorAddress: suite.programAddr.String(), // Program admin
				Name:            "Updated Name",
				Detail:          "Updated Detail",
			},
			true,
		},
		{
			"valid request - bounty admin",
			&types.MsgEditProgram{
				ProgramId:       pid,
				OperatorAddress: suite.bountyAdminAddr.String(), // Bounty admin
				Name:            "Updated Name 2",
				Detail:          "Updated Detail 2",
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.EditProgram(ctx, testCase.req)

			if testCase.expPass {
				suite.Require().NoError(err)

				// Verify changes were applied
				program, err := suite.keeper.Programs.Get(suite.ctx, pid)
				suite.Require().NoError(err)

				if testCase.req.Name != "" {
					suite.Require().Equal(testCase.req.Name, program.Name)
				}

				if testCase.req.Detail != "" {
					suite.Require().Equal(testCase.req.Detail, program.Detail)
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

// TestCloseProgram tests the CloseProgram message handler
func (suite *KeeperTestSuite) TestCloseProgram() {
	// Create and activate a program first
	pid := uuid.NewString()
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)

	testCases := []struct {
		name    string
		req     *types.MsgCloseProgram
		expPass bool
	}{
		{
			"empty request",
			&types.MsgCloseProgram{},
			false,
		},
		{
			"invalid program ID",
			&types.MsgCloseProgram{
				ProgramId:       "non-existent-id",
				OperatorAddress: suite.programAddr.String(),
			},
			false,
		},
		{
			"unauthorized operator",
			&types.MsgCloseProgram{
				ProgramId:       pid,
				OperatorAddress: suite.normalAddr.String(), // Not program admin or bounty admin
			},
			false,
		},
		{
			"valid request - program admin",
			&types.MsgCloseProgram{
				ProgramId:       pid,
				OperatorAddress: suite.programAddr.String(), // Program admin
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			// If the previous test was successful, we need to reactivate the program
			if pid == testCase.req.ProgramId {
				hasProgram, err := suite.keeper.Programs.Has(suite.ctx, pid)
				suite.Require().NoError(err)

				if hasProgram {
					program, err := suite.keeper.Programs.Get(suite.ctx, pid)
					suite.Require().NoError(err)

					if program.Status == types.ProgramStatusClosed {
						// Reactivate the program
						program.Status = types.ProgramStatusActive
						err = suite.keeper.Programs.Set(suite.ctx, pid, program)
						suite.Require().NoError(err)
					}
				}
			}

			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.CloseProgram(ctx, testCase.req)

			if testCase.expPass {
				suite.Require().NoError(err)

				// Verify the program was closed
				program, err := suite.keeper.Programs.Get(suite.ctx, pid)
				suite.Require().NoError(err)
				suite.Require().Equal(types.ProgramStatusClosed, program.Status)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSubmitFinding() {
	// Create and activate a program first
	pid := uuid.NewString()
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)

	desc, poc := "desc", "poc"
	hash := sha256.Sum256([]byte(desc + poc + suite.whiteHatAddr.String()))

	testCases := []struct {
		name    string
		req     *types.MsgSubmitFinding
		expPass bool
	}{
		{
			"empty request",
			&types.MsgSubmitFinding{},
			false,
		},
		{
			"invalid program ID",
			&types.MsgSubmitFinding{
				ProgramId:       "non-existent-id",
				FindingId:       "1",
				OperatorAddress: suite.whiteHatAddr.String(),
				SeverityLevel:   types.Critical,
			},
			false,
		},
		{
			"successful submission",
			&types.MsgSubmitFinding{
				ProgramId:       pid,
				FindingId:       uuid.NewString(),
				FindingHash:     hex.EncodeToString(hash[:]),
				OperatorAddress: suite.whiteHatAddr.String(),
				SeverityLevel:   types.Critical,
			},
			true,
		},
		{
			"duplicate finding ID",
			&types.MsgSubmitFinding{
				ProgramId:       pid,
				FindingId:       "duplicate-id",
				FindingHash:     hex.EncodeToString(hash[:]),
				OperatorAddress: suite.whiteHatAddr.String(),
				SeverityLevel:   types.Critical,
			},
			true, // First submission should pass
		},
		{
			"duplicate finding ID - second attempt",
			&types.MsgSubmitFinding{
				ProgramId:       pid,
				FindingId:       "duplicate-id",
				OperatorAddress: suite.whiteHatAddr.String(),
				SeverityLevel:   types.Critical,
			},
			false, // Second submission with same ID should fail
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.SubmitFinding(ctx, testCase.req)

			if testCase.expPass {
				suite.Require().NoError(err)

				// Verify finding was created
				finding, err := suite.keeper.Findings.Get(suite.ctx, testCase.req.FindingId)
				suite.Require().NoError(err)
				suite.Require().Equal(testCase.req.ProgramId, finding.ProgramId)
				suite.Require().Equal(testCase.req.FindingId, finding.FindingId)
				suite.Require().Equal(testCase.req.OperatorAddress, finding.SubmitterAddress)
				suite.Require().Equal(testCase.req.SeverityLevel, finding.SeverityLevel)
				suite.Require().Equal(types.FindingStatusSubmitted, finding.Status)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

// TestEditFinding tests the EditFinding message handler
func (suite *KeeperTestSuite) TestEditFinding() {
	// Create a program and submit a finding
	pid, fid := uuid.NewString(), uuid.NewString()
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)
	suite.InitSubmitFinding(pid, fid)

	testCases := []struct {
		name    string
		req     *types.MsgEditFinding
		expPass bool
	}{
		{
			"empty request",
			&types.MsgEditFinding{},
			false,
		},
		{
			"invalid finding ID",
			&types.MsgEditFinding{
				FindingId:       "non-existent-id",
				OperatorAddress: suite.whiteHatAddr.String(),
				FindingHash:     "updated-hash",
				SeverityLevel:   types.High,
			},
			false,
		},
		{
			"unauthorized operator",
			&types.MsgEditFinding{
				FindingId:       fid,
				OperatorAddress: suite.normalAddr.String(), // Not the submitter
				FindingHash:     "updated-hash",
				SeverityLevel:   types.High,
			},
			false,
		},
		{
			"valid request",
			&types.MsgEditFinding{
				FindingId:       fid,
				OperatorAddress: suite.whiteHatAddr.String(), // Original submitter
				FindingHash:     "updated-hash",
				SeverityLevel:   types.High,
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.EditFinding(ctx, testCase.req)

			if testCase.expPass {
				suite.Require().NoError(err)

				// Verify changes were applied
				finding, err := suite.keeper.Findings.Get(suite.ctx, fid)
				suite.Require().NoError(err)

				if testCase.req.FindingHash != "" {
					suite.Require().Equal(testCase.req.FindingHash, finding.FindingHash)
				}

				if testCase.req.SeverityLevel != types.Unspecified {
					suite.Require().Equal(testCase.req.SeverityLevel, finding.SeverityLevel)
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestActivateFinding() {
	pid, fid := uuid.NewString(), uuid.NewString()
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)
	suite.InitSubmitFinding(pid, fid)

	// Directly verify finding existence using collections
	_, err := suite.keeper.Findings.Get(suite.ctx, fid)
	suite.Require().NoError(err)

	testCases := []struct {
		name    string
		req     *types.MsgActivateFinding
		expPass bool
	}{
		{
			"empty request",
			&types.MsgActivateFinding{},
			false,
		},
		{
			"valid request",
			&types.MsgActivateFinding{
				FindingId:       fid,
				OperatorAddress: suite.bountyAdminAddr.String(),
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.ActivateFinding(ctx, testCase.req)

			if testCase.expPass {
				suite.Require().NoError(err)
				finding, err := suite.keeper.Findings.Get(suite.ctx, fid)
				suite.Require().NoError(err)
				suite.Require().Equal(finding.Status, types.FindingStatusActive)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestConfirmFinding() {
	pid, fid := uuid.NewString(), uuid.NewString()
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)
	suite.InitSubmitFinding(pid, fid)
	suite.InitActivateFinding(fid)

	finding, err := suite.keeper.Findings.Get(suite.ctx, fid)
	suite.Require().NoError(err)
	findingFingerPrintHash := suite.app.BountyKeeper.GetFindingFingerprintHash(&finding)

	testCases := []struct {
		name    string
		req     *types.MsgConfirmFinding
		expPass bool
	}{
		{
			"empty request",
			&types.MsgConfirmFinding{},
			false,
		},
		{
			"valid request => ",
			&types.MsgConfirmFinding{
				FindingId:       fid,
				OperatorAddress: suite.programAddr.String(),
				Fingerprint:     findingFingerPrintHash,
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.ConfirmFinding(ctx, testCase.req)

			if testCase.expPass {
				suite.Require().NoError(err)
				finding, err := suite.keeper.Findings.Get(suite.ctx, fid)
				suite.Require().NoError(err)
				suite.Require().Equal(finding.Status, types.FindingStatusConfirmed)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestConfirmFindingPaid() {
	pid, fid := uuid.NewString(), uuid.NewString()
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)
	suite.InitSubmitFinding(pid, fid)
	suite.InitActivateFinding(fid)

	finding, err := suite.keeper.Findings.Get(suite.ctx, fid)
	suite.Require().NoError(err)
	// fingerprint
	findingFingerPrintHash := suite.app.BountyKeeper.GetFindingFingerprintHash(&finding)
	suite.InitConfirmFinding(fid, findingFingerPrintHash)

	testCases := []struct {
		name    string
		req     *types.MsgConfirmFindingPaid
		expPass bool
	}{
		{
			"empty request",
			&types.MsgConfirmFindingPaid{},
			false,
		},
		{
			"valid request",
			&types.MsgConfirmFindingPaid{
				FindingId:       fid,
				OperatorAddress: suite.whiteHatAddr.String(),
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.ConfirmFindingPaid(ctx, testCase.req)

			if testCase.expPass {
				suite.Require().NoError(err)
				finding, err := suite.keeper.Findings.Get(suite.ctx, fid)
				suite.Require().NoError(err)
				suite.Require().Equal(finding.Status, types.FindingStatusPaid)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestCloseFinding() {
	pid, fid := uuid.NewString(), uuid.NewString()
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)
	suite.InitSubmitFinding(pid, fid)

	testCases := []struct {
		name    string
		req     *types.MsgCloseFinding
		expPass bool
	}{
		{
			"empty request",
			&types.MsgCloseFinding{},
			false,
		},
		{
			"valid request",
			&types.MsgCloseFinding{
				FindingId:       fid,
				OperatorAddress: suite.programAddr.String(),
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.CloseFinding(ctx, testCase.req)

			if testCase.expPass {
				suite.Require().NoError(err)
				finding, err := suite.keeper.Findings.Get(suite.ctx, fid)
				suite.Require().NoError(err)
				suite.Require().Equal(finding.Status, types.FindingStatusClosed)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestPublishConfirmFinding() {
	pid, fid := uuid.NewString(), uuid.NewString()
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)
	suite.InitSubmitFinding(pid, fid)
	suite.InitActivateFinding(fid)

	finding, err := suite.keeper.Findings.Get(suite.ctx, fid)
	suite.Require().NoError(err)
	findingFingerPrintHash := suite.app.BountyKeeper.GetFindingFingerprintHash(&finding)
	suite.InitConfirmFinding(fid, findingFingerPrintHash)

	suite.InitConfirmFindingPaid(fid)

	testCases := []struct {
		name    string
		req     *types.MsgPublishFinding
		expPass bool
	}{
		{
			"empty request",
			&types.MsgPublishFinding{},
			false,
		},
		{
			"valid request",
			&types.MsgPublishFinding{
				FindingId:       fid,
				OperatorAddress: suite.programAddr.String(),
				Description:     "desc",
				ProofOfConcept:  "poc",
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.PublishFinding(ctx, testCase.req)

			if testCase.expPass {
				suite.Require().NoError(err)
				finding, err := suite.keeper.Findings.Get(suite.ctx, fid)
				suite.Require().NoError(err)
				suite.Require().Equal(finding.Description, "desc")
				suite.Require().Equal(finding.ProofOfConcept, "poc")
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) InitCreateProgram(pid string) {
	msgCreateProgram := &types.MsgCreateProgram{
		ProgramId:       pid,
		Name:            "name",
		Detail:          "detail",
		OperatorAddress: suite.programAddr.String(),
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.CreateProgram(ctx, msgCreateProgram)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) InitActivateProgram(pid string) {
	msgCreateProgram := &types.MsgActivateProgram{
		ProgramId:       pid,
		OperatorAddress: suite.bountyAdminAddr.String(),
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.ActivateProgram(ctx, msgCreateProgram)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) InitSubmitFinding(pid, fid string) string {
	desc, poc := "desc", "poc"
	hash := sha256.Sum256([]byte(desc + poc + suite.whiteHatAddr.String()))

	msgSubmitFinding := &types.MsgSubmitFinding{
		ProgramId:       pid,
		FindingId:       fid,
		FindingHash:     hex.EncodeToString(hash[:]),
		OperatorAddress: suite.whiteHatAddr.String(),
		SeverityLevel:   types.Critical,
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.SubmitFinding(ctx, msgSubmitFinding)
	suite.Require().NoError(err)

	return msgSubmitFinding.FindingId
}

func (suite *KeeperTestSuite) InitActivateFinding(fid string) string {
	msgActivateFinding := &types.MsgActivateFinding{
		FindingId:       fid,
		OperatorAddress: suite.bountyAdminAddr.String(),
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.ActivateFinding(ctx, msgActivateFinding)
	suite.Require().NoError(err)

	return msgActivateFinding.FindingId
}

func (suite *KeeperTestSuite) InitConfirmFinding(fid, fingerprint string) string {
	msgConfirmFinding := &types.MsgConfirmFinding{
		FindingId:       fid,
		OperatorAddress: suite.programAddr.String(),
		Fingerprint:     fingerprint,
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.ConfirmFinding(ctx, msgConfirmFinding)
	suite.Require().NoError(err)

	return msgConfirmFinding.FindingId
}

func (suite *KeeperTestSuite) InitConfirmFindingPaid(fid string) string {
	msgConfirmFindingPaid := &types.MsgConfirmFindingPaid{
		FindingId:       fid,
		OperatorAddress: suite.whiteHatAddr.String(),
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.ConfirmFindingPaid(ctx, msgConfirmFindingPaid)
	suite.Require().NoError(err)

	return msgConfirmFindingPaid.FindingId
}

// TestGetProgramFindings tests the GetProgramFindings function
func (suite *KeeperTestSuite) TestGetProgramFindings() {
	// Create a program
	programID := uuid.NewString()
	suite.InitCreateProgram(programID)
	suite.InitActivateProgram(programID)

	// Create several findings for this program
	findingIDs := []string{
		uuid.NewString(),
		uuid.NewString(),
		uuid.NewString(),
	}

	// Submit each finding
	for _, findingID := range findingIDs {
		suite.InitSubmitFinding(programID, findingID)
	}

	// Query the findings using GetProgramFindings
	findings, err := suite.keeper.GetProgramFindings(suite.ctx, programID)
	suite.Require().NoError(err)

	// Verify all findings are returned
	suite.Require().Equal(len(findingIDs), len(findings), "Number of findings doesn't match")

	// Check that all expected finding IDs are in the result
	for _, expectedFindingID := range findingIDs {
		found := false
		for _, actualFindingID := range findings {
			if expectedFindingID == actualFindingID {
				found = true
				break
			}
		}
		suite.Require().True(found, "Finding %s was not returned by GetProgramFindings", expectedFindingID)
	}

	// Test with a non-existent program ID
	emptyFindings, err := suite.keeper.GetProgramFindings(suite.ctx, "non-existent-program")
	suite.Require().NoError(err)
	suite.Require().Empty(emptyFindings, "Should return empty slice for non-existent program")
}

// TestCreateTheorem tests the CreateTheorem message handler
func (suite *KeeperTestSuite) TestCreateTheorem() {
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	suite.Require().NoError(err)

	testCases := []struct {
		name    string
		req     *types.MsgCreateTheorem
		expPass bool
	}{
		{
			"empty request",
			&types.MsgCreateTheorem{},
			false,
		},
		{
			"invalid address",
			&types.MsgCreateTheorem{
				Title:       "Test Theorem",
				Description: "A test theorem description",
				Proposer:    "invalid-address",
			},
			false,
		},
		{
			"insufficient initial grant",
			&types.MsgCreateTheorem{
				Title:        "Test Theorem",
				Description:  "A test theorem description",
				Code:         "function example() { return true; }",
				InitialGrant: sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(10))), // Too small amount
				Proposer:     suite.programAddr.String(),
			},
			false,
		},
		{
			"missing title",
			&types.MsgCreateTheorem{
				Description:  "A test theorem description",
				Code:         "function example() { return true; }",
				InitialGrant: sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1e6))),
				Proposer:     suite.programAddr.String(),
			},
			false,
		},
		{
			"missing description",
			&types.MsgCreateTheorem{
				Title:        "Test Theorem",
				Code:         "function example() { return true; }",
				InitialGrant: sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1e6))),
				Proposer:     suite.programAddr.String(),
			},
			false,
		},
		{
			"valid request",
			&types.MsgCreateTheorem{
				Title:        "Test Theorem",
				Description:  "A test theorem description",
				Code:         "function example() { return true; }",
				InitialGrant: sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1e6))),
				Proposer:     suite.programAddr.String(),
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			if !testCase.expPass {
				_, err = suite.msgServer.CreateTheorem(ctx, testCase.req)
				suite.Require().Error(err)
				return
			}

			// Get initial balance before creating the theorem
			proposerAddr, err := sdk.AccAddressFromBech32(testCase.req.Proposer)
			suite.Require().NoError(err)
			initialBalance := suite.app.BankKeeper.GetBalance(suite.ctx, proposerAddr, bondDenom)

			resp, err := suite.msgServer.CreateTheorem(ctx, testCase.req)
			suite.Require().NoError(err)
			suite.Require().NotEmpty(resp.TheoremId)

			// Verify the theorem was created
			theorem, err := suite.keeper.Theorems.Get(suite.ctx, resp.TheoremId)
			suite.Require().NoError(err)
			suite.Require().Equal(testCase.req.Title, theorem.Title)
			suite.Require().Equal(testCase.req.Description, theorem.Description)
			suite.Require().Equal(testCase.req.Code, theorem.Code)
			suite.Require().Equal(types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD, theorem.Status)
			suite.Require().Equal(testCase.req.Proposer, theorem.Proposer)

			// Verify grant was created
			grant, err := suite.keeper.Grants.Get(suite.ctx, collections.Join(resp.TheoremId, sdk.MustAccAddressFromBech32(testCase.req.Proposer)))
			suite.Require().NoError(err)
			suite.Require().Equal(resp.TheoremId, grant.TheoremId)
			suite.Require().Equal(testCase.req.Proposer, grant.Grantor)
			suite.Require().True(testCase.req.InitialGrant[0].Equal(grant.Amount[0]))

			// Verify funds were transferred
			finalBalance := suite.app.BankKeeper.GetBalance(suite.ctx, proposerAddr, bondDenom)
			expectedBalance := initialBalance.Sub(testCase.req.InitialGrant[0])
			suite.Require().Equal(expectedBalance.Amount.String(), finalBalance.Amount.String())

			// Verify module account received the funds
			moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
			moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
			suite.Require().True(testCase.req.InitialGrant[0].Equal(moduleBalance))
		})
	}
}

// TestGrant tests the Grant message handler
func (suite *KeeperTestSuite) TestGrant() {
	// Create a theorem first
	theoremID := suite.InitCreateTheorem()
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	suite.Require().NoError(err)

	// Get initial module balance once for all test cases
	moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	initialModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)

	testCases := []struct {
		name    string
		req     *types.MsgGrant
		expPass bool
	}{
		{
			name:    "empty request",
			req:     &types.MsgGrant{},
			expPass: false,
		},
		{
			name: "invalid theorem ID",
			req: &types.MsgGrant{
				TheoremId: 9999, // Non-existent theorem ID
				Grantor:   suite.normalAddr.String(),
				Amount:    sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(500000))),
			},
			expPass: false,
		},
		{
			name: "invalid address",
			req: &types.MsgGrant{
				TheoremId: theoremID,
				Grantor:   "invalid-address",
				Amount:    sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(500000))),
			},
			expPass: false,
		},
		{
			name: "insufficient grant amount",
			req: &types.MsgGrant{
				TheoremId: theoremID,
				Grantor:   suite.normalAddr.String(),
				Amount:    sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(10))), // Too small
			},
			expPass: false,
		},
		{
			name: "grant with multiple coins",
			req: &types.MsgGrant{
				TheoremId: theoremID,
				Grantor:   suite.bountyAdminAddr.String(),
				Amount:    sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1e6)), sdk.NewCoin("othertoken", math.NewInt(5e5))),
			},
			expPass: false,
		},
		{
			name: "empty amount",
			req: &types.MsgGrant{
				TheoremId: theoremID,
				Grantor:   suite.normalAddr.String(),
				Amount:    sdk.NewCoins(),
			},
			expPass: false,
		},
		{
			name: "valid grant request",
			req: &types.MsgGrant{
				TheoremId: theoremID,
				Grantor:   suite.normalAddr.String(),
				Amount:    sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(5e6))),
			},
			expPass: true,
		},
		{
			name: "second valid grant from same address",
			req: &types.MsgGrant{
				TheoremId: theoremID,
				Grantor:   suite.normalAddr.String(),
				Amount:    sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(2e6))),
			},
			expPass: true,
		},
		{
			name: "grant from different address",
			req: &types.MsgGrant{
				TheoremId: theoremID,
				Grantor:   suite.whiteHatAddr.String(),
				Amount:    sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(3e6))),
			},
			expPass: true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)

			// Skip initial checks for non-passing tests
			if !testCase.expPass {
				_, err = suite.msgServer.Grant(ctx, testCase.req)
				suite.Require().Error(err)
				return
			}

			// Get initial balances and state for passing tests
			grantorAddr, err := sdk.AccAddressFromBech32(testCase.req.Grantor)
			suite.Require().NoError(err)

			initialGrantorBalance := suite.app.BankKeeper.GetBalance(suite.ctx, grantorAddr, bondDenom)

			// Get initial theorem state
			theorem, err := suite.keeper.Theorems.Get(suite.ctx, testCase.req.TheoremId)
			suite.Require().NoError(err)
			initialTheoremTotalGrant := theorem.TotalGrant

			// Check for existing grant
			var existingGrant types.Grant
			var hasExistingGrant bool
			grant, err := suite.keeper.Grants.Get(suite.ctx, collections.Join(testCase.req.TheoremId, grantorAddr))
			if err == nil {
				existingGrant = grant
				hasExistingGrant = true
			}

			// Execute the grant
			resp, err := suite.msgServer.Grant(ctx, testCase.req)
			suite.Require().NoError(err)
			suite.Require().NotNil(resp)

			// Verify grant was updated correctly
			updatedGrant, err := suite.keeper.Grants.Get(suite.ctx, collections.Join(testCase.req.TheoremId, grantorAddr))
			suite.Require().NoError(err)
			suite.Require().Equal(testCase.req.TheoremId, updatedGrant.TheoremId)
			suite.Require().Equal(testCase.req.Grantor, updatedGrant.Grantor)

			// Calculate expected grant amount based on whether there was an existing grant
			var expectedGrantCoins sdk.Coins
			if hasExistingGrant {
				// Convert existing grant amount to sdk.Coins and add new amount
				existingCoins := sdk.NewCoins(existingGrant.Amount...)
				expectedGrantCoins = existingCoins.Add(testCase.req.Amount...)
			} else {
				expectedGrantCoins = testCase.req.Amount
			}

			// Convert updatedGrant.Amount to sdk.Coins for comparison
			updatedGrantCoins := sdk.NewCoins(updatedGrant.Amount...)
			suite.Require().True(expectedGrantCoins.Equal(updatedGrantCoins))

			// Verify theorem total grant was updated
			updatedTheorem, err := suite.keeper.Theorems.Get(suite.ctx, testCase.req.TheoremId)
			suite.Require().NoError(err)

			// Convert initial theorem total grant to sdk.Coins and add new amount
			initialTheoremCoins := sdk.NewCoins(initialTheoremTotalGrant...)
			expectedTotalGrant := initialTheoremCoins.Add(testCase.req.Amount...)

			// Convert updatedTheorem.TotalGrant to sdk.Coins for comparison
			updatedTheoremCoins := sdk.NewCoins(updatedTheorem.TotalGrant...)
			suite.Require().True(expectedTotalGrant.Equal(updatedTheoremCoins))

			// Verify funds were transferred
			// Check grantor balance decreased
			finalGrantorBalance := suite.app.BankKeeper.GetBalance(suite.ctx, grantorAddr, bondDenom)
			expectedGrantorBalance := initialGrantorBalance.Sub(testCase.req.Amount[0])
			suite.Require().True(expectedGrantorBalance.Equal(finalGrantorBalance))

			// Check module balance increased
			finalModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
			expectedModuleBalance := initialModuleBalance.Add(testCase.req.Amount[0])
			suite.Require().True(expectedModuleBalance.Equal(finalModuleBalance))

			// Update the initial module balance for next test case
			initialModuleBalance = finalModuleBalance
		})
	}
}

// Helper function to create a theorem for testing
func (suite *KeeperTestSuite) InitCreateTheorem() uint64 {
	ctx := sdk.WrapSDKContext(suite.ctx)
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	suite.Require().NoError(err)

	createReq := &types.MsgCreateTheorem{
		Title:        "Test Theorem",
		Description:  "A test theorem description",
		Code:         "function test() { return true; }",
		InitialGrant: sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(1e6))),
		Proposer:     suite.programAddr.String(),
	}

	resp, err := suite.msgServer.CreateTheorem(ctx, createReq)
	suite.Require().NoError(err)
	return resp.TheoremId
}

// TestSubmitProofHash tests the SubmitProofHash message handler
func (suite *KeeperTestSuite) TestSubmitProofHash() {
	// Create a theorem first
	theoremID := suite.InitCreateTheorem()
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	suite.Require().NoError(err)

	testCases := []struct {
		name    string
		req     *types.MsgSubmitProofHash
		expPass bool
	}{
		{
			"empty request",
			&types.MsgSubmitProofHash{},
			false,
		},
		{
			"invalid theorem ID",
			&types.MsgSubmitProofHash{
				TheoremId: 9999, // Non-existent theorem ID
				Prover:    suite.whiteHatAddr.String(),
				ProofHash: "7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069",
				Deposit:   sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(500000))),
			},
			false,
		},
		{
			"invalid address",
			&types.MsgSubmitProofHash{
				TheoremId: theoremID,
				Prover:    "invalid-address",
				ProofHash: "7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069",
				Deposit:   sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(500000))),
			},
			false,
		},
		{
			"invalid hash - too short",
			&types.MsgSubmitProofHash{
				TheoremId: theoremID,
				Prover:    suite.whiteHatAddr.String(),
				ProofHash: "short",
				Deposit:   sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(500000))),
			},
			false,
		},
		{
			"invalid hash - no hex",
			&types.MsgSubmitProofHash{
				TheoremId: theoremID,
				Prover:    suite.whiteHatAddr.String(),
				ProofHash: "thisisnotahexstringthisisnotahexstringthisisnotahexstring",
				Deposit:   sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(500000))),
			},
			false,
		},
		{
			"insufficient deposit amount",
			&types.MsgSubmitProofHash{
				TheoremId: theoremID,
				Prover:    suite.whiteHatAddr.String(),
				ProofHash: "7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069",
				Deposit:   sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(10))), // Too small
			},
			false,
		},
		{
			"valid proof hash submission",
			&types.MsgSubmitProofHash{
				TheoremId: theoremID,
				Prover:    suite.whiteHatAddr.String(),
				ProofHash: "7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069",
				Deposit:   sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(500000))),
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			if !testCase.expPass {
				_, err = suite.msgServer.SubmitProofHash(ctx, testCase.req)
				suite.Require().Error(err)
				return
			}

			// Check initial balances if we expect success
			proverAddr, _ := sdk.AccAddressFromBech32(testCase.req.Prover)
			initialProverBalance := suite.app.BankKeeper.GetBalance(suite.ctx, proverAddr, bondDenom)

			moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
			initialModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)

			_, err = suite.msgServer.SubmitProofHash(ctx, testCase.req)
			suite.Require().NoError(err)
			// Verify proof was created
			proof, err := suite.keeper.Proofs.Get(suite.ctx, testCase.req.ProofHash)
			suite.Require().NoError(err)
			suite.Require().Equal(testCase.req.TheoremId, proof.TheoremId)
			suite.Require().Equal(testCase.req.ProofHash, proof.Id)
			suite.Require().Equal(testCase.req.Prover, proof.Prover)
			suite.Require().Equal(types.ProofStatus_PROOF_STATUS_HASH_LOCK_PERIOD, proof.Status)
			suite.Require().Equal(testCase.req.Deposit, proof.Deposit)

			// Verify proof-theorem relationship was established
			hasRelationship, err := suite.keeper.ProofsByTheorem.Has(suite.ctx,
				collections.Join(testCase.req.TheoremId, testCase.req.ProofHash))
			suite.Require().NoError(err)
			suite.Require().True(hasRelationship)

			// Verify funds were transferred
			finalProverBalance := suite.app.BankKeeper.GetBalance(suite.ctx, proverAddr, bondDenom)
			expectedProverBalance := initialProverBalance.SubAmount(sdk.NewCoins(testCase.req.Deposit...).AmountOf(bondDenom))
			suite.Require().Equal(expectedProverBalance.Amount.String(), finalProverBalance.Amount.String())

			// Verify module account received the funds
			finalModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
			expectedModuleBalance := initialModuleBalance.AddAmount(sdk.NewCoins(testCase.req.Deposit...).AmountOf(bondDenom))
			suite.Require().Equal(expectedModuleBalance.Amount.String(), finalModuleBalance.Amount.String())
		})
	}
}

// TestSubmitProofDetail tests the SubmitProofDetail message handler
func (suite *KeeperTestSuite) TestSubmitProofDetail() {
	// Create a theorem and submit a proof hash
	theoremID := suite.InitCreateTheorem()
	proofHash := suite.InitSubmitProofHash(theoremID)

	testCases := []struct {
		name    string
		req     *types.MsgSubmitProofDetail
		expPass bool
	}{
		{
			"empty request",
			&types.MsgSubmitProofDetail{},
			false,
		},
		{
			"invalid proof ID (non-existent)",
			&types.MsgSubmitProofDetail{
				ProofId: "invalid-hash",
				Prover:  suite.whiteHatAddr.String(),
				Detail:  "This is a valid proof detail",
			},
			false,
		},
		{
			"invalid address",
			&types.MsgSubmitProofDetail{
				ProofId: proofHash,
				Prover:  "invalid-address",
				Detail:  "This is a valid proof detail",
			},
			false,
		},
		{
			"invalid proof detail submission",
			&types.MsgSubmitProofDetail{
				ProofId: proofHash,
				Prover:  suite.whiteHatAddr.String(),
				Detail:  "This is a invalid proof detail",
			},
			false,
		},
		{
			"valid proof detail submission",
			&types.MsgSubmitProofDetail{
				ProofId: proofHash,
				Prover:  suite.whiteHatAddr.String(),
				Detail:  "This is a valid proof detail",
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			// Get the proof state before if it exists
			var initialProofStatus types.ProofStatus
			if testCase.req.ProofId != "" {
				proof, err := suite.keeper.Proofs.Get(suite.ctx, testCase.req.ProofId)
				if err == nil {
					initialProofStatus = proof.Status
				}
			}

			wrappedCtx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.SubmitProofDetail(wrappedCtx, testCase.req)

			if testCase.expPass {
				suite.Require().NoError(err)

				// Verify proof detail was updated
				proof, err := suite.keeper.Proofs.Get(suite.ctx, testCase.req.ProofId)
				suite.Require().NoError(err)

				suite.Require().Equal(testCase.req.Detail, proof.Detail)
				suite.Require().Equal(types.ProofStatus_PROOF_STATUS_HASH_DETAIL_PERIOD, proof.Status)

				// Verify proof is no longer in active proofs queue
				hasInActiveQueue, err := suite.keeper.ActiveProofsQueue.Has(suite.ctx,
					collections.Join(*proof.EndTime, proof.Id))
				suite.Require().NoError(err)
				suite.Require().False(hasInActiveQueue)
			} else {
				suite.Require().Error(err)

				// If the proof exists, verify its status wasn't changed
				if testCase.req.ProofId != "" {
					proof, err := suite.keeper.Proofs.Get(suite.ctx, testCase.req.ProofId)
					if err == nil && initialProofStatus != types.ProofStatus_PROOF_STATUS_UNSPECIFIED {
						suite.Require().Equal(initialProofStatus, proof.Status)
					}
				}
			}
		})
	}
}

// Helper function to initialize proof hash submission
func (suite *KeeperTestSuite) InitSubmitProofHash(theoremID uint64) string {
	ctx := sdk.WrapSDKContext(suite.ctx)
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	suite.Require().NoError(err)

	validHash := suite.app.BountyKeeper.GetProofHash(theoremID, suite.whiteHatAddr.String(), "This is a valid proof detail")
	submitHashReq := &types.MsgSubmitProofHash{
		TheoremId: theoremID,
		Prover:    suite.whiteHatAddr.String(),
		ProofHash: validHash,
		Deposit:   sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(500000))),
	}

	_, err = suite.msgServer.SubmitProofHash(ctx, submitHashReq)
	suite.Require().NoError(err)
	return validHash
}

// Helper function to initialize proof detail submission
func (suite *KeeperTestSuite) InitSubmitProofDetail(proofID string) {
	// Submit proof detail
	submitDetailReq := &types.MsgSubmitProofDetail{
		ProofId: proofID,
		Prover:  suite.whiteHatAddr.String(),
		Detail:  "This is a valid proof detail",
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.SubmitProofDetail(ctx, submitDetailReq)
	suite.Require().NoError(err)
}

// Helper function to initialize proof verification
func (suite *KeeperTestSuite) InitVerifyProof(proofID string, status types.ProofStatus) {
	// Verify the proof
	verifyReq := &types.MsgSubmitProofVerification{
		ProofId: proofID,
		Status:  status,
		Checker: suite.bountyAdminAddr.String(),
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.SubmitProofVerification(ctx, verifyReq)
	suite.Require().NoError(err)
}

// TestSubmitProofVerification tests the SubmitProofVerification message handler for passed proofs
func (suite *KeeperTestSuite) TestSubmitProofVerification() {
	// Create a theorem, submit a proof hash and detail
	theoremID := suite.InitCreateTheorem()
	validHash := suite.InitSubmitProofHash(theoremID)
	suite.InitSubmitProofDetail(validHash)

	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	suite.Require().NoError(err)

	testCases := []struct {
		name    string
		req     *types.MsgSubmitProofVerification
		expPass bool
	}{
		{
			"empty request",
			&types.MsgSubmitProofVerification{},
			false,
		},
		{
			"invalid proof ID",
			&types.MsgSubmitProofVerification{
				ProofId: "invalid-hash",
				Status:  types.ProofStatus_PROOF_STATUS_PASSED,
				Checker: suite.bountyAdminAddr.String(),
			},
			false,
		},
		{
			"non-existent proof ID",
			&types.MsgSubmitProofVerification{
				ProofId: "0000000000000000000000000000000000000000000000000000000000000000",
				Status:  types.ProofStatus_PROOF_STATUS_PASSED,
				Checker: suite.bountyAdminAddr.String(),
			},
			false,
		},
		{
			"unauthorized checker",
			&types.MsgSubmitProofVerification{
				ProofId: validHash,
				Status:  types.ProofStatus_PROOF_STATUS_PASSED,
				Checker: suite.normalAddr.String(), // Not a bounty admin
			},
			false,
		},
		{
			"invalid status",
			&types.MsgSubmitProofVerification{
				ProofId: validHash,
				Status:  types.ProofStatus_PROOF_STATUS_UNSPECIFIED, // Invalid status
				Checker: suite.bountyAdminAddr.String(),
			},
			false,
		},
		{
			"valid verification - passed",
			&types.MsgSubmitProofVerification{
				ProofId: validHash,
				Status:  types.ProofStatus_PROOF_STATUS_PASSED,
				Checker: suite.bountyAdminAddr.String(),
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			wrappedCtx := sdk.WrapSDKContext(suite.ctx)
			if !testCase.expPass {
				_, err := suite.msgServer.SubmitProofVerification(wrappedCtx, testCase.req)
				suite.Require().Error(err)
				return
			}

			// Get the initial balances if we expect the test to pass
			proof, err := suite.keeper.Proofs.Get(suite.ctx, testCase.req.ProofId)
			suite.Require().NoError(err)

			proverAddr, _ := sdk.AccAddressFromBech32(proof.Prover)
			initialProverBalance := suite.app.BankKeeper.GetBalance(suite.ctx, proverAddr, bondDenom)

			// Store the deposit amount for later verification
			depositAmount := sdk.NewCoins(proof.Deposit...)

			moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
			initialModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)

			_, err = suite.msgServer.SubmitProofVerification(wrappedCtx, testCase.req)
			suite.Require().NoError(err)

			// Verify proof status was updated
			proof, err = suite.keeper.Proofs.Get(suite.ctx, testCase.req.ProofId)
			suite.Require().NoError(err)
			suite.Require().Equal(testCase.req.Status, proof.Status)

			// Verify reward records were created but balances remain unchanged
			hasProverReward, err := suite.keeper.Rewards.Has(suite.ctx, proverAddr)
			suite.Require().NoError(err)
			suite.Require().True(hasProverReward, "Prover should have received a reward record")

			// Verify checker reward record was created
			checkerAddr, _ := sdk.AccAddressFromBech32(testCase.req.Checker)
			hasCheckerReward, err := suite.keeper.Rewards.Has(suite.ctx, checkerAddr)
			suite.Require().NoError(err)
			suite.Require().True(hasCheckerReward, "Checker should have received a reward record")

			// Verify deposit was returned to prover
			finalProverBalance := suite.app.BankKeeper.GetBalance(suite.ctx, proverAddr, bondDenom)
			expectedProverBalance := initialProverBalance.AddAmount(depositAmount.AmountOf(bondDenom))
			suite.Require().Equal(expectedProverBalance.Amount.String(), finalProverBalance.Amount.String(),
				"Prover should have received back the deposit when proof is passed")

			// Verify module account decreased by the deposit amount
			finalModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
			expectedModuleBalance := initialModuleBalance.SubAmount(depositAmount.AmountOf(bondDenom))
			suite.Require().Equal(expectedModuleBalance.Amount.String(), finalModuleBalance.Amount.String(),
				"Module account should have returned the deposit when proof is passed")

			// Verify reward records contain correct amounts
			proverReward, err := suite.keeper.Rewards.Get(suite.ctx, proverAddr)
			suite.Require().NoError(err)
			suite.Require().True(proverReward.Reward.IsValid(), "Prover reward should be valid")
			suite.Require().False(proverReward.Reward.IsZero(), "Prover reward should not be zero")

			checkerReward, err := suite.keeper.Rewards.Get(suite.ctx, checkerAddr)
			suite.Require().NoError(err)
			suite.Require().True(checkerReward.Reward.IsValid(), "Checker reward should be valid")
			suite.Require().False(checkerReward.Reward.IsZero(), "Checker reward should not be zero")
		})
	}
}

// TestSubmitProofVerificationFailed tests the SubmitProofVerification message handler for failed proofs
func (suite *KeeperTestSuite) TestSubmitProofVerificationFailed() {
	// Create a theorem, submit a proof hash and detail
	theoremID := suite.InitCreateTheorem()
	validHash := suite.InitSubmitProofHash(theoremID)
	suite.InitSubmitProofDetail(validHash)

	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	suite.Require().NoError(err)

	wrappedCtx := sdk.WrapSDKContext(suite.ctx)

	// Reset proof status to HASH_DETAIL_PERIOD
	proof, err := suite.keeper.Proofs.Get(suite.ctx, validHash)
	suite.Require().NoError(err)

	// Get initial balances for verification
	proverAddr, _ := sdk.AccAddressFromBech32(proof.Prover)
	initialProverBalance := suite.app.BankKeeper.GetBalance(suite.ctx, proverAddr, bondDenom)

	moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	initialModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)

	// Verify that the deposit exists in the Deposits collection before verification
	hasDeposit, err := suite.keeper.Deposits.Has(suite.ctx, collections.Join(validHash, proverAddr))
	suite.Require().NoError(err)
	suite.Require().True(hasDeposit, "Deposit should exist in Deposits collection before verification")

	// Create verification request for failed proof
	verifyReq := &types.MsgSubmitProofVerification{
		ProofId: validHash,
		Status:  types.ProofStatus_PROOF_STATUS_FAILED,
		Checker: suite.bountyAdminAddr.String(),
	}

	// Submit verification
	_, err = suite.msgServer.SubmitProofVerification(wrappedCtx, verifyReq)
	suite.Require().NoError(err)

	// Verify proof was deleted
	_, err = suite.keeper.Proofs.Get(suite.ctx, validHash)
	suite.Require().Error(err, "Proof should be deleted after failed verification")
	suite.Require().True(errors.IsOf(err, collections.ErrNotFound), "Expected ErrNotFound for deleted proof")

	// Verify the proof-theorem relationship was removed
	hasRelationship, err := suite.keeper.ProofsByTheorem.Has(suite.ctx,
		collections.Join(theoremID, validHash))
	suite.Require().NoError(err)
	suite.Require().False(hasRelationship, "Proof-theorem relationship should be removed")

	// Verify that the deposit record in Deposits collection is deleted
	hasDeposit, err = suite.keeper.Deposits.Has(suite.ctx, collections.Join(validHash, proverAddr))
	suite.Require().NoError(err)
	suite.Require().False(hasDeposit, "Deposit should be deleted from Deposits collection for failed proof")

	// Verify deposit is deleted (not returned to prover)
	proverBalance := suite.app.BankKeeper.GetBalance(suite.ctx, proverAddr, bondDenom)
	suite.Require().Equal(initialProverBalance.Amount.String(), proverBalance.Amount.String(),
		"Prover balance should remain unchanged for failed proof")

	// Verify module account balance decreased by the deposit amount
	moduleFinalBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
	suite.Require().Equal(initialModuleBalance.Amount.String(), moduleFinalBalance.Amount.String(),
		"Module account balance should remain unchanged for failed proof")

	// Verify no reward records were created
	hasProverReward, err := suite.keeper.Rewards.Has(suite.ctx, proverAddr)
	suite.Require().NoError(err)
	suite.Require().False(hasProverReward, "Prover should not have received a reward record for failed proof")

	// Verify checker did not receive a reward
	checkerAddr, _ := sdk.AccAddressFromBech32(verifyReq.Checker)
	hasCheckerReward, err := suite.keeper.Rewards.Has(suite.ctx, checkerAddr)
	suite.Require().NoError(err)
	suite.Require().False(hasCheckerReward, "Checker should not have received a reward record for failed proof")

	// Verify theorem status was not changed
	theorem, err := suite.keeper.Theorems.Get(suite.ctx, theoremID)
	suite.Require().NoError(err)
	suite.Require().Equal(types.TheoremStatus_THEOREM_STATUS_PROOF_PERIOD, theorem.Status,
		"Theorem status should not change for failed proof")
}

// TestWithdrawReward tests the WithdrawReward message handler
func (suite *KeeperTestSuite) TestWithdrawReward() {
	// Create test rewards for different addresses
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	suite.Require().NoError(err)

	// Set up small reward for whiteHat
	smallReward := types.Reward{
		Address: suite.whiteHatAddr.String(),
		Reward:  sdk.NewDecCoins(sdk.NewDecCoin(bondDenom, math.NewInt(100000))),
	}
	err = suite.keeper.Rewards.Set(suite.ctx, suite.whiteHatAddr, smallReward)
	suite.Require().NoError(err)

	// Set up large reward for program
	largeReward := types.Reward{
		Address: suite.programAddr.String(),
		Reward:  sdk.NewDecCoins(sdk.NewDecCoin(bondDenom, math.NewInt(5000000))),
	}
	err = suite.keeper.Rewards.Set(suite.ctx, suite.programAddr, largeReward)
	suite.Require().NoError(err)

	// Calculate total rewards and add to module account
	totalRewardAmount := smallReward.Reward[0].Amount.TruncateInt().Add(largeReward.Reward[0].Amount.TruncateInt())
	moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)

	// Add tokens to module account to match rewards
	err = suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(bondDenom, totalRewardAmount)))
	suite.Require().NoError(err)

	// Verify module account has sufficient funds for rewards
	moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
	suite.Require().True(moduleBalance.Amount.GTE(totalRewardAmount),
		"Module account should have sufficient funds for all rewards")

	testCases := []struct {
		name    string
		req     *types.MsgWithdrawReward
		expPass bool
	}{
		{
			"empty request",
			&types.MsgWithdrawReward{},
			false,
		},
		{
			"invalid address",
			&types.MsgWithdrawReward{
				Address: "invalid-address",
			},
			false,
		},
		{
			"non-existent reward",
			&types.MsgWithdrawReward{
				Address: suite.normalAddr.String(), // No reward exists for this address
			},
			false,
		},
		{
			"withdraw small reward",
			&types.MsgWithdrawReward{
				Address: suite.whiteHatAddr.String(), // Has small reward
			},
			true,
		},
		{
			"withdraw large reward",
			&types.MsgWithdrawReward{
				Address: suite.programAddr.String(), // Has large reward
			},
			true,
		},
		{
			"already withdrawn reward",
			&types.MsgWithdrawReward{
				Address: suite.whiteHatAddr.String(), // Reward should be gone after previous test
			},
			false,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			if !testCase.expPass {
				_, err := suite.msgServer.WithdrawReward(ctx, testCase.req)
				suite.Require().Error(err)

				// If reward existed before, check it still exists
				addrObj, err := sdk.AccAddressFromBech32(testCase.req.Address)
				if err != nil && testCase.name == "already withdrawn reward" {
					hasReward, err := suite.keeper.Rewards.Has(suite.ctx, addrObj)
					suite.Require().NoError(err)
					suite.Require().True(hasReward, "Reward should still exist after failed withdrawal")
				}
				return
			}

			// Get the current reward if any
			addrObj, err := sdk.AccAddressFromBech32(testCase.req.Address)
			suite.Require().NoError(err)
			currentReward, err := suite.keeper.Rewards.Get(suite.ctx, addrObj)

			// Get initial balance
			initialBalance := suite.app.BankKeeper.GetBalance(suite.ctx, addrObj, bondDenom)
			moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
			moduleInitialBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)

			_, err = suite.msgServer.WithdrawReward(ctx, testCase.req)
			suite.Require().NoError(err)

			// Verify the reward was withdrawn (should be removed)
			hasReward, err := suite.keeper.Rewards.Has(suite.ctx, addrObj)
			suite.Require().NoError(err)
			suite.Require().False(hasReward, "Reward should be removed after withdrawal")

			// Check that the user's balance increased by the reward amount
			finalBalance := suite.app.BankKeeper.GetBalance(suite.ctx, addrObj, bondDenom)

			// Convert DecCoins to regular Coins for comparison
			rewardCoins := sdk.NewCoins()
			for _, decCoin := range currentReward.Reward {
				rewardCoins = rewardCoins.Add(sdk.NewCoin(decCoin.Denom, decCoin.Amount.TruncateInt()))
			}

			expectedBalance := initialBalance.AddAmount(rewardCoins.AmountOf(bondDenom))
			suite.Require().Equal(expectedBalance.Amount.String(), finalBalance.Amount.String())

			// Verify module account balance decreased
			moduleFinalBalance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, bondDenom)
			moduleExpectedBalance := moduleInitialBalance.SubAmount(rewardCoins.AmountOf(bondDenom))
			suite.Require().Equal(moduleExpectedBalance.Amount.String(), moduleFinalBalance.Amount.String())
		})
	}
}
