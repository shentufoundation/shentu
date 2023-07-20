package v280

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shentufoundation/shentu/v2/common"
	"github.com/shentufoundation/shentu/v2/x/shield/types"
)

func migratePool(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, types.PoolKey)

	iterator := oldStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var oldPool types.Pool
		cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &oldPool)

		shentuSponsorAddr, err := common.PrefixToShentu(oldPool.SponsorAddr)
		if err != nil {
			return err
		}

		newPool := types.Pool{
			Id:          oldPool.Id,
			Description: oldPool.Description,
			Sponsor:     oldPool.Sponsor,
			SponsorAddr: shentuSponsorAddr,
			ShieldLimit: oldPool.ShieldLimit,
			Active:      oldPool.Active,
			Shield:      oldPool.Shield,
		}

		newProviderBz := cdc.MustMarshalLengthPrefixed(&newPool)
		oldStore.Set(iterator.Key(), newProviderBz)
	}
	return nil
}

func migrateProviders(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, types.ProviderKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var provider types.Provider
		err := cdc.UnmarshalLengthPrefixed(oldStoreIter.Value(), &provider)
		if err != nil {
			return err
		}

		shentuProviderAddr, err := common.PrefixToShentu(provider.Address)
		if err != nil {
			return err
		}

		provider.Address = shentuProviderAddr
		bz := cdc.MustMarshalLengthPrefixed(&provider)
		oldStore.Set(oldStoreIter.Key(), bz)
	}
	return nil
}

func migratePurchases(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, types.PurchaseListKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var oldPurchaseList types.PurchaseList
		err := cdc.UnmarshalLengthPrefixed(oldStoreIter.Value(), &oldPurchaseList)
		if err != nil {
			return err
		}

		purchaser, err := common.PrefixToShentu(oldPurchaseList.Purchaser)
		if err != nil {
			return err
		}
		newPurchaseList := types.PurchaseList{
			PoolId:    oldPurchaseList.PoolId,
			Purchaser: purchaser,
			Entries:   oldPurchaseList.Entries,
		}

		newPurchaseBz := cdc.MustMarshalLengthPrefixed(&newPurchaseList)
		oldStore.Set(oldStoreIter.Key(), newPurchaseBz)
	}

	return nil
}

func migrateExpiringPurchaseQueueTimeSlice(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, types.PurchaseQueueKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var oldPoolPurchaserPairs types.PoolPurchaserPairs
		err := cdc.UnmarshalLengthPrefixed(oldStoreIter.Value(), &oldPoolPurchaserPairs)
		if err != nil {
			return err
		}

		purchasePairs := make([]types.PoolPurchaser, 0, len(oldPoolPurchaserPairs.Pairs))
		for _, pair := range oldPoolPurchaserPairs.Pairs {
			purchaser, err := common.PrefixToShentu(pair.Purchaser)
			if err != nil {
				return err
			}

			newPoolPurchaser := types.PoolPurchaser{
				PoolId:    pair.PoolId,
				Purchaser: purchaser,
			}
			purchasePairs = append(purchasePairs, newPoolPurchaser)
		}

		oldPoolPurchaserPairs.Pairs = purchasePairs

		bz := cdc.MustMarshalLengthPrefixed(&oldPoolPurchaserPairs)
		oldStore.Set(oldStoreIter.Key(), bz)
	}
	return nil
}

func migrateReimbursementKey(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, types.ReimbursementKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var reimbursement types.Reimbursement
		err := cdc.UnmarshalLengthPrefixed(oldStoreIter.Value(), &reimbursement)
		if err != nil {
			return err
		}

		beneficiary, err := common.PrefixToShentu(reimbursement.Beneficiary)
		if err != nil {
			return err
		}
		reimbursement.Beneficiary = beneficiary
		bz := cdc.MustMarshalLengthPrefixed(&reimbursement)
		oldStore.Set(oldStoreIter.Key(), bz)
	}

	return nil
}

func migrateWithdraws(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, types.WithdrawQueueKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		oldWithdraws := types.Withdraws{}
		err := cdc.UnmarshalLengthPrefixed(oldStoreIter.Value(), &oldWithdraws)
		if err != nil {
			return err
		}

		newWithdraws := types.Withdraws{}
		for _, withdraw := range oldWithdraws.Withdraws {
			addr, err := common.PrefixToShentu(withdraw.Address)
			if err != nil {
				return err
			}

			newWithdraw := types.Withdraw{
				Address:        addr,
				Amount:         withdraw.Amount,
				CompletionTime: withdraw.CompletionTime,
			}
			newWithdraws.Withdraws = append(newWithdraws.Withdraws, newWithdraw)
		}

		bz := cdc.MustMarshalLengthPrefixed(&newWithdraws)
		oldStore.Set(oldStoreIter.Key(), bz)
	}

	return nil
}

func migrateStakeForShield(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldStore := prefix.NewStore(store, types.StakeForShieldKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var oldShieldStaking types.ShieldStaking
		err := cdc.UnmarshalLengthPrefixed(oldStoreIter.Value(), &oldShieldStaking)
		if err != nil {
			return err
		}

		purchaser, err := common.PrefixToShentu(oldShieldStaking.Purchaser)
		if err != nil {
			return err
		}
		oldShieldStaking.Purchaser = purchaser

		bz := cdc.MustMarshalLengthPrefixed(&oldShieldStaking)
		oldStore.Set(oldStoreIter.Key(), bz)
	}

	return nil
}

func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)

	if err := migratePool(store, cdc); err != nil {
		return err
	}

	if err := migrateProviders(store, cdc); err != nil {
		return err
	}

	if err := migratePurchases(store, cdc); err != nil {
		return err
	}

	if err := migrateExpiringPurchaseQueueTimeSlice(store, cdc); err != nil {
		return err
	}
	if err := migrateReimbursementKey(store, cdc); err != nil {
		return err
	}
	if err := migrateWithdraws(store, cdc); err != nil {
		return err
	}
	if err := migrateStakeForShield(store, cdc); err != nil {
		return err
	}
	return nil
}
