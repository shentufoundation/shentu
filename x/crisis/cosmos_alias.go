package crisis

import (
	"github.com/cosmos/cosmos-sdk/x/crisis"
)

const (
	ModuleName        = crisis.ModuleName
	DefaultParamspace = crisis.DefaultParamspace
)

var (
	// functions aliases
	DefaultGenesisState = crisis.DefaultGenesisState
	NewKeeper           = crisis.NewKeeper
	NewCosmosAppModule  = crisis.NewAppModule

	// variable aliases
	CosmosModuleCdc = crisis.ModuleCdc
)

type (
	GenesisState         = crisis.GenesisState
	MsgVerifyInvariant   = crisis.MsgVerifyInvariant
	InvarRoute           = crisis.InvarRoute
	Keeper               = crisis.Keeper
	CosmosAppModule      = crisis.AppModule
	CosmosAppModuleBasic = crisis.AppModuleBasic
)
