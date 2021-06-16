package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	certtypes "github.com/certikfoundation/shentu/x/cert/types"
)

type CertKeeper interface {
	GetCertifier(ctx sdk.Context, address sdk.AccAddress) (certtypes.Certifier, error)
}
