package utils

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/gov/client/utils"

	"github.com/certikfoundation/shentu/x/gov/internal/types"
)

// QueryProposerByTxQuery will query for a proposer of a governance proposal by ID.
func QueryProposerByTxQuery(cliCtx context.CLIContext, proposalID uint64, queryRoute string) (utils.Proposer, error) {
	res, err := utils.QueryProposalByID(proposalID, cliCtx, queryRoute)
	if err != nil {
		return utils.Proposer{}, err
	}

	var proposal types.Proposal
	if err := json.Unmarshal(res, &proposal); err != nil {
		return utils.Proposer{}, err
	}

	proposer := utils.NewProposer(proposalID, proposal.ProposerAddress.String())

	return proposer, nil
}
