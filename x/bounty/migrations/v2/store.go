package v2

import (
	"cosmossdk.io/core/store"
	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func MigrateStore(ctx sdk.Context, storeService store.KVStoreService, cdc codec.BinaryCodec) error {
	store := runtime.KVStoreAdapter(storeService.OpenKVStore(ctx))
	findingStore := prefix.NewStore(store, types.FindingKey)

	iter := findingStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var finding types.Finding
		err := cdc.Unmarshal(iter.Value(), &finding)
		if err != nil {
			return err
		}

		newFinding := types.Finding{
			ProgramId:        finding.ProgramId,
			FindingId:        finding.FindingId,
			FindingHash:      finding.FindingHash,
			SubmitterAddress: finding.SubmitterAddress,
			SeverityLevel:    finding.SeverityLevel,
			Status:           finding.Status,
			PaymentHash:      finding.PaymentHash,
			CreateTime:       finding.CreateTime,
		}

		// set the new proposal with proposer
		bz, err := cdc.Marshal(&newFinding)
		if err != nil {
			panic(err)
		}
		store.Set(types.GetFindingKey(finding.FindingId), bz)
	}

	return nil
}
