package keeper_test

import (
	"context"
	"fmt"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryProgram() {
	_, queryClient := suite.addrs, suite.queryClient

	var (
		req *types.QueryProgramRequest
		//expProgram types.Program
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
		// TODO need add a bounty account
		//{
		//	"valid request",
		//	func() {
		//		req = &types.QueryProgramRequest{ProgramId: 1}
		//
		//		// create a program
		//		decKey, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
		//		suite.Require().NoError(err)
		//		encKey := crypto.FromECDSAPub(&decKey.ExportECDSA().PublicKey)
		//
		//		deposit := sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(1e5)))
		//
		//		var sET, jET, cET time.Time
		//
		//		msg, err := types.NewMsgCreateProgram(addr[0].String(), "test", encKey, sdk.ZeroDec(), deposit, sET, jET, cET)
		//		suite.Require().NoError(err)
		//		res, err := suite.msgServer.CreateProgram(sdk.WrapSDKContext(suite.ctx), msg)
		//		suite.Require().NoError(err)
		//		suite.Require().NotNil(res.ProgramId)
		//	},
		//	true,
		//},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			programRes, err := queryClient.Program(context.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				//suite.Require().Equal(expProgram.String(), programRes.Program.String())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(programRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryPrograms() {

}
