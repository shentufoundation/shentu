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

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/auth/types"
)

var (
	denom    = common.MicroCTKDenom
	unlocker = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	pubkeys  = []cryptotypes.PubKey{
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
