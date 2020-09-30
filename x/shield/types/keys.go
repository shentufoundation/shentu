package types

import (
	"encoding/binary"
	"encoding/hex"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of this module
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
	PoolKey          = []byte{0x0}
	ShieldAdminKey   = []byte{0x1}
	NextPoolIDKey    = []byte{0x2}
	PurchaseKey      = []byte{0x3}
	ReimbursementKey = []byte{0x4}

	ParticipantKey     = []byte{0x5}
	PendingPayoutsKey  = []byte{0x6}
	WithdrawalQueueKey = []byte{0x7}
)

// GetPoolKey gets the key for the pool identified by pool ID.
func GetPoolKey(id uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, id)
	return append(PoolKey, b...)
}

// GetShieldadminKey gets the key for the shield admin.
func GetShieldAdminKey() []byte {
	return ShieldAdminKey
}

// GetNextPoolIDKey gets the key for the next pool ID.
func GetNextPoolIDKey() []byte {
	return NextPoolIDKey
}

// GetPendingPayoutsKey gets the key for pending payouts.
func GetPendingPayoutsKey(denom string) []byte {
	return append(PendingPayoutsKey, []byte(denom)...)
}

// GetPurchaseTxHashKey gets the key for a purchase.
func GetPurchaseTxHashKey(txhashStr string) []byte {
	txhash, err := hex.DecodeString(txhashStr)
	if err != nil {
		panic(err)
	}
	return append(PurchaseKey, txhash...)
}

// GetParticipantKey gets the key for the delegator's tracker.
func GetParticipantKey(addr sdk.AccAddress) []byte {
	return append(ParticipantKey, addr...)
}

// GetWithdrawalCompletionTimeKey gets a withdrawal queue key,
// which is obtained from the completion time.
func GetWithdrawalCompletionTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(WithdrawalQueueKey, bz...)
}

// GetReimbursement gets the key for a reimbursement.
func GetReimbursementKey(txhashStr string) []byte {
	txhash, err := hex.DecodeString(txhashStr)
	if err != nil {
		panic(err)
	}
	return append(ReimbursementKey, txhash...)
}
