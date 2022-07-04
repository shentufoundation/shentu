package v1beta1

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	certtypes "github.com/certikfoundation/shentu/v2/x/cert/types"
	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
)

type CertKeeper interface {
	IsCertifier(ctx sdk.Context, addr sdk.AccAddress) bool
	GetAllCertifiers(ctx sdk.Context) (certifiers certtypes.Certifiers)
	GetCertifier(ctx sdk.Context, certifierAddress sdk.AccAddress) (certtypes.Certifier, error)
	HasCertifierAlias(ctx sdk.Context, alias string) bool
	IsCertified(ctx sdk.Context, content string, certType string) bool
	GetCertifiedIdentities(ctx sdk.Context) []sdk.AccAddress
}

type UpgradeKeeper interface {
	ValidatePlan(ctx sdk.Context, plan upgradetypes.Plan) error
}

type ShieldKeeper interface {
	GetPurchase(purchaseList shieldtypes.PurchaseList, purchaseID uint64) (shieldtypes.Purchase, bool)
	GetPurchaseList(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress) (shieldtypes.PurchaseList, bool)
	GetClaimProposalParams(ctx sdk.Context) shieldtypes.ClaimProposalParams
	SecureCollaterals(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, purchaseID uint64, loss sdk.Coins, lockPeriod time.Duration) error
	RestoreShield(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, id uint64, loss sdk.Coins) error
	ClaimEnd(ctx sdk.Context, id, poolID uint64, loss sdk.Coins)
}

type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
	WithKeyTable(table paramtypes.KeyTable) paramtypes.Subspace
}

type StakingKeeper interface {
	IterateBondedValidatorsByPower(sdk.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	TotalBondedTokens(sdk.Context) sdk.Int
	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress, fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))
	BondDenom(sdk.Context) string
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
}
