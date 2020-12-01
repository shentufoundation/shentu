package auth

import (
	"github.com/cosmos/cosmos-sdk/x/auth/"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	ModuleName       = types.ModuleName
	StoreKey         = types.StoreKey
	FeeCollectorName = types.FeeCollectorName
	//DefaultParamspace = types.DefaultParamspace
)

var (
	// functions aliases
	NewAnteHandler                    = ante.NewAnteHandler
	DefaultSigVerificationGasConsumer = ante.DefaultSigVerificationGasConsumer
	NewAccountKeeper                  = keeper.NewAccountKeeper
	ProtoBaseAccount                  = types.ProtoBaseAccount
	DefaultTxDecoder                  = tx.DefaultTxDecoder
	NewCosmosAppModule                = auth.NewAppModule

	// variable aliases
	CosmosModuleCdc = types.ModuleCdc
)

type (
	AccountKeeper        = keeper.AccountKeeper
	CosmosAppModule      = auth.AppModule
	CosmosAppModuleBasic = auth.AppModuleBasic
)
