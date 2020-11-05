# Shield

CertiKShield is a decentralized pool of CTK that uses CertiK Chain on-chain governance system to reimburse lost, stolen, or inaccessible assets from any blockchain network. There are two members of the CertiKShield system: Collateral Providers and Shield Purchasers. Providers contribute cryptocurrency as collateral to fill the CertiKShield Pool. In return, they receive a portion of the fees paid by Purchasers, in addition to the usual staking rewards. Purchasers pay a recurring fee, based on their riskiness as determined by their CertiK Security Oracle score, that entitles them to submit a Claim Proposal to be reimbursed from the pool for stolen assets.

See the [whitepaper](https://www.certik.foundation/whitepaper#3-CertiKShield) for more information on CertiKShield.

## State
```go
type MixedCoins struct {
	Native  sdk.Coins
	Foreign sdk.Coins
}
```

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

// PurchaseList is a collection of purchase.
type PurchaseList struct {
	// PoolID is the id of the shield of the purchase.
	PoolID uint64 `json:"pool_id" yaml:"pool_id"`

	// Purchaser is the address making the purchase.
	Purchaser sdk.AccAddress `json:"purchaser" yaml:"purchaser"`

	// Entries stores all purchases by the purchaser in the pool.
	Entries []Purchase `json:"entries" yaml:"entries"`
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

## Messages
// MsgCreatePool defines the attributes of a create-pool transaction.
type MsgCreatePool struct {
	From        sdk.AccAddress `json:"from" yaml:"from"`
	Shield      sdk.Coins      `json:"shield" yaml:"shield"`
	Deposit     MixedCoins     `json:"deposit" yaml:"deposit"`
	Sponsor     string         `json:"sponsor" yaml:"sponsor"`
	SponsorAddr sdk.AccAddress `json:"sponsor_addr" yaml:"sponsor_addr"`
	Description string         `json:"description" yaml:"description"`
	ShieldLimit sdk.Int        `json:"shield_limit" yaml:"shield_limit"`
}
// MsgUpdatePool defines the attributes of a shield pool update transaction.
type MsgUpdatePool struct {
	From        sdk.AccAddress `json:"from" yaml:"from"`
	Shield      sdk.Coins      `json:"Shield" yaml:"Shield"`
	ServiceFees MixedCoins     `json:"service_fees" yaml:"service_fees"`
	PoolID      uint64         `json:"pool_id" yaml:"pool_id"`
	Description string         `json:"description" yaml:"description"`
	ShieldLimit sdk.Int        `json:"shield_limit" yaml:"shield_limit"`
}
// MsgPausePool defines the attributes of a pausing a shield pool.
type MsgPausePool struct {
	From   sdk.AccAddress `json:"from" yaml:"from"`
	PoolID uint64         `json:"pool_id" yaml:"pool_id"`
}
// MsgResumePool defines the attributes of a resuming a shield pool.
type MsgResumePool struct {
	From   sdk.AccAddress `json:"from" yaml:"from"`
	PoolID uint64         `json:"pool_id" yaml:"pool_id"`
}
// MsgDepositCollateral defines the attributes of a depositing collaterals.
type MsgDepositCollateral struct {
	From       sdk.AccAddress `json:"sender" yaml:"sender"`
	Collateral sdk.Coins      `json:"collateral" yaml:"collateral"`
}
// MsgDepositCollateral defines the attributes of a depositing collaterals.
type MsgDepositCollateral struct {
	From       sdk.AccAddress `json:"sender" yaml:"sender"`
	Collateral sdk.Coins      `json:"collateral" yaml:"collateral"`
}
// MsgWithdrawForeignRewards defines attribute of withdraw rewards transaction.
type MsgWithdrawRewards struct {
	From sdk.AccAddress `json:"sender" yaml:"sender"`
}
// MsgWithdrawForeignRewards defines attributes of withdraw foreign rewards transaction.
type MsgWithdrawForeignRewards struct {
	From   sdk.AccAddress `json:"sender" yaml:"sender"`
	Denom  string         `json:"denom" yaml:"denom"`
	ToAddr string         `json:"to_addr" yaml:"to_addr"`
}
// MsgClearPayouts defines attributes of clear payouts transaction.
type MsgClearPayouts struct {
	From  sdk.AccAddress `json:"sender" yaml:"sender"`
	Denom string         `json:"denom" yaml:"denom"`
}
// MsgPurchaseShield defines the attributes of purchase shield transaction.
type MsgPurchaseShield struct {
	PoolID      uint64         `json:"pool_id" yaml:"pool_id"`
	Shield      sdk.Coins      `json:"shield" yaml:"shield"`
	Description string         `json:"description" yaml:"description"`
	From        sdk.AccAddress `json:"from" yaml:"from"`
}
// MsgWithdrawReimburse defines the attributes of withdraw reimbursement transaction.
type MsgWithdrawReimbursement struct {
	ProposalID uint64         `json:"proposal_id" yaml:"proposal_id"`
	From       sdk.AccAddress `json:"from" yaml:"from"`
}
// MsgUnstakeFromShield defines the attributes of staking for purchase transaction.
type MsgUnstakeFromShield struct {
	PoolID uint64         `json:"pool_id" yaml:"pool_id"`
	Shield sdk.Coins      `json:"shield" yaml:"shield"`
	From   sdk.AccAddress `json:"from" yaml:"from"`
}

## Parameters
// PoolParams defines the parameters for the shield pool.
type PoolParams struct {
	ProtectionPeriod  time.Duration `json:"protection_period" yaml:"protection_period"`
	ShieldFeesRate    sdk.Dec       `json:"shield_fees_rate" yaml:"shield_fees_rate"`
	WithdrawPeriod    time.Duration `json:"withdraw_period" yaml:"withdraw_period"`
	PoolShieldLimit   sdk.Dec       `json:"pool_shield_limit" yaml:"pool_shield_limit"`
	MinShieldPurchase sdk.Coins     `json:"min_shield_purchase" yaml:"min_shield_purchase"`
}
// ClaimProposalParams defines the parameters for the shield claim proposals.
type ClaimProposalParams struct {
	ClaimPeriod  time.Duration `json:"claim_period" yaml:"claim_period"`
	PayoutPeriod time.Duration `json:"payout_period" yaml:"payout_period"`
	MinDeposit   sdk.Coins     `json:"min_deposit" json:"min_deposit"`
	DepositRate  sdk.Dec       `json:"deposit_rate" yaml:"deposit_rate"`
	FeesRate     sdk.Dec       `json:"fees_rate" yaml:"fees_rate"`
}

var (
	// default values for Shield pool's parameters
	DefaultProtectionPeriod  = time.Hour * 24 * 21                                                   // 21 days
	DefaultShieldFeesRate    = sdk.NewDecWithPrec(769, 5)                                            // 0.769%
	DefaultWithdrawPeriod    = time.Hour * 24 * 21                                                   // 21 days
	DefaultPoolShieldLimit   = sdk.NewDecWithPrec(50, 2)                                             // 50%
	DefaultMinShieldPurchase = sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(50000000))) // 50 CTK

	// default values for Shield claim proposal's parameters
	DefaultClaimPeriod              = time.Hour * 24 * 21                                                    // 21 days
	DefaultPayoutPeriod             = time.Hour * 24 * 56                                                    // 56 days
	DefaultMinClaimProposalDeposit  = sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(100000000))) // 100 CTK
	DefaultClaimProposalDepositRate = sdk.NewDecWithPrec(10, 2)                                              // 10%
	DefaultClaimProposalFeesRate    = sdk.NewDecWithPrec(1, 2)                                               // 1%

	// default value for staking-shield rate parameter
	DefaultStakingShieldRate = sdk.NewDec(2)
)
