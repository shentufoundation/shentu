package keeper_test

import (
	"fmt"
	"strings"
	"time"

	"cosmossdk.io/collections"
	"github.com/google/uuid"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (suite *KeeperTestSuite) TestSetGetFinding() {
	type args struct {
		findings []types.Finding
	}
	type errArgs struct {
		shouldPass bool
		contains   string
	}
	// create program
	program := types.Program{
		ProgramId:    uuid.NewString(),
		Name:         "name",
		Detail:       "detail",
		AdminAddress: suite.programAddr.String(),
		Status:       types.ProgramStatusInactive,
	}
	err := suite.keeper.Programs.Set(suite.ctx, program.ProgramId, program)
	suite.Require().NoError(err)

	storedProgram, err := suite.keeper.Programs.Get(suite.ctx, program.ProgramId)
	suite.Require().NoError(err)
	suite.Require().Equal(program.ProgramId, storedProgram.ProgramId)

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{"Finding(1)  -> Set: Simple",
			args{
				findings: []types.Finding{
					{
						ProgramId:        "1",
						FindingId:        "1",
						Title:            "title",
						Description:      "desc",
						SubmitterAddress: suite.whiteHatAddr.String(),
						CreateTime:       time.Time{},
						Status:           types.FindingStatusSubmitted,
						FindingHash:      "",
						SeverityLevel:    types.Low,
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
			for _, finding := range tc.args.findings {
				// Set finding
				err := suite.keeper.Findings.Set(suite.ctx, finding.FindingId, finding)
				suite.Require().NoError(err)

				// Set in ProgramFindings index
				err = suite.keeper.ProgramFindings.Set(suite.ctx, collections.Join(finding.ProgramId, finding.FindingId))
				suite.Require().NoError(err)

				// Get finding
				storedFinding, err := suite.keeper.Findings.Get(suite.ctx, finding.FindingId)
				suite.Require().NoError(err)

				// Get all findings
				var storedFindings []types.Finding
				err = suite.keeper.Findings.Walk(suite.ctx, nil, func(_ string, f types.Finding) (bool, error) {
					storedFindings = append(storedFindings, f)
					return false, nil
				})
				suite.Require().NoError(err)
				suite.Require().Equal(1, len(storedFindings))

				if tc.errArgs.shouldPass {
					suite.Require().Equal(finding.FindingId, storedFinding.FindingId)
				} else {
					suite.Require().NotEqual(finding.FindingId, storedFinding.FindingId)
				}
			}
		})
	}
}

// TestProgramFindingsIndex validates the ProgramFindings index functionality
func (suite *KeeperTestSuite) TestProgramFindingsIndex() {
	// Create a program
	programID := uuid.NewString()
	program := types.Program{
		ProgramId:    programID,
		Name:         "Test Program",
		Detail:       "Program for testing",
		AdminAddress: suite.programAddr.String(),
		Status:       types.ProgramStatusActive,
	}
	err := suite.keeper.Programs.Set(suite.ctx, program.ProgramId, program)
	suite.Require().NoError(err)

	// Create multiple findings for the program
	findingIDs := []string{uuid.NewString(), uuid.NewString(), uuid.NewString()}
	for i, findingID := range findingIDs {
		finding := types.Finding{
			ProgramId:        programID,
			FindingId:        findingID,
			Title:            fmt.Sprintf("Title %d", i),
			Description:      fmt.Sprintf("Description %d", i),
			SubmitterAddress: suite.whiteHatAddr.String(),
			CreateTime:       time.Now(),
			Status:           types.FindingStatusSubmitted,
			SeverityLevel:    types.Low,
		}

		err := suite.keeper.Findings.Set(suite.ctx, finding.FindingId, finding)
		suite.Require().NoError(err)

		err = suite.keeper.ProgramFindings.Set(suite.ctx, collections.Join(finding.ProgramId, finding.FindingId))
		suite.Require().NoError(err)
	}

	// Get all findings for the program
	programFindings, err := suite.keeper.GetProgramFindings(suite.ctx, programID)
	suite.Require().NoError(err)
	suite.Require().Equal(len(findingIDs), len(programFindings))

	// Verify each finding ID is in the list
	for _, findingID := range findingIDs {
		found := false
		for _, pfID := range programFindings {
			if pfID == findingID {
				found = true
				break
			}
		}
		suite.Require().True(found, "Finding ID not found in program findings: %s", findingID)
	}

	// Test non-existent program
	nonExistentProgramID := uuid.NewString()
	emptyFindings, err := suite.keeper.GetProgramFindings(suite.ctx, nonExistentProgramID)
	suite.Require().NoError(err) // Should succeed but return empty list
	suite.Require().Empty(emptyFindings)
}

// TestEdgeCases tests various edge cases for Programs and Findings
func (suite *KeeperTestSuite) TestEdgeCases() {
	// Test empty program ID
	emptyProgram := types.Program{
		ProgramId:    "", // Empty ID
		Name:         "Empty ID Program",
		Detail:       "Program with empty ID",
		AdminAddress: suite.programAddr.String(),
		Status:       types.ProgramStatusInactive,
	}
	err := suite.keeper.Programs.Set(suite.ctx, emptyProgram.ProgramId, emptyProgram)
	suite.Require().NoError(err) // Should succeed but may be bad practice

	// Retrieve the empty ID program
	retrievedEmptyProgram, err := suite.keeper.Programs.Get(suite.ctx, "")
	suite.Require().NoError(err)
	suite.Require().Equal("Empty ID Program", retrievedEmptyProgram.Name)

	// Test very long program ID
	longID := strings.Repeat("a", 1000) // Very long ID
	longProgram := types.Program{
		ProgramId:    longID,
		Name:         "Long ID Program",
		Detail:       "Program with very long ID",
		AdminAddress: suite.programAddr.String(),
		Status:       types.ProgramStatusInactive,
	}
	err = suite.keeper.Programs.Set(suite.ctx, longProgram.ProgramId, longProgram)
	suite.Require().NoError(err)

	// Retrieve the long ID program
	retrievedLongProgram, err := suite.keeper.Programs.Get(suite.ctx, longID)
	suite.Require().NoError(err)
	suite.Require().Equal("Long ID Program", retrievedLongProgram.Name)

	// Test invalid status
	invalidStatusProgram := types.Program{
		ProgramId:    uuid.NewString(),
		Name:         "Invalid Status Program",
		Detail:       "Program with invalid status",
		AdminAddress: suite.programAddr.String(),
		Status:       999, // Invalid status
	}
	err = suite.keeper.Programs.Set(suite.ctx, invalidStatusProgram.ProgramId, invalidStatusProgram)
	suite.Require().NoError(err) // Should succeed at storage level, but application logic should validate

	// Retrieve the invalid status program
	retrievedInvalidProgram, err := suite.keeper.Programs.Get(suite.ctx, invalidStatusProgram.ProgramId)
	suite.Require().NoError(err)
	suite.Require().Equal(types.ProgramStatus(999), retrievedInvalidProgram.Status)

	// Test duplicate set (overwrite)
	originalProgram := types.Program{
		ProgramId:    "duplicate-id",
		Name:         "Original Program",
		Detail:       "Original Program Detail",
		AdminAddress: suite.programAddr.String(),
		Status:       types.ProgramStatusInactive,
	}
	err = suite.keeper.Programs.Set(suite.ctx, originalProgram.ProgramId, originalProgram)
	suite.Require().NoError(err)

	// Now override with a new program with the same ID
	overwriteProgram := types.Program{
		ProgramId:    "duplicate-id",
		Name:         "Overwrite Program",
		Detail:       "Overwrite Program Detail",
		AdminAddress: suite.programAddr.String(),
		Status:       types.ProgramStatusActive,
	}
	err = suite.keeper.Programs.Set(suite.ctx, overwriteProgram.ProgramId, overwriteProgram)
	suite.Require().NoError(err)

	// Verify the overwrite was successful
	retrievedProgram, err := suite.keeper.Programs.Get(suite.ctx, "duplicate-id")
	suite.Require().NoError(err)
	suite.Require().Equal("Overwrite Program", retrievedProgram.Name)
	suite.Require().Equal(types.ProgramStatusActive, retrievedProgram.Status)
}

// TestBulkOperations tests the handling of multiple programs and findings
func (suite *KeeperTestSuite) TestBulkOperations() {
	// Create multiple programs
	numberOfPrograms := 10
	programIDs := make([]string, numberOfPrograms)

	// Insert programs in bulk
	for i := 0; i < numberOfPrograms; i++ {
		programID := uuid.NewString()
		programIDs[i] = programID

		program := types.Program{
			ProgramId:    programID,
			Name:         fmt.Sprintf("Program %d", i),
			Detail:       fmt.Sprintf("Detail for program %d", i),
			AdminAddress: suite.programAddr.String(),
			Status:       types.ProgramStatusActive,
		}

		err := suite.keeper.Programs.Set(suite.ctx, programID, program)
		suite.Require().NoError(err)
	}

	// Verify all programs were stored
	var storedPrograms []types.Program
	err := suite.keeper.Programs.Walk(suite.ctx, nil, func(_ string, program types.Program) (bool, error) {
		// Check if this program is one of our test programs
		for _, pid := range programIDs {
			if program.ProgramId == pid {
				storedPrograms = append(storedPrograms, program)
				break
			}
		}
		return false, nil
	})
	suite.Require().NoError(err)
	suite.Require().GreaterOrEqual(len(storedPrograms), numberOfPrograms)

	// Create multiple findings for the first program
	findingsPerProgram := 5
	firstProgramID := programIDs[0]
	findingIDs := make([]string, findingsPerProgram)

	for i := 0; i < findingsPerProgram; i++ {
		findingID := uuid.NewString()
		findingIDs[i] = findingID

		finding := types.Finding{
			ProgramId:        firstProgramID,
			FindingId:        findingID,
			Title:            fmt.Sprintf("Finding %d", i),
			Description:      fmt.Sprintf("Description for finding %d", i),
			SubmitterAddress: suite.whiteHatAddr.String(),
			CreateTime:       time.Now(),
			Status:           types.FindingStatusSubmitted,
			SeverityLevel:    types.Low,
		}

		err := suite.keeper.Findings.Set(suite.ctx, findingID, finding)
		suite.Require().NoError(err)

		err = suite.keeper.ProgramFindings.Set(suite.ctx, collections.Join(firstProgramID, findingID))
		suite.Require().NoError(err)
	}

	// Verify all findings were stored
	var storedFindings []types.Finding
	err = suite.keeper.Findings.Walk(suite.ctx, nil, func(_ string, finding types.Finding) (bool, error) {
		// Check if this finding is for our first program
		if finding.ProgramId == firstProgramID {
			storedFindings = append(storedFindings, finding)
		}
		return false, nil
	})
	suite.Require().NoError(err)
	suite.Require().GreaterOrEqual(len(storedFindings), findingsPerProgram)

	// Test retrieving all findings for the first program
	programFindings, err := suite.keeper.GetProgramFindings(suite.ctx, firstProgramID)
	suite.Require().NoError(err)
	suite.Require().Equal(findingsPerProgram, len(programFindings))

	// Update all findings in the first program to a new status
	for _, findingID := range findingIDs {
		finding, err := suite.keeper.Findings.Get(suite.ctx, findingID)
		suite.Require().NoError(err)

		// Update status
		finding.Status = types.FindingStatusActive
		err = suite.keeper.Findings.Set(suite.ctx, findingID, finding)
		suite.Require().NoError(err)
	}

	// Verify all findings were updated
	updatedFindings := 0
	err = suite.keeper.Findings.Walk(suite.ctx, nil, func(_ string, finding types.Finding) (bool, error) {
		if finding.ProgramId == firstProgramID && finding.Status == types.FindingStatusActive {
			updatedFindings++
		}
		return false, nil
	})
	suite.Require().NoError(err)
	suite.Require().Equal(findingsPerProgram, updatedFindings)
}
