package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

// Vote wraps a vote and corresponding txhash.
type Vote struct {
	types.Vote
	TxHash string `json:"txhash" yaml:"txhash"`
}

// NewVote creates a new Vote instance.
func NewVote(proposalID uint64, voter sdk.AccAddress, option types.VoteOption, txhash string) Vote {
	return Vote{
		types.NewVote(proposalID, voter, option),
		txhash,
	}
}

// Votes is a collection of Vote objects.
type Votes []Vote
