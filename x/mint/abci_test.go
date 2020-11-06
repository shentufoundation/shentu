package mint

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cosmosDistr "github.com/cosmos/cosmos-sdk/x/distribution"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/certikfoundation/shentu/x/bank"
	"github.com/certikfoundation/shentu/x/cert"
	"github.com/certikfoundation/shentu/x/distribution"
	"github.com/certikfoundation/shentu/x/gov"
	"github.com/certikfoundation/shentu/x/shield"
	"github.com/certikfoundation/shentu/x/slashing"
	"github.com/certikfoundation/shentu/x/staking"
	"github.com/certikfoundation/shentu/x/upgrade"
)

var (
	addrs = []sdk.AccAddress{
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
	}
)

type testInput struct {
	ctx sdk.Context
	k   Keeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	params.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	distribution.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	return cdc
}

type TestDistrKeeper struct {
	CommunityPool *sdk.Coins
}

func (tdk *TestDistrKeeper) FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error {
	coins := tdk.CommunityPool.Add(amount...)
	tdk.CommunityPool = &coins
	fmt.Println("updated: ", (*tdk.CommunityPool).String())
	return nil
}

func (tdk *TestDistrKeeper) GetFeePool(ctx sdk.Context) cosmosDistr.FeePool {
	commPool := sdk.DecCoins{}
	for _, coin := range *tdk.CommunityPool {
		decCoin := sdk.NewDecCoin(coin.Denom, coin.Amount)
		commPool = commPool.Add(decCoin)
	}
	return cosmosDistr.FeePool{commPool}
}

type fillerStoreKey string

func (sk fillerStoreKey) String() string {
	return string(sk)
}

func (sk fillerStoreKey) Name() string {
	return string(sk)
}

func createTestInput(t *testing.T) testInput {
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keyCert := sdk.NewKVStoreKey(cert.StoreKey)
	keySupp := sdk.NewKVStoreKey(supply.StoreKey)
	keyShield := sdk.NewKVStoreKey(shield.StoreKey)
	keySlashing := sdk.NewKVStoreKey(slashing.StoreKey)
	keyGov := sdk.NewKVStoreKey(gov.StoreKey)
	keyMint := sdk.NewKVStoreKey(StoreKey)

	cdc := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupp, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMint, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyCert, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySlashing, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyShield, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyGov, sdk.StoreTypeIAVL, db)
	require.NoError(t, ms.LoadLatestVersion())

	// module account permissions
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		distribution.ModuleName:   nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		shield.ModuleName:         {supply.Burner},
		gov.ModuleName:            {supply.Burner, supply.Minter},
	}

	blacklistedAddrs := map[string]bool{
		auth.FeeCollectorName: true,
	}

	paramsKeeper := params.NewKeeper(cdc, keyParams, tKeyParams)
	accKeeper := auth.NewAccountKeeper(
		cdc,
		keyMint,
		paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)
	bankKeeper := bank.NewBaseKeeper(
		accKeeper,
		paramsKeeper.Subspace(bank.DefaultParamspace),
		blacklistedAddrs,
	)

	supplyKeeper := supply.NewKeeper(cdc, keySupp, accKeeper, bankKeeper, maccPerms)
	coins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000))}
	supplyKeeper.SetSupply(ctx, supply.NewSupply(coins))
	stakingKeeper := staking.NewKeeper(cdc, keyStaking, &supplyKeeper, paramsKeeper.Subspace(staking.DefaultParamspace))
	stakingKeeper.SetParams(ctx, staking.DefaultParams())
	bondedPool := stakingKeeper.GetBondedPool(ctx)
	err := bondedPool.SetCoins(coins)
	require.Nil(t, err)
	supplyKeeper.SetModuleAccount(ctx, bondedPool)
	distrKeeper := TestDistrKeeper{&sdk.Coins{}}
	slashingKeeper := slashing.NewKeeper(cdc, keySlashing, stakingKeeper, paramsKeeper.Subspace(slashing.DefaultParamspace))
	certKeeper := cert.NewKeeper(cdc, keyCert, slashingKeeper, stakingKeeper)
	govKeeper := gov.Keeper{}
	shieldKeeper := shield.NewKeeper(cdc, keyShield, accKeeper, stakingKeeper, &govKeeper, supplyKeeper, paramsKeeper.Subspace(shield.DefaultParamSpace))

	upgradeKeeper := upgrade.NewKeeper(map[int64]bool{}, fillerStoreKey(""), cdc)
	rtr := govTypes.NewRouter().
		AddRoute(gov.RouterKey, types.ProposalHandler)
	govKeeper = gov.NewKeeper(
		cdc,
		keyGov,
		paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable()),
		supplyKeeper,
		stakingKeeper,
		certKeeper,
		shieldKeeper,
		upgradeKeeper,
		rtr,
	)
	Keeper := NewKeeper(
		cdc,
		keyMint,
		paramsKeeper.Subspace(mint.DefaultParamspace),
		stakingKeeper,
		&supplyKeeper,
		&distrKeeper,
		shieldKeeper,
		auth.FeeCollectorName,
	)

	for _, addr := range addrs {
		_, err := bankKeeper.AddCoins(ctx, addr,
			sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000))})
		require.NoError(t, err)
	}

	return testInput{ctx, Keeper}
}

func TestBeginBlocker(t *testing.T) {
	testInput := createTestInput(t)
	ctx := testInput.ctx
	k := testInput.k
	p := mint.DefaultParams()
	k.SetParams(ctx, p)
	type args struct {
		minter mint.Minter
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"normal", args{
				mint.Minter{
					Inflation:        sdk.NewDecWithPrec(12, 2),
					AnnualProvisions: sdk.NewDecWithPrec(7, 2)},
			},
		},
		{
			"zero inflation", args{
				mint.Minter{
					Inflation:        sdk.NewDecWithPrec(0, 2),
					AnnualProvisions: sdk.NewDecWithPrec(0, 2)},
			},
		},
		{
			"hundred inflation", args{
				mint.Minter{
					Inflation:        sdk.NewDecWithPrec(100, 2),
					AnnualProvisions: sdk.NewDecWithPrec(100, 2)},
			},
		},
	}
	for _, tt := range tests {
		k.SetMinter(ctx, tt.args.minter)
		t.Run(tt.name, func(t *testing.T) {
			BeginBlocker(ctx, k)
		})
	}
}
