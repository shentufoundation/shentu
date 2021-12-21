package keeper_test

import (
	"github.com/certikfoundation/shentu/v2/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (suite *KeeperTestSuite) TestQueryProposal() {
	ctx, queryClient := suite.ctx, suite.queryClient
	type proposal struct {
		title       string
		description string
	}
	tests := []struct {
		name       string
		proposal   proposal
		proposer   sdk.AccAddress
		proposalId uint64
		shouldPass bool
	}{
		{
			name: "Proposal ID exists",
			proposal: proposal{
				title:       "title0",
				description: "description0",
			},
			proposalId: 1,
			proposer:   suite.address[0],
			shouldPass: true,
		},
		{
			name: "Proposal ID does not exist",
			proposal: proposal{
				title:       "title1",
				description: "description1",
			},
			proposalId: 10,
			proposer:   suite.address[0],
			shouldPass: false,
		},
	}

	for _, tc := range tests {
		textProposalContent := govtypes.NewTextProposal(tc.proposal.title, tc.proposal.description)

		// submit a new proposal
		_, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, textProposalContent, tc.proposer)
		suite.Require().NoError(err)
	}

	for _, tc := range tests {
		queryResponse, err := queryClient.Proposal(ctx.Context(), &types.QueryProposalRequest{ProposalId: tc.proposalId})
		if tc.shouldPass {
			suite.Require().NoError(err)
			suite.Require().Equal(tc.proposalId, queryResponse.Proposal.ProposalId)
		} else {
			suite.Require().Error(err)
		}
	}
}
