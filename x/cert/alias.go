package cert

import (
	"github.com/certikfoundation/shentu/x/cert/client"
	"github.com/certikfoundation/shentu/x/cert/internal/keeper"
	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

const (
	RouterKey                   = types.RouterKey
	StoreKey                    = types.StoreKey
	ModuleName                  = types.ModuleName
	Add                         = types.Add
	Remove                      = types.Remove
	CertificateTypeCompilation  = types.CertificateTypeCompilation
	MaxTimestamp                = keeper.MaxTimestamp
	ProposalTypeCertifierUpdate = types.ProposalTypeCertifierUpdate
)

var (
	// function aliases
	NewKeeper                   = keeper.NewKeeper
	RegisterCodec               = types.RegisterCodec
	GetGenesisStateFromAppState = types.GetGenesisStateFromAppState
	NewCertifier                = types.NewCertifier
	DefaultGenesisState         = types.DefaultGenesisState
	NewGeneralCertificate       = types.NewGeneralCertificate
	NewCompilationCertificate   = types.NewCompilationCertificate
	NewCertifierUpdateProposal  = types.NewCertifierUpdateProposal

	// variable aliases
	ProposalHandler           = client.ProposalHandler
	ErrUnqualifiedCertifier   = types.ErrUnqualifiedCertifier
	ErrRepeatedAlias          = types.ErrRepeatedAlias
	ErrCertifierAlreadyExists = types.ErrCertifierAlreadyExists
)

type (
	Keeper                  = keeper.Keeper
	GenesisState            = types.GenesisState
	CertifierUpdateProposal = types.CertifierUpdateProposal
	Certifier               = types.Certifier
	Certifiers              = types.Certifiers
	Certificate             = types.Certificate
	Validator               = types.Validator
	Platform                = types.Platform
	Library                 = types.Library
	AddOrRemove             = types.AddOrRemove
)
