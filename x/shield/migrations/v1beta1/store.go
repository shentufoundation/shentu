package v231

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	types "github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
)

const (
	stakingParamsPath = "/cosmos.staking.v1beta1.Query/Params"
)

func migratePools(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, types.PoolKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var oldPool types.Pool
		err := cdc.UnmarshalLengthPrefixed(oldStoreIter.Value(), &oldPool)
		if err != nil {
			return err
		}

		newPool := types.Pool{}
	}
}
