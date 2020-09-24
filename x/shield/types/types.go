package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type Pool struct {
	// TODO: better field names
	Community []Collateral
	Coverage  sdk.Coins
	Premium   MixedCoins
	CertiK    Collateral
	Sponsor   string
}

func NewPool(coverage sdk.Coins, deposit MixedCoins, sponsor string) Pool {
	return Pool{
		Coverage: coverage,
		Premium:  deposit,
		Sponsor:  sponsor,
	}
}

type Collateral struct {
	Provider sdk.AccAddress
	Amount   sdk.Coins
}

// ForeignCoins separates sdk.Coins to shield foreign coins
type ForeignCoins sdk.Coins

type MixedCoins struct {
	Native  sdk.Coins
	Foreign ForeignCoins
}

func (mc MixedCoins) String() string {
	return append(mc.Native, mc.Foreign...).String()
}
