package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	nfttypes "github.com/irisnet/irismod/modules/nft/types"
)

const (
	errUnqualifiedCertifier uint32 = iota + 101
	errCertificateNotExists
	errCertificateGenesis
	errInvalidCertificateType
	errSourceCodeHash
	errCompiler
	errBytecodeHash
	errInvalidRequestContentType
	errUnqualifiedRevoker
)

var (
	ErrAdminNotFound = sdkerrors.Register(nfttypes.ModuleName, 13, "nft admin not found")
)

var (
	ErrUnqualifiedCertifier      = sdkerrors.Register(ModuleName, errUnqualifiedCertifier, "certifier not qualified")
	ErrCertificateNotExists      = sdkerrors.Register(ModuleName, errCertificateNotExists, "certificate id does not exist")
	ErrCertificateGenesis        = sdkerrors.Register(ModuleName, errCertificateGenesis, "invalid certificate genesis")
	ErrInvalidCertificateType    = sdkerrors.Register(ModuleName, errInvalidCertificateType, "invalid certificate type")
	ErrSourceCodeHash            = sdkerrors.Register(ModuleName, errSourceCodeHash, "invalid source code hash")
	ErrCompiler                  = sdkerrors.Register(ModuleName, errCompiler, "invalid compiler")
	ErrBytecodeHash              = sdkerrors.Register(ModuleName, errBytecodeHash, "invalid bytecode hash")
	ErrInvalidRequestContentType = sdkerrors.Register(ModuleName, errInvalidRequestContentType, "invalid request content type")
	ErrUnqualifiedRevoker        = sdkerrors.Register(ModuleName, errUnqualifiedRevoker, "only certifiers can revoke this certificate")
)
