package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"

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
						Detail:          "detail",
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
				ctx := sdk.WrapSDKContext(suite.ctx)

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

	pid := uuid.NewString()
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
						ProgramId:        pid,
						FindingId:        "1",
						Title:            "Test bug 1",
						Detail:           "detail",
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
						Detail:           "detail",
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
						Detail:           "detail",
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
				ctx := sdk.WrapSDKContext(suite.ctx)

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

func (suite *KeeperTestSuite) TestConfirmFinding() {
	pid, fid := "1", "1"
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)
	suite.InitSubmitFinding(pid, fid)

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
				OperatorAddress: suite.address[0].String(),
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.ConfirmFinding(ctx, testCase.req)

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

func (suite *KeeperTestSuite) TestCloseFinding() {
	pid, fid := "1", "1"
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)
	findingId := suite.InitSubmitFinding(pid, fid)

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
			"valid request => comment is empty",
			&types.MsgCloseFinding{
				FindingId:       "1",
				OperatorAddress: suite.address[0].String(),
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.CloseFinding(ctx, testCase.req)

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
		ProgramId:       pid,
		Name:            "name",
		Detail:          "detail",
		OperatorAddress: suite.address[0].String(),
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.CreateProgram(ctx, msgCreateProgram)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) InitActivateProgram(pid string) {
	msgCreateProgram := &types.MsgActivateProgram{
		ProgramId:       pid,
		OperatorAddress: suite.address[3].String(),
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.ActivateProgram(ctx, msgCreateProgram)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) InitSubmitFinding(pid, fid string) string {
	msgSubmitFinding := &types.MsgSubmitFinding{
		ProgramId:        pid,
		FindingId:        fid,
		Title:            "title",
		FindingHash:      "finding hash",
		SubmitterAddress: suite.address[0].String(),
		SeverityLevel:    types.Critical,
		Detail:           "detail",
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.SubmitFinding(ctx, msgSubmitFinding)
	suite.Require().NoError(err)

	return msgSubmitFinding.FindingId
}
