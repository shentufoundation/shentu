package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper - bounty keeper
type Keeper struct {
	// The codec for binary encoding/decoding.
	cdc codec.BinaryCodec

	//authKeeper authtypes.AccountKeeper
	certKeeper types.CertKeeper
	authKeeper types.AccountKeeper
	bankKeeper types.BankKeeper
	authority  string

	// The (unexposed) keys used to access the stores from the Context.
	storeService corestoretypes.KVStoreService

	Schema collections.Schema

	Params    collections.Item[types.Params]
	TheoremID collections.Sequence
	// Theorems key: TheoremID | value: Theorem
	Theorems collections.Map[uint64, types.Theorem]
	// Grants key: Grantor+TheoremID | value: Grant
	Grants collections.Map[collections.Pair[sdk.AccAddress, uint64], types.Grant]
	// Proofs key: ProofsID | value: Proof
	Proofs              collections.Map[string, types.Proof]
	ActiveTheoremsQueue collections.Map[collections.Pair[time.Time, uint64], uint64]
	// MVP don't add
	//InactiveTheoremsQueue collections.Map[collections.Pair[time.Time, uint64], uint64]
	HashLockedProofsQueue collections.Map[collections.Pair[time.Time, string], string]
	DetailProofsQueue     collections.Map[collections.Pair[time.Time, string], string]
}

// NewKeeper creates a new Keeper object
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService corestoretypes.KVStoreService,
	ak types.AccountKeeper,
	ck types.CertKeeper,
	bk types.BankKeeper,
	authority string,
) Keeper {

	if _, err := ak.AddressCodec().StringToBytes(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	sb := collections.NewSchemaBuilder(storeService)

	return Keeper{
		cdc:                   cdc,
		certKeeper:            ck,
		authKeeper:            ak,
		bankKeeper:            bk,
		authority:             authority,
		storeService:          storeService,
		Params:                collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		TheoremID:             collections.NewSequence(sb, types.TheoremIDKey, "theorem_id"),
		Theorems:              collections.NewMap(sb, types.TheoremsKeyKeyPrefix, "theorems", collections.Uint64Key, codec.CollValue[types.Theorem](cdc)),
		Grants:                collections.NewMap(sb, types.GrantsKeyPrefix, "grants", collections.PairKeyCodec(sdk.LengthPrefixedAddressKey(sdk.AccAddressKey), collections.Uint64Key), codec.CollValue[types.Grant](cdc)), // nolint: staticcheck // sdk.LengthPrefixedAddressKey is needed to retain state compatibility
		Proofs:                collections.NewMap(sb, types.ProofsKeyPrefix, "proofs", collections.StringKey, codec.CollValue[types.Proof](cdc)),
		ActiveTheoremsQueue:   collections.NewMap(sb, types.ActiveTheoremQueuePrefix, "active_theorems_queue", collections.PairKeyCodec(sdk.TimeKey, collections.Uint64Key), collections.Uint64Value),
		HashLockedProofsQueue: collections.NewMap(sb, types.HashLockProofQueuePrefix, "hash_lock_proofs_queue", collections.PairKeyCodec(sdk.TimeKey, collections.StringKey), collections.StringValue),
		DetailProofsQueue:     collections.NewMap(sb, types.DetailProofQueuePrefix, "detail_proofs_queue", collections.PairKeyCodec(sdk.TimeKey, collections.StringKey), collections.StringValue),
	}
}
