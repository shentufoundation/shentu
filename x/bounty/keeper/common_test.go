package keeper_test

import (
	"crypto/rand"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/common"
	"github.com/shentufoundation/shentu/v2/x/bounty/client/cli"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (suite *KeeperTestSuite) CreatePrograms() (uint64, *ecies.PrivateKey) {
	// create a program
	decKey, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	suite.Require().NoError(err)
	encKey := crypto.FromECDSAPub(&decKey.ExportECDSA().PublicKey)
	deposit := sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(1e5)))
	var sET time.Time

	msg, err := types.NewMsgCreateProgram(suite.address[0].String(), "test", encKey, sdk.ZeroDec(), deposit, sET)
	suite.Require().NoError(err)
	res, err := suite.msgServer.CreateProgram(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res.ProgramId)
	return res.ProgramId, decKey
}

func (suite *KeeperTestSuite) CreateSubmitFinding(programId uint64, pubKey *ecies.PublicKey, desc string) uint64 {

	randBytes, reader := cli.GetRandBytes()
	encryptedDescBytes, err := ecies.Encrypt(reader, pubKey, []byte(desc), nil, nil)
	encryptedDescBytes = append(encryptedDescBytes, randBytes...)
	encDesc := types.EciesEncryptedDesc{
		FindingDesc: encryptedDescBytes,
	}
	descAny, err := codectypes.NewAnyWithValue(&encDesc)
	if err != nil {
		return 0
	}

	msgSubmitFinding := &types.MsgSubmitFinding{
		Title:            "title",
		EncryptedDesc:    descAny,
		ProgramId:        programId,
		EncryptedPoc:     nil,
		SeverityLevel:    types.SeverityLevelCritical,
		SubmitterAddress: suite.address[0].String(),
	}

	ctx := sdk.WrapSDKContext(suite.ctx)

	findingId, err := suite.keeper.GetNextFindingID(suite.ctx)
	resp, err := suite.msgServer.SubmitFinding(ctx, msgSubmitFinding)
	suite.Require().NoError(err)
	suite.Require().Equal(findingId, resp.FindingId)

	return findingId
}

func (suite *KeeperTestSuite) GetReleaseFinding(programId uint64, pubKey *ecies.PublicKey) uint64 {

	randBytes, reader := cli.GetRandBytes()
	encryptedDescBytes, err := ecies.Encrypt(reader, pubKey, []byte("test"), nil, nil)
	encryptedDescBytes = append(encryptedDescBytes, randBytes...)
	encDesc := types.EciesEncryptedDesc{
		FindingDesc: encryptedDescBytes,
	}
	descAny, err := codectypes.NewAnyWithValue(&encDesc)
	if err != nil {
		return 0
	}

	msgSubmitFinding := &types.MsgSubmitFinding{
		Title:            "title",
		EncryptedDesc:    descAny,
		ProgramId:        programId,
		EncryptedPoc:     nil,
		SeverityLevel:    types.SeverityLevelCritical,
		SubmitterAddress: suite.address[0].String(),
	}

	ctx := sdk.WrapSDKContext(suite.ctx)

	findingId, err := suite.keeper.GetNextFindingID(suite.ctx)
	resp, err := suite.msgServer.SubmitFinding(ctx, msgSubmitFinding)
	suite.Require().NoError(err)
	suite.Require().Equal(findingId, resp.FindingId)

	return findingId
}
