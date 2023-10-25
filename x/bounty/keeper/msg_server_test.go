package keeper_test

import (
	"fmt"

	types1 "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (suite *KeeperTestSuite) TestCreateProgram() {
	type args struct {
		msgCresatePrograms []types.MsgCreateProgram
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
				msgCresatePrograms: []types.MsgCreateProgram{
					{
						Name:            "Name",
						Description:     "Desc",
						OperatorAddress: suite.address[0].String(),
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
			for _, program := range tc.args.msgCresatePrograms {
				ctx := types1.WrapSDKContext(suite.ctx)

				_, err := suite.msgServer.CreateProgram(ctx, &program)
				suite.Require().NoError(err)
				storedProgram, result := suite.keeper.GetProgram(suite.ctx, program.ProgramId)
				if tc.errArgs.shouldPass {
					suite.Require().NoError(err)
					suite.Require().True(result)
					suite.Require().Equal(storedProgram.ProgramId, program.ProgramId)
				} else {
					suite.Require().Error(err)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSubmitFinding() {
	type args struct {
		msgSubmitFindings []types.MsgSubmitFinding
	}

	type errArgs struct {
		shouldPass bool
	}

	pid := suite.InitCreateProgram()

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{"Submit finding(1)  -> submit: Simple",
			args{
				msgSubmitFindings: []types.MsgSubmitFinding{
					{
						ProgramId:        pid,
						FindingId:        "1",
						Title:            "Test bug 1",
						Description:      "Desc",
						SubmitterAddress: suite.address[0].String(),
						SeverityLevel:    types.SeverityLevelCritical,
					},
				},
			},
			errArgs{
				shouldPass: true,
			},
		},
		{"Submit finding(2)  -> submit: Simple",
			args{
				msgSubmitFindings: []types.MsgSubmitFinding{
					{
						ProgramId:        "Not exist pid",
						FindingId:        "1",
						Title:            "Test bug 1",
						Description:      "Desc",
						SubmitterAddress: suite.address[0].String(),
						SeverityLevel:    types.SeverityLevelCritical,
					},
				},
			},
			errArgs{
				shouldPass: false,
			},
		},
		{"Submit finding(3)  -> submit: Simple",
			args{
				msgSubmitFindings: []types.MsgSubmitFinding{
					{
						ProgramId:        "not exist pid",
						FindingId:        "1",
						Title:            "Test bug 1",
						Description:      "Desc",
						SubmitterAddress: "Test address",
						SeverityLevel:    types.SeverityLevelCritical,
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
				ctx := types1.WrapSDKContext(suite.ctx)

				_, err := suite.msgServer.SubmitFinding(ctx, &finding)

				if tc.errArgs.shouldPass {
					suite.Require().NoError(err)
					_, exist := suite.keeper.GetFinding(suite.ctx, finding.FindingId)
					suite.Require().True(exist)
				} else {
					suite.Require().Error(err)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) InitCreateProgram() string {

	msgCreateProgram := &types.MsgCreateProgram{
		Name:            "name",
		Description:     "create test1",
		OperatorAddress: suite.address[0].String(),
		MemberAccounts:  []string{suite.address[1].String(), suite.address[2].String()},
		ProgramId:       "1",
	}

	ctx := types1.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.CreateProgram(ctx, msgCreateProgram)
	suite.Require().NoError(err)

	return msgCreateProgram.ProgramId
}

func (suite *KeeperTestSuite) InitSubmitFinding(pid string) string {
	msgSubmitFinding := &types.MsgSubmitFinding{
		ProgramId:        pid,
		FindingId:        "1",
		Title:            "Bug title",
		Description:      "Bug desc",
		SubmitterAddress: suite.address[0].String(),
		SeverityLevel:    types.SeverityLevelCritical,
	}

	ctx := types1.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.SubmitFinding(ctx, msgSubmitFinding)
	suite.Require().NoError(err)

	return msgSubmitFinding.FindingId
}

func (suite *KeeperTestSuite) TestAcceptFinding() {
	programId := suite.InitCreateProgram()
	findingId := suite.InitSubmitFinding(programId)

	testCases := []struct {
		name    string
		req     *types.MsgModifyFindingStatus
		expPass bool
	}{
		{
			"empty request",
			&types.MsgModifyFindingStatus{},
			false,
		},
		{
			"valid request => ",
			&types.MsgModifyFindingStatus{
				FindingId:       "1",
				OperatorAddress: suite.address[0].String(),
				Status:          types.FindingStatusConfirmed,
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := types1.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.AcceptFinding(ctx, testCase.req)

			finding, _ := suite.keeper.GetFinding(suite.ctx, findingId)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(finding.Status, types.FindingStatusConfirmed)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(finding.String(), types.FindingStatusReported)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHostRejectFinding() {
	programId := suite.InitCreateProgram()
	findingId := suite.InitSubmitFinding(programId)

	testCases := []struct {
		name    string
		req     *types.MsgModifyFindingStatus
		expPass bool
	}{
		{
			"empty request",
			&types.MsgModifyFindingStatus{},
			false,
		},
		{
			"valid request => comment is empty",
			&types.MsgModifyFindingStatus{
				FindingId:       "1",
				OperatorAddress: suite.address[0].String(),
				Status:          types.FindingStatusClosed,
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := types1.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.RejectFinding(ctx, testCase.req)

			finding, _ := suite.keeper.GetFinding(suite.ctx, findingId)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(finding.Status, types.FindingStatusClosed)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(finding.Status, types.FindingStatusReported)
			}
		})
	}
}
