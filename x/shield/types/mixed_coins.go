package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MixedCoins defines the struct for mixed coins with native and foreign coins.
type MixedCoins struct {
	Native  sdk.Coins
	Foreign sdk.Coins
}

// Add implements the add method of MixedCoins.
func (mc MixedCoins) Add(a MixedCoins) MixedCoins {
	native := mc.Native.Add(a.Native...)
	foreign := mc.Foreign.Add(a.Foreign...)
	return MixedCoins{
		Native:  native,
		Foreign: foreign,
	}
}

// String implements the Stringer for MixedCoins.
func (mc MixedCoins) String() string {
	return append(mc.Native, mc.Foreign...).String()
}

// MixedDecCoins defines the struct for mixed coins in decimal with native and foreign decimal coins.
type MixedDecCoins struct {
	Native  sdk.DecCoins `json:"native" yaml:"native"`
	Foreign sdk.DecCoins `json:"foreign" yaml:"foreign"`
}

// InitMixedDecCoins initialize an empty mixed decimal coins instance.
func InitMixedDecCoins() MixedDecCoins {
	return MixedDecCoins{
		Native:  sdk.DecCoins{},
		Foreign: sdk.DecCoins{},
	}
}

// NewMixedDecCoins returns a new mixed decimal coins instance.
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

// String implements the Stringer for MixedDecCoins.
func (mdc MixedDecCoins) String() string {
	return append(mdc.Native, mdc.Foreign...).String()
}
