package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	"github.com/certikfoundation/shentu/x/cert"
	"github.com/certikfoundation/shentu/x/shield"
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

type ShieldKeeper interface {
	GetPurchaseList(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress) (shield.PurchaseList, bool)
	GetClaimProposalParams(ctx sdk.Context) shield.ClaimProposalParams
	ClaimLock(ctx sdk.Context, proposalID, poolID uint64, purchaser sdk.AccAddress, purchaseID uint64, loss sdk.Coins, lockPeriod time.Duration) error
	ClaimUnlock(ctx sdk.Context, proposalID, poolID uint64, loss sdk.Coins) error
	RestoreShield(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, id uint64, loss sdk.Coins) error
}

type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
	WithKeyTable(table subspace.KeyTable) subspace.Subspace
}
