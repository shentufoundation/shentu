package vesting

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtime "github.com/tendermint/tendermint/types/time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting"
)

var (
	denom  = "uctk"
	denom2 = "uctk2"
	addrs  = []sdk.AccAddress{
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
	}
)

func TestManualVestingAcc(t *testing.T) {
	// Account setup
	origCoins := sdk.Coins{sdk.NewInt64Coin(denom, 1000)}
	ba := authtypes.NewBaseAccount(addrs[0], origCoins, nil, 0, 0)
	bva, err := authvesting.NewBaseVestingAccount(ba, origCoins, 0)
	require.Nil(t, err)

	mva := NewManualVestingAccountRaw(bva, sdk.NewCoins(), addrs[1])

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

	bz1, err := mva.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, string(bz1), string(bz))

	var a ManualVestingAccount
	require.NoError(t, json.Unmarshal(bz, &a))
	require.Equal(t, mva.String(), a.String())

	// New account setup
	origCoins = sdk.Coins{sdk.NewInt64Coin(denom, 1000)}
	origVesting := sdk.Coins{sdk.NewInt64Coin(denom, 300)}

	ba2 := authtypes.NewBaseAccount(addrs[0], origCoins, nil, 0, 0)
	bva2, err := authvesting.NewBaseVestingAccount(ba2, origVesting, 0)
	require.Nil(t, err)
	mva2 := NewManualVestingAccountRaw(bva2, sdk.NewCoins(), addrs[1])

	// Test SpendableCoins
	spendableCoins := mva2.SpendableCoins(now)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 700)}, spendableCoins)

	coinToUnlock = sdk.NewCoin(denom, sdk.NewInt(150))
	mva2.VestedCoins = mva2.VestedCoins.Add(coinToUnlock)

	spendableCoins = mva2.SpendableCoins(now)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 850)}, spendableCoins)

	// TODO: Test delegation, undelegation, genesis validation
}
