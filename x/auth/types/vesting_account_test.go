package types_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	"github.com/certikfoundation/shentu/v2/common"
	"github.com/certikfoundation/shentu/v2/simapp"
	"github.com/certikfoundation/shentu/v2/x/auth/types"
)

var (
	denom    = common.MicroCTKDenom
	unlocker = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	pubkeys  = []cryptotypes.PubKey{
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
	}
)

func TestManualVestingAcc(t *testing.T) {
	origAmt := sdk.NewInt(1000)
	origCoins := sdk.Coins{sdk.NewCoin(denom, origAmt)}

	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})

	// Set up an MVA with all its base coins vesting
	simapp.AddTestAddrsFromPubKeys(app, ctx, pubkeys, origAmt)
	ba := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(pubkeys[0].Address()))
	bva := authvesting.NewBaseVestingAccount(ba, origCoins, 0)
	mva := types.NewManualVestingAccountRaw(bva, sdk.NewCoins(), unlocker)

	// Test GetVestedCoins and GetVestingCoins
	now := tmtime.Now()
	vestedCoins := mva.GetVestedCoins(now)
	require.Nil(t, vestedCoins)
	vestingCoins := mva.GetVestingCoins(now)
	require.Equal(t, origCoins, vestingCoins)
	lockedCoins := mva.LockedCoins(now)
	require.Equal(t, vestingCoins, lockedCoins)

	coinToUnlock := sdk.NewCoin(denom, sdk.NewInt(700))
	mva.VestedCoins = mva.VestedCoins.Add(coinToUnlock)

	vestedCoins = mva.GetVestedCoins(now)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 700)}, vestedCoins)
	vestingCoins = mva.GetVestingCoins(now)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 300)}, vestingCoins)
	lockedCoins = mva.LockedCoins(now)
	require.Equal(t, vestingCoins, lockedCoins)

	// Test JSON (un)marshal
	bz, err := json.Marshal(mva)
	require.NoError(t, err)

	var a types.ManualVestingAccount
	require.NoError(t, json.Unmarshal(bz, &a))
	require.Equal(t, mva.String(), a.String())

	// Set up an MVA with 300 out of 1000 base coin vesting
	origVesting := sdk.Coins{sdk.NewInt64Coin(denom, 300)}

	ba2 := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(pubkeys[1].Address()))
	bva2 := authvesting.NewBaseVestingAccount(ba2, origVesting, 0)
	mva2 := types.NewManualVestingAccountRaw(bva2, sdk.NewCoins(), unlocker)
	app.AccountKeeper.SetAccount(ctx, mva2)

	lockedCoins = mva2.LockedCoins(now)
	require.Equal(t, origVesting, lockedCoins)

	// Test SpendableCoins
	spendableCoins := app.BankKeeper.SpendableCoins(ctx, mva2.GetAddress())
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 700)}, spendableCoins)

	// Test SpendableCoins after unlocking 150
	coinToUnlock = sdk.NewCoin(denom, sdk.NewInt(150))
	mva2.VestedCoins = mva2.VestedCoins.Add(coinToUnlock)
	app.AccountKeeper.SetAccount(ctx, mva2)

	spendableCoins = app.BankKeeper.SpendableCoins(ctx, mva2.GetAddress())
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 850)}, spendableCoins)

	// TODO: Test delegation, undelegation, genesis validation
}

func TestTrackDelegation(t *testing.T) {
	now := tmtime.Now()

	origCoins := sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(100))}
	ba := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(pubkeys[2].Address()))
	bva := authvesting.NewBaseVestingAccount(ba, origCoins, 0)

	// require the ability to delegate all vesting coins
	mva := types.NewManualVestingAccountRaw(bva, sdk.NewCoins(), unlocker)
	mva.TrackDelegation(now, origCoins, origCoins)
	require.Equal(t, origCoins, mva.DelegatedVesting) // 100uctk
	require.Empty(t, mva.DelegatedFree)               // 0uctk

	// require the ability to delegate all vested coins
	bva = authvesting.NewBaseVestingAccount(ba, origCoins, 0)
	mva = types.NewManualVestingAccountRaw(bva, origCoins, unlocker)
	mva.TrackDelegation(now, origCoins, origCoins)
	require.Empty(t, mva.DelegatedVesting)         // 0uctk
	require.Equal(t, origCoins, mva.DelegatedFree) // 100uctk

	// create account with 50 vested and 50 vesting coins
	bva = authvesting.NewBaseVestingAccount(ba, origCoins, 0)
	mva = types.NewManualVestingAccountRaw(bva, sdk.NewCoins(), unlocker)
	mva.LockedCoins(now)
	coinToUnlock := sdk.NewCoin(denom, sdk.NewInt(50))
	mva.VestedCoins = mva.VestedCoins.Add(coinToUnlock)
	// require the ability to delegate all vesting coins (50%) and all vested coins (50%)
	mva.TrackDelegation(now, origCoins, sdk.Coins{sdk.NewInt64Coin(denom, 50)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 50)}, mva.DelegatedVesting) // 50uctk
	require.Empty(t, mva.DelegatedFree)                                            // 0uctk

	mva.TrackDelegation(now.Add(12*time.Hour), origCoins, sdk.Coins{sdk.NewInt64Coin(denom, 50)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 50)}, mva.DelegatedVesting) // 50uctk
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 50)}, mva.DelegatedFree)    // 50uctk

	// require panic when delegation amount is zero or not enough funds
	bva = authvesting.NewBaseVestingAccount(ba, origCoins, 0)
	mva = types.NewManualVestingAccountRaw(bva, sdk.NewCoins(), unlocker)
	require.Panics(t, func() {
		mva.TrackDelegation(now, origCoins, sdk.Coins{sdk.NewInt64Coin(denom, 1000000)})
	})
}

func TestTrackUndelegation(t *testing.T) {
	now := tmtime.Now()

	origCoins := sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(100))}
	ba := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(pubkeys[2].Address()))
	bva := authvesting.NewBaseVestingAccount(ba, origCoins, 0)

	// require the ability to undelegate all vesting coins
	mva := types.NewManualVestingAccountRaw(bva, sdk.NewCoins(), unlocker)
	mva.TrackDelegation(now, origCoins, origCoins)
	mva.TrackUndelegation(origCoins)
	require.Empty(t, mva.DelegatedFree)
	require.Empty(t, mva.DelegatedVesting)

	// require the ability to undelegate all vested coins
	bva = authvesting.NewBaseVestingAccount(ba, origCoins, 0)
	mva = types.NewManualVestingAccountRaw(bva, origCoins, unlocker)
	mva.TrackDelegation(now, origCoins, origCoins)
	mva.TrackUndelegation(origCoins)
	require.Empty(t, mva.DelegatedFree)
	require.Empty(t, mva.DelegatedVesting)

	// create account with equal vested and vesting coins
	bva = authvesting.NewBaseVestingAccount(ba, origCoins, 0)
	mva = types.NewManualVestingAccountRaw(bva, sdk.NewCoins(), unlocker)
	mva.LockedCoins(now)
	coinToUnlock := sdk.NewCoin(denom, sdk.NewInt(50))
	mva.VestedCoins = mva.VestedCoins.Add(coinToUnlock)
	// vest 50% and delegate to two validators
	mva.TrackDelegation(now, origCoins, sdk.Coins{sdk.NewInt64Coin(denom, 50)})
	mva.TrackDelegation(now, origCoins, sdk.Coins{sdk.NewInt64Coin(denom, 50)})

	// undelegate from one validator that got slashed 50%
	mva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(denom, 25)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 25)}, mva.DelegatedFree)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 50)}, mva.DelegatedVesting)

	// undelegate from the other validator that did not get slashed
	mva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(denom, 50)})
	require.Empty(t, mva.DelegatedFree)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 25)}, mva.DelegatedVesting)
}
