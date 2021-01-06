package types_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/auth/types"
)

var (
	denom    = "stake"
	unlocker = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	pubkeys  = []crypto.PubKey{
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
	}
)

func TestManualVestingAcc(t *testing.T) {
	origAmt := sdk.NewInt(1000)
	origCoins := sdk.Coins{sdk.NewCoin(denom, origAmt)}

	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})

	// Account setup
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

	coinToUnlock := sdk.NewCoin(denom, sdk.NewInt(700))
	mva.VestedCoins = mva.VestedCoins.Add(coinToUnlock)

	vestedCoins = mva.GetVestedCoins(now)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 700)}, vestedCoins)
	vestingCoins = mva.GetVestingCoins(now)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 300)}, vestingCoins)

	// Test JSON (un)marshal
	bz, err := json.Marshal(mva)
	require.NoError(t, err)

	var a types.ManualVestingAccount
	require.NoError(t, json.Unmarshal(bz, &a))
	require.Equal(t, mva.String(), a.String())

	// New account setup
	origCoins = sdk.Coins{sdk.NewInt64Coin(denom, 1000)}
	origVesting := sdk.Coins{sdk.NewInt64Coin(denom, 300)}

	//ba2 := authtypes.NewBaseAccount(addrs[0], origCoins, nil, 0, 0)
	ba2 := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(pubkeys[0].Address()))
	bva2 := authvesting.NewBaseVestingAccount(ba2, origVesting, 0)
	mva2 := types.NewManualVestingAccountRaw(bva2, sdk.NewCoins(), unlocker)

	// Test SpendableCoins
	//spendableCoins := mva2.SpendableCoins(now)
	spendableCoins := app.BankKeeper.SpendableCoins(ctx, mva2.GetAddress())
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 700)}, spendableCoins)

	coinToUnlock = sdk.NewCoin(denom, sdk.NewInt(150))
	mva2.VestedCoins = mva2.VestedCoins.Add(coinToUnlock)

	//spendableCoins = mva2.SpendableCoins(now)
	spendableCoins = app.BankKeeper.SpendableCoins(ctx, mva2.GetAddress())
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 850)}, spendableCoins)

	// spendableCoins := mva2.SpendableCoins(now)
	// require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 700)}, spendableCoins)

	// coinToUnlock = sdk.NewCoin(denom, sdk.NewInt(150))
	// mva2.VestedCoins = mva2.VestedCoins.Add(coinToUnlock)

	// spendableCoins = mva2.SpendableCoins(now)
	// require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 850)}, spendableCoins)

	// TODO: Test delegation, undelegation, genesis validation
}
