package v1beta1

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	v1beta1 "github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
	v1alpha1 "github.com/certikfoundation/shentu/v2/x/shield/types/v1alpha1"
)

const (
	stakingParamsPath = "/cosmos.staking.v1beta1.Query/Params"
)

func migratePools(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, v1beta1.PoolKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var oldPool v1alpha1.Pool
		err := cdc.UnmarshalLengthPrefixed(oldStoreIter.Value(), &oldPool)
		if err != nil {
			return err
		}

		newPool := v1beta1.Pool{
			Id: oldPool.Id,
			Description: oldPool.Description,
			Sponsor: oldPool.Sponsor,
			SponsorAddr: oldPool.SponsorAddr,
			ShieldLimit: oldPool.ShieldLimit,
			Active: oldPool.Active,
			Shield: oldPool.Shield,
		}

		oldStore.Delete(oldStoreIter.Key())
		newPoolBz := cdc.MustMarshal(&newPool)
		oldStore.Set(oldStoreIter.Key(), newPoolBz)
	}

	return nil
}

func migrateProviders(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, v1beta1.ProviderKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var oldProvider v1alpha1.Provider
		err := cdc.UnmarshalLengthPrefixed(oldStoreIter.Value(), &oldProvider)
		if err != nil {
			return err;
		}

		newProvider := v1beta1.Provider{
			Address: oldProvider.Address,
			DelegationBonded: oldProvider.DelegationBonded,
			Collateral: oldProvider.Collateral,
			TotalLocked: oldProvider.TotalLocked,
			Withdrawing: oldProvider.Withdrawing,
			Rewards: oldProvider.Rewards.Native.Add(oldProvider.Rewards.Foreign...),
		}

		oldStore.Delete(oldStoreIter.Key())
		newProviderBz := cdc.MustMarshal(&newProvider)
		oldStore.Set(oldStoreIter.Key(), newProviderBz)
	}

	return nil
}

func migratePurchases(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, v1beta1.PurchaseListKey)

	oldStoreIter := oldStore.Iterator(nil, nil) 
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var oldPurchase v1alpha1.Purchase
		err := cdc.UnmarshalLengthPrefixed(oldStoreIter.Value(), &oldPurchase)
		if err != nil {
			return err;
		}

		newPurchase := v1beta1.Purchase{
			PurchaseId: oldPurchase.PurchaseId,
			ProtectionEndTime: oldPurchase.ProtectionEndTime,
			DeletionTime: oldPurchase.DeletionTime,
			Description: oldPurchase.Description,
			Shield: oldPurchase.Shield,
			Fees: oldPurchase.ServiceFees.Native.Add(oldPurchase.ServiceFees.Foreign...),
		}

		oldStore.Delete(oldStoreIter.Key())
		newPurchaseBz := cdc.MustMarshal(&newPurchase)
		oldStore.Set(oldStoreIter.Key(), newPurchaseBz)
	}

	return nil
}

func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)

	poolMigrationErr := migratePools(store, cdc) 
	if poolMigrationErr != nil {
		return poolMigrationErr
	}

	providerMigrationErr := migrateProviders(store, cdc)
	if providerMigrationErr != nil {
		return providerMigrationErr
	}

	purchaseMigrationErr := migratePurchases(store, cdc)
	if purchaseMigrationErr != nil {
		return purchaseMigrationErr
	}

	return nil
}