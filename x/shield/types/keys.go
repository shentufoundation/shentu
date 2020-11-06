package types

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of this module.
	ModuleName = "shield"

	// RouterKey is used to route messages.
	RouterKey = ModuleName

	// StoreKey is the prefix under which we store this module's data.
	StoreKey = ModuleName

	// QuerierRoute is used to handle abci_query requests.
	QuerierRoute = ModuleName

	// DefaultParamspace is the default name for parameter store.
	DefaultParamspace = ModuleName
)

var (
	ShieldAdminKey              = []byte{0x00}
	TotalCollateralKey          = []byte{0x01}
	TotalWithdrawingKey         = []byte{0x02}
	TotalShieldKey              = []byte{0x03}
	TotalClaimedKey             = []byte{0x04}
	ServiceFeesKey              = []byte{0x05}
	RemainingServiceFeesKey     = []byte{0x06}
	PoolKey                     = []byte{0x07}
	NextPoolIDKey               = []byte{0x08}
	NextPurchaseIDKey           = []byte{0x09}
	PurchaseListKey             = []byte{0x0A}
	PurchaseQueueKey            = []byte{0x0B}
	ProviderKey                 = []byte{0x0C}
	WithdrawQueueKey            = []byte{0x0D}
	LastUpdateTimeKey           = []byte{0x0E}
	GlobalStakeForShieldPoolKey = []byte{0x0F}
	StakeForShieldKey           = []byte{0x11}
	BlockServiceFeesKey         = []byte{0x12}
	OriginalStakingKey          = []byte{0x13}
	ReimbursementKey            = []byte{0x14}
)

func GetTotalCollateralKey() []byte {
	return TotalCollateralKey
}

func GetTotalWithdrawingKey() []byte {
	return TotalWithdrawingKey
}

func GetTotalShieldKey() []byte {
	return TotalShieldKey
}

func GetTotalClaimedKey() []byte {
	return TotalClaimedKey
}

func GetServiceFeesKey() []byte {
	return ServiceFeesKey
}

func GetBlockServiceFeesKey() []byte {
	return BlockServiceFeesKey
}

func GetRemainingServiceFeesKey() []byte {
	return RemainingServiceFeesKey
}

// GetPoolKey gets the key for the pool identified by pool ID.
func GetPoolKey(id uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, id)
	return append(PoolKey, b...)
}

// GetShieldAdminKey gets the key for the shield admin.
func GetShieldAdminKey() []byte {
	return ShieldAdminKey
}

// GetNextPoolIDKey gets the key for the next pool ID.
func GetNextPoolIDKey() []byte {
	return NextPoolIDKey
}

// GetNextPurchaseIDKey gets the key for the next pool ID.
func GetNextPurchaseIDKey() []byte {
	return NextPurchaseIDKey
}

// GetPurchaseTxHashKey gets the key for a purchase.
func GetPurchaseListKey(id uint64, purchaser sdk.AccAddress) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	return append(PurchaseListKey, append(bz, purchaser.Bytes()...)...)
}

// GetProviderKey gets the key for the delegator's tracker.
func GetProviderKey(addr sdk.AccAddress) []byte {
	return append(ProviderKey, addr...)
}

// GetWithdrawCompletionTimeKey gets a withdraw queue key,
// which is obtained from the completion time.
func GetWithdrawCompletionTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(WithdrawQueueKey, bz...)
}

// GetPurchaseExpirationTimeKey gets a withdraw queue key,
// which is obtained from the expiration time.
func GetPurchaseExpirationTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(PurchaseQueueKey, bz...)
}

// GetLastUpdateTimeKey gets the key for the last update time.
func GetLastUpdateTimeKey() []byte {
	return LastUpdateTimeKey
}

func GetGlobalStakeForShieldPoolKey() []byte {
	return GlobalStakeForShieldPoolKey
}

func GetStakeForShieldKey(poolID uint64, purchaser sdk.AccAddress) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, poolID)
	return append(StakeForShieldKey, append(bz, purchaser...)...)
}

func GetOriginalStakingKey(purchaseID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, purchaseID)
	return append(OriginalStakingKey, bz...)
}

// GetReimbursementKey gets the key for a reimbursement.
func GetReimbursementKey(proposalID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, proposalID)
	return append(ReimbursementKey, bz...)
}
