package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Library is a type for certified libraries.
type Library struct {
	Address   sdk.AccAddress
	Publisher sdk.AccAddress
}

// Libraries is a collection of Library objects.
type Libraries []Library
