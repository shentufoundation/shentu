package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

	// Premium is the undistributed pool
	Premium MixedDecCoins `json:"premium" yaml:"premium"`

	// TotalCollateral is the amount of all collaterals in the pool.
	TotalCollateral sdk.Int `json:"total_collateral" yaml:"total_collateral"`

	// Available is the amount of collaterals available to be purchased.
	Available sdk.Int `json:"available" yaml:"available"`

	// TotalLocked is the amount of
	TotalLocked sdk.Int `json:"total_locked" yaml:"total_locked"`

	// Shiled is the amount of unused shield for the pool sponsor.
	Shield sdk.Coins `json:"shield" yaml:"shield"`

	// EndTime is the time pool maintainence ends.
	EndTime time.Time `json:"end_time" yaml:"end_time"`
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

// Collateral records the collaterals provided by a provider in a shield pool.
type Collateral struct {
	// PoolID is the id of the shield pool.
	PoolID uint64 `json:"pool_id" yaml:"pool_id"`

	// Provider is the chain address of the provider.
	Provider sdk.AccAddress `json:"provider" yaml:"provider"`

	// Amount is the collateral amount, excluding withdrawing or locked.
	Amount sdk.Int `json:"amount" yaml:"amount"`

	// Withdrawing is the amount of collateral in withdrawing process.
	Withdrawing sdk.Int `json:"withdrawing" yaml:"withdrawing"`

	// TotalLocked is the amount of collateral locked up for pending claims.
	TotalLocked sdk.Int `json:"total_locked" yaml:"total_locked"`

	// LockedCollaterals stores collaterals locked up for claims against the
	// provider and the pool.
	LockedCollaterals []LockedCollateral `json:"locked_collaterals" yaml:"locked_collaterals"`
}

// NewCollateral creates a new collateral object.
func NewCollateral(pool Pool, provider sdk.AccAddress, amount sdk.Int) Collateral {
	return Collateral{
		PoolID:   pool.PoolID,
		Provider: provider,
		Amount:   amount,
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

	// Available is the amount of staked CTK available to be deposited.
	Available sdk.Int `json:"available" yaml:"available"`

	// Withdrawing is the amount of collateral in withdraw queues.
	Withdrawing sdk.Int `json:"withrawal" yaml:"withdraw"`

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
	// PurchaseID is the purchase_id
	PurchaseID uint64 `json:"purchase_id" yaml:"purchase_id"`

	// Description is the information about the protected asset.
	Description string `json:"description" yaml:"description"`

	// Shield is the unused amount of sheild purchased.
	Shield sdk.Coins `json:"shield" yaml:"shield"`

	// StartBlockHeight is the purchasing block height.
	StartBlockHeight int64 `json:"start_block_height" yaml:"start_block_height"`

	// ProtectionEndTime is the time when the protection of the shield ends.
	ProtectionEndTime time.Time `json:"protection_end_time" yaml:"protection_end_time"`

	// ClaimPeriodEndTime is the time when any claims to the shield must be
	// filed before.
	ClaimPeriodEndTime time.Time `json:"claim_period_end_time" yaml:"claim_period_end_time"`

	// ExpirationTime is the time when the purchase is scheduled to be
	// deleted.
	ExpirationTime time.Time `json:"expiration_time" yaml:"expiration_time"`
}

// NewPurchase creates a new purhase object.
func NewPurchase(purchaseID uint64, shield sdk.Coins, startBlockHeight int64, protectionEndTime, claimPeriodEndTime, expirationTime time.Time, description string) Purchase {
	return Purchase{
		PurchaseID:         purchaseID,
		Description:        description,
		Shield:             shield,
		StartBlockHeight:   startBlockHeight,
		ProtectionEndTime:  protectionEndTime,
		ClaimPeriodEndTime: claimPeriodEndTime,
		ExpirationTime:     expirationTime,
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
	// PoolID is the id of the shield withdrawing collateral from.
	PoolID uint64 `json:"pool_id" yaml:"pool_id"`

	// Address is the chain address of the provider withdrawing.
	Address sdk.AccAddress `json:"address" yaml:"address"`

	// Amount is the amount of withdraw.
	Amount sdk.Int `json:"amount" yaml:"amount"`

	// CompletionTime is the scheduled withdraw completion time.
	CompletionTime time.Time `json:"completion_time" yaml:"completion_time"`
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
