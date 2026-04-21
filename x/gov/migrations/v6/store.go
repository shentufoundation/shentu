package v6

import (
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var CertVotesKeyPrefix = []byte("certvote")

func MigrateStore(ctx sdk.Context, storeService corestoretypes.KVStoreService) error {
	kv := storeService.OpenKVStore(ctx)
	certVotes := prefix.NewStore(runtime.KVStoreAdapter(kv), CertVotesKeyPrefix)

	it := certVotes.Iterator(nil, nil)
	defer it.Close()

	var keys [][]byte
	for ; it.Valid(); it.Next() {
		keys = append(keys, append([]byte(nil), it.Key()...))
	}
	for _, k := range keys {
		certVotes.Delete(k)
	}
	return nil
}
