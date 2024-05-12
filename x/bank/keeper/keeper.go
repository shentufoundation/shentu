// Package keeper implements custom bank keeper through CVM.
package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/hyperledger/burrow/crypto"

	"github.com/shentufoundation/shentu/v2/x/bank/types"
)

// Keeper is a wrapper of the basekeeper with CVM keeper.
type Keeper struct {
	bankKeeper.BaseKeeper
	cvmk types.CVMKeeper
	ak   types.AccountKeeper
}

// NewKeeper returns a new Keeper.
func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, ak types.AccountKeeper, cvmk types.CVMKeeper, paramSpace paramsTypes.Subspace,
	blockedAddrs map[string]bool) Keeper {
	bk := bankKeeper.NewBaseKeeper(cdc, storeKey, ak, paramSpace, blockedAddrs)
	return Keeper{
		BaseKeeper: bk,
		cvmk:       cvmk,
		ak:         ak,
	}
}

// GetCode retrieves VM code from an account.
func (k Keeper) GetCode(ctx sdk.Context, addr sdk.AccAddress) ([]byte, error) {
	vmAddress := crypto.MustAddressFromBytes(addr)
	return k.cvmk.GetCode(ctx, vmAddress)
}

// SendCoins checks if there is code in the receiver account, and wires the send through CVM if it does.
func (k Keeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	code, err := k.GetCode(ctx, toAddr)
	if err != nil {
		return err
	}
	if len(code) > 0 {
		return k.cvmk.Send(ctx, fromAddr, toAddr, amt)
	}
	return k.BaseKeeper.SendCoins(ctx, fromAddr, toAddr, amt)
}

// InputOutputCoins handles multisend logic.
func (k Keeper) InputOutputCoins(ctx sdk.Context, inputs []bankTypes.Input, outputs []bankTypes.Output) error {
	for _, out := range outputs {
		outAddr, err := sdk.AccAddressFromBech32(out.Address)
		if err != nil {
			return err
		}
		code, err := k.GetCode(ctx, outAddr)
		if err != nil {
			return err
		}
		if len(code) > 0 {
			return types.ErrCodeExists
		}
	}
	return k.BaseKeeper.InputOutputCoins(ctx, inputs, outputs)
}
