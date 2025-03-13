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

// Keeper defines the bounty keeper
type Keeper struct {
	cdc codec.BinaryCodec // The codec for binary encoding/decoding

	certKeeper types.CertKeeper
	authKeeper types.AccountKeeper
	bankKeeper types.BankKeeper
	authority  string

	storeService corestoretypes.KVStoreService // The store service for accessing KV stores

	Schema collections.Schema

	// State
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
	// ProofsByTheorem key: ProofID+TheoremID | value: none used (index key for proofs by theorem index)
	ProofsByTheorem collections.Map[collections.Pair[uint64, string], []byte]
	// TheoremProofList key: TheoremID | value: ProofID
	TheoremProof collections.Map[uint64, string]
	// ActiveTheoremsQueue key: EndTime+TheoremID | value: TheoremID
	ActiveTheoremsQueue collections.Map[collections.Pair[time.Time, uint64], uint64]
	// ActiveProofsQueue key: EndTime+ProofID | value: Proof
	ActiveProofsQueue collections.Map[collections.Pair[time.Time, string], types.Proof]
}

// NewKeeper creates a new bounty Keeper instance
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

	k := Keeper{
		cdc:                 cdc,
		certKeeper:          ck,
		authKeeper:          ak,
		bankKeeper:          bk,
		authority:           authority,
		storeService:        storeService,
		Params:              collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		TheoremID:           collections.NewSequence(sb, types.TheoremIDKey, "theorem_id"),
		Theorems:            collections.NewMap(sb, types.TheoremKeyPrefix, "theorems", collections.Uint64Key, codec.CollValue[types.Theorem](cdc)),
		Grants:              collections.NewMap(sb, types.GrantKeyPrefix, "grants", collections.PairKeyCodec(collections.Uint64Key, sdk.LengthPrefixedAddressKey(sdk.AccAddressKey)), codec.CollValue[types.Grant](cdc)),
		Rewards:             collections.NewMap(sb, types.RewardKeyPrefix, "rewards", sdk.LengthPrefixedAddressKey(sdk.AccAddressKey), codec.CollValue[types.Reward](cdc)),
		Deposits:            collections.NewMap(sb, types.DepositKeyPrefix, "deposits", collections.PairKeyCodec(collections.StringKey, sdk.LengthPrefixedAddressKey(sdk.AccAddressKey)), codec.CollValue[types.Deposit](cdc)),
		Proofs:              collections.NewMap(sb, types.ProofKeyPrefix, "proofs", collections.StringKey, codec.CollValue[types.Proof](cdc)),
		ProofsByTheorem:     collections.NewMap(sb, types.ProofByTheoremPrefix, "proofs_by_theorem", collections.PairKeyCodec(collections.Uint64Key, collections.StringKey), collections.BytesValue),
		TheoremProof:        collections.NewMap(sb, types.TheoremProofPrefix, "theorem_proof", collections.Uint64Key, collections.StringValue),
		ActiveTheoremsQueue: collections.NewMap(sb, types.ActiveTheoremQueueKey, "active_theorems_queue", collections.PairKeyCodec(sdk.TimeKey, collections.Uint64Key), collections.Uint64Value),
		ActiveProofsQueue:   collections.NewMap(sb, types.ActiveProofQueueKey, "active_proofs_queue", collections.PairKeyCodec(sdk.TimeKey, collections.StringKey), codec.CollValue[types.Proof](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}
