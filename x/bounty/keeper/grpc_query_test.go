package keeper_test

import (
	"context"
	"fmt"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryProgram() {
	queryClient := suite.queryClient

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

				// create programs
				suite.CreatePrograms()
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
	queryClient := suite.queryClient

	var (
		req *types.QueryProgramsRequest
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryProgramsRequest{
					FindingAddress: "",
					Pagination:     nil,
				}

				// create programs
				suite.CreatePrograms()
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
				suite.Require().Equal(len(programRes.Programs), SIZE)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(programRes)
			}
		})
	}
}
