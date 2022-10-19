package bank

import (
	"context"
	"encoding/json"
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
	sdkbank "github.com/cosmos/cosmos-sdk/x/bank"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	v040 "github.com/cosmos/cosmos-sdk/x/bank/legacy/v040"
	banksim "github.com/cosmos/cosmos-sdk/x/bank/simulation"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/shentufoundation/shentu/v2/x/bank/client/cli"
	"github.com/shentufoundation/shentu/v2/x/bank/client/rest"
	"github.com/shentufoundation/shentu/v2/x/bank/keeper"
	"github.com/shentufoundation/shentu/v2/x/bank/simulation"
	"github.com/shentufoundation/shentu/v2/x/bank/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the bank module.
type AppModuleBasic struct{}

// Name returns the bank module's name.
func (AppModuleBasic) Name() string {
	return sdkbank.AppModuleBasic{}.Name()
}

// RegisterLegacyAminoCodec registers the bank module's types on the LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
	*banktypes.ModuleCdc = *types.ModuleCdc
}

// DefaultGenesis returns default genesis state as raw bytes for the bank
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(banktypes.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the bank module.
func (am AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	return sdkbank.AppModuleBasic{}.ValidateGenesis(cdc, config, bz)
}

// RegisterRESTRoutes registers the REST routes for the bank module.
func (AppModuleBasic) RegisterRESTRoutes(cliCtx client.Context, route *mux.Router) {
	rest.RegisterRoutes(cliCtx, route)
	sdkbank.AppModuleBasic{}.RegisterRESTRoutes(cliCtx, route)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the bank module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	banktypes.RegisterQueryHandlerClient(context.Background(), mux, banktypes.NewQueryClient(clientCtx))
}

// GetTxCmd returns the root tx command for the bank module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	cmds := bankcli.NewTxCmd()
	cmds.AddCommand(cli.LockedSendTxCmd())
	return cmds
}

// GetQueryCmd returns the root query command for the bank module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return bankcli.GetQueryCmd()
}

// RegisterInterfaces registers interfaces and implementations of the bank module.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
	banktypes.RegisterInterfaces(registry)

	// Register legacy interfaces for migration scripts.
	v040.RegisterInterfaces(registry)
}

// ___________________________

// AppModule implements an application module for the bank module.
type AppModule struct {
	AppModuleBasic
	cosmosAppModule sdkbank.AppModule
	keeper          keeper.Keeper
	accountKeeper   types.AccountKeeper
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	am.cosmosAppModule.RegisterServices(cfg)
}

// NewAppModule creates a new AppModule object.
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, accountKeeper types.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic:  AppModuleBasic{},
		cosmosAppModule: sdkbank.NewAppModule(cdc, keeper.BaseKeeper, accountKeeper),
		keeper:          keeper,
		accountKeeper:   accountKeeper,
	}
}

// Name returns the bank module's name.
func (am AppModule) Name() string {
	return am.cosmosAppModule.Name()
}

// RegisterInvariants registers the bank module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	am.cosmosAppModule.RegisterInvariants(ir)
}

// Route returns the message routing key for the bank module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(banktypes.RouterKey, NewHandler(am.keeper))
}

// QuerierRoute returns the bank module's querier route name.
func (am AppModule) QuerierRoute() string { return am.cosmosAppModule.QuerierRoute() }

// LegacyQuerierHandler returns the bank module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return bankkeeper.NewQuerier(am.keeper.BaseKeeper, legacyQuerierCdc)
}

// InitGenesis performs genesis initialization for the bank module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	return am.cosmosAppModule.InitGenesis(ctx, cdc, data)
}

// ExportGenesis returns the exported genesis state as raw bytes for the bank module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return am.cosmosAppModule.ExportGenesis(ctx, cdc)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (am AppModule) ConsensusVersion() uint64 { return am.cosmosAppModule.ConsensusVersion() }

//____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the bank module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	banksim.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized bank param changes for the simulator.
func (AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	return banksim.ParamChanges(r)
}

// RegisterStoreDecoder performs a no-op.
func (AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.accountKeeper, am.keeper)
}
