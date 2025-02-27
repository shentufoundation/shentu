package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (k Keeper) SetDeposit(ctx context.Context, deposit types.Deposit) error {
	depositor, err := k.authKeeper.AddressCodec().StringToBytes(deposit.Depositor)
	if err != nil {
		return err
	}
	return k.Deposits.Set(ctx, collections.Join(deposit.ProofId, sdk.AccAddress(depositor)), deposit)
}

// validateMinDeposit validates if deposit is greater than or equal to the minimum
// required at the time of proof submission. Returns nil on success, error otherwise.
func (k Keeper) validateMinDeposit(ctx context.Context, params types.Params, initialDeposit sdk.Coins) error {
	// Check if the initial deposit is valid and has no negative amount
	if !initialDeposit.IsValid() || initialDeposit.IsAnyNegative() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, initialDeposit.String())
	}

	// Check if the initial deposit meets the minimum required grant amount
	if !initialDeposit.IsAllGTE(params.MinDeposit) {
		return errors.Wrapf(types.ErrMinDepositTooSmall, "was (%s), need (%s)", initialDeposit, params.MinGrant)
	}
	return nil
}
