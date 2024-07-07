package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// ParamSubspace defines the expected Subspace interface for parameters (noalias)
type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
}

// AccountKeeper defines the expected account keeper.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
	IterateAccounts(ctx sdk.Context, process func(authtypes.AccountI) (stop bool))
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
}

// BankKeeper defines the expected bank keeper.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	DelegateCoins(ctx sdk.Context, fromAdd, toAddr sdk.AccAddress, amt sdk.Coins) error
	UndelegateCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error

	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
}
