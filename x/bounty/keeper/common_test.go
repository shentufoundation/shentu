package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (suite *KeeperTestSuite) CreatePrograms(pid string) {
	// create a program
	msg, err := types.NewMsgCreateProgram("name", "desc", pid, suite.address[0], nil, nil)

	suite.Require().NoError(err)
	_, err = suite.msgServer.CreateProgram(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) CreateSubmitFinding(pid, fid string) {
	ctx := sdk.WrapSDKContext(suite.ctx)

	msgSubmitFinding := &types.MsgSubmitFinding{
		ProgramId:        pid,
		FindingId:        fid,
		Title:            "title",
		Description:      "desc",
		SubmitterAddress: suite.address[0].String(),
		SeverityLevel:    types.SeverityLevelHigh,
	}

	_, err := suite.msgServer.SubmitFinding(ctx, msgSubmitFinding)
	suite.Require().NoError(err)
}

//func (suite *KeeperTestSuite) GetReleaseFinding(programId uint64, pubKey *ecies.PublicKey) uint64 {
//
//	randBytes, reader := cli.GetRandBytes()
//	encryptedDescBytes, err := ecies.Encrypt(reader, pubKey, []byte("test"), nil, nil)
//	encryptedDescBytes = append(encryptedDescBytes, randBytes...)
//	encDesc := types.EciesEncryptedDesc{
//		FindingDesc: encryptedDescBytes,
//	}
//	descAny, err := codectypes.NewAnyWithValue(&encDesc)
//	if err != nil {
//		return 0
//	}
//
//	msgSubmitFinding := &types.MsgSubmitFinding{
//		Title:            "title",
//		EncryptedDesc:    descAny,
//		ProgramId:        programId,
//		EncryptedPoc:     nil,
//		SeverityLevel:    types.SeverityLevelCritical,
//		SubmitterAddress: suite.address[0].String(),
//	}
//
//	ctx := sdk.WrapSDKContext(suite.ctx)
//
//	findingId, err := suite.keeper.GetNextFindingID(suite.ctx)
//	resp, err := suite.msgServer.SubmitFinding(ctx, msgSubmitFinding)
//	suite.Require().NoError(err)
//	suite.Require().Equal(findingId, resp.FindingId)
//
//	return findingId
//}
