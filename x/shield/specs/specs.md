# Shield

CertiKShield is a decentralized pool of CTK that uses CertiK Chain on-chain governance system to reimburse lost, stolen, or inaccessible assets from any blockchain network. There are two members of the CertiKShield system: Collateral Providers and Shield Purchasers. Providers contribute cryptocurrency as collateral to fill the CertiKShield Pool. In return, they receive a portion of the fees paid by Purchasers, in addition to the usual staking rewards. Purchasers pay a recurring fee, based on their riskiness as determined by their CertiK Security Oracle score, that entitles them to submit a Claim Proposal to be reimbursed from the pool for stolen assets.

See the [whitepaper](https://www.certik.foundation/whitepaper#3-CertiKShield) for more information on CertiKShield.

## State

### Admins

`Admin` represents the Shield admin account address.

- Admin: `0x0 -> sdk.AccAddress`

### Pools

Every project that wants to buy a Shield needs to have a `Pool` created. Then, a project can purchase Shields (a.k.a. `Purchase`s) up to their `ShieldLimit`.

- Pool: `0x7 | LittleEndian(Id) -> amino(pool)`
- GetNextPoolId: `0x8 -> LittleEndian(Id)`

```go
// Pool contains a shield project pool's data.
type Pool struct {
    Id          uint64  `json:"id" yaml:"id"`
    Description string  `json:"description" yaml:"description"`
    Sponsor     string  `json:"sponsor" yaml:"sponsor"`
    SponsorAddr string  `json:"sponsor_addr" yaml:"sponsor_addr"`
    ShieldLimit sdk.Int `json:"shield_limit" yaml:"shield_limit"`
    Active      bool    `json:"active" yaml:"active"`
    Shield      sdk.Int `json:"shield" yaml:"shield"`
}
```

Relevant states that are tracked along with the pool are wrapped in `sdk.IntProto` objects.

- TotalCollateral: `0x1 -> amino(totalCollateral)`
- TotalWithdrawing: `0x2 -> amino(totalWithdrawing)`
- TotalShield: `0x3 -> amino(totalShield)`
- TotalClaimed: `0x4 -> amino(totalClaimed)`

`ServiceFees` are small fees charged for each purchase, and are stored as `MixDecCoins` objects. `MixedDecCoins` keeps track of native and foreign decimal tokens together.

- ServiceFees: `0x5 -> amino(serviceFees)`
- RemainingServiceFees: `0x6 -> amino(serviceFees)`
- BlockServiceFees: `0x13 -> amino(serviceFees)`

```go
// MixedDecCoins defines the struct for mixed coins in decimal with native and foreign decimal coins.
type MixedDecCoins struct {
    Native  sdk.DecCoins    `json:"native"`
    Foreign sdk.DecCoins    `json:"foreign"`
}
```

### Providers

`Provider` tracks total delegation, total collateral, and rewards of a provider. A collateral provider contributes their assets as the collateral, which is used to pay out approved reimbursement requests from Shield purchasers. The provider earns staking rewards for their staked assets as well as a portion of the fees paid by purchasers.

- Provider: `0xC | Address -> amino(provider)`

```go
// Provider tracks total delegation, total collateral, and rewards of a provider.
type Provider struct {
    // Address is the address of the provider.
    Address             string          `json:"address" yaml:"address"`
    // DelegationBonded is the amount of bonded delegation.
    DelegationBonded    sdk.Int         `json:"delegation_bonded" yaml:"provider"`
    // Collateral is amount of all collaterals for the provider, including
    // those in withdraw queue but excluding those currently locked, in all
    // pools.
    Collateral          sdk.Int         `json:"collateral" yaml:"collateral"`
    // TotalLocked is the amount locked for pending claims.
    TotalLocked         sdk.Int         `json:"total_locked" yaml:"total_locked"`
    // Withdrawing is the amount of collateral in withdraw queues.
    Withdrawing         sdk.Int         `json:"withdrawing" yaml:"withdrawing"`
    // Rewards is the pooling rewards to be collected.
    Rewards             MixedDecCoins   `json:"rewards" yaml:"rewards"`
}
```

### Purchases

`Purchase` records an individual purchase. Purchases are stored in the store as `PurchaseList` objects.

- PurchaseList: `0xA | LittleEndian(Id) | Purchaser -> amino(purchaseList)`
- NextPurchaseId: `0x9 -> LittleEndian(Id)`

```go
// Purchase record an individual purchase.
type Purchase struct {
    // PurchaseID is the purchase_id.
    PurchaseId          uint64          `json:"purchase_id" yaml:"purchase_id"`
    // ProtectionEndTime is the time when the protection of the shield ends.
    ProtectionEndTime   time.Time       `json:"protection_end_time" yaml:"protection_end_time"`
    // DeletionTime is the time when the purchase should be deleted.
    DeletionTime        time.Time       `json:"deletion_time" yaml:"deletion_time"`
    // Description is the information about the protected asset.
    Description         string          `json:"description" yaml:"description"`
    // Shield is the unused amount of shield purchased.
    Shield              sdk.Int         `json:"shield" yaml:"shield"`
    // ServiceFees is the service fees paid by this purchase.
    ServiceFees         MixedDecCoins   `json:"service_fees" yaml:"service_fees"`
}
```

```go
// PurchaseList is a collection of purchase.
type PurchaseList struct {
    // PoolID is the id of the shield of the purchase.
    PoolId      uint64      `json:"pool_id,omitempty" yaml:"pool_id"`
    // Purchaser is the address making the purchase.
    Purchaser   string      `json:"purchaser,omitempty" yaml:"purchaser"`
    // Entries stores all purchases by the purchaser in the pool.
    Entries     []Purchase  `json:"entries" yaml:"entries"`
}
```

Purchases are queued as `(PoolId, Purchaser)` pairs according to their expiration timestamps. Shield purchasers have their assets protected upon purchase of Shields until the expiration timestamp.

- PurchaseExpirationTime: `0xB | Timestamp -> amino(poolPurchaserPairs)`

```go
// PoolPurchase is a pair of pool id and purchaser.
type PoolPurchaser struct {
    // PoolID is the id of the shield pool.
    PoolId      uint64  `json:"pool_id" yaml:"pool_id"`
    // Purchaser is the chain address of the purchaser.
    Purchaser   string  `json:"purchaser" yaml:"purchaser"`
}
```

`LastUpdateTime` is set when the purchase is made or service fees are distributed.

- LastUpdateTime: `0xE -> amino(time)`

### Withdraws

`Withdraw` stores an ongoing withdraw of pool collateral. Withdraws are queued according to their completion timestamps.

- WithdrawQueue: `0xD -> amino([]withdraw)`

```go
// Withdraw stores an ongoing withdraw of pool collateral.
type Withdraw struct {
    // Address is the chain address of the provider withdrawing.
    Address         string      `json:"address" yaml:"address"`
    // Amount is the amount of withdraw.
    Amount          sdk.Int     `json:"amount" yaml:"amount"`
    // CompletionTime is the scheduled withdraw completion time.
    CompletionTime  time.Time   `json:"completion_time" yaml:"completion_time"`
}
```

### Staking Purchases

Collateral providers can stake on Shield pool, which is a higher-risk, higher-reward staking alternative to CertiK Node staking. Providers can stake assets as collaterals on `GlobalShieldStakingPool`, from which purchases are then stored as `StakeForShield` and `OriginalStaking`, which keep track of purchases and staking amounts, respectively. Shield staking purchases are stored as `ShieldStaking` objects.

- GlobalStakeForShieldPool: `0xF -> amino(pool)`
- StakeForShield: `0x11 | LittleEndian(PoolId) | Purchaser -> amino(purchase)`
- OriginalStaking: `0x13 | LittleEndian(PurchaseId) -> amino(stakingAmt)`

```go
type ShieldStaking struct {
    PoolId              uint64  `json:"pool_id" yaml:"pool_id"`
    Purchaser           string  `json:"purchaser" yaml:"purchaser"`
    Amount              sdk.Int `json:"amount" yaml:"amount"`
    WithdrawRequested   sdk.Int `json:"withdraw_requested" yaml:"withdraw_requested"`
}
```

### Reimbursements

`Reimbursement` tracks relevant information for the payout granted upon approval of `ShieldClaimProposal`, which is a proposal Shield purchasers submit when they lose protected assets.

- Reimbursement: `0x14 | LittleEndian(ProposalId) -> amino(reimbursement)`

```go
type ShieldClaimProposal struct {
    ProposalId  uint64      `json:"proposal_id" yaml:"proposal_id"`
    PoolId      uint64      `json:"pool_id" yaml:"pool_id"`
    Loss        sdk.Coins   `json:"loss" yaml:"loss"`
    Evidence    string      `json:"evidence" yaml:"evidence"`
    Description string      `json:"description" yaml:"description"`
    Proposer    string      `json:"proposer" yaml:"proposer"`
}

type Reimbursement struct {
    Amount      sdk.Coins   `json:"amount"`
    Beneficiary string      `json:"beneficiary" yaml:"beneficiary"`
    PayoutTime  time.Time   `json:"payout_time" yaml:"payout_time"`
}
```

## Messages

### Pools

`MsgCreatePool` creates a new pool for a project. Once created, the project can purchase Shields via `MsgPurchaseShield`. The pool can be updated with `MsgUpdatePool`--for example, a project's `ShieldLimit` could be increased.

```go
// MsgCreatePool defines the attributes of a create-pool transaction.
type MsgCreatePool struct {
    From        string      `json:"from" yaml:"from"`
    Shield      sdk.Coins   `json:"shield"`
    Deposit     MixedCoins  `json:"deposit" yaml:"deposit"`
    Sponsor     string      `json:"sponsor" yaml:"sponsor"`
    SponsorAddr string      `json:"sponsor_addr" yaml:"sponsor_addr"`
    Description string      `json:"description" yaml:"description"`
    ShieldLimit sdk.Int     `json:"shield_limit"`
}

// MsgUpdatePool defines the attributes of a shield pool update transaction.
type MsgUpdatePool struct {
    From        string      `json:"from" yaml:"from"`
    Shield      sdk.Coins   `json:"shield"`
    ServiceFees MixedCoins  `json:"service_fees" yaml:"service_fees"`
    PoolId      uint64      `json:"pool_id" yaml:"pool_id"`
    Description string      `json:"description" yaml:"description"`
    ShieldLimit sdk.Int     `json:"shield_limit"`
}
```

`MsgPausePool` sets the pool's `Active` to `false`; `MsgResumePool` sets it to `true`. While inactive, new Shields cannot be purchased.

```go
// MsgPausePool defines the attributes of a pausing a shield pool.
type MsgPausePool struct {
    From    string  `json:"from" yaml:"from"`
    PoolId  uint64  `json:"pool_id" yaml:"pool_id"`
}

// MsgResumePool defines the attributes of a resuming a shield pool.
type MsgResumePool struct {
    From    string  `json:"from" yaml:"from"`
    PoolId  uint64  `json:"pool_id" yaml:"pool_id"`
}
```

Projects with a `Pool` can use `MsgPurchaseShield` to purchase a new Shield.

```go
// MsgPurchaseShield defines the attributes of purchase shield transaction.
type MsgPurchaseShield struct {
    PoolId      uint64      `json:"pool_id" yaml:"pool_id"`
    Shield      sdk.Coins   `json:"shield"`
    Description string      `json:"description" yaml:"description"`
    From        string      `json:"from" yaml:"from"`
}
```

### Deposits

`MsgDepositCollateral` creates a new provider with the given `Collateral`, or it adds `Collateral` to an existing provider's collateral. There's no `MsgCreateProvider` because this message has that functionality.

```go
// MsgDepositCollateral defines the attributes of a depositing collaterals.
type MsgDepositCollateral struct {
    From       string       `json:"from" yaml:"from"`
    Collateral sdk.Coins    `json:"collateral"`
}
```

### Withdraws

`MsgWithdrawCollateral` inserts a collateral withdraw to the withdraw queue.

```go
// MsgWithdrawCollateral defines the attributes of a withdrawing collaterals.
type MsgWithdrawCollateral struct {
    From       string       `json:"from" yaml:"from"`
    Collateral sdk.Coins    `json:"collateral"`
}
```

`MsgWithdrawRewards` pays out pending CTK rewards. Currently, `MsgWithdrawForeignRewards` and `MsgClearPayouts` are not callable or implemented.

```go
// MsgWithdrawRewards defines attribute of withdraw rewards transaction.
type MsgWithdrawRewards struct {
    From    string  `json:"from" yaml:"from"`
}

// MsgWithdrawForeignRewards defines attributes of withdraw foreign rewards transaction.
type MsgWithdrawForeignRewards struct {
    From    string  `json:"from" yaml:"from"`
    Denom   string  `json:"denom" yaml:"denom"`
    ToAddr  string  `json:"to_addr" yaml:"to_addr"`
}

`MsgWithdrawReimbursement` withdraws a reimbursement made for a beneficiary.

```go
// MsgWithdrawReimbursement defines the attributes of withdraw reimbursement transaction.
type MsgWithdrawReimbursement struct {
    ProposalId  uint64  `json:"proposal_id" yaml:"proposal_id"`
    From        string  `json:"from" yaml:"from"`
}
```

`MsgStakeForShield` purchases a Shield and pays for the fees using staking.

```go
// MsgStakeForShield defines the attributes of staking for purchase transaction.
type MsgStakeForShield struct {
    PoolId      uint64      `json:"pool_id" yaml:"pool_id"`
    Shield      sdk.Coins   `json:"shield"`
    Description string      `json:"description" yaml:"description"`
    From        string      `json:"from" yaml:"from"`
}

// MsgUnstakeFromShield defines the attributes of staking for purchase transaction.
type MsgUnstakeFromShield struct {
    PoolId  uint64      `json:"pool_id,omitempty" yaml:"pool_id"`
    Shield  sdk.Coins   `json:"shield"`
    From    string      `json:"from" yaml:"from"`
}
```

`MsgUpdateSponsor` updates the sponsor information of a given pool specified by `PoolID`.
```go
// MsgUpdateSponsor defines the attributes of a update-sponsor transaction.
type MsgUpdateSponsor struct {
    PoolId      uint64  `json:"pool_id" yaml:"pool_id"`
    Sponsor     string  `json:"sponsor" yaml:"from"`
    SponsorAddr string  `json:"sponsor_addr" yaml:"sponsor_addr"`
    From        string  `json:"from" yaml:"from"`
}
```

## Parameters
| Parameter           | Info                                                                          | Default |
|---------------------|-------------------------------------------------------------------------------|---------|
| `ProtectionPeriod`  | how long a Shield lasts                                                       | 21 days |
| `ShieldFeesRate`    | percentage of protected assets to be paid as fee                              | 0.769%  |
| `WithdrawPeriod`    | how long a pending withdraw sits in the queue                                 | 21 days |
| `PoolShieldLimit`   | percentage of total collateral that a single Shield can protect               | 50%     |
| `MinShieldPurchase` | smallest allowed Shield purchase amount                                       | 50 CTK  |
| `ClaimPeriod`       |                              _(currently unused)_                             | 21 days |
| `PayoutPeriod`      |                              _(currently unused)_                             | 56 days |
| `MinDeposit`        |                              _(currently unused)_                             | 100 CTK |
| `DepositRate`       |                              _(currently unused)_                             | 10%     |
| `FeesRate`          |                              _(currently unused)_                             | 1%      |
| `StakingShieldRate` | multiple of Shield's protected assets that purchaser can stake in lieu of fee | 2       |
