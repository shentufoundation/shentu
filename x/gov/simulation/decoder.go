package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/v2/x/gov/types"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding gov type
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], govtypes.ProposalsKeyPrefix):
			var proposalA, proposalB types.Proposal
			cdc.MustUnmarshal(kvA.Value, &proposalA)
			cdc.MustUnmarshal(kvB.Value, &proposalB)
			return fmt.Sprintf("%v\n%v", proposalA, proposalB)

		case bytes.Equal(kvA.Key[:1], govtypes.ActiveProposalQueuePrefix),
			bytes.Equal(kvA.Key[:1], govtypes.InactiveProposalQueuePrefix),
			bytes.Equal(kvA.Key[:1], govtypes.ProposalIDKey):
			proposalIDA := binary.LittleEndian.Uint64(kvA.Value)
			proposalIDB := binary.LittleEndian.Uint64(kvB.Value)
			return fmt.Sprintf("%d\n%d", proposalIDA, proposalIDB)

<<<<<<< HEAD
		case bytes.Equal(kvA.Key[:1], govTypes.DepositsKeyPrefix):
			var depositA, depositB types.Deposit
=======
		case bytes.Equal(kvA.Key[:1], govtypes.DepositsKeyPrefix):
			var depositA, depositB govtypes.Deposit
>>>>>>> 6f4b45bce5f277e193c4116dbea18212f40e242a
			cdc.MustUnmarshal(kvA.Value, &depositA)
			cdc.MustUnmarshal(kvB.Value, &depositB)
			return fmt.Sprintf("%v\n%v", depositA, depositB)

<<<<<<< HEAD
		case bytes.Equal(kvA.Key[:1], govTypes.VotesKeyPrefix):
			var voteA, voteB types.Vote
=======
		case bytes.Equal(kvA.Key[:1], govtypes.VotesKeyPrefix):
			var voteA, voteB govtypes.Vote
>>>>>>> 6f4b45bce5f277e193c4116dbea18212f40e242a
			cdc.MustUnmarshal(kvA.Value, &voteA)
			cdc.MustUnmarshal(kvB.Value, &voteB)
			return fmt.Sprintf("%v\n%v", voteA, voteB)

		default:
			panic(fmt.Sprintf("invalid governance key prefix %X", kvA.Key[:1]))
		}
	}
}
