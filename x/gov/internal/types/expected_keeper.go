package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	"github.com/certikfoundation/shentu/x/cert"
)

type CertKeeper interface {
	IsCertifier(ctx sdk.Context, addr sdk.AccAddress) bool
	GetAllCertifiers(ctx sdk.Context) (certifiers cert.Certifiers)
	GetCertifier(ctx sdk.Context, certifierAddress sdk.AccAddress) (cert.Certifier, error)
	HasCertifierAlias(ctx sdk.Context, alias string) bool
}

type UpgradeKeeper interface {
	ValidatePlan(ctx sdk.Context, plan upgrade.Plan) error
}

type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
	WithKeyTable(table subspace.KeyTable) subspace.Subspace
}
