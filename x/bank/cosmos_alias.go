package bank

import (
	"github.com/cosmos/cosmos-sdk/x/bank"
)

const (
	ModuleName        = bank.ModuleName
	RouterKey         = bank.RouterKey
	DefaultParamspace = bank.DefaultParamspace
)

var (
	// functions aliases
	NewBaseKeeper      = bank.NewBaseKeeper
	NewCosmosAppModule = bank.NewAppModule

	// variable aliases
	CosmosModuleCdc = bank.ModuleCdc
)

type (
	BaseKeeper           = bank.BaseKeeper // ibc module depends on this
	MsgSend              = bank.MsgSend
	MsgMultiSend         = bank.MsgMultiSend
	Input                = bank.Input
	Output               = bank.Output
	CosmosAppModule      = bank.AppModule
	CosmosAppModuleBasic = bank.AppModuleBasic
)
