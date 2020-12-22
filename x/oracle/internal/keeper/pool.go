package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

// CollectBounty collects task bounty from the operator.
func (k Keeper) CollectBounty(ctx sdk.Context, value sdk.Coins, creator sdk.AccAddress) error {
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, value); err != nil {
		return err
	}
	return nil
}

// SetTotalCollateral sets total collateral to store.
func (k Keeper) SetTotalCollateral(ctx sdk.Context, collateral sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(&types.CoinsProto{Coins: collateral})
	store.Set(types.TotalCollateralKey(), bz)
}

// GetTotalCollateral gets total collateral from store.
func (k Keeper) GetTotalCollateral(ctx sdk.Context) (sdk.Coins, error) {
	store := ctx.KVStore(k.storeKey)
	opBz := store.Get(types.TotalCollateralKey())
	if opBz != nil {
		var coinsProto types.CoinsProto
		k.cdc.MustUnmarshalBinaryLengthPrefixed(opBz, &coinsProto)
		return coinsProto.Coins, nil
	}
	return sdk.Coins{}, types.ErrNoTotalCollateralFound
}

// AddTotalCollateral increases total collateral.
func (k Keeper) AddTotalCollateral(ctx sdk.Context, increment sdk.Coins) error {
	currentCollateral, err := k.GetTotalCollateral(ctx)
	if err != nil {
		return err
	}
	currentCollateral = currentCollateral.Add(increment...)
	k.SetTotalCollateral(ctx, currentCollateral)
	return nil
}

// ReduceTotalCollateral reduces total collateral.
func (k Keeper) ReduceTotalCollateral(ctx sdk.Context, decrement sdk.Coins) error {
	currentCollateral, err := k.GetTotalCollateral(ctx)
	if err != nil {
		return err
	}
	currentCollateral = currentCollateral.Sub(decrement)
	if currentCollateral.IsAnyNegative() {
		return types.ErrNoEnoughTotalCollateral
	}
	k.SetTotalCollateral(ctx, currentCollateral)
	return nil
}

// FundCommunityPool transfers money from module account to community pool.
func (k Keeper) FundCommunityPool(ctx sdk.Context, amount sdk.Coins) error {
	macc := k.supplyKeeper.GetModuleAddress(types.ModuleName)
	return k.distrKeeper.FundCommunityPool(ctx, amount, macc)
}

// FinalizeMatureWithdraws finishes mature (unlocked) withdrawals and removes them.
func (k Keeper) FinalizeMatureWithdraws(ctx sdk.Context) {
	k.IterateMatureWithdraws(ctx, func(withdraw types.Withdraw) bool {
		withdrawAddr, err := sdk.AccAddressFromBech32(withdraw.Address)
		if err != nil {
			panic(err)
		}
		if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, withdraw.Amount); err != nil {
			panic(err)
		}
		if err := k.DeleteWithdraw(ctx, withdrawAddr, withdraw.DueBlock); err != nil {
			panic(err)
		}
		return false
	})
}
