package mint

import (
	"github.com/cosmos/cosmos-sdk/x/mint"
)

const (
	ModuleName        = mint.ModuleName
	StoreKey          = mint.StoreKey
	DefaultParamspace = mint.DefaultParamspace
)

var (
	// function aliases
	NewQuerier    = mint.NewQuerier
	InitGenesis   = mint.InitGenesis
	ExportGenesis = mint.ExportGenesis

	// variable aliases
	ModuleCdc = mint.ModuleCdc
)

type (
	CosmosAppModule      = mint.AppModule
	CosmosAppModuleBasic = mint.AppModuleBasic
	GenesisState         = mint.GenesisState
)
