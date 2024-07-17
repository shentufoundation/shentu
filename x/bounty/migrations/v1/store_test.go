package v1_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/shentufoundation/shentu/v2/x/bounty"
	v1 "github.com/shentufoundation/shentu/v2/x/bounty/migrations/v1"
	bountytypes "github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func Test_MigratelStore(t *testing.T) {
	bKey := sdk.NewKVStoreKey(bountytypes.StoreKey)
	ctx := testutil.DefaultContext(bKey, sdk.NewTransientStoreKey("transient_test"))
	cdc := moduletestutil.MakeTestEncodingConfig(bounty.AppModuleBasic{}).Codec

	store := ctx.KVStore(bKey)

	findings := []struct {
		finding bountytypes.Finding
	}{
		{
			bountytypes.Finding{
				ProgramId:        uuid.New().String(),
				FindingId:        uuid.New().String(),
				Title:            "title",
				Description:      "desc",
				ProofOfConcept:   "poc",
				FindingHash:      uuid.New().String(),
				SubmitterAddress: "addr",
				SeverityLevel:    1,
				Status:           1,
				Detail:           "detail",
				PaymentHash:      uuid.New().String(),
				CreateTime:       time.Now(),
			},
		},
	}

	for _, test := range findings {
		bz, err := cdc.Marshal(&test.finding)
		require.NoError(t, err)
		store.Set(bountytypes.GetFindingKey(test.finding.FindingId), bz)
	}

	findingStore := prefix.NewStore(store, bountytypes.FindingKey)
	iter := findingStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var finding bountytypes.Finding
		err := cdc.Unmarshal(iter.Value(), &finding)
		require.NoError(t, err)
		t.Log(finding.FindingId, finding.FindingHash, finding.Status, finding.Title, finding.Description)
	}

	err := v1.MigrateStore(ctx, bKey, cdc)
	require.NoError(t, err)

	iter = findingStore.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var finding bountytypes.Finding
		err := cdc.Unmarshal(iter.Value(), &finding)
		require.NoError(t, err)
		t.Log(finding.FindingId, finding.FindingHash, finding.Status, finding.Title, finding.Description)
	}
}
