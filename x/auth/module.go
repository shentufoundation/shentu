package auth

import (
	"encoding/json"
	"fmt"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"cosmossdk.io/core/appmodule"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	cosmosauth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/shentufoundation/shentu/v2/x/auth/keeper"
	"github.com/shentufoundation/shentu/v2/x/auth/simulation"
	"github.com/shentufoundation/shentu/v2/x/auth/types"
)

var (
	_ module.AppModuleBasic      = AppModule{}
	_ module.AppModuleSimulation = AppModule{}
	_ module.HasGenesis          = AppModule{}
	_ module.HasServices         = AppModule{}

	_ appmodule.AppModule = AppModule{}
)

// AppModuleBasic defines the basic application module used by the auth module.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the auth module's name.
func (AppModuleBasic) Name() string {
	return authtypes.ModuleName
}

// RegisterLegacyAminoCodec registers the module's types with the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
	authtypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interfaces and implementations with
// the given interface registry.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
	authtypes.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the auth module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cosmosauth.AppModuleBasic{}.DefaultGenesis(cdc)
}

// ValidateGenesis performs genesis state validation for the auth module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	return cosmosauth.AppModuleBasic{}.ValidateGenesis(cdc, config, bz)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the auth module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	cosmosauth.AppModuleBasic{}.RegisterGRPCGatewayRoutes(clientCtx, mux)
}

//____________________________________________________________________________

// AppModule implements an application module for the auth module.
type AppModule struct {
	AppModuleBasic
	cosmosAppModule cosmosauth.AppModule

	keeper     keeper.Keeper
	authKeeper authkeeper.AccountKeeper
	bankKeeper types.BankKeeper
	certKeeper types.CertKeeper

	// legacySubspace is used solely for migration of x/params managed parameters
	legacySubspace exported.Subspace
}

func (am AppModule) IsOnePerModuleType() {}

func (am AppModule) IsAppModule() {}

// NewAppModule creates a new AppModule object.
func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	ak authkeeper.AccountKeeper,
	bk types.BankKeeper,
	ck types.CertKeeper,
	randGenAccountsFn authtypes.RandomGenesisAccountsFn,
	ss exported.Subspace,
) AppModule {
	return AppModule{
		AppModuleBasic:  AppModuleBasic{cdc: cdc},
		cosmosAppModule: cosmosauth.NewAppModule(cdc, ak, randGenAccountsFn, ss),
		keeper:          keeper,
		authKeeper:      ak,
		bankKeeper:      bk,
		certKeeper:      ck,
		legacySubspace:  ss,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	authtypes.RegisterQueryServer(cfg.QueryServer(), authkeeper.NewQueryServer(am.authKeeper))

	m := keeper.NewMigrator(am.authKeeper, am.keeper, cfg.QueryServer(), am.legacySubspace)
	err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2)
	if err != nil {
		panic(err)
	}

	err = cfg.RegisterMigration(types.ModuleName, 2, m.Migrate2to3)
	if err != nil {
		panic(err)
	}

	if err := cfg.RegisterMigration(types.ModuleName, 3, m.Migrate3to4); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 3 to 4: %v", types.ModuleName, err))
	}

	if err := cfg.RegisterMigration(types.ModuleName, 4, m.Migrate4To5); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 4 to 5", types.ModuleName))
	}
}

// InitGenesis performs genesis initialization for the auth module. It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	am.cosmosAppModule.InitGenesis(ctx, cdc, data)
}

// ExportGenesis returns the exported genesis state as raw bytes for the auth module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return am.cosmosAppModule.ExportGenesis(ctx, cdc)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (am AppModule) ConsensusVersion() uint64 { return am.cosmosAppModule.ConsensusVersion() }

//____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the auth module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// RegisterStoreDecoder registers a decoder for auth module's types.
func (am AppModule) RegisterStoreDecoder(sdr simtypes.StoreDecoderRegistry) {
	am.cosmosAppModule.RegisterStoreDecoder(sdr)
}

// WeightedOperations returns auth operations for use in simulations.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return nil
}
