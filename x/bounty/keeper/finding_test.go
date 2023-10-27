package keeper_test

import (
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
	"time"
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
		ProgramId:    "1",
		Name:         "name",
		Description:  "desc",
		AdminAddress: suite.address[0].String(),
		Status:       types.ProgramStatusInactive,
	}
	suite.keeper.SetProgram(suite.ctx, program)
	storedProgram, isExist := suite.keeper.GetProgram(suite.ctx, program.ProgramId)
	suite.Require().Equal(true, isExist)
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
						SubmitterAddress: suite.address[0].String(),
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
				suite.keeper.SetFinding(suite.ctx, finding)
				storedFinding, isExist := suite.keeper.GetFinding(suite.ctx, finding.FindingId)
				suite.Require().Equal(true, isExist)

				storedFindings := suite.keeper.GetAllFindings(suite.ctx)
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
