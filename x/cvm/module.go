// Package cvm defines the cvm module.
package cvm

import (
	"encoding/json"
	"math/rand"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/cvm/client/cli"
	"github.com/certikfoundation/shentu/x/cvm/client/rest"
	"github.com/certikfoundation/shentu/x/cvm/simulation"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic specifies the app module basics object.
type AppModuleBasic struct {
	common.AppModuleBasic
}

// NewAppModuleBasic create a new AppModuleBasic object in cvm module.
func NewAppModuleBasic() AppModuleBasic {
	return AppModuleBasic{
		common.NewAppModuleBasic(
			ModuleName,
			RegisterCodec,
			ModuleCdc,
			DefaultGenesisState(),
			ValidateGenesis,
			StoreKey,
			rest.RegisterRoutes,
			cli.GetQueryCmd,
			cli.GetTxCmd,
		),
	}
}

// AppModule specifies the app module object.
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(k Keeper) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(),
		keeper:         k,
	}
}

// Route returns the module's route key.
func (AppModule) Route() string {
	return RouterKey
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
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.keeper, genesisState)
}

// ExportGenesis initializes the module export genesis.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// NewHandler returns a new module handler.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// NewQuerierHandler returns a new querier module handler.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// QuerierRoute returns the module querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// RegisterInvariants registers the module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// GenerateGenesisState creates a randomized GenState of this module.
func (AppModuleBasic) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// RegisterStoreDecoder registers a decoder for cvm module.
func (AppModuleBasic) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[StoreKey] = simulation.DecodeStore
}

// WeightedOperations returns cvm operations for use in simulations.
func (am AppModule) WeightedOperations(simState module.SimulationState) []sim.WeightedOperation {
	return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.keeper)
}

// ProposalContents returns functions that generate gov proposals for the module.
func (AppModule) ProposalContents(_ module.SimulationState) []sim.WeightedProposalContent {
	return nil
}

// RandomizedParams returns functions that generate params for the module.
func (AppModuleBasic) RandomizedParams(_ *rand.Rand) []sim.ParamChange {
	return nil
}
