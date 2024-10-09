package types

import (
	"context"

	addresscodec "cosmossdk.io/core/address"
	"cosmossdk.io/math"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	AddressCodec() addresscodec.Codec

	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI

	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx context.Context, name string) sdk.ModuleAccountI

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(context.Context, sdk.ModuleAccountI)
}

type CertKeeper interface {
	IsCertifier(ctx context.Context, addr sdk.AccAddress) (bool, error)
	//GetAllCertifiers(ctx context.Context) (certifiers certtypes.Certifiers, err error)
	GetCertifier(ctx context.Context, certifierAddress sdk.AccAddress) (certtypes.Certifier, error)
	HasCertifierAlias(ctx context.Context, alias string) (bool, error)
	//IsCertified(ctx context.Context, content string, certType string) bool
	//GetCertifiedIdentities(ctx context.Context) []sdk.AccAddress
}

// StakingKeeper expected staking keeper (Validator and Delegator sets) (noalias)
type StakingKeeper interface {
	ValidatorAddressCodec() addresscodec.Codec
	// iterate through bonded validators by operator address, execute func for each validator
	IterateBondedValidatorsByPower(
		context.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool),
	) error

	TotalBondedTokens(context.Context) (math.Int, error) // total bonded tokens within the validator set
	IterateDelegations(
		ctx context.Context, delegator sdk.AccAddress,
		fn func(index int64, delegation stakingtypes.DelegationI) (stop bool),
	) error
}

// DistributionKeeper defines the expected distribution keeper.
type DistributionKeeper interface {
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
}
