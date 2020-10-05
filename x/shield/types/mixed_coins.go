package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MixedCoins struct {
	Native  sdk.Coins
	Foreign sdk.Coins
}

func (mc MixedCoins) Add(a MixedCoins) MixedCoins {
	native := mc.Native.Add(a.Native...)
	foreign := mc.Foreign.Add(a.Foreign...)
	return MixedCoins{
		Native:  native,
		Foreign: foreign,
	}
}

func (mc MixedCoins) String() string {
	return append(mc.Native, mc.Foreign...).String()
}

type MixedDecCoins struct {
	Native  sdk.DecCoins
	Foreign sdk.DecCoins
}

func InitMixedDecCoins() MixedDecCoins {
	return MixedDecCoins{
		Native:  sdk.DecCoins{},
		Foreign: sdk.DecCoins{},
	}
}

func NewMixedDecCoins(native, foreign sdk.DecCoins) MixedDecCoins {
	return MixedDecCoins{
		Native:  native,
		Foreign: foreign,
	}
}

// MixedDecCoinsFromMixedCoins converts MixedCoins to MixedDecCoins.
func MixedDecCoinsFromMixedCoins(mc MixedCoins) MixedDecCoins {
	return MixedDecCoins{
		Native:  sdk.NewDecCoinsFromCoins(mc.Native...),
		Foreign: sdk.NewDecCoinsFromCoins(mc.Foreign...),
	}
}

// Add adds two MixedDecCoins type coins together.
func (mdc MixedDecCoins) Add(a MixedDecCoins) MixedDecCoins {
	return MixedDecCoins{
		Native:  mdc.Native.Add(a.Native...),
		Foreign: mdc.Foreign.Add(a.Foreign...),
	}
}

// MulDec multiplies native and foreign coins by a decimal.
func (mdc MixedDecCoins) MulDec(d sdk.Dec) MixedDecCoins {
	return MixedDecCoins{
		Native:  mdc.Native.MulDec(d),
		Foreign: mdc.Foreign.MulDec(d),
	}
}

// QuoDec divides native and foreign coins by a decimal.
func (mdc MixedDecCoins) QuoDec(d sdk.Dec) MixedDecCoins {
	return MixedDecCoins{
		Native:  mdc.Native.QuoDec(d),
		Foreign: mdc.Foreign.QuoDec(d),
	}
}

func (mdc MixedDecCoins) String() string {
	return append(mdc.Native, mdc.Foreign...).String()
}
