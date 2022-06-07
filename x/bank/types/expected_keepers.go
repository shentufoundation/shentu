// Package types adds AccountKeeper and CVMKeeper expected keeper.
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/hyperledger/burrow/crypto"
)

// AccountKeeper defines the account contract that must be fulfilled when creating a x/bank keeper.
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	NewAccount(ctx sdk.Context, acc types.AccountI) types.AccountI

	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetAllAccounts(ctx sdk.Context) []types.AccountI
	SetAccount(ctx sdk.Context, acc types.AccountI)
	HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool
	IterateAccounts(ctx sdk.Context, process func(types.AccountI) bool)

	ValidatePermissions(macc types.ModuleAccountI) error

	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAddressAndPermissions(moduleName string) (addr sdk.AccAddress, permissions []string)
	GetModuleAccountAndPermissions(ctx sdk.Context, moduleName string) (types.ModuleAccountI, []string)
	GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI
	SetModuleAccount(ctx sdk.Context, macc types.ModuleAccountI)
}

// CVMKeeper defines the CVM interface that must be fulfilled when wrapping the basekeeper.
type CVMKeeper interface {
	Send(ctx sdk.Context, from, to sdk.AccAddress, value sdk.Coins) error
	GetCode(ctx sdk.Context, addr crypto.Address) ([]byte, error)
}
