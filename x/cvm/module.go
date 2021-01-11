// Package cvm defines the cvm module.
package cvm

import (
	"context"
	"encoding/json"
	"math/rand"

	sim "github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/certikfoundation/shentu/x/cvm/client/cli"
	"github.com/certikfoundation/shentu/x/cvm/client/rest"
	"github.com/certikfoundation/shentu/x/cvm/keeper"
	"github.com/certikfoundation/shentu/x/cvm/simulation"
	"github.com/certikfoundation/shentu/x/cvm/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic specifies the app module basics object.
type AppModuleBasic struct {
	cdc codec.Marshaler
}

// NewAppModuleBasic create a new AppModuleBasic object in cvm module.
func NewAppModuleBasic() AppModuleBasic {
	return AppModuleBasic{}
}

// Name returns the gov module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the gov module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the gov
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the gov module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, config client.TxEncodingConfig, bz json.RawMessage) error {
	return types.ValidateGenesis(bz)
}

// RegisterRESTRoutes registers the REST routes for the cvm module.
func (a AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	rest.RegisterHandlers(clientCtx, rtr)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the gov module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}

// GetTxCmd returns the root tx command for the gov module.
func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns the root query command for the gov module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterInterfaces implements InterfaceModule.RegisterInterfaces
func (a AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// AppModule specifies the app module object.
type AppModule struct {
	AppModuleBasic
	keeper     keeper.Keeper
	bankkeeper types.BankKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(k keeper.Keeper, bk types.BankKeeper) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(),
		keeper:         k,
		bankkeeper:     bk,
	}
}

// Route returns the module's route key.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(am.keeper))
}

// BeginBlock implements the Cosmos SDK BeginBlock module function.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock implements the Cosmos SDK EndBlock module function.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}

// InitGenesis initializes the module genesis.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.keeper, genesisState)
}

// ExportGenesis initializes the module export genesis.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// NewHandler returns a new module handler.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// LegacyQuerierHandler returns a new querier module handler.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return keeper.NewQuerier(am.keeper, legacyQuerierCdc)
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	querier := keeper.Querier{Keeper: am.keeper}
	types.RegisterQueryServer(cfg.QueryServer(), querier)
}

// QuerierRoute returns the module querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// RegisterInvariants registers the module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// TODO Simulation

// GenerateGenesisState creates a randomized GenState of this module.
func (AppModuleBasic) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// RegisterStoreDecoder registers a decoder for cvm module.
func (am AppModuleBasic) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[types.StoreKey] = simulation.NewDecodeStore(am.cdc)
}

// WeightedOperations returns cvm operations for use in simulations.
func (am AppModule) WeightedOperations(simState module.SimulationState) sim.WeightedOperations {
	return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.keeper, am.bankkeeper)
}

// ProposalContents returns functions that generate gov proposals for the module.
func (AppModule) ProposalContents(_ module.SimulationState) []sim.WeightedProposalContent {
	return nil
}

// RandomizedParams returns functions that generate params for the module.
func (AppModuleBasic) RandomizedParams(_ *rand.Rand) []sim.ParamChange {
	return nil
}
