// Package staking defines the staking module.
package staking

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	sdksimulation "github.com/cosmos/cosmos-sdk/x/staking/simulation"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/certikfoundation/shentu/v2/common"
	"github.com/certikfoundation/shentu/v2/x/staking/keeper"
	"github.com/certikfoundation/shentu/v2/x/staking/simulation"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic is the basic app module.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the staking module's name.
func (AppModuleBasic) Name() string {
	return staking.AppModuleBasic{}.Name()
}

// RegisterCodec registers the staking module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	stakingtypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (b AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	stakingtypes.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the staking module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	defaultGenesis := stakingtypes.DefaultGenesisState()
	defaultGenesis.Params.BondDenom = common.MicroCTKDenom
	return cdc.MustMarshalJSON(defaultGenesis)
}

// ValidateGenesis performs genesis state validation for the staking module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	return staking.AppModuleBasic{}.ValidateGenesis(cdc, config, bz)
}

// RegisterRESTRoutes registers the REST routes for the staking module.
func (AppModuleBasic) RegisterRESTRoutes(cliCtx client.Context, route *mux.Router) {
	staking.AppModuleBasic{}.RegisterRESTRoutes(cliCtx, route)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the staking module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	staking.AppModuleBasic{}.RegisterGRPCGatewayRoutes(clientCtx, mux)
}

// GetTxCmd returns no root tx command for the staking module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns the root query command for the staking module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// AppModule implements an application module for the staking module.
type AppModule struct {
	AppModuleBasic
	cosmosAppModule staking.AppModule
	authKeeper      stakingtypes.AccountKeeper
	bankKeeper      stakingtypes.BankKeeper
	keeper          keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, stakingKeeper keeper.Keeper, ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper) AppModule {
	return AppModule{
		AppModuleBasic:  AppModuleBasic{},
		cosmosAppModule: staking.NewAppModule(cdc, stakingKeeper.Keeper, ak, bk),
		authKeeper:      ak,
		bankKeeper:      bk,
		keeper:          stakingKeeper,
	}
}

// Name returns the module name.
func (AppModule) Name() string {
	return stakingtypes.ModuleName
}

// RegisterInvariants registers module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the staking module.
func (am AppModule) Route() sdk.Route {
	return am.cosmosAppModule.Route()
}

// QuerierRoute returns the module query route.
func (AppModule) QuerierRoute() string {
	return stakingtypes.QuerierRoute
}

// NewQuerierHandler create new query handler.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return stakingkeeper.NewQuerier(am.keeper.Keeper, legacyQuerierCdc)
}

// RegisterServices registers a gRPC query service to respond to the
// module-specific gRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	am.cosmosAppModule.RegisterServices(cfg)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
<<<<<<< HEAD
func (AppModule) ConsensusVersion() uint64 { return 1 }
=======
func (am AppModule) ConsensusVersion() uint64 { return am.cosmosAppModule.ConsensusVersion() }
>>>>>>> 6f4b45bce5f277e193c4116dbea18212f40e242a

// BeginBlock implements the Cosmos SDK BeginBlock module function.
func (am AppModule) BeginBlock(ctx sdk.Context, rbb abci.RequestBeginBlock) {
	am.cosmosAppModule.BeginBlock(ctx, rbb)
}

// InitGenesis performs genesis initialization for the staking module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	return am.cosmosAppModule.InitGenesis(ctx, cdc, data)
}

// ExportGenesis exports genesis state data.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := staking.ExportGenesis(ctx, am.keeper.Keeper)
	return cdc.MustMarshalJSON(gs)
}

// EndBlock processes module beginblock.
func (am AppModule) EndBlock(ctx sdk.Context, reb abci.RequestEndBlock) []abci.ValidatorUpdate {
	return am.cosmosAppModule.EndBlock(ctx, reb)
}

//____________________________________________________________________________

// AppModuleSimulation functions

// WeightedOperations returns staking operations for use in simulations.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.authKeeper, am.bankKeeper, am.keeper.Keeper)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams returns functions that generate params for the module.
func (AppModuleBasic) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{}
	//return stakingSim.ParamChanges(r)
}

// GenerateGenesisState creates a randomized GenState of this module.
func (AppModuleBasic) GenerateGenesisState(simState *module.SimulationState) {
	sdksimulation.RandomizedGenState(simState)
}

// RegisterStoreDecoder registers a decoder for this module.
func (am AppModuleBasic) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[stakingtypes.StoreKey] = sdksimulation.NewDecodeStore(am.cdc)
}
