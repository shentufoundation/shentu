package v280

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shentufoundation/shentu/v2/common"
	"github.com/shentufoundation/shentu/v2/x/shield/types"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

const (
	DENOM           = "uctk"
	AMOUNT_PROVIDER = 10e5
	SHIELD          = 1e8
	SERVICE_FEES    = 1e3
	PURCHASE_NUM    = 8
	PURCHASE_COIN   = 1000
)

type ProviderAndAddr struct {
	addr sdk.AccAddress
	pvd  types.Provider
}

func makeProvider() ProviderAndAddr {
	_, _, addr := testdata.KeyTestPubAddr()
	newAddr, _ := common.PrefixToCertik(addr.String())
	return ProviderAndAddr{
		addr,
		types.Provider{
			Address:          newAddr,
			DelegationBonded: sdk.NewInt(10),
			Collateral:       sdk.NewInt(3),
			TotalLocked:      sdk.NewInt(20),
			Withdrawing:      sdk.NewInt(7),
			Rewards:          sdk.NewDecCoins(sdk.NewDecCoin(DENOM, sdk.NewInt(AMOUNT_PROVIDER))),
		},
	}
}

func makePurchases() (entries []types.Purchase) {
	for i := 0; i < PURCHASE_NUM; i++ {
		entries = append(entries, types.Purchase{
			PurchaseId:        rand.Uint64(),
			ProtectionEndTime: time.Now().Add(10 * time.Minute),
			DeletionTime:      time.Now().Add(time.Hour),
			Description:       "--",
			Shield:            sdk.NewInt(SHIELD),
			ServiceFees:       sdk.NewDecCoins(sdk.NewDecCoin(DENOM, sdk.NewInt(SERVICE_FEES))),
		})
	}
	return
}

func makePool() types.Pool {
	_, _, addr := testdata.KeyTestPubAddr()
	sponsorAddr, _ := common.PrefixToCertik(addr.String())

	pool := types.Pool{
		Id:          rand.Uint64(),
		Description: "for_test",
		Sponsor:     "for_test",
		SponsorAddr: sponsorAddr,
		ShieldLimit: sdk.NewInt(PURCHASE_NUM),
	}
	return pool
}

func makePurchaseList(poolId uint64) (sdk.AccAddress, types.PurchaseList) {
	_, _, addr := testdata.KeyTestPubAddr()
	purchaserAddr, _ := common.PrefixToCertik(addr.String())
	return addr, types.PurchaseList{
		PoolId:    poolId,
		Purchaser: purchaserAddr,
		Entries:   makePurchases(),
	}
}

func makePoolPurchaser() (poolPurchases []types.PoolPurchaser) {
	for i := 0; i < PURCHASE_NUM; i++ {
		_, _, addr := testdata.KeyTestPubAddr()
		purchaserAddr, _ := common.PrefixToCertik(addr.String())
		poolPurchases = append(poolPurchases, types.PoolPurchaser{
			PoolId:    rand.Uint64(),
			Purchaser: purchaserAddr,
		})
	}
	return
}

func makeReimbursement() types.Reimbursement {
	_, _, addr := testdata.KeyTestPubAddr()
	beneficiaryAddr, _ := common.PrefixToCertik(addr.String())
	return types.Reimbursement{
		Amount:      sdk.NewCoins(sdk.NewCoin(DENOM, sdk.NewInt(PURCHASE_COIN))),
		PayoutTime:  time.Now().Add(time.Hour),
		Beneficiary: beneficiaryAddr,
	}
}

func makeWithdraws() (withdraws []types.Withdraw) {
	for i := 0; i < 5; i++ {
		_, _, addr := testdata.KeyTestPubAddr()
		withdrawAddr, _ := common.PrefixToCertik(addr.String())

		withdraws = append(withdraws, types.Withdraw{
			Address:        withdrawAddr,
			Amount:         sdk.NewInt(rand.Int63()),
			CompletionTime: time.Now().Add(time.Hour),
		})
	}
	return
}

func makeShieldStaking(poolId uint64) types.ShieldStaking {
	_, _, addr := testdata.KeyTestPubAddr()
	purchaser, _ := common.PrefixToCertik(addr.String())

	return types.ShieldStaking{
		PoolId:            poolId,
		Purchaser:         purchaser,
		Amount:            sdk.NewInt(rand.Int63()),
		WithdrawRequested: sdk.NewInt(rand.Int63()),
	}
}

func TestMigrateStore(t *testing.T) {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(common.Bech32PrefixAccAddr, common.Bech32PrefixAccPub)
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	shieldKey := sdk.NewKVStoreKey("shield")
	ctx := testutil.DefaultContext(shieldKey, sdk.NewTransientStoreKey("transient_test"))
	store := ctx.KVStore(shieldKey)

	providerCases := []ProviderAndAddr{
		makeProvider(),
		makeProvider(),
	}
	for _, pc := range providerCases {
		bz := cdc.MustMarshalLengthPrefixed(&pc.pvd)
		store.Set(types.GetProviderKey(pc.addr), bz)
	}

	//set data of PurchaseList in the store
	type PLCase struct {
		poolID    uint64
		key       []byte
		purchaser string
	}

	type poolCase struct {
		poolID      uint64
		sponsorAddr string
	}
	var poolPurchaserPairs types.PoolPurchaserPairs
	plCases := make([]PLCase, 0, 5)
	poolCases := make([]poolCase, 0, 5)
	for i := 0; i < 5; i++ {
		pool := makePool()
		bz := cdc.MustMarshalLengthPrefixed(&pool)
		store.Set(types.GetPoolKey(pool.Id), bz)
		poolCases = append(poolCases, poolCase{poolID: pool.Id, sponsorAddr: pool.SponsorAddr})

		addr, purchaseList := makePurchaseList(pool.Id)
		key := types.GetPurchaseListKey(pool.Id, addr)
		bz = cdc.MustMarshalLengthPrefixed(&purchaseList)
		store.Set(key, bz)
		plCases = append(plCases, PLCase{poolID: pool.Id, key: key, purchaser: purchaseList.Purchaser})

		poolPurchaserPairs.Pairs = append(poolPurchaserPairs.Pairs, types.PoolPurchaser{
			PoolId:    pool.Id,
			Purchaser: purchaseList.Purchaser,
		})
	}

	oneHour := time.Now().Add(time.Hour)
	{
		//set PurchaseExpirationTimeKey value
		bz := cdc.MustMarshalLengthPrefixed(&poolPurchaserPairs)
		store.Set(types.GetPurchaseExpirationTimeKey(oneHour), bz)
	}

	proposalID := rand.Uint64()
	reimbursement := makeReimbursement()
	{
		bz := cdc.MustMarshalLengthPrefixed(&reimbursement)
		store.Set(types.GetReimbursementKey(proposalID), bz)
	}

	withdraws := makeWithdraws()
	{
		bz := cdc.MustMarshalLengthPrefixed(&types.Withdraws{Withdraws: withdraws})
		store.Set(types.GetWithdrawCompletionTimeKey(oneHour), bz)
	}

	shieldStaking := makeShieldStaking(plCases[0].poolID)
	{
		purchaseAddr, _ := sdk.AccAddressFromBech32(shieldStaking.Purchaser)
		bz := cdc.MustMarshalLengthPrefixed(&shieldStaking)
		store.Set(types.GetStakeForShieldKey(plCases[0].poolID, purchaseAddr), bz)

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

			shentuAddr, _ := common.PrefixToShentu(item.pvd.Address)
			require.Equal(t, pvd.Address, shentuAddr)
			dc := pvd.Rewards[0]
			require.Equal(t, sdk.NewInt64DecCoin(DENOM, AMOUNT_PROVIDER), dc)
		})
	}

	//check for PurchaseList
	for _, item := range plCases {
		t.Run("PurchaseList", func(t *testing.T) {
			var pl types.PurchaseList
			bz := store.Get(item.key)
			cdc.MustUnmarshalLengthPrefixed(bz, &pl)
			require.Equal(t, item.poolID, pl.PoolId)
			shentuPurchaser, _ := common.PrefixToShentu(item.purchaser)
			require.Equal(t, shentuPurchaser, pl.Purchaser)
			require.Equal(t, PURCHASE_NUM, len(pl.Entries))
			require.Equal(t, sdk.NewInt(SHIELD), pl.Entries[0].Shield)
			require.Equal(t,
				sdk.DecCoins{sdk.NewInt64DecCoin(DENOM, SERVICE_FEES)},
				pl.Entries[1].ServiceFees)
		})
	}

	//check for pool
	for _, item := range poolCases {
		t.Run("pool", func(t *testing.T) {
			var pool types.Pool
			bz := store.Get(types.GetPoolKey(item.poolID))
			cdc.MustUnmarshalLengthPrefixed(bz, &pool)
			shentuSponsor, _ := common.PrefixToShentu(item.sponsorAddr)
			require.Equal(t, pool.SponsorAddr, shentuSponsor)
		})
	}

	//check for Expiration PurchaseQueue
	{
		bz := store.Get(types.GetPurchaseExpirationTimeKey(oneHour))
		require.NotNil(t, bz)
		var ppPairs types.PoolPurchaserPairs
		cdc.MustUnmarshalLengthPrefixed(bz, &ppPairs)

		old := poolPurchaserPairs.Pairs
		for i, item := range ppPairs.Pairs {
			shentuAddr, _ := common.PrefixToShentu(old[i].Purchaser)
			require.Equal(t, item.Purchaser, shentuAddr)
		}
	}

	// check for Reimbursement
	{
		bz := store.Get(types.GetReimbursementKey(proposalID))
		var newReimbursement types.Reimbursement
		cdc.MustUnmarshalLengthPrefixed(bz, &newReimbursement)

		beneficiary, _ := common.PrefixToShentu(reimbursement.Beneficiary)
		require.Equal(t, beneficiary, newReimbursement.Beneficiary)
	}

	//check for Withdraws
	{
		bz := store.Get(types.GetWithdrawCompletionTimeKey(oneHour))
		var newWithdraws types.Withdraws
		cdc.MustUnmarshalLengthPrefixed(bz, &newWithdraws)

		for i, withdraw := range newWithdraws.Withdraws {
			shentuAddr, _ := common.PrefixToShentu(withdraws[i].Address)
			require.Equal(t, shentuAddr, withdraw.Address)
		}
	}

	{
		purchaser, err := sdk.AccAddressFromBech32(shieldStaking.Purchaser)
		require.NoError(t, err)

		bz := store.Get(types.GetStakeForShieldKey(shieldStaking.PoolId, purchaser))
		var purchase types.ShieldStaking
		cdc.MustUnmarshalLengthPrefixed(bz, &purchase)

		shentuPurchaser, _ := common.PrefixToShentu(shieldStaking.Purchaser)
		require.Equal(t, shentuPurchaser, purchase.Purchaser)
	}
}
