package types

import (
	errorsmod "cosmossdk.io/errors"
)

// Error Code Enums
const (
	errUnqualifiedCertifier uint32 = iota + 101
	errCertifierAlreadyExists
	errCertifierNotExists
	errUnqualifiedProposer
	errEmptyCertifier
	errAddOrRemove
	errOnlyOneCertifier
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

// [1xx] Certifier
var (
	ErrUnqualifiedCertifier   = errorsmod.Register(ModuleName, errUnqualifiedCertifier, "certifier not qualified")
	ErrCertifierAlreadyExists = errorsmod.Register(ModuleName, errCertifierAlreadyExists, "certifier already added")
	ErrCertifierNotExists     = errorsmod.Register(ModuleName, errCertifierNotExists, "certifier does not exist")
	ErrUnqualifiedProposer    = errorsmod.Register(ModuleName, errUnqualifiedProposer, "proposer not qualified")
	ErrEmptyCertifier         = errorsmod.Register(ModuleName, errEmptyCertifier, "certifier address empty")
	ErrAddOrRemove            = errorsmod.Register(ModuleName, errAddOrRemove, "must be `add` or `remove`")
	ErrOnlyOneCertifier       = errorsmod.Register(ModuleName, errOnlyOneCertifier, "cannot remove only certifier")
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
