package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// [1xx] Certifier
var (
	ErrUnqualifiedCertifier   = sdkerrors.Register(ModuleName, 101, "certifier not qualified")
	ErrCertifierAlreadyExists = sdkerrors.Register(ModuleName, 102, "certifier already added")
	ErrCertifierNotExists     = sdkerrors.Register(ModuleName, 103, "certifier does not exist")
	ErrRepeatedAlias          = sdkerrors.Register(ModuleName, 104, "the alias has been used by other certifiers")
	ErrUnqualifiedProposer    = sdkerrors.Register(ModuleName, 105, "proposer not qualified")
	ErrEmptyCertifier         = sdkerrors.Register(ModuleName, 106, "certifier address empty")
	ErrAddOrRemove            = sdkerrors.Register(ModuleName, 107, "must be `add` or `remove`")
	ErrInvalidCertifierAlias  = sdkerrors.Register(ModuleName, 108, "invalid certifier alias`")
	ErrOnlyOneCertifier       = sdkerrors.Register(ModuleName, 109, "cannot remove only certifier")
)

// [2xx] Validator
var (
	ErrRejectedValidator    = sdkerrors.Register(ModuleName, 201, "only certifiers can certify or de-certify validators")
	ErrValidatorCertified   = sdkerrors.Register(ModuleName, 202, "validator has already been certified")
	ErrValidatorUncertified = sdkerrors.Register(ModuleName, 203, "validator has not been certified")
	ErrTombstonedValidator  = sdkerrors.Register(ModuleName, 204, "validator has already been tombstoned")
	ErrMissingValidator     = sdkerrors.Register(ModuleName, 205, "validator missing from staking store")
)

// [3xx] Certificate
var (
	ErrCertificateNotExists      = sdkerrors.Register(ModuleName, 301, "certificate id does not exist")
	ErrCertificateGenesis        = sdkerrors.Register(ModuleName, 302, "invalid certificate genesis")
	ErrInvalidCertificateType    = sdkerrors.Register(ModuleName, 303, "invalid certificate type")
	ErrSourceCodeHash            = sdkerrors.Register(ModuleName, 304, "invalid source code hash")
	ErrCompiler                  = sdkerrors.Register(ModuleName, 305, "invalid compiler")
	ErrBytecodeHash              = sdkerrors.Register(ModuleName, 306, "invalid bytecode hash")
	ErrInvalidRequestContentType = sdkerrors.Register(ModuleName, 307, "invalid request content type")
	ErrUnqualifiedRevoker        = sdkerrors.Register(ModuleName, 308, "only certifiers can revoke this certificate")
)

// [4xx] Library
var (
	ErrLibraryNotExists     = sdkerrors.Register(ModuleName, 401, "library does not exist")
	ErrLibraryAlreadyExists = sdkerrors.Register(ModuleName, 402, "library already exists")
)
