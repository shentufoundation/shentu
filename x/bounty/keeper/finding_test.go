package keeper_test

import "github.com/shentufoundation/shentu/v2/x/bounty/types"

func (suite *KeeperTestSuite) TestFindingList_GetSet() {
	findIDs := []uint64{10, 20, 30, 40}
	var pid uint64
	pid = 2
	err := suite.keeper.SetPidFindingIDList(suite.ctx, pid, findIDs)
	suite.Require().NoError(err)

	findIDs2, err := suite.keeper.GetPidFindingIDList(suite.ctx, pid)
	suite.Require().NoError(err)
	suite.Require().Equal(findIDs, findIDs2)
}

func (suite *KeeperTestSuite) TestFinding_GetSet() {
	type args struct {
		finding []types.Finding
	}

	type errArgs struct {
		shouldPass bool
	}

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{"Finding(1)  -> Set: Simple",
			args{
				finding: []types.Finding{
					{
						FindingId:        1,
						Title:            "test finding",
						ProgramId:        1,
						SeverityLevel:    types.SeverityLevelCritical,
						SubmitterAddress: suite.address[0].String(),
					},
				},
			},
			errArgs{
				shouldPass: true,
			},
		},
		{"Finding(2)  -> Set: Simple",
			args{
				finding: []types.Finding{
					{
						FindingId:        3,
						Title:            "test findingv3",
						ProgramId:        3,
						SeverityLevel:    types.SeverityLevelCritical,
						SubmitterAddress: suite.address[0].String(),
					},
				},
			},
			errArgs{
				shouldPass: true,
			},
		},
		{"Finding(3)  -> get: Simple",
			args{
				finding: []types.Finding{
					{
						FindingId:        30,
						Title:            "",
						ProgramId:        3,
						SeverityLevel:    types.SeverityLevelCritical,
						SubmitterAddress: suite.address[0].String(),
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
			for _, finding := range tc.args.finding {
				if tc.errArgs.shouldPass {
					suite.keeper.SetFinding(suite.ctx, finding)
					findingResult, result := suite.keeper.GetFinding(suite.ctx, finding.FindingId)

					suite.Require().True(result)
					suite.Require().Equal(findingResult.FindingId, finding.FindingId)
				} else {
					_, result := suite.keeper.GetFinding(suite.ctx, finding.FindingId)
					suite.Require().False(result)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestNextFindingID_GetSet() {
	var findingID uint64 = 2
	suite.keeper.SetNextFindingID(suite.ctx, findingID)
	nextFindingID := suite.keeper.GetNextFindingID(suite.ctx)

	suite.Require().Equal(findingID, nextFindingID)
}
