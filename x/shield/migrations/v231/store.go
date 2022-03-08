package v231

import (
	"github.com/gogo/protobuf/grpc"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1alpha1"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
)

const (
	stakingParamsPath = "/cosmos.staking.v1beta1.Query/Params"
)

func migratePools(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, types.PoolKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var oldPool v1alpha1.Pool
		cdc.MustUnmarshal(oldStoreIter.Value(), &oldPool)

		newPool := v1beta1.Pool{
			Id:          oldPool.Id,
			Description: oldPool.Description,
			SponsorAddr: oldPool.SponsorAddr,
			Active:      oldPool.Active,
			Shield:      oldPool.Shield,
			ShieldRate:  v1beta1.DefaultShieldRate,
		}

		newPoolBz := cdc.MustMarshal(&newPool)
		store.Set(oldStoreIter.Key(), newPoolBz)
	}
	return nil
}

func deleteUnusedStores(store sdk.KVStore) error {
	store.Delete(types.GetNextPurchaseIDKey())
	store.Delete(types.GetLastUpdateTimeKey())
	store.Delete(types.PurchaseKey)
	store.Delete(types.PurchaseQueueKey)
	return nil
}

func migrateparams(store sdk.KVStore, cdc codec.BinaryCodec, ps types.ParamSubspace) error {
	return nil
}

func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec, paramSpace types.ParamSubspace, queryServer grpc.Server) error {
	store := ctx.KVStore(storeKey)
	err := migratePools(store, cdc)
	if err != nil {
		return err
	}

	err = migrateparams(store, cdc, paramSpace)
	if err != nil {
		return err
	}

	return deleteUnusedStores(store)
}
