package utils

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/gov/client/utils"

	"github.com/certikfoundation/shentu/x/gov/types"
)

// QueryProposerByTxQuery will query for a proposer of a governance proposal by ID.
func QueryProposerByTxQuery(cliCtx client.Context, proposalID uint64, queryRoute string) (utils.Proposer, error) {
	res, err := utils.QueryProposalByID(proposalID, cliCtx, queryRoute)
	if err != nil {
		return utils.Proposer{}, err
	}

	var proposal types.Proposal
	if err := json.Unmarshal(res, &proposal); err != nil {
		return utils.Proposer{}, err
	}

	proposer := utils.NewProposer(proposalID, proposal.ProposerAddress)

	return proposer, nil
}
