package v3

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func MigrateStore(ctx sdk.Context, storeService store.KVStoreService, cdc codec.BinaryCodec) error {
	sb := collections.NewSchemaBuilder(storeService)
	programFindings := collections.NewKeySet(sb, types.ProgramFindingListKey, "program_findings", collections.PairKeyCodec(collections.StringKey, collections.StringKey))

	kvStore := runtime.KVStoreAdapter(storeService.OpenKVStore(ctx))
	ProgramFindingsStore := prefix.NewStore(kvStore, types.ProgramFindingListKey)
	findingStore := prefix.NewStore(kvStore, types.FindingKeyPrefix)

	// delete old ProgramFindingsStore
	programIter := ProgramFindingsStore.Iterator(nil, nil)
	defer programIter.Close()
	for ; programIter.Valid(); programIter.Next() {
		ProgramFindingsStore.Delete(programIter.Key())
	}

	// migrate new ProgramFindingsStore
	findingIter := findingStore.Iterator(nil, nil)
	defer findingIter.Close()
	for ; findingIter.Valid(); findingIter.Next() {
		var finding types.Finding

		err := cdc.Unmarshal(findingIter.Value(), &finding)
		if err != nil {
			return err
		}
		err = programFindings.Set(ctx, collections.Join(finding.ProgramId, finding.FindingId))
		if err != nil {
			return err
		}
	}

	return nil
}
