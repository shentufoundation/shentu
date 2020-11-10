package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingexported "github.com/cosmos/cosmos-sdk/x/staking/exported"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	"github.com/certikfoundation/shentu/x/cert"
	"github.com/certikfoundation/shentu/x/shield"
)

type CertKeeper interface {
	IsCertifier(ctx sdk.Context, addr sdk.AccAddress) bool
	GetAllCertifiers(ctx sdk.Context) (certifiers cert.Certifiers)
	GetCertifier(ctx sdk.Context, certifierAddress sdk.AccAddress) (cert.Certifier, error)
	HasCertifierAlias(ctx sdk.Context, alias string) bool
	IsCertified(ctx sdk.Context, requestContentType string, content string, certType string) bool
	GetCertifiedIdentities(ctx sdk.Context) []sdk.AccAddress
}

type UpgradeKeeper interface {
	ValidatePlan(ctx sdk.Context, plan upgrade.Plan) error
}

type ShieldKeeper interface {
	GetPurchaseList(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress) (shield.PurchaseList, bool)
	GetClaimProposalParams(ctx sdk.Context) shield.ClaimProposalParams
	SecureCollaterals(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, purchaseID uint64, loss sdk.Coins, lockPeriod time.Duration) error
	RestoreShield(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, id uint64, loss sdk.Coins) error
	ClaimEnd(ctx sdk.Context, id, poolID uint64, loss sdk.Coins)
}

type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
	WithKeyTable(table subspace.KeyTable) subspace.Subspace
}

type StakingKeeper interface {
	IterateBondedValidatorsByPower(sdk.Context, func(index int64, validator stakingexported.ValidatorI) (stop bool))
	TotalBondedTokens(sdk.Context) sdk.Int
	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress, fn func(index int64, delegation stakingexported.DelegationI) (stop bool))
	BondDenom(sdk.Context) string
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator staking.Validator, found bool)
}
