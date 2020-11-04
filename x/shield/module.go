package shield

import (
	"encoding/json"
	"math/rand"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/shield/client/cli"
	"github.com/certikfoundation/shentu/x/shield/client/rest"
	"github.com/certikfoundation/shentu/x/shield/simulation"
	"github.com/certikfoundation/shentu/x/shield/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic specifies the app module basics object.
type AppModuleBasic struct {
	common.AppModuleBasic
}

// NewAppModuleBasic creates a new AppModuleBasic object in shield module.
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

// AppModule implements an application module for the shield module.
type AppModule struct {
	AppModuleBasic

	keeper        Keeper
	accountKeeper types.AccountKeeper
	stakingKeeper types.StakingKeeper
	supplyKeeper  types.SupplyKeeper
}

// NewAppModule creates a new AppModule object.
func NewAppModule(keeper Keeper, ak types.AccountKeeper, stk types.StakingKeeper, sk types.SupplyKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		accountKeeper:  ak,
		stakingKeeper:  stk,
		supplyKeeper:   sk,
	}
}

// Name returns the shield module's name.
func (am AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the shield module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	RegisterInvariants(ir, am.keeper)
}

// Route returns the message routing key for the shield module.
func (am AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the shield module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the shield module's querier route name.
func (am AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the shield module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the slashing module.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.keeper, genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the shield module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the shield module.
func (am AppModule) BeginBlock(ctx sdk.Context, rbb abci.RequestBeginBlock) {
	BeginBlock(ctx, rbb, am.keeper)
}

// EndBlock returns the end blocker for the shield module.
func (am AppModule) EndBlock(ctx sdk.Context, rbb abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}

// TODO: Simulation functions
// WeightedOperations returns shield operations for use in simulations.
func (am AppModule) WeightedOperations(simState module.SimulationState) []sim.WeightedOperation {
	return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.keeper, am.accountKeeper, am.stakingKeeper)
}

// ProposalContents returns functions that generate gov proposals for the module.
func (am AppModule) ProposalContents(_ module.SimulationState) []sim.WeightedProposalContent {
	return simulation.ProposalContents(am.keeper, am.stakingKeeper)
}

// RandomizedParams returns functions that generate params for the module.
func (AppModuleBasic) RandomizedParams(r *rand.Rand) []sim.ParamChange {
	return simulation.ParamChanges(r)
}

// GenerateGenesisState creates a randomized GenState of this module.
func (AppModuleBasic) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// RegisterStoreDecoder registers a decoder for this module.
func (AppModuleBasic) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[StoreKey] = simulation.DecodeStore
}

// ValidateGenesis validates the module genesis.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return ValidateGenesis(bz)
}
