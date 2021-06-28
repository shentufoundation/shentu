package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	nfttypes "github.com/irisnet/irismod/modules/nft/types"
)

const (
	errUnqualifiedCertifier uint32 = iota + 101
	errCertificateNotExists
	errCertificateGenesis
	errInvalidDenomID
	errUnqualifiedRevoker
)

var (
	ErrAdminNotFound = sdkerrors.Register(nfttypes.ModuleName, 13, "nft admin not found")
)

var (
	ErrUnqualifiedCertifier = sdkerrors.Register(ModuleName, errUnqualifiedCertifier, "certifier not qualified")
	ErrCertificateNotExists = sdkerrors.Register(ModuleName, errCertificateNotExists, "certificate id does not exist")
	ErrCertificateGenesis   = sdkerrors.Register(ModuleName, errCertificateGenesis, "invalid certificate genesis")
	ErrInvalidDenomID       = sdkerrors.Register(ModuleName, errInvalidDenomID, "invalid certificate denom ID")
	ErrUnqualifiedRevoker   = sdkerrors.Register(ModuleName, errUnqualifiedRevoker, "only certifiers can revoke this certificate")
)
