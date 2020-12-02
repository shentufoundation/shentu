package bank

import (
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

const (
	ModuleName        = bankTypes.ModuleName
	RouterKey         = bankTypes.RouterKey
	DefaultParamspace = bankTypes.ModuleName
)

var (
	// functions aliases
	NewBaseKeeper      = bankKeeper.NewBaseKeeper
	NewCosmosAppModule = bank.NewAppModule

	// variable aliases
	CosmosModuleCdc = bankTypes.ModuleCdc
)

type (
	BaseKeeper           = bankKeeper.BaseKeeper // ibc module depends on this
	MsgSend              = bankTypes.MsgSend
	MsgMultiSend         = bankTypes.MsgMultiSend
	Input                = bankTypes.Input
	Output               = bankTypes.Output
	CosmosAppModule      = bank.AppModule
	CosmosAppModuleBasic = bank.AppModuleBasic
)
