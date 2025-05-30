package keeper_test

import (
	"fmt"

	"github.com/google/uuid"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

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
				req = &types.QueryProgramRequest{ProgramId: "3"}
			},
			false,
		},
		{
			"zero program id request",
			func() {
				req = &types.QueryProgramRequest{ProgramId: "0"}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryProgramRequest{ProgramId: "1"}
				// create programs
				suite.InitCreateProgram("1")
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()
			ctx := sdk.WrapSDKContext(suite.ctx)
			programRes, err := queryClient.Program(ctx, req)

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
					Pagination: nil,
				}

				// create two programs
				suite.InitCreateProgram("1")
				suite.InitCreateProgram("2")
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()
			ctx := sdk.WrapSDKContext(suite.ctx)
			programRes, err := queryClient.Programs(ctx, req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(len(programRes.Programs), 2)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(programRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryFinding() {
	queryClient := suite.queryClient

	// create programs
	pid := uuid.NewString()
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)
	var (
		req *types.QueryFindingRequest
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryFindingRequest{}
			},
			false,
		},
		{
			"non existing finding id request",
			func() {
				req = &types.QueryFindingRequest{FindingId: "100"}
			},
			false,
		},
		{
			"zero finding id request",
			func() {
				req = &types.QueryFindingRequest{FindingId: "1"}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryFindingRequest{FindingId: "1"}
				suite.InitSubmitFinding(pid, "1")
			},
			true,
		},
		{
			"valid request",
			func() {
				req = &types.QueryFindingRequest{FindingId: "2"}
				suite.InitSubmitFinding(pid, "2")

				ctx := sdk.WrapSDKContext(suite.ctx)
				suite.msgServer.PublishFinding(ctx, types.NewMsgPublishFinding("2", "desc", "poc", suite.address[0]))
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()
			ctx := sdk.WrapSDKContext(suite.ctx)
			findingRes, err := queryClient.Finding(ctx, req)
			if testCase.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(findingRes)
			}
		})
	}

}

func (suite *KeeperTestSuite) TestGRPCQueryFindings() {
	queryClient := suite.queryClient

	pid, fid := uuid.NewString(), uuid.NewString()
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)
	suite.InitSubmitFinding(pid, fid)
	suite.InitSubmitFinding(pid, uuid.NewString())

	testCases := []struct {
		msg          string
		req          *types.QueryFindingsRequest
		expResultLen int
		expPass      bool
	}{
		{
			"invalid request",
			&types.QueryFindingsRequest{},
			0,
			false,
		},
		{
			"valid request => piq and submitter address",
			&types.QueryFindingsRequest{ProgramId: pid, SubmitterAddress: suite.whiteHatAddr.String()},
			2,
			true,
		},
		{
			"valid request => piq",
			&types.QueryFindingsRequest{ProgramId: pid},
			2,
			true,
		},
		{
			"valid request => submitter address",
			&types.QueryFindingsRequest{SubmitterAddress: suite.whiteHatAddr.String()},
			2,
			true,
		},
		{
			"valid request => invalid pid",
			&types.QueryFindingsRequest{ProgramId: "not exist"},
			0,
			true,
		},
		{
			"valid request => invalid submitter address",
			&types.QueryFindingsRequest{SubmitterAddress: suite.normalAddr.String()},
			0,
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			findingRes, err := queryClient.Findings(ctx, testCase.req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(len(findingRes.Findings), testCase.expResultLen)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(findingRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryFindingFingerprint() {
	queryClient := suite.queryClient

	// Create a program, activate it, and submit a finding
	pid := uuid.NewString()
	fid := uuid.NewString()
	suite.InitCreateProgram(pid)
	suite.InitActivateProgram(pid)
	suite.InitSubmitFinding(pid, fid)

	testCases := []struct {
		name    string
		req     *types.QueryFindingFingerprintRequest
		expPass bool
	}{
		{
			"empty request",
			&types.QueryFindingFingerprintRequest{},
			false,
		},
		{
			"non-existent finding ID",
			&types.QueryFindingFingerprintRequest{FindingId: "non-existent"},
			false,
		},
		{
			"valid request",
			&types.QueryFindingFingerprintRequest{FindingId: fid},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := queryClient.FindingFingerprint(ctx, tc.req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotEmpty(res.Fingerprint)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryProgramFingerprint() {
	queryClient := suite.queryClient

	// Create a program
	pid := uuid.NewString()
	suite.InitCreateProgram(pid)

	testCases := []struct {
		name    string
		req     *types.QueryProgramFingerprintRequest
		expPass bool
	}{
		{
			"empty request",
			&types.QueryProgramFingerprintRequest{},
			false,
		},
		{
			"non-existent program ID",
			&types.QueryProgramFingerprintRequest{ProgramId: "non-existent"},
			false,
		},
		{
			"valid request",
			&types.QueryProgramFingerprintRequest{ProgramId: pid},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := queryClient.ProgramFingerprint(ctx, tc.req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotEmpty(res.Fingerprint)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryTheorems() {
	queryClient := suite.queryClient

	// Create some theorems
	_ = suite.InitCreateTheorem()
	_ = suite.InitCreateTheorem()

	testCases := []struct {
		name          string
		req           *types.QueryTheoremsRequest
		expResultsLen int
		expPass       bool
	}{
		{
			"empty request (valid)",
			&types.QueryTheoremsRequest{},
			2, // We created 2 theorems above
			true,
		},
		{
			"with pagination",
			&types.QueryTheoremsRequest{
				Pagination: &query.PageRequest{
					Limit: 1,
				},
			},
			1, // Should return only 1 theorem due to pagination limit
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := queryClient.Theorems(ctx, tc.req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Len(res.Theorems, tc.expResultsLen)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryTheorem() {
	queryClient := suite.queryClient

	// Create a theorem
	theoremID := suite.InitCreateTheorem()

	testCases := []struct {
		name    string
		req     *types.QueryTheoremRequest
		expPass bool
	}{
		{
			"empty request",
			&types.QueryTheoremRequest{},
			false,
		},
		{
			"zero theorem ID",
			&types.QueryTheoremRequest{TheoremId: 0},
			false,
		},
		{
			"non-existent theorem ID",
			&types.QueryTheoremRequest{TheoremId: 9999},
			false,
		},
		{
			"valid request",
			&types.QueryTheoremRequest{TheoremId: theoremID},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := queryClient.Theorem(ctx, tc.req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(tc.req.TheoremId, res.Theorem.Id)
				suite.Require().Equal("Test Theorem", res.Theorem.Title)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryProof() {
	queryClient := suite.queryClient

	// Create a theorem and submit a proof hash
	theoremID := suite.InitCreateTheorem()
	proofHash := suite.InitSubmitProofHash(theoremID)

	testCases := []struct {
		name    string
		req     *types.QueryProofRequest
		expPass bool
	}{
		{
			"empty request",
			&types.QueryProofRequest{},
			false,
		},
		{
			"empty proof ID",
			&types.QueryProofRequest{ProofId: ""},
			false,
		},
		{
			"non-existent proof ID",
			&types.QueryProofRequest{ProofId: "non-existent"},
			false,
		},
		{
			"valid request",
			&types.QueryProofRequest{ProofId: proofHash},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := queryClient.Proof(ctx, tc.req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(tc.req.ProofId, res.Proof.Id)
				suite.Require().Equal(theoremID, res.Proof.TheoremId)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryProofs() {
	queryClient := suite.queryClient

	// Create a theorem and submit multiple proofs
	theoremID := suite.InitCreateTheorem()

	// Submit a proof hashes for the same theorem
	_ = suite.InitSubmitProofHash(theoremID)

	testCases := []struct {
		name          string
		req           *types.QueryProofsRequest
		expResultsLen int
		expPass       bool
	}{
		{
			"empty request",
			&types.QueryProofsRequest{},
			0,
			false,
		},
		{
			"zero theorem ID",
			&types.QueryProofsRequest{TheoremId: 0},
			0,
			false,
		},
		{
			"non-existent theorem ID",
			&types.QueryProofsRequest{TheoremId: 9999},
			0,
			true, // Should pass but return empty results
		},
		{
			"valid request",
			&types.QueryProofsRequest{TheoremId: theoremID},
			1, // We submitted one proof
			true,
		},
		{
			"valid request with pagination",
			&types.QueryProofsRequest{
				TheoremId: theoremID,
				Pagination: &query.PageRequest{
					Limit: 2,
				},
			},
			1,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := queryClient.Proofs(ctx, tc.req)

			if tc.expPass {
				suite.Require().NoError(err)
				if tc.expResultsLen > 0 {
					suite.Require().NotNil(res)
					suite.Require().Len(res.Proofs, tc.expResultsLen)
				}
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryReward() {
	queryClient := suite.queryClient

	// Create a theorem and submit a proof
	theoremID := suite.InitCreateTheorem()
	proofHash := suite.InitSubmitProofHash(theoremID)

	// Submit proof detail
	suite.InitSubmitProofDetail(proofHash)

	// Verify the proof to generate a reward
	suite.InitVerifyProof(proofHash, types.ProofStatus_PROOF_STATUS_PASSED)

	testCases := []struct {
		name    string
		req     *types.QueryRewardsRequest
		expPass bool
	}{
		{
			"empty request",
			&types.QueryRewardsRequest{},
			false,
		},
		{
			"invalid address",
			&types.QueryRewardsRequest{Address: "invalid-address"},
			false,
		},
		{
			"address with no rewards",
			&types.QueryRewardsRequest{Address: suite.normalAddr.String()},
			false,
		},
		{
			"valid request - prover",
			&types.QueryRewardsRequest{Address: suite.whiteHatAddr.String()},
			true,
		},
		{
			"valid request - checker",
			&types.QueryRewardsRequest{Address: suite.bountyAdminAddr.String()},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := queryClient.Reward(ctx, tc.req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().NotEmpty(res.Rewards)
				suite.Require().False(res.Rewards[0].Amount.IsZero())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryParams() {
	queryClient := suite.queryClient

	testCases := []struct {
		name    string
		req     *types.QueryParamsRequest
		expPass bool
	}{
		{
			"valid request",
			&types.QueryParamsRequest{},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := queryClient.Params(ctx, tc.req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().NotNil(res.Params)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryGrants() {
	queryClient := suite.queryClient

	// Create a theorem and add some grants
	theoremID := suite.InitCreateTheorem()

	// Add another grant from a different address
	bondDenom, err := suite.app.StakingKeeper.BondDenom(suite.ctx)
	suite.Require().NoError(err)

	grantMsg := &types.MsgGrant{
		TheoremId: theoremID,
		Grantor:   suite.normalAddr.String(),
		Amount:    sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(2000000))),
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err = suite.msgServer.Grant(ctx, grantMsg)
	suite.Require().NoError(err)

	testCases := []struct {
		name          string
		req           *types.QueryGrantsRequest
		expResultsLen int
		expPass       bool
	}{
		{
			"empty request",
			&types.QueryGrantsRequest{},
			0,
			false,
		},
		{
			"zero theorem ID",
			&types.QueryGrantsRequest{TheoremId: 0},
			0,
			false,
		},
		{
			"non-existent theorem ID",
			&types.QueryGrantsRequest{TheoremId: 9999},
			0,
			true, // Should pass but return empty results
		},
		{
			"valid request",
			&types.QueryGrantsRequest{TheoremId: theoremID},
			2, // The original grant plus the one we just added
			true,
		},
		{
			"valid request with pagination",
			&types.QueryGrantsRequest{
				TheoremId: theoremID,
				Pagination: &query.PageRequest{
					Limit: 1,
				},
			},
			1, // Pagination limit is 1
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := queryClient.Grants(ctx, tc.req)

			if tc.expPass {
				suite.Require().NoError(err)
				if tc.expResultsLen > 0 {
					suite.Require().NotNil(res)
					suite.Require().Len(res.Grants, tc.expResultsLen)
				}
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}
