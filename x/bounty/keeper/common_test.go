package keeper_test

import (
	"crypto/rand"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/common"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (suite *KeeperTestSuite) CreatePrograms() uint64 {
	// create a program
	decKey, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	suite.Require().NoError(err)
	encKey := crypto.FromECDSAPub(&decKey.ExportECDSA().PublicKey)
	deposit := sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(1e5)))
	var sET, jET, cET time.Time

	msg, err := types.NewMsgCreateProgram(suite.address[0].String(), "test", encKey, sdk.ZeroDec(), deposit, sET, jET, cET)
	suite.Require().NoError(err)
	res, err := suite.msgServer.CreateProgram(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res.ProgramId)
	return res.ProgramId
}

func (suite *KeeperTestSuite) CreateSubmitFinding(proposalId uint64) uint64 {

	msgSubmitFinding := &types.MsgSubmitFinding{
		Title:            "title",
		Desc:             "desc",
		ProgramId:        proposalId,
		Poc:              "poc",
		SeverityLevel:    types.SeverityLevelCritical,
		SubmitterAddress: suite.address[0].String(),
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	findingId := suite.keeper.GetNextFindingID(suite.ctx)
	resp, err := suite.msgServer.SubmitFinding(ctx, msgSubmitFinding)
	suite.Require().NoError(err)
	suite.Require().Equal(findingId, resp.FindingId)

	return findingId
}
