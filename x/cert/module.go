// Package cert defines the cert module.
package cert

import (
	"context"
	"encoding/json"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/shentufoundation/shentu/v2/x/cert/client/cli"
	"github.com/shentufoundation/shentu/v2/x/cert/keeper"
	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic specifies the app module basics object.
type AppModuleBasic struct{}

// NewAppModuleBasic create a new AppModuleBasic object in cert module
func NewAppModuleBasic() AppModuleBasic {
	return AppModuleBasic{}
}

// Name returns the cert module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the cert module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the cert
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the cert module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	return types.ValidateGenesis(bz)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the cert module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd returns the root tx command for the cert module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns the root query command for the cert module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// AppModule specifies the app module object.
type AppModule struct {
	AppModuleBasic
	moduleKeeper keeper.Keeper
	authKeeper   types.AccountKeeper
	bankKeeper   types.BankKeeper
}

func (am AppModule) IsOnePerModuleType() {}

func (am AppModule) IsAppModule() {}

// NewAppModule creates a new AppModule object
func NewAppModule(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(),
		moduleKeeper:   k,
		authKeeper:     ak,
		bankKeeper:     bk,
	}
}

// Name returns the cert module's name.
func (am AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers the module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.moduleKeeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.Querier{Keeper: am.moduleKeeper})

	m := keeper.NewMigrator(am.moduleKeeper)
	err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2)
	if err != nil {
		panic(err)
	}
}

// InitGenesis initializes the module genesis.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.moduleKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis initializes the module export genesis.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.moduleKeeper)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 2 }

//____________________________________________________________________________

// AppModuleSimulation functions

func (am AppModule) GenerateGenesisState(input *module.SimulationState) {
}

// WeightedOperations returns cert operations for use in simulations.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	// return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.authKeeper, am.bankKeeper, am.moduleKeeper)
	return nil
}

func (am AppModule) RegisterStoreDecoder(registry simtypes.StoreDecoderRegistry) {
}

// ProposalMsgs returns functions that generate gov proposals for the module
func (am AppModule) ProposalMsgs(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	// return simulation.ProposalContents(am.moduleKeeper)
	return nil
}
