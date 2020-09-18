package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CertKeeper interface {
	IsCertifier(ctx sdk.Context, addr sdk.AccAddress) bool
}
