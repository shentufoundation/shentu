package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Global variables
// 
// 	// TotalCollateral is the amount of all collaterals in the pool.
// 	TotalCollateral sdk.Int `json:"total_collateral" yaml:"total_collateral"`
//
//	// TotalWithdrawing is the amount of collateral in withdraw queues.
//	TotalWithdrawing sdk.Int `json:"withdrawing" yaml:"withdrawing"`
//
//	// Shield is the amount of all active purchased shields.
//	Shield sdk.Coins `json:"shield" yaml:"shield"`
//
//	// Premium is the undistributed pool premium from the sponsor.
//	Premium MixedDecCoins `json:"premium" yaml:"premium"`
//
//	// TotalLocked is the amount of collaterals locked for pending claims.
//	TotalLocked sdk.Int `json:"total_locked" yaml:"total_locked"`

// Pool contains a shield pool's data.
type Pool struct {
	// PoolID is the id of the pool.
	PoolID uint64 `json:"pool_id" yaml:"pool_id"`

	// Description is the term of the pool.
	Description string `json:"description" yaml:"description"`

	// Sponsor is the project owner of the pool.
	Sponsor string `json:"sponsor" yaml:"sponsor"`

	// SponsorAddr is the CertiK Chain address of the sponsor.
	SponsorAddr sdk.AccAddress `json:"sponsor_address" yaml:"sponsor_address"`

	// Active means new purchases are allowed.
	Active bool `json:"active" yaml:"active"`

	// Shield is the amount of all active purchased shields.
	Shield sdk.Coins `json:"shield" yaml:"shield"`
}

// NewPool creates a new shield pool.
func NewPool(shield sdk.Coins, totalCollateral sdk.Int, deposit MixedDecCoins, sponsor string, sponsorAddr sdk.AccAddress, endTime time.Time, id uint64) Pool {
	return Pool{
		Shield:          shield,
		Premium:         deposit,
		Sponsor:         sponsor,
		SponsorAddr:     sponsorAddr,
		Active:          true,
		TotalCollateral: totalCollateral,
		EndTime:         endTime,
		PoolID:          id,
	}
}

// To purchase
//
// SUM(Collateral) - SUM(Withdrawing) >= SHIELD

// To withdraw
//
// Min(Collateral - Withdrawing, SUM(Collateral) - SUM(Withdrawing) - SUM(SHIELD)) >= WITHDRAW

// Pool 100 0  A1 50 0  A2 50 0 Shield 0
// P1 => purchase 50
// Pool 100 0  A1 50 0  A2 50 0 Shield 50 // A1 and A2 both get rewards from P1
// Day 1 A1 withdraw from unbond 50
// Pool 100 50 A1 50 50 A2 50 0 Shield 50
// Day 21 P1 claims loss of 50
// Day 22
// Pool 50 0 A1 0 0 A2 50 0 Shield 0 Locked 50
// Day 25 P1 claim of 50 passed
// Pool 0 0 A1 0 0 A2 0 0 Shield 0 Locked 0

// Provider tracks total delegation, total collateral, and rewards of a provider.
type Provider struct {
	// Address is the address of the provider.
	Address sdk.AccAddress `json:"address" yaml:"address"`

	// DelegationBonded is the amount of bonded delegation.
	DelegationBonded sdk.Int `json:"delegation_bonded" yaml:"delegation_bonded"`

	// B 100 UB 0 Collateral 100 Withdraw 0
	// Redelegate 10
	// N: Unbond 10
	// B 90 UB 10 Collateral 100 Withdraw 10
	// N: Bond 10
	// B 100 UB 10 Collateral 100 Withdraw 10
	
	// Collateral is amount of all collaterals for the provider, including
	// those in withdraw queue but excluding those currently locked, in all
	// pools.
	Collateral sdk.Int `json:"collateral" yaml:"collateral"`

	// Withdrawing is the amount of collateral in withdraw queues.
	Withdrawing sdk.Int `json:"withdrawing" yaml:"withdrawing"`

	// TotalLocked is the amount locked for pending claims.
	TotalLocked sdk.Int `json:"total_locked" yaml:"total_locked"`

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
		Available:        sdk.ZeroInt(),
		Withdrawing:      sdk.ZeroInt(),
	}
}

// Purchase record an individual purchase.
type Purchase struct {
	// PurchaseID is the purchase_id.
	PurchaseID uint64 `json:"purchase_id" yaml:"purchase_id"`

	// Description is the information about the protected asset.
	Description string `json:"description" yaml:"description"`

	// Shield is the unused amount of shield purchased.
	Shield sdk.Coins `json:"shield" yaml:"shield"`

	// ProtectionEndTime is the time when the protection of the shield ends.
	ProtectionEndTime time.Time `json:"protection_end_time" yaml:"protection_end_time"`
}

// NewPurchase creates a new purchase object.
func NewPurchase(purchaseID uint64, shield sdk.Coins, startBlockHeight int64, protectionEndTime, claimPeriodEndTime, deleteTime time.Time, description string) Purchase {
	return Purchase{
		PurchaseID:         purchaseID,
		Description:        description,
		Shield:             shield,
		StartBlockHeight:   startBlockHeight,
		ProtectionEndTime:  protectionEndTime,
		ClaimPeriodEndTime: claimPeriodEndTime,
		DeleteTime:         deleteTime,
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

// PoolPurchaser is a pair of pool id and purchaser.
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
}

// NewWithdraw creates a new withdraw object.
func NewWithdraw(poolID uint64, addr sdk.AccAddress, amount sdk.Int, completionTime time.Time) Withdraw {
	return Withdraw{
		PoolID:         poolID,
		Address:        addr,
		Amount:         amount,
		CompletionTime: completionTime,
	}
}

// Withdraws contains multiple withdraws.
type Withdraws []Withdraw
