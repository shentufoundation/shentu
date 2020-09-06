package slashing

import (
	"github.com/cosmos/cosmos-sdk/x/slashing"
)

const (
	ModuleName        = slashing.ModuleName
	StoreKey          = slashing.StoreKey
	DefaultParamspace = slashing.DefaultParamspace
)

var (
	// functions aliases
	NewCosmosAppModule = slashing.NewAppModule
	NewKeeper          = slashing.NewKeeper

	// variable aliases
	CosmosModuleCdc = slashing.ModuleCdc
)

type (
	GenesisState            = slashing.GenesisState
	MissedBlock             = slashing.MissedBlock
	MsgUnjail               = slashing.MsgUnjail
	Params                  = slashing.Params
	QuerySigningInfoParams  = slashing.QuerySigningInfoParams
	QuerySigningInfosParams = slashing.QuerySigningInfosParams
	ValidatorSigningInfo    = slashing.ValidatorSigningInfo
	Keeper                  = slashing.Keeper
	CosmosAppModule         = slashing.AppModule
	CosmosAppModuleBasic    = slashing.AppModuleBasic
)
