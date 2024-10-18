package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

// CollectBounty collects task bounty from the operator.
func (k Keeper) CollectBounty(ctx context.Context, value sdk.Coins, creator sdk.AccAddress) error {
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, value); err != nil {
		return err
	}
	return nil
}

// SetTotalCollateral sets total collateral to store.
func (k Keeper) SetTotalCollateral(ctx context.Context, collateral sdk.Coins) error {
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshalLengthPrefixed(&types.CoinsProto{Coins: collateral})
	return store.Set(types.TotalCollateralKey(), bz)
}

// GetTotalCollateral gets total collateral from store.
func (k Keeper) GetTotalCollateral(ctx context.Context) (sdk.Coins, error) {
	store := k.storeService.OpenKVStore(ctx)
	opBz, err := store.Get(types.TotalCollateralKey())
	if err != nil {
		return nil, err
	}
	if opBz != nil {
		var coinsProto types.CoinsProto
		k.cdc.MustUnmarshalLengthPrefixed(opBz, &coinsProto)
		return coinsProto.Coins, nil
	}
	return nil, types.ErrNoTotalCollateralFound
}

// AddTotalCollateral increases total collateral.
func (k Keeper) AddTotalCollateral(ctx context.Context, increment sdk.Coins) error {
	currentCollateral, err := k.GetTotalCollateral(ctx)
	if err != nil {
		return err
	}
	currentCollateral = currentCollateral.Add(increment...)
	return k.SetTotalCollateral(ctx, currentCollateral)
}

// ReduceTotalCollateral reduces total collateral.
func (k Keeper) ReduceTotalCollateral(ctx context.Context, decrement sdk.Coins) error {
	currentCollateral, err := k.GetTotalCollateral(ctx)
	if err != nil {
		return err
	}
	currentCollateral = currentCollateral.Sub(decrement...)
	k.SetTotalCollateral(ctx, currentCollateral)
	return nil
}

// FundCommunityPool transfers money from module account to community pool.
func (k Keeper) FundCommunityPool(ctx context.Context, amount sdk.Coins) error {
	macc := k.accountKeeper.GetModuleAddress(types.ModuleName)
	return k.distrKeeper.FundCommunityPool(ctx, amount, macc)
}

// FinalizeMatureWithdraws finishes mature (unlocked) withdrawals and removes them.
func (k Keeper) FinalizeMatureWithdraws(ctx context.Context) error {
	k.IterateMatureWithdraws(ctx, func(withdraw types.Withdraw) bool {
		withdrawAddr := sdk.MustAccAddressFromBech32(withdraw.Address)
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, withdraw.Amount); err != nil {
			return false
		}
		if err := k.DeleteWithdraw(ctx, withdrawAddr, withdraw.DueBlock); err != nil {
			panic(err)
		}
		return false
	})

	return nil
}
