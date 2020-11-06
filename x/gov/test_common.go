package gov

import (
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
	"github.com/cosmos/cosmos-sdk/x/bank"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/cert"
	distr "github.com/certikfoundation/shentu/x/distribution"
	"github.com/certikfoundation/shentu/x/gov/internal/keeper"
	"github.com/certikfoundation/shentu/x/gov/internal/types"
	"github.com/certikfoundation/shentu/x/shield"
	"github.com/certikfoundation/shentu/x/staking"
	"github.com/certikfoundation/shentu/x/upgrade"
)

var (
	uCTKAmount = sdk.NewInt(1005).MulRaw(common.MicroUnit)

	addrs = []sdk.AccAddress{
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
	}
)

type testInput struct {
	ctx        sdk.Context
	govKeeper  keeper.Keeper
	bankKeeper bank.Keeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()

	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	cert.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	params.RegisterCodec(cdc)
	govTypes.RegisterCodec(cdc)

	return cdc
}

func createTestInput(t *testing.T) testInput {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyCert := sdk.NewKVStoreKey(cert.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keySlashing := sdk.NewKVStoreKey(slashing.StoreKey)
	tKeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keyDistr := sdk.NewKVStoreKey(distr.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyGov := sdk.NewKVStoreKey(StoreKey)
	keyShield := sdk.NewKVStoreKey(shield.StoreKey)

	cdc := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyCert, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyGov, sdk.StoreTypeIAVL, db)
	require.NoError(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		distr.ModuleName:          nil,
		ModuleName:                {supply.Burner, supply.Minter},
	}
	blacklistedAddrs := map[string]bool{
		// TODO
	}

	paramsKeeper := params.NewKeeper(cdc, keyParams, tKeyParams)
	accKeeper := auth.NewAccountKeeper(
		cdc,
		keyAcc,
		paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)
	bankKeeper := bank.NewBaseKeeper(
		accKeeper,
		paramsKeeper.Subspace(bank.DefaultParamspace),
		blacklistedAddrs,
	)

	supplyKeeper := supply.NewKeeper(cdc, keySupply, accKeeper, bankKeeper, maccPerms)
	totalSupply := sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, uCTKAmount.MulRaw(int64(len(addrs)))))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	stakingKeeper := staking.NewKeeper(
		cdc,
		keyStaking,
		supplyKeeper,
		paramsKeeper.Subspace(staking.DefaultParamspace),
	)
	genesis := stakingTypes.DefaultGenesisState()
	_ = staking.InitGenesis(ctx, stakingKeeper.Keeper, accKeeper, supplyKeeper, genesis)

	distrKeeper := distr.NewKeeper(
		cdc,
		keyDistr,
		paramsKeeper.Subspace(distr.DefaultParamspace),
		stakingKeeper,
		supplyKeeper,
		auth.FeeCollectorName,
		blacklistedAddrs,
	)
	distrKeeper.SetFeePool(ctx, distrTypes.InitialFeePool())

	slashingKeeper := slashing.NewKeeper(cdc, keySlashing, stakingKeeper, paramsKeeper.Subspace(slashing.DefaultParamspace))
	certKeeper := cert.NewKeeper(cdc, keyCert, slashingKeeper, stakingKeeper)
	govKeeper := keeper.Keeper{}
	shieldKeeper := shield.NewKeeper(cdc, keyShield, accKeeper, stakingKeeper, &govKeeper, supplyKeeper, paramsKeeper.Subspace(shield.DefaultParamSpace))

	upgradeKeeper := upgrade.NewKeeper(map[int64]bool{}, fillerStoreKey(""), cdc)

	rtr := govTypes.NewRouter().
		AddRoute(RouterKey, types.ProposalHandler)

	govKeeper = keeper.NewKeeper(
		cdc,
		keyGov,
		paramsKeeper.Subspace(govTypes.DefaultParamspace).WithKeyTable(ParamKeyTable()),
		supplyKeeper,
		stakingKeeper,
		certKeeper,
		shieldKeeper,
		upgradeKeeper,
		rtr,
	)

	InitGenesis(ctx, govKeeper, supplyKeeper, DefaultGenesisState())

	for _, addr := range addrs {
		_, err := bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(10000))})
		require.NoError(t, err)
	}

	return testInput{ctx, govKeeper, bankKeeper}
}

type fillerStoreKey string

func (sk fillerStoreKey) String() string {
	return string(sk)
}

func (sk fillerStoreKey) Name() string {
	return string(sk)
}
