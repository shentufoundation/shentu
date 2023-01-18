package keeper_test

import (
	"crypto/rand"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

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

	deposit1 := types1.NewInt(10000)
	dd, _ := time.ParseDuration("24h")

	decKey, _ := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	encPubKey := crypto.FromECDSAPub(&decKey.ExportECDSA().PublicKey)

	encKeyMsg := types.EciesPubKey{
		EncryptionKey: encPubKey,
	}

	encAny, _ := codectypes.NewAnyWithValue(&encKeyMsg)

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{"Program(1)  -> Set: Simple",
			args{
				msgCresatePrograms: []types.MsgCreateProgram{
					{
						Description:       "create test1",
						CommissionRate:    types1.NewDec(1),
						SubmissionEndTime: time.Now().Add(dd),
						CreatorAddress:    suite.address[0].String(),
						EncryptionKey:     encAny,
						Deposit: []types1.Coin{
							{
								Denom:  "uctk",
								Amount: deposit1,
							},
						},
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
				programID := suite.keeper.GetNextProgramID(suite.ctx)
				resp, err := suite.msgServer.CreateProgram(ctx, &program)
				storedProgram, result := suite.keeper.GetProgram(suite.ctx, programID)
				if tc.errArgs.shouldPass {
					suite.Require().NoError(err)
					suite.Require().True(result)
					suite.Require().Equal(storedProgram.ProgramId, resp.ProgramId)
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

	programId := suite.InitCreateProgram()
	errorProgramId := suite.InitCreateErrorProgram()

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{"Submit finding(1)  -> submit: Simple",
			args{
				msgSubmitFindings: []types.MsgSubmitFinding{
					{
						Title:            "Test bug 1",
						Desc:             "This is real bug 1",
						ProgramId:        programId,
						Poc:              "bug1",
						SeverityLevel:    types.SeverityLevelCritical,
						SubmitterAddress: suite.address[0].String(),
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
						Title:            "Test bug 2",
						Desc:             "This is real bug 2",
						ProgramId:        200,
						Poc:              "bug2",
						SeverityLevel:    types.SeverityLevelCritical,
						SubmitterAddress: suite.address[0].String(),
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
						Title:            "Test bug 2",
						Desc:             "This is real bug 2",
						ProgramId:        200,
						Poc:              "bug2",
						SeverityLevel:    types.SeverityLevelCritical,
						SubmitterAddress: "test address",
					},
				},
			},
			errArgs{
				shouldPass: false,
			},
		},
		{"Submit finding(4)  -> submit: Simple",
			args{
				msgSubmitFindings: []types.MsgSubmitFinding{
					{
						Title:            "Test bug 2",
						Desc:             "This is real bug 2",
						ProgramId:        errorProgramId,
						Poc:              "bug2",
						SeverityLevel:    types.SeverityLevelCritical,
						SubmitterAddress: "test address",
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

				findingID := suite.keeper.GetNextFindingID(suite.ctx)
				resp, err := suite.msgServer.SubmitFinding(ctx, &finding)

				if tc.errArgs.shouldPass {
					suite.Require().NoError(err)
					f, exist := suite.keeper.GetFinding(suite.ctx, findingID)
					suite.Require().True(exist)
					suite.Require().Equal(f.FindingId, resp.FindingId)
				} else {
					suite.Require().Error(err)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) InitCreateProgram() uint64 {
	dd, _ := time.ParseDuration("24h")
	decKey, _ := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	encPubKey := crypto.FromECDSAPub(&decKey.ExportECDSA().PublicKey)

	encKeyMsg := types.EciesPubKey{
		EncryptionKey: encPubKey,
	}
	encAny, _ := codectypes.NewAnyWithValue(&encKeyMsg)

	deposit1 := types1.NewInt(10000)
	msgCreateProgram := types.MsgCreateProgram{
		Description:       "create test1",
		CommissionRate:    types1.NewDec(1),
		SubmissionEndTime: time.Now().Add(dd),
		CreatorAddress:    suite.address[0].String(),
		EncryptionKey:     encAny,
		Deposit: []types1.Coin{
			{
				Denom:  "uctk",
				Amount: deposit1,
			},
		},
	}

	ctx := types1.WrapSDKContext(suite.ctx)
	resp, err := suite.msgServer.CreateProgram(ctx, &msgCreateProgram)
	suite.Require().NoError(err)

	return resp.ProgramId
}

func (suite *KeeperTestSuite) InitCreateErrorProgram() uint64 {
	dd, _ := time.ParseDuration("24h")

	encKeyMsg := types.EciesPubKey{
		EncryptionKey: []byte{
			1, 2, 3, 5,
		},
	}
	encAny, _ := codectypes.NewAnyWithValue(&encKeyMsg)

	deposit1 := types1.NewInt(10000)
	msgCreateProgram := types.MsgCreateProgram{
		Description:       "create test1",
		CommissionRate:    types1.NewDec(1),
		SubmissionEndTime: time.Now().Add(dd),
		CreatorAddress:    suite.address[0].String(),
		EncryptionKey:     encAny,
		Deposit: []types1.Coin{
			{
				Denom:  "uctk",
				Amount: deposit1,
			},
		},
	}

	ctx := types1.WrapSDKContext(suite.ctx)
	resp, err := suite.msgServer.CreateProgram(ctx, &msgCreateProgram)
	suite.Require().NoError(err)

	return resp.ProgramId
}
