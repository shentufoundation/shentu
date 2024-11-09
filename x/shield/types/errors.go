package types

import (
	errorsmod "cosmossdk.io/errors"
)

// Error Code Enums
const (
	errEmptySender uint32 = iota + 101
	errProviderNotFound
)

var (
	ErrEmptySender      = errorsmod.Register(ModuleName, errEmptySender, "no sender provided")
	ErrProviderNotFound = errorsmod.Register(ModuleName, errProviderNotFound, "provider is not found")
)
