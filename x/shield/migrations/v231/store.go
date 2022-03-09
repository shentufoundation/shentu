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

func resolvePurchases(store sdk.KVStore) error {
	oldStore := prefix.NewStore(store, types.PurchaseKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		oldStore.Delete(oldStoreIter.Key())
	}

	queueStore := prefix.NewStore(store, types.PurchaseQueueKey)

	queueStoreIter := queueStore.Iterator(nil, nil)
	defer queueStoreIter.Close()

	for ; queueStoreIter.Valid(); queueStoreIter.Next() {
		oldStore.Delete(queueStoreIter.Key())
	}

	return nil
}

func deleteUnusedStores(store sdk.KVStore) error {
	store.Delete(types.GetNextPurchaseIDKey())
	store.Delete(types.GetLastUpdateTimeKey())
	return nil
}

func resolveReimbursements(store sdk.KVStore, cdc codec.BinaryCodec, bondDenom string) error {
	var reserve v1beta1.Reserve
	reserveBz := store.Get(types.ReserveKey)
	cdc.MustUnmarshal(reserveBz, &reserve)

	oldStore := prefix.NewStore(store, types.ReimbursementKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()
	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var reimbursement v1alpha1.Reimbursement
		cdc.MustUnmarshal(oldStoreIter.Value(), &reimbursement)

		reserve.Amount = reserve.Amount.Add(reimbursement.Amount.AmountOf(bondDenom))
		oldStore.Delete(oldStoreIter.Key())
	}
	reserveBz = cdc.MustMarshal(&reserve)
	store.Set(types.ReserveKey, reserveBz)
	return nil
}

func migrateParams(ctx sdk.Context, ps types.ParamSubspace) error {
	var poolParamsV1 v1alpha1.PoolParams
	ps.Get(ctx, v1alpha1.ParamStoreKeyPoolParams, &poolParamsV1)
	// fmt.Println(poolParamsV1.String())
	poolParamsV2 := v1beta1.PoolParams{
		ProtectionPeriod:  poolParamsV1.ProtectionPeriod,
		ShieldFeesRate:    poolParamsV1.ShieldFeesRate,
		WithdrawPeriod:    poolParamsV1.WithdrawPeriod,
		PoolShieldLimit:   poolParamsV1.PoolShieldLimit,
		MinShieldPurchase: poolParamsV1.MinShieldPurchase,
		CooldownPeriod:    v1beta1.DefaultCooldownPeriod,
		WithdrawFeesRate:  v1beta1.DefaultWithdrawFeesRate,
	}
	ps.Set(ctx, v1beta1.ParamStoreKeyPoolParams, &poolParamsV2)

	// Claim proposal params didn't change, do nothing.
	// var claimProposal v1beta1.ClaimProposalParams
	// ps.Get(ctx, v1beta1.ParamStoreKeyClaimProposalParams, &claimProposal)
	// fmt.Println(claimProposal.String())

	// Staking shield rate didn't change, do nothing.
	// var stakingShieldRate sdk.Dec
	// ps.Get(ctx, v1beta1.ParamStoreKeyStakingShieldRate, &stakingShieldRate)
	// fmt.Println(stakingShieldRate.String())

	blockRewardParams := v1beta1.DefaultBlockRewardParams()
	ps.Set(ctx, v1beta1.ParamStoreKeyBlockRewardParams, &blockRewardParams)

	return nil
}

func initReserve(store sdk.KVStore, cdc codec.BinaryCodec) error {
	reserve := v1beta1.InitialReserve()
	bz := cdc.MustMarshalLengthPrefixed(&reserve)
	store.Set(types.GetReserveKey(), bz)
	return nil
}

func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec, paramSpace types.ParamSubspace, queryServer grpc.Server) error {
	bondDenom, err := getBondDenom(ctx, queryServer)
	if err != nil {
		return err
	}
	store := ctx.KVStore(storeKey)
	err = migratePools(store, cdc)
	if err != nil {
		return err
	}

	err = resolvePurchases(store)
	if err != nil {
		return err
	}

	err = resolveReimbursements(store, cdc, bondDenom)
	if err != nil {
		return err
	}

	err = migrateParams(ctx, paramSpace)
	if err != nil {
		return err
	}

	err = initReserve(store, cdc)
	if err != nil {
		return err
	}

	return deleteUnusedStores(store)
}
