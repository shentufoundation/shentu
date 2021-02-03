package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/types"
)

// SetOperator sets an operator to store.
func (k Keeper) SetOperator(ctx sdk.Context, operator types.Operator) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(operator)
	store.Set(types.OperatorStoreKey(operator.Address), bz)
}

// GetOperator gets an operator from store.
func (k Keeper) GetOperator(ctx sdk.Context, address sdk.AccAddress) (types.Operator, error) {
	store := ctx.KVStore(k.storeKey)
	opBz := store.Get(types.OperatorStoreKey(address))
	if opBz != nil {
		var operator types.Operator
		k.cdc.MustUnmarshalBinaryLengthPrefixed(opBz, &operator)
		return operator, nil
	}
	return types.Operator{}, types.ErrNoOperatorFound
}

// IsOperator determines if an address belongs to an operator.
func (k Keeper) IsOperator(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.OperatorStoreKey(address))
}

// DeleteOperators deletes an operator from store.
func (k Keeper) DeleteOperator(ctx sdk.Context, address sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.OperatorStoreKey(address))
	return nil
}

// IterateAllOperators iterates all operators.
func (k Keeper) IterateAllOperators(ctx sdk.Context, callback func(operator types.Operator) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OperatorStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var operator types.Operator
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &operator)
		if callback(operator) {
			break
		}
	}
}

// IsBelowMinCollateral determines if collateral is below the minimum requirement.
func (k Keeper) IsBelowMinCollateral(ctx sdk.Context, currentCollateral sdk.Coins) bool {
	params := k.GetLockedPoolParams(ctx)
	return currentCollateral.AmountOf(k.stakingKeeper.BondDenom(ctx)).LT(sdk.NewInt(params.MinimumCollateral))
}

// CreateOperator creates an operator and deposits collateral.
func (k Keeper) CreateOperator(ctx sdk.Context, address sdk.AccAddress, collateral sdk.Coins, proposer sdk.AccAddress, name string) error {
	if k.IsOperator(ctx, address) {
		return types.ErrOperatorAlreadyExists
	}
	if k.IsBelowMinCollateral(ctx, collateral) {
		return types.ErrNoEnoughCollateral
	}
	operator := types.NewOperator(address, proposer, collateral, nil, name)
	k.SetOperator(ctx, operator)
	if err := k.AddTotalCollateral(ctx, collateral); err != nil {
		return err
	}
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, collateral); err != nil {
		return err
	}
	return nil
}

// RemoveOperator removes an operator, creates an withdrawal for collateral and gives back rewards immediately.
func (k Keeper) RemoveOperator(ctx sdk.Context, address sdk.AccAddress) error {
	if !k.IsOperator(ctx, address) {
		return types.ErrNoOperatorFound
	}
	operator, err := k.GetOperator(ctx, address)
	if err != nil {
		return nil
	}
	if err := k.ReduceTotalCollateral(ctx, operator.Collateral); err != nil {
		return err
	}
	if err := k.CreateWithdraw(ctx, operator.Address, operator.Collateral); err != nil {
		return err
	}
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address,
		operator.AccumulatedRewards); err != nil {
		return err
	}
	return k.DeleteOperator(ctx, address)
}

// GetAllOperators gets all operators.
func (k Keeper) GetAllOperators(ctx sdk.Context) types.Operators {
	operators := types.Operators{}
	k.IterateAllOperators(ctx, func(operator types.Operator) bool {
		operators = append(operators, operator)
		return false
	})
	return operators
}

// AddCollateral increases an operator's collateral, effective immediately.
func (k Keeper) AddCollateral(ctx sdk.Context, address sdk.AccAddress, increment sdk.Coins) error {
	if !k.IsOperator(ctx, address) {
		return types.ErrNoOperatorFound
	}
	operator, err := k.GetOperator(ctx, address)
	if err != nil {
		return err
	}
	operator.Collateral = operator.Collateral.Add(increment...)
	k.SetOperator(ctx, operator)
	if err := k.AddTotalCollateral(ctx, increment); err != nil {
		return err
	}
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, increment); err != nil {
		return err
	}
	return nil
}

// ReduceCollateral reduces an operator's collateral and creates a withdrawal for it.
func (k Keeper) ReduceCollateral(ctx sdk.Context, address sdk.AccAddress, decrement sdk.Coins) error {
	if !k.IsOperator(ctx, address) {
		return types.ErrNoOperatorFound
	}
	operator, err := k.GetOperator(ctx, address)
	if err != nil {
		return err
	}
	operator.Collateral = operator.Collateral.Sub(decrement)
	if k.IsBelowMinCollateral(ctx, operator.Collateral) {
		return types.ErrNoEnoughCollateral
	}
	k.SetOperator(ctx, operator)
	if err := k.ReduceTotalCollateral(ctx, decrement); err != nil {
		return err
	}
	if err := k.CreateWithdraw(ctx, operator.Address, decrement); err != nil {
		return err
	}
	return nil
}

// AddReward increases an operators accumulated rewards.
func (k Keeper) AddReward(ctx sdk.Context, address sdk.AccAddress, increment sdk.Coins) error {
	if !k.IsOperator(ctx, address) {
		return types.ErrNoOperatorFound
	}
	operator, err := k.GetOperator(ctx, address)
	if err != nil {
		return err
	}
	operator.AccumulatedRewards = operator.AccumulatedRewards.Add(increment...)
	k.SetOperator(ctx, operator)
	return nil
}

// WithdrawAllReward gives back all rewards of an operator.
func (k Keeper) WithdrawAllReward(ctx sdk.Context, address sdk.AccAddress) (sdk.Coins, error) {
	if !k.IsOperator(ctx, address) {
		return nil, types.ErrNoOperatorFound
	}
	operator, err := k.GetOperator(ctx, address)
	if err != nil {
		return nil, err
	}
	reward := operator.AccumulatedRewards
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, reward); err != nil {
		return nil, err
	}
	operator.AccumulatedRewards = nil
	k.SetOperator(ctx, operator)
	return reward, nil
}

// GetCollateralAmount gets an operator's collateral.
func (k Keeper) GetCollateralAmount(ctx sdk.Context, address sdk.AccAddress) (sdk.Int, error) {
	operator, err := k.GetOperator(ctx, address)
	if err != nil {
		return sdk.NewInt(0), err
	}
	return operator.Collateral[0].Amount, nil
}
