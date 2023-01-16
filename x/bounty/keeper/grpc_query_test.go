package keeper_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/common"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryProgram() {
	addr, queryClient := suite.addrs, suite.queryClient

	var (
		req *types.QueryProgramRequest
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryProgramRequest{}
			},
			false,
		},
		{
			"non existing program request",
			func() {
				req = &types.QueryProgramRequest{ProgramId: 3}
			},
			false,
		},
		{
			"zero program id request",
			func() {
				req = &types.QueryProgramRequest{ProgramId: 0}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryProgramRequest{ProgramId: 1}

				// create a program
				decKey, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
				suite.Require().NoError(err)
				encKey := crypto.FromECDSAPub(&decKey.ExportECDSA().PublicKey)
				deposit := sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(1e5)))
				var sET, jET, cET time.Time

				msg, err := types.NewMsgCreateProgram(addr[0].String(), "test", encKey, sdk.ZeroDec(), deposit, sET, jET, cET)
				suite.Require().NoError(err)
				res, err := suite.msgServer.CreateProgram(sdk.WrapSDKContext(suite.ctx), msg)
				suite.Require().NoError(err)
				suite.Require().NotNil(res.ProgramId)
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			programRes, err := queryClient.Program(context.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(programRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryPrograms() {
	addr, queryClient := suite.addrs, suite.queryClient

	var (
		req  *types.QueryProgramsRequest
		size = 5
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryProgramsRequest{}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryProgramsRequest{
					FindingAddress: "",
					Pagination:     nil,
				}

				// create a program
				decKey, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
				suite.Require().NoError(err)
				encKey := crypto.FromECDSAPub(&decKey.ExportECDSA().PublicKey)
				deposit := sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(1e5)))
				var sET, jET, cET time.Time

				for i := 0; i < size; i++ {
					msg, err := types.NewMsgCreateProgram(addr[0].String(), "test", encKey, sdk.ZeroDec(), deposit, sET, jET, cET)
					suite.Require().NoError(err)
					res, err := suite.msgServer.CreateProgram(sdk.WrapSDKContext(suite.ctx), msg)
					suite.Require().NoError(err)
					suite.Require().NotNil(res.ProgramId)
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			programRes, err := queryClient.Programs(context.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Len(len(programRes.Programs), size)
			}
		})
	}

}
