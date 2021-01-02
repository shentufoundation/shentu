package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

// NewVote creates a new Vote instance.
func NewVote(proposalID uint64, voter sdk.AccAddress, option types.VoteOption, txhash string) Vote {
	vote := types.NewVote(proposalID, voter, option)
	return Vote{
		&vote,
		txhash,
	}
}

// Votes is a collection of Vote objects.
type Votes []Vote
