package v220

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// MigrateStore performs migration of votes. Specifically, it performs:
// - Conversion of votes from custom type to Cosmos SDK type
// nolint
func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, govtypes.VotesKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var oldVote Vote
		if err := cdc.Unmarshal(iterator.Value(), &oldVote); err != nil {
			return err
		}

		newVote := govtypesv1beta1.Vote{ProposalId: oldVote.ProposalId, Voter: oldVote.Voter, Option: oldVote.Option}
		bz, err := cdc.Marshal(&newVote)
		if err != nil {
			return err
		}

		store.Set(iterator.Key(), bz)
	}

	return nil
}
