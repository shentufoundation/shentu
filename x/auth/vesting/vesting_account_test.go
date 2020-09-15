package vesting

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtime "github.com/tendermint/tendermint/types/time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

var (
	denom  = "uctk"
	denom2 = "uctk2"
	addr   = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
)

func TestGetVestedCoinsTriggeredVestingAcc(t *testing.T) {
	// Account setup
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)
	periods := types.Periods{
		types.Period{Length: int64(12 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 500)}},
		types.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 250)}},
		types.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 250)}},
	}

	origCoins := sdk.Coins{sdk.NewInt64Coin(denom, 1000)}
	bacc := authtypes.NewBaseAccount(addr, origCoins, nil, 0, 0)

	baseVestingAccount, err := authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)

	// Initialize without triggering vesting
	tva := NewTriggeredVestingAccountRaw(baseVestingAccount, now.Unix(), periods, false)

	// No coins are vested at any point in time because the vesting has not been triggered.
	vestedCoins := tva.GetVestedCoins(now)
	require.Nil(t, vestedCoins)
	vestedCoins = tva.GetVestedCoins(endTime)
	require.Nil(t, vestedCoins)
	vestedCoins = tva.GetVestedCoins(now.Add(6 * time.Hour))
	require.Nil(t, vestedCoins)
	vestedCoins = tva.GetVestedCoins(now.Add(12 * time.Hour))
	require.Nil(t, vestedCoins)
	vestedCoins = tva.GetVestedCoins(now.Add(15 * time.Hour))
	require.Nil(t, vestedCoins)
	vestedCoins = tva.GetVestedCoins(now.Add(18 * time.Hour))
	require.Nil(t, vestedCoins)
	vestedCoins = tva.GetVestedCoins(now.Add(48 * time.Hour))
	require.Nil(t, vestedCoins)

	// Trigger vesting
	tva.Activated = true

	// require no coins vested at the beginning of the vesting schedule
	vestedCoins = tva.GetVestedCoins(now)
	require.Nil(t, vestedCoins)
	// require all coins vested at the end of the vesting schedule
	vestedCoins = tva.GetVestedCoins(endTime)
	require.Equal(t, origCoins, vestedCoins)
	// require no coins vested during first vesting period
	vestedCoins = tva.GetVestedCoins(now.Add(6 * time.Hour))
	require.Nil(t, vestedCoins)
	// require 50% of coins vested after period 1
	vestedCoins = tva.GetVestedCoins(now.Add(12 * time.Hour))
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 500)}, vestedCoins)
	// require period 2 coins don't vest until period is over
	vestedCoins = tva.GetVestedCoins(now.Add(15 * time.Hour))
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 500)}, vestedCoins)
	// require 75% of coins vested after period 2
	vestedCoins = tva.GetVestedCoins(now.Add(18 * time.Hour))
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 750)}, vestedCoins)
	// require 100% of coins vested
	vestedCoins = tva.GetVestedCoins(now.Add(48 * time.Hour))
	require.Equal(t, origCoins, vestedCoins)
}

func TestGetVestingCoinsTriggeredVestingAcc(t *testing.T) {
	// Account setup
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)
	periods := types.Periods{
		types.Period{Length: int64(12 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 500)}},
		types.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 250)}},
		types.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 250)}},
	}

	origCoins := sdk.Coins{sdk.NewInt64Coin(denom, 1000)}
	bacc := authtypes.NewBaseAccount(addr, origCoins, nil, 0, 0)

	baseVestingAccount, err := authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)

	// Initialize without triggering vesting
	tva := NewTriggeredVestingAccountRaw(baseVestingAccount, now.Unix(), periods, false)

	// All coins are vesting at all times because the vesting has not been triggered.
	vestingCoins := tva.GetVestingCoins(now)
	require.Equal(t, origCoins, vestingCoins)
	vestingCoins = tva.GetVestingCoins(endTime)
	require.Equal(t, origCoins, vestingCoins)
	vestingCoins = tva.GetVestingCoins(now.Add(12 * time.Hour))
	require.Equal(t, origCoins, vestingCoins)
	vestingCoins = tva.GetVestingCoins(now.Add(15 * time.Hour))
	require.Equal(t, origCoins, vestingCoins)
	vestingCoins = tva.GetVestingCoins(now.Add(18 * time.Hour))
	require.Equal(t, origCoins, vestingCoins)
	vestingCoins = tva.GetVestingCoins(now.Add(48 * time.Hour))
	require.Equal(t, origCoins, vestingCoins)

	// Trigger vesting
	tva.Activated = true

	// require all coins vesting at the beginning of the vesting schedule
	vestingCoins = tva.GetVestingCoins(now)
	require.Equal(t, origCoins, vestingCoins)
	// require no coins vesting at the end of the vesting schedule
	vestingCoins = tva.GetVestingCoins(endTime)
	require.Nil(t, vestingCoins)
	// require 50% of coins vesting
	vestingCoins = tva.GetVestingCoins(now.Add(12 * time.Hour))
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 500)}, vestingCoins)
	// require 50% of coins vesting after period 1, but before period 2 completes.
	vestingCoins = tva.GetVestingCoins(now.Add(15 * time.Hour))
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 500)}, vestingCoins)
	// require 25% of coins vesting after period 2
	vestingCoins = tva.GetVestingCoins(now.Add(18 * time.Hour))
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 250)}, vestingCoins)
	// require 0% of coins vesting after vesting complete
	vestingCoins = tva.GetVestingCoins(now.Add(48 * time.Hour))
	require.Nil(t, vestingCoins)
}

func TestSpendableCoinsTriggeredVestingAcc(t *testing.T) {
	// Account setup
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)
	periods := types.Periods{
		types.Period{Length: int64(12 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 500)}},
		types.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 250)}},
		types.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 250)}},
	}

	origCoins := sdk.Coins{sdk.NewInt64Coin(denom, 1000)}
	bacc := authtypes.NewBaseAccount(addr, origCoins, nil, 0, 0)

	baseVestingAccount, err := authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)

	// Initialize without triggering vesting
	tva := NewTriggeredVestingAccountRaw(baseVestingAccount, now.Unix(), periods, false)

	spendableCoins := tva.SpendableCoins(now)
	require.Nil(t, spendableCoins)
	spendableCoins = tva.SpendableCoins(endTime)
	require.Nil(t, spendableCoins)
	spendableCoins = tva.SpendableCoins(now.Add(12 * time.Hour))
	require.Nil(t, spendableCoins)

	recvAmt := sdk.Coins{sdk.NewInt64Coin(denom, 300)}
	tva.SetCoins(tva.GetCoins().Add(recvAmt...))
	spendableCoins = tva.SpendableCoins(now.Add(12 * time.Hour))
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 300)}, spendableCoins)
	tva.SetCoins(tva.GetCoins().Sub(spendableCoins))

	spendableCoins = tva.SpendableCoins(now.Add(12 * time.Hour))
	require.Nil(t, spendableCoins)

	// Trigger vesting
	tva.Activated = true

	// require that there exist no spendable coins at the beginning of the
	// vesting schedule
	spendableCoins = tva.SpendableCoins(now)
	require.Nil(t, spendableCoins)

	// require that all original coins are spendable at the end of the vesting
	// schedule
	spendableCoins = tva.SpendableCoins(endTime)
	require.Equal(t, origCoins, spendableCoins)

	// require that all vested coins (50%) are spendable
	spendableCoins = tva.SpendableCoins(now.Add(12 * time.Hour))
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 500)}, spendableCoins)

	// receive some coins
	recvAmt = sdk.Coins{sdk.NewInt64Coin(denom, 300)}
	tva.SetCoins(tva.GetCoins().Add(recvAmt...))

	// require that all vested coins (50%) are spendable plus any received
	spendableCoins = tva.SpendableCoins(now.Add(12 * time.Hour))
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 800)}, spendableCoins)

	// spend all spendable coins
	tva.SetCoins(tva.GetCoins().Sub(spendableCoins))

	// require that no more coins are spendable
	spendableCoins = tva.SpendableCoins(now.Add(12 * time.Hour))
	require.Nil(t, spendableCoins)
}

func TestTrackDelegationTriggeredVestingAcc(t *testing.T) {
	// Account setup
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)
	periods := types.Periods{
		types.Period{Length: int64(12 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 500), sdk.NewInt64Coin(denom2, 50)}},
		types.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 250), sdk.NewInt64Coin(denom2, 25)}},
		types.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 250), sdk.NewInt64Coin(denom2, 25)}},
	}

	origCoins := sdk.Coins{sdk.NewInt64Coin(denom, 1000), sdk.NewInt64Coin(denom2, 100)}

	bacc := authtypes.NewBaseAccount(addr, origCoins, nil, 0, 0)
	bvacc, err := authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)

	// Initialize without triggering vesting
	tva := NewTriggeredVestingAccountRaw(bvacc, now.Unix(), periods, true)

	// require the ability to delegate all vesting coins
	tva.TrackDelegation(now, origCoins)
	require.Equal(t, origCoins, tva.DelegatedVesting)
	require.Equal(t, sdk.Coins{}, tva.DelegatedFree)

	// require the ability to delegate all vested coins
	bacc.SetCoins(origCoins)
	bvacc, err = authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)
	tva = NewTriggeredVestingAccountRaw(bvacc, now.Unix(), periods, true)

	tva.TrackDelegation(endTime, origCoins)
	require.Equal(t, sdk.Coins{}, tva.DelegatedVesting)
	require.Equal(t, origCoins, tva.DelegatedFree)

	// delegate half of vesting coins
	bacc.SetCoins(origCoins)
	bvacc, err = authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)
	tva = NewTriggeredVestingAccountRaw(bvacc, now.Unix(), periods, true)
	tva.TrackDelegation(now, periods[0].Amount)
	// require that all delegated coins are delegated vesting
	require.Equal(t, tva.DelegatedVesting, periods[0].Amount)
	require.Equal(t, sdk.Coins{}, tva.DelegatedFree)

	// delegate 75% of coins, split between vested and vesting
	bacc.SetCoins(origCoins)
	bvacc, err = authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)
	tva = NewTriggeredVestingAccountRaw(bvacc, now.Unix(), periods, true)
	tva.TrackDelegation(now.Add(12*time.Hour), periods[0].Amount.Add(periods[1].Amount...))
	// require that the maximum possible amount of vesting coins are chosen for delegation.
	require.Equal(t, tva.DelegatedFree, periods[1].Amount)
	require.Equal(t, tva.DelegatedVesting, periods[0].Amount)

	// require the ability to delegate all vesting coins (50%) and all vested coins (50%)
	bacc.SetCoins(origCoins)
	bvacc, err = authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)
	tva = NewTriggeredVestingAccountRaw(bvacc, now.Unix(), periods, true)
	tva.TrackDelegation(now.Add(12*time.Hour), sdk.Coins{sdk.NewInt64Coin(denom2, 50)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom2, 50)}, tva.DelegatedVesting)
	require.Equal(t, sdk.Coins{}, tva.DelegatedFree)

	tva.TrackDelegation(now.Add(12*time.Hour), sdk.Coins{sdk.NewInt64Coin(denom2, 50)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom2, 50)}, tva.DelegatedVesting)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom2, 50)}, tva.DelegatedFree)

	// require no modifications when delegation amount is zero or not enough funds
	bacc.SetCoins(origCoins)
	bvacc, err = authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)
	tva = NewTriggeredVestingAccountRaw(bvacc, now.Unix(), periods, true)
	require.Panics(t, func() {
		tva.TrackDelegation(endTime, sdk.Coins{sdk.NewInt64Coin(denom2, 1000000)})
	})
	require.Equal(t, sdk.Coins{}, tva.DelegatedVesting)
	require.Equal(t, sdk.Coins{}, tva.DelegatedFree)
}

// Comparable test as TestTrackDelegationTriggeredVestingAcc but on a TriggeredVestingAccount
// that has not been triggered.
func TestTrackDelegationTriggeredVestingAcc2(t *testing.T) {
	// Account setup
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)
	periods := types.Periods{
		types.Period{Length: int64(12 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 500), sdk.NewInt64Coin(denom2, 50)}},
		types.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 250), sdk.NewInt64Coin(denom2, 25)}},
		types.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 250), sdk.NewInt64Coin(denom2, 25)}},
	}

	origCoins := sdk.Coins{sdk.NewInt64Coin(denom, 1000), sdk.NewInt64Coin(denom2, 100)}

	bacc := authtypes.NewBaseAccount(addr, origCoins, nil, 0, 0)
	bvacc, err := authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)

	// Initialize without triggering vesting
	tva := NewTriggeredVestingAccountRaw(bvacc, now.Unix(), periods, false)

	// require the ability to delegate all vesting coins
	tva.TrackDelegation(now, origCoins)
	require.Equal(t, origCoins, tva.DelegatedVesting)
	require.Equal(t, sdk.Coins{}, tva.DelegatedFree)

	// require the ability to delegate all vested coins
	bacc.SetCoins(origCoins)
	bvacc, err = authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)
	tva = NewTriggeredVestingAccountRaw(bvacc, now.Unix(), periods, false)

	tva.TrackDelegation(endTime, origCoins)
	require.Equal(t, origCoins, tva.DelegatedVesting)
	require.Equal(t, sdk.Coins{}, tva.DelegatedFree)

	// delegate 75% of coins, split between vested and vesting
	bacc.SetCoins(origCoins)
	bvacc, err = authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)
	tva = NewTriggeredVestingAccountRaw(bvacc, now.Unix(), periods, false)
	tva.TrackDelegation(now.Add(12*time.Hour), periods[0].Amount.Add(periods[1].Amount...))
	require.Equal(t, sdk.Coins{}, tva.DelegatedFree)
	require.Equal(t, periods[0].Amount.Add(periods[1].Amount...), tva.DelegatedVesting)
}

func TestTrackUndelegationTriggeredVestingAcc(t *testing.T) {
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)
	periods := types.Periods{
		types.Period{Length: int64(12 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 500), sdk.NewInt64Coin(denom2, 50)}},
		types.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 250), sdk.NewInt64Coin(denom2, 25)}},
		types.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 250), sdk.NewInt64Coin(denom2, 25)}},
	}

	origCoins := sdk.Coins{sdk.NewInt64Coin(denom, 1000), sdk.NewInt64Coin(denom2, 100)}

	bacc := authtypes.NewBaseAccount(addr, origCoins, nil, 0, 0)
	bvacc, err := authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)

	// Initialize without triggering vesting
	tva := NewTriggeredVestingAccountRaw(bvacc, now.Unix(), periods, true)

	// require the ability to undelegate all vesting coins at the beginning of vesting
	tva.TrackDelegation(now, origCoins)
	tva.TrackUndelegation(origCoins)
	require.Equal(t, sdk.Coins{}, tva.DelegatedFree)
	require.Nil(t, tva.DelegatedVesting)

	// require the ability to undelegate all vested coins at the end of vesting
	bacc.SetCoins(origCoins)
	bvacc, err = authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)
	tva = NewTriggeredVestingAccountRaw(bvacc, now.Unix(), periods, true)

	tva.TrackDelegation(endTime, origCoins)
	tva.TrackUndelegation(origCoins)
	require.Nil(t, tva.DelegatedFree)
	require.Equal(t, sdk.Coins{}, tva.DelegatedVesting)

	// require the ability to undelegate half of coins
	bacc.SetCoins(origCoins)
	bvacc, err = authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)
	tva = NewTriggeredVestingAccountRaw(bvacc, now.Unix(), periods, true)
	tva.TrackDelegation(endTime, periods[0].Amount)
	tva.TrackUndelegation(periods[0].Amount)
	require.Nil(t, tva.DelegatedFree)
	require.Equal(t, sdk.Coins{}, tva.DelegatedVesting)

	// require no modifications when the undelegation amount is zero
	bacc.SetCoins(origCoins)
	bvacc, err = authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)
	tva = NewTriggeredVestingAccountRaw(bvacc, now.Unix(), periods, true)

	require.Panics(t, func() {
		tva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(denom2, 0)})
	})
	require.Equal(t, sdk.Coins{}, tva.DelegatedFree)
	require.Equal(t, sdk.Coins{}, tva.DelegatedVesting)

	// vest 50% and delegate to two validators
	tva = NewTriggeredVestingAccountRaw(bvacc, now.Unix(), periods, true)
	tva.TrackDelegation(now.Add(12*time.Hour), sdk.Coins{sdk.NewInt64Coin(denom2, 50)})
	tva.TrackDelegation(now.Add(12*time.Hour), sdk.Coins{sdk.NewInt64Coin(denom2, 50)})

	// undelegate from one validator that got slashed 50%
	tva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(denom2, 25)})
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom2, 25)}, tva.DelegatedFree)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom2, 50)}, tva.DelegatedVesting)

	// undelegate from the other validator that did not get slashed
	tva.TrackUndelegation(sdk.Coins{sdk.NewInt64Coin(denom2, 50)})
	require.Nil(t, tva.DelegatedFree)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom2, 25)}, tva.DelegatedVesting)
}

func TestGenesisAccountValidate(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	baseAccWithCoins := authtypes.NewBaseAccount(addr, nil, pubkey, 0, 0)
	baseAccWithCoins.SetCoins(sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 50)})
	baseVestingWithCoins, _ := authvesting.NewBaseVestingAccount(baseAccWithCoins, baseAccWithCoins.Coins, 100)

	tests := []struct {
		name   string
		acc    authexported.GenesisAccount
		expErr error
	}{
		{
			"invalid vesting period lengths",
			NewTriggeredVestingAccountRaw(
				baseVestingWithCoins,
				0, types.Periods{types.Period{Length: int64(50), Amount: sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 50)}}}, false),
			errors.New("vesting end time does not match length of all vesting periods"),
		},
		{
			"invalid vesting period amounts",
			NewTriggeredVestingAccountRaw(
				baseVestingWithCoins,
				0, types.Periods{types.Period{Length: int64(100), Amount: sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 25)}}}, false),
			errors.New("original vesting coins does not match the sum of all coins in vesting periods"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.acc.Validate()
			require.Equal(t, tt.expErr, err)
		})
	}
}

func TestTriggeredVestingAccountJSON(t *testing.T) {
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)
	periods := types.Periods{
		types.Period{Length: int64(12 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 500)}},
		types.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 250)}},
		types.Period{Length: int64(6 * 60 * 60), Amount: sdk.Coins{sdk.NewInt64Coin(denom, 250)}},
	}

	origCoins := sdk.Coins{sdk.NewInt64Coin(denom, 1000)}
	bacc := authtypes.NewBaseAccount(addr, origCoins, nil, 0, 0)

	baseVestingAccount, err := authvesting.NewBaseVestingAccount(bacc, origCoins, endTime.Unix())
	require.Nil(t, err)

	acc := NewTriggeredVestingAccountRaw(baseVestingAccount, now.Unix(), periods, false)

	bz, err := json.Marshal(acc)
	require.NoError(t, err)

	bz1, err := acc.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, string(bz1), string(bz))

	var a TriggeredVestingAccount
	require.NoError(t, json.Unmarshal(bz, &a))
	require.Equal(t, acc.String(), a.String())

	acc2 := NewTriggeredVestingAccountRaw(baseVestingAccount, now.Unix(), periods, true)

	bz2, err := json.Marshal(acc2)
	require.NoError(t, err)

	bz12, err := acc2.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, string(bz12), string(bz2))

	var a2 TriggeredVestingAccount
	require.NoError(t, json.Unmarshal(bz2, &a2))
	require.Equal(t, acc2.String(), a2.String())
}

func TestManualVestingAcc(t *testing.T) {
	// Account setup
	origCoins := sdk.Coins{sdk.NewInt64Coin(denom, 1000)}
	ba := authtypes.NewBaseAccount(addr, origCoins, nil, 0, 0)
	bva, err := authvesting.NewBaseVestingAccount(ba, origCoins, 0)
	require.Nil(t, err)

	mva := NewManualVestingAccountRaw(bva, sdk.NewCoins())

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

	ba2 := authtypes.NewBaseAccount(addr, origCoins, nil, 0, 0)
	bva2, err := authvesting.NewBaseVestingAccount(ba2, origVesting, 0)
	require.Nil(t, err)
	mva2 := NewManualVestingAccountRaw(bva2, sdk.NewCoins())

	// Test SpendableCoins
	spendableCoins := mva2.SpendableCoins(now)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 700)}, spendableCoins)
	
	coinToUnlock = sdk.NewCoin(denom, sdk.NewInt(150))
	mva2.VestedCoins = mva2.VestedCoins.Add(coinToUnlock)

	spendableCoins = mva2.SpendableCoins(now)
	require.Equal(t, sdk.Coins{sdk.NewInt64Coin(denom, 850)}, spendableCoins)

	// TODO: Test delegation, undelegation, genesis validation

}
