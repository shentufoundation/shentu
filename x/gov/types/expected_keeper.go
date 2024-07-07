package types

import (
	"cosmossdk.io/math"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type CertKeeper interface {
	IsCertifier(ctx sdk.Context, addr sdk.AccAddress) bool
	GetAllCertifiers(ctx sdk.Context) (certifiers certtypes.Certifiers)
	GetCertifier(ctx sdk.Context, certifierAddress sdk.AccAddress) (certtypes.Certifier, error)
	HasCertifierAlias(ctx sdk.Context, alias string) bool
	IsCertified(ctx sdk.Context, content string, certType string) bool
	GetCertifiedIdentities(ctx sdk.Context) []sdk.AccAddress
}

type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
	GetRaw(ctx sdk.Context, key []byte) []byte
}

type StakingKeeper interface {
	IterateBondedValidatorsByPower(sdk.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	TotalBondedTokens(sdk.Context) math.Int
	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress, fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))
	BondDenom(sdk.Context) string
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
}
