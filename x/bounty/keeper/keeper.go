package keeper

import (
	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// Keeper - bounty keeper
type Keeper struct {
	storeService store.KVStoreService
	cdc          codec.BinaryCodec
	paramSpace   paramtypes.Subspace

	//authKeeper authtypes.AccountKeeper
	certKeeper types.CertKeeper

	//Params       collections.Item[Params]
	//TheoremIDKey collections.Sequence
	////TheoremID|Theorem
	//Theorems collections.Map[uint64, Theorem]
	////TheoremID,Grantor|Grant
	//Grants collections.Map[collections.Pair[uint64, sdk.AccAddress], Grant]
	////TheoremID,ProofId|Proof
	//Proofs                collections.Map[collections.Pair[uint64, String], Proof]
	//ActiveTheoremsQueue   collections.Map[collections.Pair[time.Time, uint64], uint64]
	//InactiveTheoremsQueue collections.Map[collections.Pair[time.Time, uint64], uint64]
	//ProofPeriodTheorems   collections.Map[uint64, []byte]
}

// NewKeeper creates a new Keeper object
func NewKeeper(
	cdc codec.BinaryCodec, storeService store.KVStoreService, ck types.CertKeeper, paramSpace paramtypes.Subspace,
) Keeper {

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		paramSpace:   paramSpace,
		certKeeper:   ck,
	}
}
