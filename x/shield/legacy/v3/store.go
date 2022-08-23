package v3

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	v2 "github.com/certikfoundation/shentu/v2/x/shield/legacy/v2"
	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

func migrateProviders(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, types.ProviderKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var oldProvider v2.Provider
		err := cdc.UnmarshalLengthPrefixed(oldStoreIter.Value(), &oldProvider)
		if err != nil {
			return err
		}

		newProvider := types.Provider{
			Address:          oldProvider.Address,
			DelegationBonded: oldProvider.DelegationBonded,
			Collateral:       oldProvider.Collateral,
			TotalLocked:      oldProvider.TotalLocked,
			Withdrawing:      oldProvider.Withdrawing,
			Rewards:          oldProvider.Rewards.Native.Add(oldProvider.Rewards.Foreign...),
		}

		newProviderBz := cdc.MustMarshal(&newProvider)
		oldStore.Set(oldStoreIter.Key(), newProviderBz)
	}

	return nil
}

func migratePurchases(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, types.PurchaseListKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var oldPurchaseList v2.PurchaseList
		err := cdc.UnmarshalLengthPrefixed(oldStoreIter.Value(), &oldPurchaseList)
		if err != nil {
			return err
		}

		entries := make([]types.Purchase, len(oldPurchaseList.Entries))
		for _, op := range oldPurchaseList.Entries {
			newPurchase := types.Purchase{
				PurchaseId:        op.PurchaseId,
				ProtectionEndTime: op.ProtectionEndTime,
				DeletionTime:      op.DeletionTime,
				Description:       op.Description,
				Shield:            op.Shield,
				ServiceFees:       op.ServiceFees.Native.Add(op.ServiceFees.Foreign...),
			}
			entries = append(entries, newPurchase)
		}
		newPurchaseList := types.PurchaseList{
			PoolId:    oldPurchaseList.PoolId,
			Purchaser: oldPurchaseList.Purchaser,
			Entries:   entries,
		}

		newPurchaseBz := cdc.MustMarshalLengthPrefixed(&newPurchaseList)
		oldStore.Set(oldStoreIter.Key(), newPurchaseBz)
	}

	return nil
}

func MigrateFees(store sdk.KVStore, cdc codec.BinaryCodec, key []byte) error {
	bz := store.Get(key)
	if bz == nil {
		panic("MigrateFees: key not found")
	}
	var oldFees v2.MixedDecCoins
	cdc.MustUnmarshalLengthPrefixed(bz, &oldFees)
	newFees := types.Fees{
		Fees: oldFees.Native.Add(oldFees.Foreign...),
	}
	bz = cdc.MustMarshalLengthPrefixed(&newFees)
	store.Set(key, bz)
	return nil
}

func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)

	providerMigrationErr := migrateProviders(store, cdc)
	if providerMigrationErr != nil {
		return providerMigrationErr
	}

	purchaseMigrationErr := migratePurchases(store, cdc)
	if purchaseMigrationErr != nil {
		return purchaseMigrationErr
	}

	feesKeys := [][]byte{
		types.GetServiceFeesKey(),
		types.GetBlockServiceFeesKey(),
		types.GetRemainingServiceFeesKey(),
	}
	for _, key := range feesKeys {
		MigrateFees(store, cdc, key)
	}

	return nil
}
