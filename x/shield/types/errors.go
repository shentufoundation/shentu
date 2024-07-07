package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

// Error Code Enums
const (
	errEmptySender uint32 = iota + 101
	errProviderNotFound
)

var (
	ErrEmptySender      = sdkerrors.Register(ModuleName, errEmptySender, "no sender provided")
	ErrProviderNotFound = sdkerrors.Register(ModuleName, errProviderNotFound, "provider is not found")
)
