package types

import (
	errorsmod "cosmossdk.io/errors"
)

// Error Code Enums
const (
	errUnqualifiedCertifier uint32 = iota + 101
	errCertifierAlreadyExists
	errCertifierNotExists
	errRepeatedAlias
	errUnqualifiedProposer
	errEmptyCertifier
	errAddOrRemove
	errInvalidCertifierAlias
	errOnlyOneCertifier
)

const (
	errRejectedValidator uint32 = iota + 201
	errValidatorCertified
	errValidatorUncertified
	errTombstonedValidator
	errMissingValidator
)

const (
	errCertificateNotExists uint32 = iota + 301
	errCertificateGenesis
	errInvalidCertificateType
	errSourceCodeHash
	errCompiler
	errBytecodeHash
	errInvalidRequestContentType
	errUnqualifiedRevoker
)

const (
	errLibraryNotExists uint32 = iota + 401
	errLibraryAlreadyExists
)

// [1xx] Certifier
var (
	ErrUnqualifiedCertifier   = errorsmod.Register(ModuleName, errUnqualifiedCertifier, "certifier not qualified")
	ErrCertifierAlreadyExists = errorsmod.Register(ModuleName, errCertifierAlreadyExists, "certifier already added")
	ErrCertifierNotExists     = errorsmod.Register(ModuleName, errCertifierNotExists, "certifier does not exist")
	ErrRepeatedAlias          = errorsmod.Register(ModuleName, errRepeatedAlias, "the alias has been used by other certifiers")
	ErrUnqualifiedProposer    = errorsmod.Register(ModuleName, errUnqualifiedProposer, "proposer not qualified")
	ErrEmptyCertifier         = errorsmod.Register(ModuleName, errEmptyCertifier, "certifier address empty")
	ErrAddOrRemove            = errorsmod.Register(ModuleName, errAddOrRemove, "must be `add` or `remove`")
	ErrInvalidCertifierAlias  = errorsmod.Register(ModuleName, errInvalidCertifierAlias, "invalid certifier alias`")
	ErrOnlyOneCertifier       = errorsmod.Register(ModuleName, errOnlyOneCertifier, "cannot remove only certifier")
)

// [2xx] Validator
var (
	ErrRejectedValidator    = errorsmod.Register(ModuleName, errRejectedValidator, "only certifiers can certify or de-certify validators")
	ErrValidatorCertified   = errorsmod.Register(ModuleName, errValidatorCertified, "validator has already been certified")
	ErrValidatorUncertified = errorsmod.Register(ModuleName, errValidatorUncertified, "validator has not been certified")
	ErrTombstonedValidator  = errorsmod.Register(ModuleName, errTombstonedValidator, "validator has already been tombstoned")
	ErrMissingValidator     = errorsmod.Register(ModuleName, errMissingValidator, "validator missing from staking store")
)

// [3xx] Certificate
var (
	ErrCertificateNotExists      = errorsmod.Register(ModuleName, errCertificateNotExists, "certificate id does not exist")
	ErrCertificateGenesis        = errorsmod.Register(ModuleName, errCertificateGenesis, "invalid certificate genesis")
	ErrInvalidCertificateType    = errorsmod.Register(ModuleName, errInvalidCertificateType, "invalid certificate type")
	ErrSourceCodeHash            = errorsmod.Register(ModuleName, errSourceCodeHash, "invalid source code hash")
	ErrCompiler                  = errorsmod.Register(ModuleName, errCompiler, "invalid compiler")
	ErrBytecodeHash              = errorsmod.Register(ModuleName, errBytecodeHash, "invalid bytecode hash")
	ErrInvalidRequestContentType = errorsmod.Register(ModuleName, errInvalidRequestContentType, "invalid request content type")
	ErrUnqualifiedRevoker        = errorsmod.Register(ModuleName, errUnqualifiedRevoker, "only certifiers can revoke this certificate")
)

// [4xx] Library
var (
	ErrLibraryNotExists     = errorsmod.Register(ModuleName, errLibraryNotExists, "library does not exist")
	ErrLibraryAlreadyExists = errorsmod.Register(ModuleName, errLibraryAlreadyExists, "library already exists")
)
