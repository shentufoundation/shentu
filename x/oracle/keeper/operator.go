package keeper

import (
	"context"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

// SetOperator sets an operator to store.
func (k Keeper) SetOperator(ctx context.Context, operator types.Operator) error {
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshalLengthPrefixed(&operator)
	addr := sdk.MustAccAddressFromBech32(operator.Address)
	return store.Set(types.OperatorStoreKey(addr), bz)
}

// GetOperator gets an operator from store.
func (k Keeper) GetOperator(ctx context.Context, address sdk.AccAddress) (types.Operator, error) {
	store := k.storeService.OpenKVStore(ctx)
	opBz, err := store.Get(types.OperatorStoreKey(address))
	if err != nil {
		return types.Operator{}, err
	}
	if len(opBz) == 0 {
		return types.Operator{}, types.ErrNoOperatorFound
	}
	var operator types.Operator
	k.cdc.MustUnmarshalLengthPrefixed(opBz, &operator)
	return operator, nil
}

// IsOperator determines if an address belongs to an operator.
func (k Keeper) IsOperator(ctx context.Context, address sdk.AccAddress) (bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	return store.Has(types.OperatorStoreKey(address))
}

// DeleteOperator deletes an operator from store.
func (k Keeper) DeleteOperator(ctx context.Context, address sdk.AccAddress) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Delete(types.OperatorStoreKey(address))
}

// IterateAllOperators iterates all operators.
func (k Keeper) IterateAllOperators(ctx context.Context, callback func(operator types.Operator) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.OperatorStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var operator types.Operator
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &operator)
		if callback(operator) {
			break
		}
	}
}

// IsBelowMinCollateral determines if collateral is below the minimum requirement.
func (k Keeper) IsBelowMinCollateral(ctx context.Context, currentCollateral sdk.Coins) bool {
	params := k.GetLockedPoolParams(ctx)
	denom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return false
	}
	return currentCollateral.AmountOf(denom).LT(math.NewInt(params.MinimumCollateral))
}

// CreateOperator creates an operator and deposits collateral.
func (k Keeper) CreateOperator(ctx context.Context, address sdk.AccAddress, collateral sdk.Coins, proposer sdk.AccAddress, name string) error {
	_, err := k.IsOperator(ctx, address)
	if err != nil {
		return err
	}
	if k.IsBelowMinCollateral(ctx, collateral) {
		return types.ErrNoEnoughCollateral
	}
	isCertifier, err := k.CertKeeper.IsCertifier(ctx, proposer)
	if err != nil {
		return err
	}
	if !isCertifier {
		return types.ErrUnqualifiedRemover
	}
	operator := types.NewOperator(address, proposer, collateral, nil, name)
	if err = k.SetOperator(ctx, operator); err != nil {
		return err
	}
	if err := k.AddTotalCollateral(ctx, collateral); err != nil {
		return err
	}
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, collateral); err != nil {
		return err
	}
	return nil
}

// RemoveOperator removes an operator, creates an withdrawal for collateral and gives back rewards immediately.
func (k Keeper) RemoveOperator(ctx context.Context, operatorAddress, proposerAddress string) error {
	// Ensure that the sender of the tx is either the operator to be removed itself or a certifier.
	proposerAddr, err := sdk.AccAddressFromBech32(proposerAddress)
	if err != nil {
		return err
	}
	isCertifier, err := k.CertKeeper.IsCertifier(ctx, proposerAddr)
	if err != nil {
		return err
	}
	if operatorAddress != proposerAddress && !isCertifier {
		return types.ErrUnqualifiedRemover
	}

	operatorAddr, err := sdk.AccAddressFromBech32(operatorAddress)
	if err != nil {
		return err
	}
	if _, err = k.IsOperator(ctx, operatorAddr); err != nil {
		return err
	}
	operator, err := k.GetOperator(ctx, operatorAddr)
	if err != nil {
		return nil
	}
	if err := k.ReduceTotalCollateral(ctx, operator.Collateral); err != nil {
		return err
	}
	if err := k.CreateWithdraw(ctx, operatorAddr, operator.Collateral); err != nil {
		return err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, operatorAddr,
		operator.AccumulatedRewards); err != nil {
		return err
	}
	return k.DeleteOperator(ctx, operatorAddr)
}

// GetAllOperators gets all operators.
func (k Keeper) GetAllOperators(ctx context.Context) types.Operators {
	operators := types.Operators{}
	k.IterateAllOperators(ctx, func(operator types.Operator) bool {
		operators = append(operators, operator)
		return false
	})
	return operators
}

// AddCollateral increases an operator's collateral, effective immediately.
func (k Keeper) AddCollateral(ctx context.Context, address sdk.AccAddress, increment sdk.Coins) error {
	if _, err := k.IsOperator(ctx, address); err != nil {
		return err
	}
	operator, err := k.GetOperator(ctx, address)
	if err != nil {
		return err
	}
	operator.Collateral = operator.Collateral.Add(increment...)
	if err = k.SetOperator(ctx, operator); err != nil {
		return err
	}
	if err := k.AddTotalCollateral(ctx, increment); err != nil {
		return err
	}
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, increment); err != nil {
		return err
	}
	return nil
}

// ReduceCollateral reduces an operator's collateral and creates a withdrawal for it.
func (k Keeper) ReduceCollateral(ctx context.Context, address sdk.AccAddress, decrement sdk.Coins) error {
	if _, err := k.IsOperator(ctx, address); err != nil {
		return err
	}
	operator, err := k.GetOperator(ctx, address)
	if err != nil {
		return err
	}
	operator.Collateral = operator.Collateral.Sub(decrement...)
	if k.IsBelowMinCollateral(ctx, operator.Collateral) {
		return types.ErrNoEnoughCollateral
	}
	if err = k.SetOperator(ctx, operator); err != nil {
		return err
	}
	if err := k.ReduceTotalCollateral(ctx, decrement); err != nil {
		return err
	}
	operatorAddr := sdk.MustAccAddressFromBech32(operator.Address)
	if err := k.CreateWithdraw(ctx, operatorAddr, decrement); err != nil {
		return err
	}
	return nil
}

// AddReward increases an operators accumulated rewards.
func (k Keeper) AddReward(ctx context.Context, address sdk.AccAddress, increment sdk.Coins) error {
	if _, err := k.IsOperator(ctx, address); err != nil {
		return err
	}
	operator, err := k.GetOperator(ctx, address)
	if err != nil {
		return err
	}
	operator.AccumulatedRewards = operator.AccumulatedRewards.Add(increment...)
	if err = k.SetOperator(ctx, operator); err != nil {
		return err
	}
	return nil
}

// WithdrawAllReward gives back all rewards of an operator.
func (k Keeper) WithdrawAllReward(ctx context.Context, address sdk.AccAddress) (sdk.Coins, error) {
	if _, err := k.IsOperator(ctx, address); err != nil {
		return nil, err
	}
	operator, err := k.GetOperator(ctx, address)
	if err != nil {
		return nil, err
	}
	reward := operator.AccumulatedRewards
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, reward); err != nil {
		return nil, err
	}
	operator.AccumulatedRewards = nil
	if err = k.SetOperator(ctx, operator); err != nil {
		return nil, err
	}
	return reward, nil
}

// GetCollateralAmount gets an operator's collateral.
func (k Keeper) GetCollateralAmount(ctx context.Context, address sdk.AccAddress) (math.Int, error) {
	operator, err := k.GetOperator(ctx, address)
	if err != nil {
		return math.NewInt(0), err
	}
	return operator.Collateral[0].Amount, nil
}
