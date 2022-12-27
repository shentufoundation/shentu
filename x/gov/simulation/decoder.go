package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding gov type
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], govtypes.ProposalsKeyPrefix):
			var proposalA, proposalB govtypes.Proposal
			cdc.MustUnmarshal(kvA.Value, &proposalA)
			cdc.MustUnmarshal(kvB.Value, &proposalB)
			return fmt.Sprintf("%v\n%v", proposalA, proposalB)

		case bytes.Equal(kvA.Key[:1], govtypes.ActiveProposalQueuePrefix),
			bytes.Equal(kvA.Key[:1], govtypes.InactiveProposalQueuePrefix),
			bytes.Equal(kvA.Key[:1], govtypes.ProposalIDKey):
			proposalIDA := binary.LittleEndian.Uint64(kvA.Value)
			proposalIDB := binary.LittleEndian.Uint64(kvB.Value)
			return fmt.Sprintf("%d\n%d", proposalIDA, proposalIDB)

		case bytes.Equal(kvA.Key[:1], govtypes.DepositsKeyPrefix):
			var depositA, depositB govtypes.Deposit
			cdc.MustUnmarshal(kvA.Value, &depositA)
			cdc.MustUnmarshal(kvB.Value, &depositB)
			return fmt.Sprintf("%v\n%v", depositA, depositB)

		case bytes.Equal(kvA.Key[:1], govtypes.VotesKeyPrefix):
			var voteA, voteB govtypes.Vote
			cdc.MustUnmarshal(kvA.Value, &voteA)
			cdc.MustUnmarshal(kvB.Value, &voteB)
			return fmt.Sprintf("%v\n%v", voteA, voteB)

		default:
			panic(fmt.Sprintf("invalid governance key prefix %X", kvA.Key[:1]))
		}
	}
}
