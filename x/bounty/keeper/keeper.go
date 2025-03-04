package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
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
	// Grants key: TheoremID+Grantor | value: Grant
	Grants collections.Map[collections.Pair[uint64, sdk.AccAddress], types.Grant]
	// Deposits key: ProofID+Depositor | value: Deposit
	Deposits collections.Map[collections.Pair[string, sdk.AccAddress], types.Deposit]
	// Rewards key: address | value: Reward
	Rewards collections.Map[sdk.AccAddress, types.Reward]
	// Proofs key: ProofID | value: Proof
	Proofs collections.Map[string, types.Proof]
	// TheoremProofList key: TheoremID | value: ProofID
	TheoremProof collections.Map[uint64, string]
	// ActiveTheoremsQueue key: EndTime+TheoremID | value: TheoremID
	ActiveTheoremsQueue collections.Map[collections.Pair[time.Time, uint64], uint64]
	// ActiveProofsQueue key: EndTime+ProofID | value: Proof
	ActiveProofsQueue collections.Map[collections.Pair[time.Time, string], types.Proof]
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
		cdc:                 cdc,
		certKeeper:          ck,
		authKeeper:          ak,
		bankKeeper:          bk,
		authority:           authority,
		storeService:        storeService,
		Params:              collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		TheoremID:           collections.NewSequence(sb, types.TheoremIDKey, "theorem_id"),
		Theorems:            collections.NewMap(sb, types.TheoremsKeyKeyPrefix, "theorems", collections.Uint64Key, codec.CollValue[types.Theorem](cdc)),
		Grants:              collections.NewMap(sb, types.GrantsKeyPrefix, "grants", collections.PairKeyCodec(collections.Uint64Key, sdk.LengthPrefixedAddressKey(sdk.AccAddressKey)), codec.CollValue[types.Grant](cdc)),
		Rewards:             collections.NewMap(sb, types.RewardsKeyPrefix, "reward", sdk.LengthPrefixedAddressKey(sdk.AccAddressKey), codec.CollValue[types.Reward](cdc)),
		Deposits:            collections.NewMap(sb, types.DepositsKeyPrefix, "deposits", collections.PairKeyCodec(collections.StringKey, sdk.LengthPrefixedAddressKey(sdk.AccAddressKey)), codec.CollValue[types.Deposit](cdc)),
		Proofs:              collections.NewMap(sb, types.ProofsKeyPrefix, "proofs", collections.StringKey, codec.CollValue[types.Proof](cdc)),
		TheoremProof:        collections.NewMap(sb, types.TheoremProofPrefix, "theorem_proof", collections.Uint64Key, collections.StringValue),
		ActiveTheoremsQueue: collections.NewMap(sb, types.ActiveTheoremQueuePrefix, "active_theorems_queue", collections.PairKeyCodec(sdk.TimeKey, collections.Uint64Key), collections.Uint64Value),
		ActiveProofsQueue:   collections.NewMap(sb, types.HashLockProofQueuePrefix, "active_proofs_queue", collections.PairKeyCodec(sdk.TimeKey, collections.StringKey), codec.CollValue[types.Proof](cdc)),
	}
}
