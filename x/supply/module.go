package supply

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/supply/client/cli"
	"github.com/cosmos/cosmos-sdk/x/supply/client/rest"
	supplySim "github.com/cosmos/cosmos-sdk/x/supply/simulation"

	"github.com/certikfoundation/shentu/x/supply/internal/types"
	"github.com/certikfoundation/shentu/x/supply/simulation"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the supply module.
type AppModuleBasic struct{}

// Name returns the supply module's name.
func (AppModuleBasic) Name() string {
	return supply.ModuleName
}

// RegisterCodec registers the supply module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	supply.RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the supply module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return supply.ModuleCdc.MustMarshalJSON(supply.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the supply module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data supply.GenesisState
	if err := supply.ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", supply.ModuleName, err)
	}

	return supply.ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the supply module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the supply module.
func (AppModuleBasic) GetTxCmd(_ *codec.Codec) *cobra.Command { return nil }

// GetQueryCmd returns no root query command for the supply module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

//____________________________________________________________________________

// AppModule implements an application module for the supply module.
type AppModule struct {
	AppModuleBasic

	keeper supply.Keeper
	ak     types.AccountKeeper
}

// NewAppModule creates a new AppModule object.
func NewAppModule(keeper supply.Keeper, ak types.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		ak:             ak,
	}
}

// Name returns the supply module's name.
func (AppModule) Name() string {
	return supply.ModuleName
}

// RegisterInvariants registers the supply module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	supply.RegisterInvariants(ir, am.keeper)
}

// Route returns the message routing key for the supply module.
func (AppModule) Route() string {
	return supply.RouterKey
}

// NewHandler returns an sdk.Handler for the supply module.
func (am AppModule) NewHandler() sdk.Handler { return nil }

// QuerierRoute returns the supply module's querier route name.
func (AppModule) QuerierRoute() string {
	return supply.QuerierRoute
}

// NewQuerierHandler returns the supply module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return supply.NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the supply module. It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState supply.GenesisState
	supply.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	supply.InitGenesis(ctx, am.keeper, am.ak, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the supply module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := supply.ExportGenesis(ctx, am.keeper)
	return supply.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the supply module.
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the supply module. It returns no validator updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

//____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the supply module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []sim.WeightedProposalContent {
	return nil
}

// RandomizedParams doesn't create any randomized supply param changes for the simulator.
func (AppModule) RandomizedParams(_ *rand.Rand) []sim.ParamChange {
	return nil
}

// RegisterStoreDecoder registers a decoder for supply module's types.
func (AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[supply.StoreKey] = supplySim.DecodeStore
}

// WeightedOperations doesn't return any operation for the nft module.
func (AppModule) WeightedOperations(_ module.SimulationState) []sim.WeightedOperation {
	return nil
}
