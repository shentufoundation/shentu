package types_test

//
//import (
//	"encoding/json"
//	"testing"
//	"time"
//
//	sdkmath "cosmossdk.io/math"
//	"github.com/stretchr/testify/require"
//
//	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
//	tmtime "github.com/cometbft/cometbft/types/time"
//
//	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
//	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
//	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
//
//	shentuapp "github.com/shentufoundation/shentu/v2/app"
//	"github.com/shentufoundation/shentu/v2/common"
//	"github.com/shentufoundation/shentu/v2/x/auth/types"
//)
//
//var (
//	denom    = common.MicroCTKDenom
//	unlocker = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
//	pubkeys  = []cryptotypes.PubKey{
//		secp256k1.GenPrivKey().PubKey(),
//		secp256k1.GenPrivKey().PubKey(),
//		secp256k1.GenPrivKey().PubKey(),
//	}
//)
//
//func TestManualVestingAcc(t *testing.T) {
//	origAmt := sdkmath.NewInt(1000)
//	origCoins := sdk.Coins{sdk.NewCoin(denom, origAmt)}
//
//	app := shentuapp.Setup(t, false)
//	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
//
//	// Set up an MVA with all its base coins vesting
//	shentuapp.AddTestAddrsFromPubKeys(app, ctx, pubkeys, origAmt)
//	ba := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(pubkeys[0].Address()))
//	bva := authvesting.NewBaseVestingAccount(ba, origCoins, 0)
//	mva := types.NewManualVestingAccountRaw(bva, sdk.NewCoins(), unlocker)
//
//	// Test GetVestedCoins and GetVestingCoins
//	now := tmtime.Now()
//	vestedCoins := mva.GetVestedCoins(now)
//	require.Nil(t, vestedCoins)
//	vestingCoins := mva.GetVestingCoins(now)
//	require.Equal(t, origCoins, vestingCoins)
//	lockedCoins := mva.LockedCoins(now)
//	require.Equal(t, vestingCoins, lockedCoins)
//
//	coinToUnlock := sdk.NewCoin(denom, sdk.NewInt(700))
//	mva.VestedCoins = mva.VestedCoins.Add(coinToUnlock)
//
//	vestedCoins = mva.GetVestedCoins(now)
//	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 700)}, vestedCoins)
//	vestingCoins = mva.GetVestingCoins(now)
//	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 300)}, vestingCoins)
//	lockedCoins = mva.LockedCoins(now)
//	require.Equal(t, vestingCoins, lockedCoins)
//
//	// Test JSON (un)marshal
//	bz, err := json.Marshal(mva)
//	require.NoError(t, err)
//
//	var a types.ManualVestingAccount
//	require.NoError(t, json.Unmarshal(bz, &a))
//	require.Equal(t, mva.String(), a.String())
//
//	// Set up an MVA with 300 out of 1000 base coin vesting
//	origVesting := sdk.Coins{sdk.NewInt64Coin(denom, 300)}
//
//	ba2 := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(pubkeys[1].Address()))
//	bva2 := authvesting.NewBaseVestingAccount(ba2, origVesting, 0)
//	mva2 := types.NewManualVestingAccountRaw(bva2, sdk.NewCoins(), unlocker)
//	app.AccountKeeper.SetAccount(ctx, mva2)
//
//	lockedCoins = mva2.LockedCoins(now)
//	require.Equal(t, origVesting, lockedCoins)
//
//	// Test SpendableCoins
//	spendableCoins := app.BankKeeper.SpendableCoins(ctx, mva2.GetAddress())
//	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 700)}, spendableCoins)
//
//	// Test SpendableCoins after unlocking 150
//	coinToUnlock = sdk.NewCoin(denom, sdk.NewInt(150))
//	mva2.VestedCoins = mva2.VestedCoins.Add(coinToUnlock)
//	app.AccountKeeper.SetAccount(ctx, mva2)
//
//	spendableCoins = app.BankKeeper.SpendableCoins(ctx, mva2.GetAddress())
//	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 850)}, spendableCoins)
//}
//
//func TestTrackDelegation(t *testing.T) {
//	now := tmtime.Now()
//	origCoins := sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(100))}
//	ba := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(pubkeys[2].Address()))
//
//	tests := []struct {
//		name        string
//		vestedCoins sdk.Coins
//		amount      sdk.Coins
//		expVesting  sdk.Coins
//		expDfree    sdk.Coins
//		doubleDel   bool
//		willPanic   bool
//	}{
//		{
//			"require the ability to delegate all vesting coins",
//			sdk.NewCoins(),
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(100))},
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(100))},
//			sdk.NewCoins(),
//			false,
//			false,
//		},
//		{
//			"require the ability to delegate all vested coins",
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(100))},
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(100))},
//			sdk.NewCoins(),
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(100))},
//			false,
//			false,
//		},
//		{
//			"require the ability to delegate all vesting coins (50%) and all vested coins (50%)",
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(50))},
//			sdk.Coins{sdk.NewInt64Coin(denom, 50)},
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(50))},
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(50))},
//			true,
//			false,
//		},
//		{
//			"panic if attempting to delegate coins that exceed the vesting coins",
//			sdk.NewCoins(),
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(1000000))},
//			sdk.NewCoins(),
//			sdk.NewCoins(),
//			false,
//			true,
//		},
//	}
//
//	for _, tt := range tests {
//		tt := tt // pin variable
//
//		t.Run(tt.name, func(t *testing.T) {
//			// create manual vesting account with desired vested amount
//			bva := authvesting.NewBaseVestingAccount(ba, origCoins, 0)
//			mva := types.NewManualVestingAccountRaw(bva, sdk.NewCoins(), unlocker)
//			mva.LockedCoins(now)
//			mva.VestedCoins = tt.vestedCoins
//
//			if tt.willPanic {
//				require.Panics(t, func() {
//					mva.TrackDelegation(now, origCoins, tt.amount)
//				})
//			} else {
//				mva.TrackDelegation(now, origCoins, tt.amount)
//				if tt.doubleDel {
//					mva.TrackDelegation(now, origCoins, tt.amount)
//				}
//				require.Equal(t, tt.expDfree, mva.DelegatedFree)
//				require.Equal(t, tt.expVesting, mva.DelegatedVesting)
//			}
//		})
//	}
//}
//
//func TestTrackUndelegation(t *testing.T) {
//	now := tmtime.Now()
//	origCoins := sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(100))}
//	ba := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(pubkeys[2].Address()))
//
//	tests := []struct {
//		name        string
//		vestedCoins sdk.Coins
//		delAmount   sdk.Coins
//		unDelAmount sdk.Coins
//		doubleDel   bool
//	}{
//		{
//			"require the ability to undelegate all vesting coins",
//			sdk.NewCoins(),
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(100))},
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(100))},
//			false,
//		},
//		{
//			"require the ability to undelegate all vested coins",
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(100))},
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(100))},
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(100))},
//			false,
//		},
//		{
//			"undelegate 50% vested coins and 50% vesting coins",
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(50))},
//			sdk.Coins{sdk.NewInt64Coin(denom, 50)},
//			sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(50))},
//			true,
//		},
//	}
//
//	for _, tt := range tests {
//		tt := tt // pin variable
//
//		t.Run(tt.name, func(t *testing.T) {
//			// create manual vesting account with desired vested amount
//			bva := authvesting.NewBaseVestingAccount(ba, origCoins, 0)
//			mva := types.NewManualVestingAccountRaw(bva, sdk.NewCoins(), unlocker)
//			mva.LockedCoins(now)
//			mva.VestedCoins = tt.vestedCoins
//
//			mva.TrackDelegation(now, origCoins, tt.delAmount)
//			if tt.doubleDel {
//				mva.TrackDelegation(now, origCoins, tt.delAmount)
//			}
//			mva.TrackUndelegation(tt.unDelAmount)
//			if tt.doubleDel {
//				mva.TrackUndelegation(tt.unDelAmount)
//			}
//			require.Empty(t, mva.DelegatedFree)
//			require.Empty(t, mva.DelegatedVesting)
//		})
//	}
//}
//
//func TestGenesisAccountValidate(t *testing.T) {
//	initialVesting := sdk.Coins{sdk.NewInt64Coin(denom, 100)}
//	pubkey := secp256k1.GenPrivKey().PubKey()
//	addr := sdk.AccAddress(pubkey.Address())
//	ba := authtypes.NewBaseAccount(addr, pubkey, 0, 0)
//	bva := authvesting.NewBaseVestingAccount(ba, initialVesting, 0)
//
//	tests := []struct {
//		name   string
//		acc    authtypes.GenesisAccount
//		expErr bool
//	}{
//		{
//			"valid base account",
//			ba,
//			false,
//		},
//		{
//			"invalid base account",
//			authtypes.NewBaseAccount(addr, secp256k1.GenPrivKey().PubKey(), 0, 0),
//			true,
//		},
//		{
//			"valid base vesting account",
//			bva,
//			false,
//		},
//		{
//			"valid continuous vesting account",
//			types.NewManualVestingAccountRaw(bva, initialVesting, unlocker),
//			false,
//		},
//		{
//			"valid manual vesting amount with no vested coins",
//			types.NewManualVestingAccountRaw(bva, sdk.NewCoins(), unlocker),
//			false,
//		},
//		{
//			"valid manual vesting amount with vested coins",
//			types.NewManualVestingAccountRaw(bva, sdk.NewCoins(sdk.NewInt64Coin(denom, 100)), unlocker),
//			false,
//		},
//		{
//			"invalid vesting amount with invalid vested coins",
//			types.NewManualVestingAccountRaw(bva, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(101))), unlocker),
//			true,
//		},
//	}
//
//	for _, tt := range tests {
//		tt := tt
//
//		t.Run(tt.name, func(t *testing.T) {
//			t.Logf("acc: %+v", tt.acc)
//			require.Equal(t, tt.expErr, tt.acc.Validate() != nil)
//		})
//	}
//}
