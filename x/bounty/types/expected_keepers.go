package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the expected supply keeper (noalias)
type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// ParamSubspace defines the expected Subspace interface for parameters (noalias)
type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
	HasKeyTable() bool
}

type CertKeeper interface {
	IsCertifier(ctx sdk.Context, addr sdk.AccAddress) bool
}
