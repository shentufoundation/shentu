package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Pool contains a shield project pool's data.
type Pool struct {
	// ID is the id of the pool.
	ID uint64 `json:"id" yaml:"id"`

	// Description is the term of the pool.
	Description string `json:"description" yaml:"description"`

	// Sponsor is the project owner of the pool.
	Sponsor string `json:"sponsor" yaml:"sponsor"`

	// SponsorAddress is the CertiK Chain address of the sponsor.
	SponsorAddress sdk.AccAddress `json:"sponsor_address" yaml:"sponsor_address"`

	// ShieldLimit is the maximum shield can be purchased for the pool.
	ShieldLimit sdk.Int `json:"shield_limit" yaml:"shield_limit"`

	// Active means new purchases are allowed.
	Active bool `json:"active" yaml:"active"`

	// Shield is the amount of all active purchased shields.
	Shield sdk.Int `json:"shield" yaml:"shield"`
}

// NewPool creates a new project pool.
func NewPool(id uint64, description, sponsor string, sponsorAddress sdk.AccAddress, shieldLimit sdk.Int, shield sdk.Int) Pool {
	return Pool{
		ID:             id,
		Description:    description,
		Sponsor:        sponsor,
		SponsorAddress: sponsorAddress,
		ShieldLimit:    shieldLimit,
		Active:         true,
		Shield:         shield,
	}
}

// Provider tracks total delegation, total collateral, and rewards of a provider.
type Provider struct {
	// Address is the address of the provider.
	Address sdk.AccAddress `json:"address" yaml:"address"`

	// DelegationBonded is the amount of bonded delegation.
	DelegationBonded sdk.Int `json:"delegation_bonded" yaml:"delegation_bonded"`

	// Collateral is amount of all collaterals for the provider, including
	// those in withdraw queue but excluding those currently locked, in all
	// pools.
	Collateral sdk.Int `json:"collateral" yaml:"collateral"`

	// TotalLocked is the amount locked for pending claims.
	TotalLocked sdk.Int `json:"total_locked" yaml:"total_locked"`

	// Withdrawing is the amount of collateral in withdraw queues.
	Withdrawing sdk.Int `json:"withdrawing" yaml:"withdrawing"`

	// Rewards is the pooling rewards to be collected.
	Rewards MixedDecCoins `json:"rewards" yaml:"rewards"`
}

// NewProvider creates a new provider object.
func NewProvider(addr sdk.AccAddress) Provider {
	return Provider{
		Address:          addr,
		DelegationBonded: sdk.ZeroInt(),
		Collateral:       sdk.ZeroInt(),
		TotalLocked:      sdk.ZeroInt(),
		Withdrawing:      sdk.ZeroInt(),
	}
}

// Purchase record an individual purchase.
type Purchase struct {
	// PurchaseID is the purchase_id.
	PurchaseID uint64 `json:"purchase_id" yaml:"purchase_id"`

	// ProtectionEndTime is the time when the protection of the shield ends.
	ProtectionEndTime time.Time `json:"protection_end_time" yaml:"protection_end_time"`

	// DeletionTime is the time when the purchase should be deleted.
	DeletionTime time.Time `json:"deletion_time" yaml:"deletion_time"`

	// Description is the information about the protected asset.
	Description string `json:"description" yaml:"description"`

	// Shield is the unused amount of shield purchased.
	Shield sdk.Int `json:"shield" yaml:"shield"`

	// ServiceFees is the service fees paid by this purchase.
	ServiceFees MixedDecCoins `json:"service_fees" yaml:"service_fees"`
}

// NewPurchase creates a new purchase object.
func NewPurchase(purchaseID uint64, protectionEndTime, deletionTime time.Time, description string, shield sdk.Int, serviceFees MixedDecCoins) Purchase {
	return Purchase{
		PurchaseID:        purchaseID,
		ProtectionEndTime: protectionEndTime,
		DeletionTime:      deletionTime,
		Description:       description,
		Shield:            shield,
		ServiceFees:       serviceFees,
	}
}

// PurchaseList is a collection of purchase.
type PurchaseList struct {
	// PoolID is the id of the shield of the purchase.
	PoolID uint64 `json:"pool_id" yaml:"pool_id"`

	// Purchaser is the address making the purchase.
	Purchaser sdk.AccAddress `json:"purchaser" yaml:"purchaser"`

	// Entries stores all purchases by the purchaser in the pool.
	Entries []Purchase `json:"entries" yaml:"entries"`
}

// NewPurchaseList creates a new purchase list.
func NewPurchaseList(poolID uint64, purchaser sdk.AccAddress, purchases []Purchase) PurchaseList {
	return PurchaseList{
		PoolID:    poolID,
		Purchaser: purchaser,
		Entries:   purchases,
	}
}

// PoolPurchase is a pair of pool id and purchaser.
type PoolPurchaser struct {
	// PoolID is the id of the shield pool.
	PoolID uint64

	// Purchaser is the chain address of the purchaser.
	Purchaser sdk.AccAddress
}

// Withdraw stores an ongoing withdraw of pool collateral.
type Withdraw struct {
	// Address is the chain address of the provider withdrawing.
	Address sdk.AccAddress `json:"address" yaml:"address"`

	// Amount is the amount of withdraw.
	Amount sdk.Int `json:"amount" yaml:"amount"`

	// CompletionTime is the scheduled withdraw completion time.
	CompletionTime time.Time `json:"completion_time" yaml:"completion_time"`
}

// NewWithdraw creates a new withdraw object.
func NewWithdraw(addr sdk.AccAddress, amount sdk.Int, completionTime time.Time) Withdraw {
	return Withdraw{
		Address:        addr,
		Amount:         amount,
		CompletionTime: completionTime,
	}
}

// Withdraws contains multiple withdraws.
type Withdraws []Withdraw

type ShieldStaking struct {
	PoolID            uint64         `json:"pool_id" yaml:"pool_id"`
	Purchaser         sdk.AccAddress `json:"purchaser" yaml:"purchaser"`
	Amount            sdk.Int        `json:"amount" yaml:"amount"`
	WithdrawRequested sdk.Int        `json:"withdraw_requested" yaml:"withdraw_requested"`
}

func NewShieldStaking(poolID uint64, purchaser sdk.AccAddress, amount sdk.Int) ShieldStaking {
	return ShieldStaking{
		PoolID:            poolID,
		Purchaser:         purchaser,
		Amount:            amount,
		WithdrawRequested: sdk.NewInt(0),
	}
}
