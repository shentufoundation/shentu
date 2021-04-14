package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	ErrUnqualifiedCertifier   = sdkerrors.Register(ModuleName, errUnqualifiedCertifier, "certifier not qualified")
	ErrCertifierAlreadyExists = sdkerrors.Register(ModuleName, errCertifierAlreadyExists, "certifier already added")
	ErrCertifierNotExists     = sdkerrors.Register(ModuleName, errCertifierNotExists, "certifier does not exist")
	ErrRepeatedAlias          = sdkerrors.Register(ModuleName, errRepeatedAlias, "the alias has been used by other certifiers")
	ErrUnqualifiedProposer    = sdkerrors.Register(ModuleName, errUnqualifiedProposer, "proposer not qualified")
	ErrEmptyCertifier         = sdkerrors.Register(ModuleName, errEmptyCertifier, "certifier address empty")
	ErrAddOrRemove            = sdkerrors.Register(ModuleName, errAddOrRemove, "must be `add` or `remove`")
	ErrInvalidCertifierAlias  = sdkerrors.Register(ModuleName, errInvalidCertifierAlias, "invalid certifier alias`")
	ErrOnlyOneCertifier       = sdkerrors.Register(ModuleName, errOnlyOneCertifier, "cannot remove only certifier")
)

// [2xx] Validator
var (
	ErrRejectedValidator    = sdkerrors.Register(ModuleName, errRejectedValidator, "only certifiers can certify or de-certify validators")
	ErrValidatorCertified   = sdkerrors.Register(ModuleName, errValidatorCertified, "validator has already been certified")
	ErrValidatorUncertified = sdkerrors.Register(ModuleName, errValidatorUncertified, "validator has not been certified")
	ErrTombstonedValidator  = sdkerrors.Register(ModuleName, errTombstonedValidator, "validator has already been tombstoned")
	ErrMissingValidator     = sdkerrors.Register(ModuleName, errMissingValidator, "validator missing from staking store")
)

// [3xx] Certificate
var (
	ErrCertificateNotExists      = sdkerrors.Register(ModuleName, errCertificateNotExists, "certificate id does not exist")
	ErrCertificateGenesis        = sdkerrors.Register(ModuleName, errCertificateGenesis, "invalid certificate genesis")
	ErrInvalidCertificateType    = sdkerrors.Register(ModuleName, errInvalidCertificateType, "invalid certificate type")
	ErrSourceCodeHash            = sdkerrors.Register(ModuleName, errSourceCodeHash, "invalid source code hash")
	ErrCompiler                  = sdkerrors.Register(ModuleName, errCompiler, "invalid compiler")
	ErrBytecodeHash              = sdkerrors.Register(ModuleName, errBytecodeHash, "invalid bytecode hash")
	ErrInvalidRequestContentType = sdkerrors.Register(ModuleName, errInvalidRequestContentType, "invalid request content type")
	ErrUnqualifiedRevoker        = sdkerrors.Register(ModuleName, errUnqualifiedRevoker, "only certifiers can revoke this certificate")
)

// [4xx] Library
var (
	ErrLibraryNotExists     = sdkerrors.Register(ModuleName, errLibraryNotExists, "library does not exist")
	ErrLibraryAlreadyExists = sdkerrors.Register(ModuleName, errLibraryAlreadyExists, "library already exists")
)
