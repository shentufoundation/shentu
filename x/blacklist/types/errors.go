package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	// ErrAddressBlacklisted is returned when a blacklisted address attempts
	// to sign (or otherwise authorize) a transaction.
	ErrAddressBlacklisted = errorsmod.Register(ModuleName, 2, "address is blacklisted")

	// ErrEntryNotFound is returned when operating on an address that is not
	// on the blacklist.
	ErrEntryNotFound = errorsmod.Register(ModuleName, 3, "blacklist entry not found")

	// ErrInvalidCategory is returned for unrecognized entry categories.
	ErrInvalidCategory = errorsmod.Register(ModuleName, 4, "invalid category")

	// ErrModuleAccount is returned when attempting to blacklist a module
	// account.
	ErrModuleAccount = errorsmod.Register(ModuleName, 6, "cannot blacklist a module account")
)
