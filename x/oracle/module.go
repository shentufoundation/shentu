package oracle

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/oracle/client/cli"
	"github.com/certikfoundation/shentu/x/oracle/client/rest"
	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic specifies the app module basics object.
type AppModuleBasic struct {
	common.AppModuleBasic
}

// NewAppModuleBasic creates a new AppModuleBasic object in cert module.
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

// AppModule is the main ctk module app type.
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	oracleKeeper Keeper,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         oracleKeeper,
	}
}

// Name returns the oracle module's name.
func (am AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the this module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the oracle module.
func (am AppModule) Route() string {
	return types.RouterKey
}

// NewHandler returns an sdk.Handler for the oracle module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the oracle module's querier route name.
func (am AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the oracle module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the oracle module.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the oracle module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock implements the Cosmos SDK BeginBlock module function.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock implements the Cosmos SDK EndBlock module function.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}
