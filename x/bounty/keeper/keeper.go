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

// Keeper defines the bounty module's keeper, which handles all state management and business logic
type Keeper struct {
	cdc codec.BinaryCodec // Codec for binary encoding/decoding

	// External keepers
	certKeeper types.CertKeeper
	authKeeper types.AccountKeeper
	bankKeeper types.BankKeeper
	authority  string

	storeService corestoretypes.KVStoreService // KVStore service for state persistence
	Schema       collections.Schema            // Collections schema

	// State collections
	Params collections.Item[types.Params] // Global module parameters

	// OpenBountuy
	Programs        collections.Map[string, types.Program]
	Findings        collections.Map[string, types.Finding]
	ProgramFindings collections.KeySet[collections.Pair[string, string]] // ProgramFindings key: (programID, findingID)

	// OpenMath
	TheoremID           collections.Sequence
	Theorems            collections.Map[uint64, types.Theorem]                                   // Theorems key: TheoremID | value: Theorem
	Grants              collections.Map[collections.Pair[uint64, sdk.AccAddress], types.Grant]   // Grants key: TheoremID+Grantor | value: Grant
	Deposits            collections.Map[collections.Pair[string, sdk.AccAddress], types.Deposit] // Deposits key: ProofID+Depositor | value: Deposit
	Rewards             collections.Map[sdk.AccAddress, types.Reward]                            // Rewards key: address | value: Reward
	Proofs              collections.Map[string, types.Proof]                                     // Proofs key: ProofID | value: Proof
	ProofsByTheorem     collections.Map[collections.Pair[uint64, string], []byte]                // ProofsByTheorem key: ProofID+TheoremID | value: none used (index key for proofs by theorem index)
	ActiveTheoremsQueue collections.Map[collections.Pair[time.Time, uint64], uint64]             // ActiveTheoremsQueue key: EndTime+TheoremID | value: TheoremID
	ActiveProofsQueue   collections.Map[collections.Pair[time.Time, string], types.Proof]        // ActiveProofsQueue key: EndTime+ProofID | value: Proof
}

// NewKeeper creates and initializes a new Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService corestoretypes.KVStoreService,
	ak types.AccountKeeper,
	ck types.CertKeeper,
	bk types.BankKeeper,
	authority string,
) Keeper {
	// Validate authority address
	if _, err := ak.AddressCodec().StringToBytes(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	// Initialize schema builder
	sb := collections.NewSchemaBuilder(storeService)

	// Create keeper instance with all collections
	k := Keeper{
		cdc:                 cdc,
		certKeeper:          ck,
		authKeeper:          ak,
		bankKeeper:          bk,
		authority:           authority,
		storeService:        storeService,
		Params:              collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Programs:            collections.NewMap(sb, types.ProgramKeyPrefix, "programs", collections.StringKey, codec.CollValue[types.Program](cdc)),
		Findings:            collections.NewMap(sb, types.FindingKeyPrefix, "findings", collections.StringKey, codec.CollValue[types.Finding](cdc)),
		ProgramFindings:     collections.NewKeySet(sb, types.ProgramFindingListKey, "program_findings", collections.PairKeyCodec(collections.StringKey, collections.StringKey)),
		TheoremID:           collections.NewSequence(sb, types.TheoremIDKey, "theorem_id"),
		Theorems:            collections.NewMap(sb, types.TheoremKeyPrefix, "theorems", collections.Uint64Key, codec.CollValue[types.Theorem](cdc)),
		Grants:              collections.NewMap(sb, types.GrantKeyPrefix, "grants", collections.PairKeyCodec(collections.Uint64Key, sdk.LengthPrefixedAddressKey(sdk.AccAddressKey)), codec.CollValue[types.Grant](cdc)),
		Rewards:             collections.NewMap(sb, types.RewardKeyPrefix, "rewards", sdk.LengthPrefixedAddressKey(sdk.AccAddressKey), codec.CollValue[types.Reward](cdc)),
		Deposits:            collections.NewMap(sb, types.DepositKeyPrefix, "deposits", collections.PairKeyCodec(collections.StringKey, sdk.LengthPrefixedAddressKey(sdk.AccAddressKey)), codec.CollValue[types.Deposit](cdc)),
		Proofs:              collections.NewMap(sb, types.ProofKeyPrefix, "proofs", collections.StringKey, codec.CollValue[types.Proof](cdc)),
		ProofsByTheorem:     collections.NewMap(sb, types.ProofByTheoremPrefix, "proofs_by_theorem", collections.PairKeyCodec(collections.Uint64Key, collections.StringKey), collections.BytesValue),
		ActiveTheoremsQueue: collections.NewMap(sb, types.ActiveTheoremQueueKey, "active_theorems_queue", collections.PairKeyCodec(sdk.TimeKey, collections.Uint64Key), collections.Uint64Value),
		ActiveProofsQueue:   collections.NewMap(sb, types.ActiveProofQueueKey, "active_proofs_queue", collections.PairKeyCodec(sdk.TimeKey, collections.StringKey), codec.CollValue[types.Proof](cdc)),
	}

	// Build and validate schema
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}
