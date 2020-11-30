// Package types adds AccountKeeper and CVMKeeper expected keeper.
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/hyperledger/burrow/crypto"
)

// AccountKeeper defines the account contract that must be fulfilled when creating a x/bank keeper.
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) exported.Account
	NewAccount(ctx sdk.Context, acc exported.Account) exported.Account

	GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account
	GetAllAccounts(ctx sdk.Context) []exported.Account
	SetAccount(ctx sdk.Context, acc exported.Account)

	IterateAccounts(ctx sdk.Context, process func(exported.Account) bool)
}

// CVMKeeper defines the CVM interface that must be fulfilled when wrapping the basekeeper.
type CVMKeeper interface {
	Send(ctx sdk.Context, from, to sdk.AccAddress, value sdk.Coins) error
	GetCode(ctx sdk.Context, addr crypto.Address) ([]byte, error)
}
