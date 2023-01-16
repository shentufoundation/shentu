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

const SIZE = 5

func (suite *KeeperTestSuite) CreatePrograms() {
	// create a program
	decKey, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	suite.Require().NoError(err)
	encKey := crypto.FromECDSAPub(&decKey.ExportECDSA().PublicKey)
	deposit := sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(1e5)))
	var sET, jET, cET time.Time

	for i := 0; i < SIZE; i++ {
		msg, err := types.NewMsgCreateProgram(suite.address[0].String(), "test", encKey, sdk.ZeroDec(), deposit, sET, jET, cET)
		suite.Require().NoError(err)
		res, err := suite.msgServer.CreateProgram(sdk.WrapSDKContext(suite.ctx), msg)
		suite.Require().NoError(err)
		suite.Require().NotNil(res.ProgramId)
	}
}
