package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	nfttypes "github.com/irisnet/irismod/modules/nft/types"
)

const (
	errAdminNotFound uint32 = iota + 13
	errUnqualifiedCertifier
	errCertificateGenesis
	errInvalidDenomID
	errUnqualifiedRevoker
)

var (
	ErrAdminNotFound        = sdkerrors.Register(nfttypes.ModuleName, errAdminNotFound, "nft admin not found")
	ErrUnqualifiedCertifier = sdkerrors.Register(ModuleName, errUnqualifiedCertifier, "certifier not qualified")
	ErrCertificateGenesis   = sdkerrors.Register(ModuleName, errCertificateGenesis, "invalid certificate genesis")
	ErrInvalidDenomID       = sdkerrors.Register(ModuleName, errInvalidDenomID, "invalid certificate denom ID")
	ErrUnqualifiedRevoker   = sdkerrors.Register(ModuleName, errUnqualifiedRevoker, "only certifiers can revoke this certificate")
)
