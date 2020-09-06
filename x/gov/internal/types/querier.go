package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// QueryProposalsParams defines data structure for querying 'custom/gov/proposals'.
type QueryProposalsParams struct {
	Page           int
	Limit          int
	Voter          sdk.AccAddress
	Depositor      sdk.AccAddress
	ProposalStatus ProposalStatus
}

// NewQueryProposalsParams creates a new instance of QueryProposalsParams.
func NewQueryProposalsParams(page, limit int, status ProposalStatus, voter, depositor sdk.AccAddress) QueryProposalsParams {
	return QueryProposalsParams{
		Page:           page,
		Limit:          limit,
		Voter:          voter,
		Depositor:      depositor,
		ProposalStatus: status,
	}
}
