package keeper_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
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
