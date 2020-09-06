package cert

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
	govTypes "github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/cert/internal/types"
	distr "github.com/certikfoundation/shentu/x/distribution"
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
	certKeeper Keeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()

	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	RegisterCodec(cdc)
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
	keyCert := sdk.NewKVStoreKey(StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keySlashing := sdk.NewKVStoreKey(slashing.StoreKey)
	tKeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keyDistr := sdk.NewKVStoreKey(distr.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

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
	genesis := staking.DefaultGenesisState()
	_ = staking.InitGenesis(ctx, stakingKeeper, accKeeper, supplyKeeper, genesis)
	slashingKeeper := slashing.NewKeeper(cdc, keySlashing, stakingKeeper, paramsKeeper.Subspace(slashing.DefaultParamspace))
	certKeeper := NewKeeper(cdc, keyCert, slashingKeeper, stakingKeeper)

	InitGenesis(ctx, certKeeper, DefaultGenesisState())

	for _, addr := range addrs {
		_, err := bankKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(10000))})
		require.NoError(t, err)

		certKeeper.SetCertifier(ctx, types.NewCertifier(addr, "", addr, ""))
	}

	return testInput{ctx, certKeeper}
}

type fillerStoreKey string

func (sk fillerStoreKey) String() string {
	return string(sk)
}

func (sk fillerStoreKey) Name() string {
	return string(sk)
}
