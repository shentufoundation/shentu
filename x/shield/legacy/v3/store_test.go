package v3

import (
	v2 "github.com/shentufoundation/shentu/v2/x/shield/legacy/v2"
	"github.com/shentufoundation/shentu/v2/x/shield/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

const (
	DENOM           = "CTK"
	AMOUNT_PROVIDER = 10e5
	AMOUNT_FEE      = 2e8
	PURCHASER       = "for_test"
	SHIELD          = 1e8
	SERVICE_FEES    = 1e3
	PURCHASE_NUM    = 8
)

type ProviderAndAddr struct {
	addr sdk.AccAddress
	pvd  v2.Provider
}

func makeMixDecCoins(amt int64) v2.MixedDecCoins {
	return v2.MixedDecCoins{
		Native: sdk.DecCoins{
			sdk.NewInt64DecCoin(DENOM, amt),
		},
		Foreign: sdk.DecCoins{},
	}
}

func makeProvider() ProviderAndAddr {
	_, _, addr := testdata.KeyTestPubAddr()
	return ProviderAndAddr{
		addr,
		v2.Provider{
			Address:          addr.String(),
			DelegationBonded: sdk.NewInt(10),
			Collateral:       sdk.NewInt(3),
			TotalLocked:      sdk.NewInt(20),
			Withdrawing:      sdk.NewInt(7),
			Rewards:          makeMixDecCoins(AMOUNT_PROVIDER),
		},
	}
}

func makePurchases() (entries []v2.Purchase) {
	for i := 0; i < PURCHASE_NUM; i++ {
		entries = append(entries, v2.Purchase{
			PurchaseId:        rand.Uint64(),
			ProtectionEndTime: time.Now().Add(10 * time.Minute),
			DeletionTime:      time.Now().Add(time.Hour),
			Description:       "--",
			Shield:            sdk.NewInt(SHIELD),
			ServiceFees:       makeMixDecCoins(SERVICE_FEES),
		})
	}
	return
}

func makePurchaseList(poolId uint64) (sdk.AccAddress, v2.PurchaseList) {
	_, _, addr := testdata.KeyTestPubAddr()
	return addr, v2.PurchaseList{
		PoolId:    poolId,
		Purchaser: PURCHASER,
		Entries:   makePurchases(),
	}
}

func TestMigrateStore(t *testing.T) {
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	shieldKey := sdk.NewKVStoreKey("shield")
	ctx := testutil.DefaultContext(shieldKey, sdk.NewTransientStoreKey("transient_test"))
	store := ctx.KVStore(shieldKey)

	//set data of v2.Provider in the store
	providerCases := []ProviderAndAddr{
		makeProvider(),
		makeProvider(),
	}
	for _, pc := range providerCases {
		bz := cdc.MustMarshalLengthPrefixed(&pc.pvd)
		store.Set(types.GetProviderKey(pc.addr), bz)
	}

	//set data of *ServiceFees in the store
	mc := makeMixDecCoins(AMOUNT_FEE)
	bz := cdc.MustMarshalLengthPrefixed(&mc)
	feeKeys := [][]byte{
		types.GetServiceFeesKey(),
		types.GetBlockServiceFeesKey(),
		types.GetRemainingServiceFeesKey(),
	}
	for _, fk := range feeKeys {
		store.Set(fk, bz)
	}

	//set data of PurchaseList in the store
	type PLCase struct {
		poolID uint64
		key    []byte
	}
	plCases := make([]PLCase, 0, 5)
	for i := 0; i < 5; i++ {
		poolID := 92688 + uint64(i)
		addr, purchaseList := makePurchaseList(poolID)
		key := types.GetPurchaseListKey(poolID, addr)
		bz := cdc.MustMarshalLengthPrefixed(&purchaseList)
		store.Set(key, bz)
		plCases = append(plCases, PLCase{poolID, key})
	}

	err := MigrateStore(ctx, shieldKey, cdc)
	require.NoError(t, err)

	//check for Provider
	for _, pc := range providerCases {
		item := pc
		t.Run("provider", func(t *testing.T) {
			var pvd types.Provider
			bz := store.Get(types.GetProviderKey(item.addr))
			cdc.MustUnmarshalLengthPrefixed(bz, &pvd)
			addr, err := sdk.AccAddressFromBech32(pvd.Address)
			require.NoError(t, err)
			require.Equal(t, item.addr, addr)
			dc := pvd.Rewards[0]
			require.Equal(t, sdk.NewInt64DecCoin(DENOM, AMOUNT_PROVIDER), dc)
		})
	}

	//check for *serviceFees
	for _, fk := range feeKeys {
		item := fk
		t.Run("Fee", func(t *testing.T) {
			var fees types.Fees
			bz := store.Get(item)
			cdc.MustUnmarshalLengthPrefixed(bz, &fees)
			require.Equal(t, sdk.NewInt64DecCoin(DENOM, AMOUNT_FEE), fees.Fees[0])
		})
	}

	//check for PurchaseList
	for _, plc := range plCases {
		item := plc
		t.Run("PurchaseList", func(t *testing.T) {
			var pl types.PurchaseList
			bz := store.Get(item.key)
			cdc.MustUnmarshalLengthPrefixed(bz, &pl)
			require.Equal(t, item.poolID, pl.PoolId)
			require.Equal(t, PURCHASER, pl.Purchaser)
			require.Equal(t, PURCHASE_NUM, len(pl.Entries))
			require.Equal(t, sdk.NewInt(SHIELD), pl.Entries[0].Shield)
			require.Equal(t,
				sdk.DecCoins{sdk.NewInt64DecCoin(DENOM, SERVICE_FEES)},
				pl.Entries[1].ServiceFees)
		})
	}
}
