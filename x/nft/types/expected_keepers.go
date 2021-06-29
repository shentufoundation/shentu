package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	certtypes "github.com/certikfoundation/shentu/x/cert/types"
)

type CertKeeper interface {
	IsCertifier(ctx sdk.Context, address sdk.AccAddress) bool
	GetCertifier(ctx sdk.Context, address sdk.AccAddress) (certtypes.Certifier, error)
	GetAllCertifiers(ctx sdk.Context) certtypes.Certifiers
}
