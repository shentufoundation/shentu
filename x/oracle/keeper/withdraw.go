package keeper

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/types"
)

// SetWithdraw sets a withdrawal in store.
func (k Keeper) SetWithdraw(ctx sdk.Context, withdraw types.Withdraw) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(&withdraw)
	withdrawAddr, err := sdk.AccAddressFromBech32(withdraw.Address)
	if err != nil {
		panic(err)
	}
	store.Set(types.WithdrawStoreKey(withdrawAddr, withdraw.DueBlock), bz)
}

// DeleteWithdraw deletes a withdrawal from store.
func (k Keeper) DeleteWithdraw(ctx sdk.Context, address sdk.AccAddress, startTime int64) error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.WithdrawStoreKey(address, startTime))
	return nil
}

// IterateAllWithdraws iterates all withdrawals in store.
func (k Keeper) IterateAllWithdraws(ctx sdk.Context, callback func(withdraw types.Withdraw) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.WithdrawStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var withdraw types.Withdraw
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &withdraw)
		if callback(withdraw) {
			break
		}
	}
}

// IterateMatureWithdraws iterates all mature (unlocked) withdrawals in store.
func (k Keeper) IterateMatureWithdraws(ctx sdk.Context, callback func(withdraw types.Withdraw) (stop bool)) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(ctx.BlockHeight()))

	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.WithdrawStoreKeyPrefix,
		sdk.PrefixEndBytes(append(types.WithdrawStoreKeyPrefix, b...)))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var withdraw types.Withdraw
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &withdraw)
		if callback(withdraw) {
			break
		}
	}
}

// CreateWithdraw creates a withdrawal.
func (k Keeper) CreateWithdraw(ctx sdk.Context, address sdk.AccAddress, amount sdk.Coins) error {
	params := k.GetLockedPoolParams(ctx)
	dueBlock := ctx.BlockHeight() + params.LockedInBlocks
	withdraw := types.NewWithdraw(address, amount, dueBlock)
	k.SetWithdraw(ctx, withdraw)
	return nil
}

// GetAllWithdraws gets all withdrawals from store.
func (k Keeper) GetAllWithdraws(ctx sdk.Context) types.Withdraws {
	var withdraws types.Withdraws
	k.IterateAllWithdraws(ctx, func(withdraw types.Withdraw) bool {
		withdraws = append(withdraws, withdraw)
		return false
	})
	return withdraws
}

// GetAllWithdrawsForExport gets all withdrawals from store and adjusts DueBlock value for import-export.
func (k Keeper) GetAllWithdrawsForExport(ctx sdk.Context) types.Withdraws {
	var withdraws types.Withdraws
	k.IterateAllWithdraws(ctx, func(withdraw types.Withdraw) bool {
		withdraw.DueBlock = withdraw.DueBlock - ctx.BlockHeight()
		withdraws = append(withdraws, withdraw)
		return false
	})
	return withdraws
}
