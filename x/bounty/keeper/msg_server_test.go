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

	pid := "1"
	suite.InitCreateProgram(pid)
	suite.InitOpenProgram(pid)

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
						SeverityLevel:    types.Critical,
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
						SeverityLevel:    types.Critical,
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
						SeverityLevel:    types.Critical,
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

func (suite *KeeperTestSuite) TestAcceptFinding() {
	pid, fid := "1", "1"
	suite.InitCreateProgram(pid)
	suite.InitOpenProgram(pid)
	suite.InitSubmitFinding(pid, fid)

	testCases := []struct {
		name    string
		req     *types.MsgAcceptFinding
		expPass bool
	}{
		{
			"empty request",
			&types.MsgAcceptFinding{},
			false,
		},
		{
			"valid request => ",
			&types.MsgAcceptFinding{
				FindingId:       fid,
				OperatorAddress: suite.address[0].String(),
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := types1.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.AcceptFinding(ctx, testCase.req)

			finding, _ := suite.keeper.GetFinding(suite.ctx, fid)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(finding.Status, types.FindingStatusConfirmed)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(finding.Status, types.FindingStatusSubmitted)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestRejectFinding() {
	pid, fid := "1", "1"
	suite.InitCreateProgram(pid)
	suite.InitOpenProgram(pid)
	findingId := suite.InitSubmitFinding(pid, fid)

	testCases := []struct {
		name    string
		req     *types.MsgRejectFinding
		expPass bool
	}{
		{
			"empty request",
			&types.MsgRejectFinding{},
			false,
		},
		{
			"valid request => comment is empty",
			&types.MsgRejectFinding{
				FindingId:       "1",
				OperatorAddress: suite.address[0].String(),
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
				suite.Require().Equal(finding.Status, types.FindingStatusSubmitted)
			}
		})
	}
}

func (suite *KeeperTestSuite) InitCreateProgram(pid string) {

	msgCreateProgram := &types.MsgCreateProgram{
		Name:            "name",
		Description:     "create test1",
		OperatorAddress: suite.address[0].String(),
		MemberAccounts:  []string{suite.address[1].String(), suite.address[2].String()},
		ProgramId:       pid,
	}

	ctx := types1.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.CreateProgram(ctx, msgCreateProgram)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) InitOpenProgram(pid string) {
	msgCreateProgram := &types.MsgOpenProgram{
		ProgramId:       pid,
		OperatorAddress: suite.address[3].String(),
	}

	ctx := types1.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.OpenProgram(ctx, msgCreateProgram)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) InitSubmitFinding(pid, fid string) string {
	msgSubmitFinding := &types.MsgSubmitFinding{
		ProgramId:        pid,
		FindingId:        fid,
		Title:            "Bug title",
		Description:      "Bug desc",
		SubmitterAddress: suite.address[0].String(),
		SeverityLevel:    types.Critical,
	}

	ctx := types1.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.SubmitFinding(ctx, msgSubmitFinding)
	suite.Require().NoError(err)

	return msgSubmitFinding.FindingId
}
