package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type Pool struct {
	// TODO: better field names
	Community []Collateral
	Coverage  sdk.Coins
	Premium   MixedCoins
	CertiK    Collateral
	Sponsor   string
	Active bool
}

func NewPool(coverage sdk.Coins, deposit MixedCoins, sponsor string) Pool {
	return Pool{
		Coverage: coverage,
		Premium:  deposit,
		Sponsor:  sponsor,
		Active: true,
	}
}

type Collateral struct {
	Provider sdk.AccAddress
	Amount   sdk.Coins
}

type MixedCoins struct {
	Native  sdk.Coins
	Foreign sdk.Coins
}

func (mc MixedCoins) String() string {
	return append(mc.Native, mc.Foreign...).String()
}

func (mc MixedCoins) Add(a MixedCoins) MixedCoins {
	native := mc.Native.Add(a.Native...)
	foreign := mc.Foreign.Add(a.Foreign...)
	return MixedCoins{
		Native: native,
		Foreign: foreign,
	}
}