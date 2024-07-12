package v5

import (
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/shield/types"
)

func migratePool(store sdk.KVStore, cdc codec.BinaryCodec) error {
	poolStore := prefix.NewStore(store, PoolKey)
	iterator := poolStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		poolStore.Delete(iterator.Key())
	}

	store.Delete(PoolKey)
	return nil
}

func migratePurchaseListKey(store sdk.KVStore, cdc codec.BinaryCodec) error {
	s := prefix.NewStore(store, PurchaseListKey)

	iterator := s.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		s.Delete(iterator.Key())
	}
	s.Delete(PurchaseListKey)
	return nil
}

func migratePurchaseQueueKey(store sdk.KVStore, cdc codec.BinaryCodec) error {
	s := prefix.NewStore(store, PurchaseQueueKey)

	iterator := s.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		s.Delete(iterator.Key())
	}

	s.Delete(PurchaseQueueKey)
	return nil
}

func migrateProviders(store sdk.KVStore, cdc codec.BinaryCodec) error {
	providerStore := prefix.NewStore(store, types.ProviderKey)

	iterator := providerStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var provider types.Provider
		err := cdc.UnmarshalLengthPrefixed(iterator.Value(), &provider)
		if err != nil {
			return err
		}
		provider.Withdrawing = math.NewInt(0)
		bz := cdc.MustMarshalLengthPrefixed(&provider)
		providerStore.Set(iterator.Key(), bz)
	}
	return nil
}

func migrateWithdraws(store sdk.KVStore, cdc codec.BinaryCodec) error {
	s := prefix.NewStore(store, WithdrawQueueKey)

	iterator := s.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		s.Delete(iterator.Key())
	}

	s.Delete(WithdrawQueueKey)
	return nil
}

func migrateStakeForShield(store sdk.KVStore, cdc codec.BinaryCodec) error {
	s := prefix.NewStore(store, StakeForShieldKey)

	iterator := s.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		s.Delete(iterator.Key())
	}

	s.Delete(StakeForShieldKey)

	return nil
}

func migrateOriginalStakingKey(store sdk.KVStore, cdc codec.BinaryCodec) error {
	iter := sdk.KVStoreReversePrefixIterator(store, OriginalStakingKey)
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}

	store.Delete(OriginalStakingKey)
	return nil
}

func migrateReimbursementKey(store sdk.KVStore, cdc codec.BinaryCodec) error {
	iter := sdk.KVStoreReversePrefixIterator(store, ReimbursementKey)
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}

	store.Delete(ReimbursementKey)
	return nil
}

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	store.Delete(ShieldAdminKey)
	store.Delete(TotalCollateralKey)
	store.Delete(TotalWithdrawingKey)
	store.Delete(TotalShieldKey)
	store.Delete(TotalClaimedKey)
	store.Delete(ServiceFeesKey)
	if err := migratePool(store, cdc); err != nil {
		return err
	}
	store.Delete(NextPoolIDKey)
	store.Delete(NextPurchaseIDKey)
	if err := migratePurchaseListKey(store, cdc); err != nil {
		return err
	}
	if err := migratePurchaseQueueKey(store, cdc); err != nil {
		return err
	}
	if err := migrateProviders(store, cdc); err != nil {
		return err
	}
	if err := migrateWithdraws(store, cdc); err != nil {
		return err
	}
	store.Delete(LastUpdateTimeKey)
	store.Delete(GlobalStakeForShieldPoolKey)
	if err := migrateStakeForShield(store, cdc); err != nil {
		return err
	}
	if err := migrateOriginalStakingKey(store, cdc); err != nil {
		return err
	}
	if err := migrateReimbursementKey(store, cdc); err != nil {
		return err
	}
	return nil
}
