package keeper

// import (
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/require"

// 	abci "github.com/tendermint/tendermint/abci/types"
// 	"github.com/tendermint/tendermint/crypto/secp256k1"
// 	"github.com/tendermint/tendermint/libs/log"

// 	"github.com/cosmos/cosmos-sdk/codec"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/x/auth"
// 	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
// 	"github.com/cosmos/cosmos-sdk/x/bank"
// 	"github.com/cosmos/cosmos-sdk/x/params"
// 	"github.com/cosmos/cosmos-sdk/x/slashing"
// 	"github.com/cosmos/cosmos-sdk/x/staking"

// 	"github.com/certikfoundation/shentu/common"
// 	"github.com/certikfoundation/shentu/common/tests"
// 	"github.com/certikfoundation/shentu/x/bank"
// 	"github.com/certikfoundation/shentu/x/cert"
// 	"github.com/certikfoundation/shentu/x/cvm/types"
// 	distr "github.com/certikfoundation/shentu/x/distribution"
// )

// var (
// 	uCTKAmount = sdk.NewInt(1005).MulRaw(common.MicroUnit)

// 	Addrs = []sdk.AccAddress{
// 		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
// 		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
// 		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
// 		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
// 		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
// 	}
// )

// type TestInput struct {
// 	Cdc           *codec.Codec
// 	Ctx           sdk.Context
// 	CvmKeeper     Keeper
// 	AccountKeeper authkeeper.AccountKeeper
// 	BankKeeper    bank.Keeper
// 	DistrKeeper   *TestDistrKeeper
// 	CertKeeper    cert.Keeper
// }

// type TestDistrKeeper struct {
// 	CommunityPool *sdk.Coins
// }

// func (tdk *TestDistrKeeper) FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error {
// 	coins := tdk.CommunityPool.Add(amount...)
// 	tdk.CommunityPool = &coins
// 	fmt.Println("updated: ", (*tdk.CommunityPool).String())
// 	return nil
// }

// func NewGasMeter(limit uint64) sdk.GasMeter {
// 	return sdk.NewGasMeter(limit)
// }

// func CreateTestInput(t *testing.T) TestInput {
// 	config := sdk.GetConfig()
// 	config.SetBech32PrefixForAccount(common.Bech32PrefixAccAddr, common.Bech32PrefixAccPub)
// 	config.SetBech32PrefixForValidator(common.Bech32PrefixValAddr, common.Bech32PrefixValPub)
// 	config.SetBech32PrefixForConsensusNode(common.Bech32PrefixConsAddr, common.Bech32PrefixConsPub)

// 	cdc := tests.MakeTestCodec(
// 		[]tests.CodecRegister{
// 			bank.RegisterCodec,
// 			auth.RegisterCodec,
// 			params.RegisterCodec,
// 			types.RegisterCodec,
// 			cert.RegisterCodec,
// 		},
// 	)
// 	keys := sdk.NewKVStoreKeys([]string{params.StoreKey, auth.StoreKey, types.StoreKey, cert.StoreKey, staking.StoreKey, supply.StoreKey}...)
// 	tkeys := sdk.NewTransientStoreKeys([]string{params.TStoreKey}...)
// 	db := tests.MakeTestDB()
// 	ms := tests.MakeTestStore(db, keys, tkeys)
// 	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

// 	maccPerms := map[string][]string{
// 		auth.FeeCollectorName:     nil,
// 		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
// 		staking.BondedPoolName:    {supply.Burner, supply.Staking},
// 		distr.ModuleName:          nil,
// 		types.ModuleName:          {supply.Burner, supply.Minter},
// 	}
// 	blacklistedAddrs := map[string]bool{
// 		auth.FeeCollectorName: true,
// 	}

// 	paramsKeeper := params.NewKeeper(cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
// 	accKeeper := auth.NewAccountKeeper(
// 		cdc,
// 		keys[auth.StoreKey],
// 		paramsKeeper.Subspace(auth.DefaultParamspace),
// 		auth.ProtoBaseAccount,
// 	)
// 	distrKeeper := TestDistrKeeper{&sdk.Coins{}}

// 	var cvmKeeper Keeper
// 	bankKeeper := bank.NewKeeper(
// 		accKeeper,
// 		&cvmKeeper,
// 		paramsKeeper.Subspace(bank.DefaultParamspace),
// 		blacklistedAddrs,
// 	)
// 	supplyKeeper := supply.NewKeeper(cdc, keys[supply.StoreKey], accKeeper, bankKeeper, maccPerms)
// 	totalSupply := sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, uCTKAmount.MulRaw(int64(len(Addrs)))))
// 	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

// 	stakingKeeper := staking.NewKeeper(
// 		cdc,
// 		keys[staking.StoreKey],
// 		supplyKeeper,
// 		paramsKeeper.Subspace(staking.DefaultParamspace),
// 	)
// 	genesis := staking.DefaultGenesisState()
// 	_ = staking.InitGenesis(ctx, stakingKeeper, accKeeper, supplyKeeper, genesis)

// 	slashingKeeper := slashing.NewKeeper(cdc, keys[slashing.StoreKey], stakingKeeper, paramsKeeper.Subspace(slashing.DefaultParamspace))
// 	certKeeper := cert.NewKeeper(cdc, keys[cert.StoreKey], slashingKeeper, stakingKeeper)
// 	cert.InitDefaultGenesis(ctx, certKeeper)

// 	cvmKeeper = NewKeeper(cdc, keys[types.StoreKey], accKeeper, &distrKeeper, certKeeper, paramsKeeper.Subspace(types.DefaultParamspace))

// 	for _, addr := range Addrs {
// 		_, err := bankKeeper.AddCoins(ctx, addr,
// 			sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(10000))})
// 		require.NoError(t, err)
// 	}

// 	cvmKeeper.SetGasRate(ctx, 1)
// 	for _, addr := range Addrs {
// 		acc := accKeeper.GetAccount(ctx, addr)
// 		_ = acc.SetSequence(1)
// 		accKeeper.SetAccount(ctx, acc)
// 	}
// 	ctx = ctx.WithGasMeter(NewGasMeter(10000000000000))
// 	RegisterGlobalPermissionAcc(ctx, cvmKeeper)
// 	return TestInput{cdc, ctx, cvmKeeper, accKeeper, bankKeeper, &distrKeeper, certKeeper}
// }
