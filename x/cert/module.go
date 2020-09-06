// Package cert defines the cert module.
package cert

import (
	"encoding/json"
	"math/rand"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/cert/client/cli"
	"github.com/certikfoundation/shentu/x/cert/client/rest"
	"github.com/certikfoundation/shentu/x/cert/internal/keeper"
	"github.com/certikfoundation/shentu/x/cert/internal/types"
	"github.com/certikfoundation/shentu/x/cert/simulation"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic specifies the app module basics object.
type AppModuleBasic struct {
	common.AppModuleBasic
}

// NewAppModuleBasic create a new AppModuleBasic object in cert module
func NewAppModuleBasic() AppModuleBasic {
	return AppModuleBasic{
		common.NewAppModuleBasic(
			types.ModuleName,
			types.RegisterCodec,
			types.ModuleCdc,
			types.DefaultGenesisState(),
			types.ValidateGenesis,
			types.StoreKey,
			rest.RegisterRoutes,
			cli.GetQueryCmd,
			cli.GetTxCmd,
		),
	}
}

// AppModule specifies the app module object.
type AppModule struct {
	AppModuleBasic
	moduleKeeper keeper.Keeper
	authKeeper   types.AccountKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(k keeper.Keeper, ak types.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(),
		moduleKeeper:   k,
		authKeeper:     ak,
	}
}

// Route returns the module's route key.
func (AppModule) Route() string {
	return types.RouterKey
}

// BeginBlock implements the Cosmos SDK BeginBlock module function.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
}

// EndBlock implements the Cosmos SDK EndBlock module function.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// InitGenesis initializes the module genesis.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.moduleKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis initializes the module export genesis.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.moduleKeeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// NewHandler returns a new module handler.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.moduleKeeper)
}

// NewQuerierHandler returns a new querier module handler.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(am.moduleKeeper)
}

// QuerierRoute returns the module querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// RegisterInvariants registers the module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
}

// GenerateGenesisState creates a randomized GenState of this module.
func (AppModuleBasic) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// RegisterStoreDecoder registers a decoder for cert module.
func (AppModuleBasic) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[StoreKey] = simulation.DecodeStore
}

// WeightedOperations returns cert operations for use in simulations.
func (am AppModule) WeightedOperations(simState module.SimulationState) []sim.WeightedOperation {
	return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.authKeeper, am.moduleKeeper)
}

// ProposalContents returns functions that generate gov proposals for the module
func (AppModule) ProposalContents(_ module.SimulationState) []sim.WeightedProposalContent {
	return nil
}

// RandomizedParams returns functions that generate params for the module
func (AppModuleBasic) RandomizedParams(_ *rand.Rand) []sim.ParamChange {
	return nil
}
