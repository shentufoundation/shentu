package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

type Keeper struct {
	upgrade.Keeper
}

// NewKeeper constructs an upgrade Keeper.
func NewKeeper(skipUpgradeHeights map[int64]bool, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		upgrade.NewKeeper(skipUpgradeHeights, storeKey, cdc),
	}
}

// ValidatePlan validates an upgrade plan.
func (k Keeper) ValidatePlan(ctx sdk.Context, plan upgrade.Plan) error {
	if err := plan.ValidateBasic(); err != nil {
		return err
	}

	if !plan.Time.IsZero() {
		if !plan.Time.After(ctx.BlockHeader().Time) {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "upgrade cannot be scheduled in the past")
		}
	} else if plan.Height <= ctx.BlockHeight() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "upgrade cannot be scheduled in the past")
	}

	return nil
}
