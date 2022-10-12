// Package gov defines the gov module.
package gov

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govsim "github.com/cosmos/cosmos-sdk/x/gov/simulation"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/shentufoundation/shentu/v2/x/gov/client/cli"
	"github.com/shentufoundation/shentu/v2/x/gov/client/rest"
	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
	"github.com/shentufoundation/shentu/v2/x/gov/simulation"
	"github.com/shentufoundation/shentu/v2/x/gov/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic is the app module basics object.
type AppModuleBasic struct {
	cdc              codec.Codec
	proposalHandlers []govclient.ProposalHandler // proposal handlers which live in governance cli and rest
}

// NewAppModuleBasic creates a new AppModuleBasic object.
func NewAppModuleBasic(proposalHandlers ...govclient.ProposalHandler) AppModuleBasic {
	return AppModuleBasic{
		proposalHandlers: proposalHandlers,
	}
}

// Name returns the gov module's name.
func (AppModuleBasic) Name() string {
	return gov.AppModuleBasic{}.Name()
}

// RegisterLegacyAminoCodec registers the gov module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	govtypes.RegisterLegacyAminoCodec(cdc)
}

// DefaultGenesis returns the default genesis state.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis validates the module genesis.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", govtypes.ModuleName, err)
	}

	return types.ValidateGenesis(&data)
}

// RegisterRESTRoutes registers REST routes for the module.
func (a AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
	proposalRESTHandlers := make([]govrest.ProposalRESTHandler, 0, len(a.proposalHandlers))
	for _, proposalHandler := range a.proposalHandlers {
		proposalRESTHandlers = append(proposalRESTHandlers, proposalHandler.RESTHandler(ctx))
	}

	rest.RegisterRoutes(ctx, rtr, proposalRESTHandlers)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the gov module.
func (a AppModuleBasic) RegisterGRPCGatewayRoutes(ctx client.Context, mux *runtime.ServeMux) {
	types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(ctx))
}

// GetTxCmd gets the root tx command of this module.
func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	proposalCLIHandlers := make([]*cobra.Command, 0, len(a.proposalHandlers))
	for _, proposalHandler := range a.proposalHandlers {
		proposalCLIHandlers = append(proposalCLIHandlers, proposalHandler.CLIHandler())
	}

	return cli.NewTxCmd(proposalCLIHandlers)
}

// GetQueryCmd gets the root query command of this module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterInterfaces implements InterfaceModule.RegisterInterfaces
func (a AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	govtypes.RegisterInterfaces(registry)
}

// AppModule is the main ctk module app type.
type AppModule struct {
	AppModuleBasic

	keeper        keeper.Keeper
	accountKeeper govtypes.AccountKeeper
	bankKeeper    govtypes.BankKeeper
}

// NewAppModule creates a new AppModule object.
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, ak govtypes.AccountKeeper, bk govtypes.BankKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
		accountKeeper:  ak,
		bankKeeper:     bk,
	}
}

// Name returns the governance module's name.
func (am AppModule) Name() string {
	return govtypes.ModuleName
}

// RegisterInvariants registers the governance module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	// TODO: Register cosmos invariant?
}

// Route returns the message routing key for the governance module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(govtypes.RouterKey, NewHandler(am.keeper))
}

// QuerierRoute returns the governance module's querier route name.
func (am AppModule) QuerierRoute() string {
	return govtypes.QuerierRoute
}

// LegacyQuerierHandler returns the governance module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return keeper.NewQuerier(am.keeper, legacyQuerierCdc)
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	//govtypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)

	m := keeper.NewMigrator(am.keeper)
	err := cfg.RegisterMigration(govtypes.ModuleName, 1, m.Migrate1to2)
	if err != nil {
		panic(err)
	}
}

// InitGenesis performs genesis initialization for the governance module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, am.accountKeeper, am.bankKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the governance module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (am AppModule) ConsensusVersion() uint64 { return 2 }

// BeginBlock implements the Cosmos SDK BeginBlock module function.
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock implements the Cosmos SDK EndBlock module function.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}

//____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of this module.
func (AppModuleBasic) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// ProposalContents returns all the gov content functions used to
// simulate governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return govsim.ProposalContents()
}

// RandomizedParams creates randomized gov param changes for the simulator.
func (AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	return simulation.ParamChanges(r)
}

// RegisterStoreDecoder registers a decoder for gov module.
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[govtypes.StoreKey] = simulation.NewDecodeStore(am.cdc)
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.accountKeeper, am.bankKeeper, am.keeper.CertKeeper, am.keeper, simState.Contents)
}
