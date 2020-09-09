package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/x/gov/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding gov type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], govTypes.ProposalsKeyPrefix):
		var proposalA, proposalB types.Proposal
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &proposalA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &proposalB)
		return fmt.Sprintf("%v\n%v", proposalA, proposalB)

	case bytes.Equal(kvA.Key[:1], govTypes.ActiveProposalQueuePrefix),
		bytes.Equal(kvA.Key[:1], govTypes.InactiveProposalQueuePrefix),
		bytes.Equal(kvA.Key[:1], govTypes.ProposalIDKey):
		proposalIDA := binary.LittleEndian.Uint64(kvA.Value)
		proposalIDB := binary.LittleEndian.Uint64(kvB.Value)
		return fmt.Sprintf("%d\n%d", proposalIDA, proposalIDB)

	case bytes.Equal(kvA.Key[:1], govTypes.DepositsKeyPrefix):
		var depositA, depositB types.Deposit
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &depositA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &depositB)
		return fmt.Sprintf("%v\n%v", depositA, depositB)

	case bytes.Equal(kvA.Key[:1], govTypes.VotesKeyPrefix):
		var voteA, voteB types.Vote
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &voteA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &voteB)
		return fmt.Sprintf("%v\n%v", voteA, voteB)

	default:
		panic(fmt.Sprintf("invalid governance key prefix %X", kvA.Key[:1]))
	}
}
