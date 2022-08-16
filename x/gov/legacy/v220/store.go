package v220

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// MigrateStore performs migration of votes. Specifically, it performs:
// - Conversion of votes from custom type to Cosmos SDK type
//nolint
func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, govtypes.VotesKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var oldVote Vote
		if err := cdc.Unmarshal(iterator.Value(), &oldVote); err != nil {
			return err
		}

		newVote := govtypes.Vote{ProposalId: oldVote.ProposalId, Voter: oldVote.Voter, Option: oldVote.Option}
		bz, err := cdc.Marshal(&newVote)
		if err != nil {
			return err
		}

		store.Set(iterator.Key(), bz)
	}

	return nil
}
