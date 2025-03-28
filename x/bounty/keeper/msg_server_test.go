package keeper_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"

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

func (suite *KeeperTestSuite) TestSubmitFinding() {
	type args struct {
		msgSubmitFindings []types.MsgSubmitFinding
	}

	type errArgs struct {
		shouldPass bool
	}

	pid := uuid.NewString()
	fid := uuid.NewString()
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{"Submit finding(1)  -> submit: Simple",
			args{
				msgSubmitFindings: []types.MsgSubmitFinding{
					{
						ProgramId:       pid,
						FindingId:       fid,
						OperatorAddress: suite.whiteHatAddr.String(),
						SeverityLevel:   types.Critical,
					},
				},
			},
			errArgs{
				shouldPass: true,
			},
		},
		{"Submit finding(2)  -> fid repeat",
			args{
				msgSubmitFindings: []types.MsgSubmitFinding{
					{
						ProgramId:       pid,
						FindingId:       fid,
						OperatorAddress: suite.whiteHatAddr.String(),
						SeverityLevel:   types.Critical,
					},
				},
			},
			errArgs{
				shouldPass: false,
			},
		},
		{"Submit finding(3)  -> pid not exist",
			args{
				msgSubmitFindings: []types.MsgSubmitFinding{
					{
						ProgramId:       "not exist pid",
						FindingId:       "1",
						OperatorAddress: suite.whiteHatAddr.String(),
						SeverityLevel:   types.Critical,
					},
				},
			},
			errArgs{
				shouldPass: false,
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			for _, finding := range tc.args.msgSubmitFindings {
				ctx := sdk.WrapSDKContext(suite.ctx)

				_, err := suite.msgServer.SubmitFinding(ctx, &finding)

				if tc.errArgs.shouldPass {
					suite.Require().NoError(err)
					// Directly verify finding existence using collections
					_, err := suite.keeper.Findings.Get(suite.ctx, finding.FindingId)
					suite.Require().NoError(err)
				} else {
					suite.Require().Error(err)
				}
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

			finding, err := suite.keeper.Findings.Get(suite.ctx, fid)
			suite.Require().NoError(err)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(finding.Status, types.FindingStatusActive)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(finding.Status, types.FindingStatusSubmitted)
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

			finding, err := suite.keeper.Findings.Get(suite.ctx, fid)
			suite.Require().NoError(err)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(finding.Status, types.FindingStatusConfirmed)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(finding.Status, types.FindingStatusActive)
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
			finding, err := suite.keeper.Findings.Get(suite.ctx, fid)
			suite.Require().NoError(err)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(finding.Status, types.FindingStatusPaid)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(finding.Status, types.FindingStatusConfirmed)
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
			finding, err := suite.keeper.Findings.Get(suite.ctx, fid)
			suite.Require().NoError(err)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(finding.Status, types.FindingStatusClosed)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(finding.Status, types.FindingStatusSubmitted)
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

			finding, err := suite.keeper.Findings.Get(suite.ctx, fid)
			suite.Require().NoError(err)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(finding.Description, "desc")
				suite.Require().Equal(finding.ProofOfConcept, "poc")
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(finding.Status, types.FindingStatusPaid)
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

// TestErrorHandling tests various error scenarios in the message server
func (suite *KeeperTestSuite) TestErrorHandling() {
	// Test create program with invalid operator address
	suite.Run("Create Program with Invalid Address", func() {
		ctx := sdk.WrapSDKContext(suite.ctx)

		// Use an invalid address format
		invalidMsg := &types.MsgCreateProgram{
			ProgramId:       uuid.NewString(),
			Name:            "Invalid Address Program",
			Detail:          "This should fail due to invalid address",
			OperatorAddress: "invalid-address", // Invalid address format
		}

		// This should fail due to invalid address
		_, err := suite.msgServer.CreateProgram(ctx, invalidMsg)
		suite.Require().Error(err)
	})

	// Test activating a non-existent program
	suite.Run("Activate Non-existent Program", func() {
		ctx := sdk.WrapSDKContext(suite.ctx)

		nonExistentProgramID := uuid.NewString()

		// Try to activate a program that doesn't exist
		activateMsg := &types.MsgActivateProgram{
			ProgramId:       nonExistentProgramID,
			OperatorAddress: suite.bountyAdminAddr.String(),
		}

		// This should fail because the program doesn't exist
		_, err := suite.msgServer.ActivateProgram(ctx, activateMsg)
		suite.Require().Error(err)
	})	
	// Test submitting a finding for a non-existent program
	suite.Run("Submit Finding for Non-existent Program", func() {					
		ctx := sdk.WrapSDKContext(suite.ctx)

		nonExistentProgramID := uuid.NewString()

		// Try to submit a finding for a program that doesn't exist
		submitMsg := &types.MsgSubmitFinding{
			ProgramId:       nonExistentProgramID,
			FindingId:       uuid.NewString(),
			OperatorAddress: suite.whiteHatAddr.String(),
			SeverityLevel:   types.Critical,
		}

		// This should fail because the program doesn't exist
		_, err := suite.msgServer.SubmitFinding(ctx, submitMsg)
		suite.Require().Error(err)
	})

	// Test activating a non-existent finding
	suite.Run("Activate Non-existent Finding", func() {
		ctx := sdk.WrapSDKContext(suite.ctx)

		nonExistentFindingID := uuid.NewString()

		// Try to activate a finding that doesn't exist
		activateMsg := &types.MsgActivateFinding{
			FindingId:       nonExistentFindingID,
			OperatorAddress: suite.bountyAdminAddr.String(),
		}

		// This should fail because the finding doesn't exist
		_, err := suite.msgServer.ActivateFinding(ctx, activateMsg)
		suite.Require().Error(err)
	})
}

// TestConcurrentOperations tests concurrent operations on the keeper
func (suite *KeeperTestSuite) TestConcurrentOperations() {
	// Create a program to work with
	programID := uuid.NewString()
	program := types.Program{
		ProgramId:    programID,
		Name:         "Concurrent Test Program",
		Detail:       "Program for concurrent testing",
		AdminAddress: suite.programAddr.String(),
		Status:       types.ProgramStatusActive,
	}
	err := suite.keeper.Programs.Set(suite.ctx, program.ProgramId, program)
	suite.Require().NoError(err)

	// Number of concurrent operations
	concurrentCount := 10

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(concurrentCount)

	// Create a slice to hold all finding IDs
	findingIDs := make([]string, concurrentCount)

	// Use a mutex to protect access to shared data
	var mutex sync.Mutex

	// Launch concurrent submission of findings
	for i := 0; i < concurrentCount; i++ {
		go func(index int) {
			defer wg.Done()

			// Create a new finding ID
			findingID := uuid.NewString()

			// Store the finding ID in the slice (protected by mutex)
			mutex.Lock()
			findingIDs[index] = findingID
			mutex.Unlock()

			// Create a finding with this ID
			finding := types.Finding{
				ProgramId:        programID,
				FindingId:        findingID,
				Title:            fmt.Sprintf("Concurrent Finding %d", index),
				Description:      fmt.Sprintf("Description for concurrent finding %d", index),
				SubmitterAddress: suite.whiteHatAddr.String(),
				CreateTime:       suite.ctx.BlockTime(),
				Status:           types.FindingStatusSubmitted,
				SeverityLevel:    types.Low,
			}

			// Set the finding in the store
			err := suite.keeper.Findings.Set(suite.ctx, findingID, finding)

			// Set the program-finding relationship
			if err == nil {
				err = suite.keeper.ProgramFindings.Set(suite.ctx, collections.Join(programID, findingID))
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Verify all findings were stored
	for _, findingID := range findingIDs {
		_, err := suite.keeper.Findings.Get(suite.ctx, findingID)
		suite.Require().NoError(err, "Finding %s was not stored correctly", findingID)
	}

	// Verify program-finding relationships
	programFindings, err := suite.keeper.GetProgramFindings(suite.ctx, programID)
	suite.Require().NoError(err)

	// Check that all our findings are in the program findings
	foundCount := 0
	for _, pfID := range programFindings {
		for _, findingID := range findingIDs {
			if pfID == findingID {
				foundCount++
				break
			}
		}
	}

	// We should have found all our findings
	suite.Require().Equal(concurrentCount, foundCount, "Not all findings were properly associated with the program")
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

// TestCreateTheorem tests the CreateTheorem message handler
func (suite *KeeperTestSuite) TestCreateTheorem() {
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
			"valid request",
			&types.MsgCreateTheorem{
				Title:       "Test Theorem",
				Description: "A test theorem description",
				Proposer:    suite.programAddr.String(),
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			resp, err := suite.msgServer.CreateTheorem(ctx, testCase.req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().NotEmpty(resp.TheoremId)

				// Verify the theorem was created
				theorem, err := suite.keeper.Theorems.Get(suite.ctx, resp.TheoremId)
				suite.Require().NoError(err)
				suite.Require().Equal(testCase.req.Title, theorem.Title)
				suite.Require().Equal(testCase.req.Description, theorem.Description)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

// TestWithdrawReward tests the WithdrawReward message handler
func (suite *KeeperTestSuite) TestWithdrawReward() {
	// Create a test reward for the white hat address
	reward := types.Reward{
		Address: suite.whiteHatAddr.String(),
		Reward:  sdk.NewDecCoins(sdk.NewDecCoin("uctk", math.NewInt(1000000))),
	}
	err := suite.keeper.Rewards.Set(suite.ctx, suite.whiteHatAddr, reward)
	suite.Require().NoError(err)

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
			"valid request",
			&types.MsgWithdrawReward{
				Address: suite.whiteHatAddr.String(), // Has reward
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.WithdrawReward(ctx, testCase.req)

			if testCase.expPass {
				suite.Require().NoError(err)

				// Verify the reward was withdrawn (should be removed)
				hasReward, err := suite.keeper.Rewards.Has(suite.ctx, suite.whiteHatAddr)
				suite.Require().NoError(err)
				suite.Require().False(hasReward, "Reward should be removed after withdrawal")
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
