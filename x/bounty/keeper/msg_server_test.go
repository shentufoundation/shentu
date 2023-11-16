package keeper_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"

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
			for _, program := range tc.args.msgCreatePrograms {
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
						Title:           "Test bug 1",
						Detail:          "detail",
						OperatorAddress: suite.address[0].String(),
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
						Title:           "Test bug 1",
						Detail:          "detail",
						OperatorAddress: suite.address[0].String(),
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
						Title:           "Test bug 1",
						Detail:          "detail",
						OperatorAddress: "Test address",
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
					_, exist := suite.keeper.GetFinding(suite.ctx, finding.FindingId)
					suite.Require().True(exist)
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

	_, found := suite.keeper.GetFinding(suite.ctx, fid)
	suite.Require().True(found)

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
				OperatorAddress: suite.address[0].String(),
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.ActivateFinding(ctx, testCase.req)

			finding, _ := suite.keeper.GetFinding(suite.ctx, fid)

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

	finding, found := suite.keeper.GetFinding(suite.ctx, fid)
	suite.Require().True(found)
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
				OperatorAddress: suite.address[0].String(),
				Fingerprint:     findingFingerPrintHash,
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

func (suite *KeeperTestSuite) TestConfirmFindingPaid() {
	pid, fid := uuid.NewString(), uuid.NewString()
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)
	suite.InitSubmitFinding(pid, fid)

	finding, found := suite.keeper.GetFinding(suite.ctx, fid)
	suite.Require().True(found)
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
				OperatorAddress: suite.address[0].String(),
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.ConfirmFindingPaid(ctx, testCase.req)
			finding, _ := suite.keeper.GetFinding(suite.ctx, fid)

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
				OperatorAddress: suite.address[0].String(),
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.CloseFinding(ctx, testCase.req)
			finding, _ := suite.keeper.GetFinding(suite.ctx, fid)

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

func (suite *KeeperTestSuite) TestReleaseConfirmFinding() {
	pid, fid := uuid.NewString(), uuid.NewString()
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)
	suite.InitSubmitFinding(pid, fid)

	finding, found := suite.keeper.GetFinding(suite.ctx, fid)
	suite.Require().True(found)
	findingFingerPrintHash := suite.app.BountyKeeper.GetFindingFingerprintHash(&finding)
	suite.InitConfirmFinding(fid, findingFingerPrintHash)

	suite.InitConfirmFindingPaid(fid)

	testCases := []struct {
		name    string
		req     *types.MsgReleaseFinding
		expPass bool
	}{
		{
			"empty request",
			&types.MsgReleaseFinding{},
			false,
		},
		{
			"valid request",
			&types.MsgReleaseFinding{
				FindingId:       fid,
				OperatorAddress: suite.address[0].String(),
				Description:     "desc",
				ProofOfConcept:  "poc",
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.msgServer.ReleaseFinding(ctx, testCase.req)

			finding, _ := suite.keeper.GetFinding(suite.ctx, fid)

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
	desc, poc := "desc", "poc"
	hash := sha256.Sum256([]byte(desc + poc + suite.address[0].String()))

	msgSubmitFinding := &types.MsgSubmitFinding{
		ProgramId:       pid,
		FindingId:       fid,
		Title:           "title",
		FindingHash:     hex.EncodeToString(hash[:]),
		OperatorAddress: suite.address[0].String(),
		SeverityLevel:   types.Critical,
		Detail:          "detail",
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.SubmitFinding(ctx, msgSubmitFinding)
	suite.Require().NoError(err)

	return msgSubmitFinding.FindingId
}

func (suite *KeeperTestSuite) InitConfirmFinding(fid, fingerprint string) string {
	msgConfirmFinding := &types.MsgConfirmFinding{
		FindingId:       fid,
		OperatorAddress: suite.address[0].String(),
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
		OperatorAddress: suite.address[0].String(),
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.msgServer.ConfirmFindingPaid(ctx, msgConfirmFindingPaid)
	suite.Require().NoError(err)

	return msgConfirmFindingPaid.FindingId
}
