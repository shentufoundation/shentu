package mint

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/mint"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is the basic app module.
type AppModuleBasic struct{}

// Name returns the staking module's name.
func (AppModuleBasic) Name() string {
	return CosmosAppModuleBasic{}.Name()
}

// RegisterCodec registers the staking module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {}

// DefaultGenesis returns default genesis state as raw bytes for the mint module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the mint module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return CosmosAppModuleBasic{}.ValidateGenesis(bz)
}

// RegisterRESTRoutes registers the REST routes for the mint module.
func (AppModuleBasic) RegisterRESTRoutes(cliCtx context.CLIContext, route *mux.Router) {
	CosmosAppModuleBasic{}.RegisterRESTRoutes(cliCtx, route)
}

// GetTxCmd returns the root tx command for the mint module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return CosmosAppModuleBasic{}.GetTxCmd(cdc)
}

// GetQueryCmd returns the root query command for the mint module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return CosmosAppModuleBasic{}.GetQueryCmd(cdc)
}

// AppModule implements an application module for the mint module.
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule object.
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// Name returns the module name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route routes message routes.
func (AppModule) Route() string { return "" }

// NewHandler creates new module handler.
func (am AppModule) NewHandler() sdk.Handler { return nil }

// QuerierRoute returns the module query route.
func (AppModule) QuerierRoute() string {
	return mint.QuerierRoute
}

// NewQuerierHandler create new query handler.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper.Keeper)
}

// InitGenesis initializes genesis state from data.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper.Keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports genesis state data.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper.Keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock processes module beginblock.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock processes module beginblock.
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
