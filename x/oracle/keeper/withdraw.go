package keeper

import (
	"context"
	"encoding/binary"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

// SetWithdraw sets a withdrawal in store.
func (k Keeper) SetWithdraw(ctx context.Context, withdraw types.Withdraw) error {
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshalLengthPrefixed(&withdraw)
	withdrawAddr := sdk.MustAccAddressFromBech32(withdraw.Address)
	return store.Set(types.WithdrawStoreKey(withdrawAddr, withdraw.DueBlock), bz)
}

// DeleteWithdraw deletes a withdrawal from store.
func (k Keeper) DeleteWithdraw(ctx context.Context, address sdk.AccAddress, startTime int64) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Delete(types.WithdrawStoreKey(address, startTime))
}

// IterateAllWithdraws iterates all withdrawals in store.
func (k Keeper) IterateAllWithdraws(ctx context.Context, callback func(withdraw types.Withdraw) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.WithdrawStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var withdraw types.Withdraw
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &withdraw)
		if callback(withdraw) {
			break
		}
	}
}

// IterateMatureWithdraws iterates all mature (unlocked) withdrawals in store.
func (k Keeper) IterateMatureWithdraws(ctx context.Context, callback func(withdraw types.Withdraw) (stop bool)) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(sdkCtx.BlockHeight()))

	store := k.storeService.OpenKVStore(ctx)
	iterator, _ := store.Iterator(types.WithdrawStoreKeyPrefix,
		storetypes.PrefixEndBytes(append(types.WithdrawStoreKeyPrefix, b...)))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var withdraw types.Withdraw
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &withdraw)
		if callback(withdraw) {
			break
		}
	}
}

// CreateWithdraw creates a withdrawal.
func (k Keeper) CreateWithdraw(ctx context.Context, address sdk.AccAddress, amount sdk.Coins) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetLockedPoolParams(ctx)
	dueBlock := sdkCtx.BlockHeight() + params.LockedInBlocks
	withdraw := types.NewWithdraw(address, amount, dueBlock)
	return k.SetWithdraw(ctx, withdraw)
}

// GetAllWithdraws gets all withdrawals from store.
func (k Keeper) GetAllWithdraws(ctx context.Context) types.Withdraws {
	var withdraws types.Withdraws
	k.IterateAllWithdraws(ctx, func(withdraw types.Withdraw) bool {
		withdraws = append(withdraws, withdraw)
		return false
	})
	return withdraws
}

// GetAllWithdrawsForExport gets all withdrawals from store and adjusts DueBlock value for import-export.
func (k Keeper) GetAllWithdrawsForExport(ctx context.Context) types.Withdraws {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var withdraws types.Withdraws
	k.IterateAllWithdraws(ctx, func(withdraw types.Withdraw) bool {
		withdraw.DueBlock -= sdkCtx.BlockHeight()
		withdraws = append(withdraws, withdraw)
		return false
	})
	return withdraws
}
