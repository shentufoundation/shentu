// Package gov defines the gov module.
package gov

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	modulev1 "cosmossdk.io/api/cosmos/gov/module/v1"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govsim "github.com/cosmos/cosmos-sdk/x/gov/simulation"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/shentufoundation/shentu/v2/x/gov/client/cli"
	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
	"github.com/shentufoundation/shentu/v2/x/gov/simulation"
	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

const ConsensusVersion = 5

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic is the app module basics object.
type AppModuleBasic struct {
	cdc                    codec.Codec
	legacyProposalHandlers []govclient.ProposalHandler // proposal handlers which live in governance cli and rest
}

// NewAppModuleBasic creates a new AppModuleBasic object.
func NewAppModuleBasic(legacyProposalHandlers []govclient.ProposalHandler) AppModuleBasic {
	return AppModuleBasic{
		legacyProposalHandlers: legacyProposalHandlers,
	}
}

// Name returns the gov module's name.
func (AppModuleBasic) Name() string {
	return gov.AppModuleBasic{}.Name()
}

// RegisterLegacyAminoCodec registers the gov module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	govtypesv1beta1.RegisterLegacyAminoCodec(cdc)
	govtypesv1.RegisterLegacyAminoCodec(cdc)
}

// DefaultGenesis returns the default genesis state.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(typesv1.DefaultGenesisState())
}

// ValidateGenesis validates the module genesis.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data typesv1.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", govtypes.ModuleName, err)
	}

	return typesv1.ValidateGenesis(&data)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the gov module.
func (a AppModuleBasic) RegisterGRPCGatewayRoutes(ctx client.Context, mux *runtime.ServeMux) {
	err := typesv1.RegisterQueryHandlerClient(context.Background(), mux, typesv1.NewQueryClient(ctx))
	if err != nil {
		panic(err)
	}
	if err := govtypesv1beta1.RegisterQueryHandlerClient(context.Background(), mux, govtypesv1beta1.NewQueryClient(ctx)); err != nil {
		panic(err)
	}

}

// GetTxCmd gets the root tx command of this module.
func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	proposalCLIHandlers := make([]*cobra.Command, 0, len(a.legacyProposalHandlers))
	for _, proposalHandler := range a.legacyProposalHandlers {
		proposalCLIHandlers = append(proposalCLIHandlers, proposalHandler.CLIHandler())
	}

	return govcli.NewTxCmd(proposalCLIHandlers)
}

// GetQueryCmd gets the root query command of this module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterInterfaces implements InterfaceModule.RegisterInterfaces
func (a AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	govtypesv1beta1.RegisterInterfaces(registry)
	govtypesv1.RegisterInterfaces(registry)
}

// AppModule is the main ctk module app type.
type AppModule struct {
	AppModuleBasic

	keeper        keeper.Keeper
	accountKeeper govtypes.AccountKeeper
	bankKeeper    govtypes.BankKeeper

	// legacySubspace is used solely for migration of x/params managed parameters
	legacySubspace govtypes.ParamSubspace
}

var _ appmodule.AppModule = AppModule{}

// NewAppModule creates a new AppModule object.
func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper, ak govtypes.AccountKeeper, bk govtypes.BankKeeper,
	ss govtypes.ParamSubspace) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
		accountKeeper:  ak,
		bankKeeper:     bk,
		legacySubspace: ss,
	}
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

func init() {
	appmodule.Register(
		&modulev1.Module{},
		appmodule.Provide(gov.ProvideModule, typesv1.ParamKeyTable),
		appmodule.Invoke(gov.InvokeAddRoutes, gov.InvokeSetHooks))
}

// Name returns the governance module's name.
func (am AppModule) Name() string {
	return govtypes.ModuleName
}

// RegisterInvariants registers the governance module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	// TODO: Register cosmos invariant?
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	msgServer := keeper.NewMsgServerImpl(am.keeper)
	govtypesv1.RegisterMsgServer(cfg.MsgServer(), msgServer)
	govtypesv1beta1.RegisterMsgServer(cfg.MsgServer(), keeper.NewLegacyMsgServerImpl(am.accountKeeper.GetModuleAddress(govtypes.ModuleName).String(), msgServer, am.keeper))

	typesv1.RegisterQueryServer(cfg.QueryServer(), am.keeper)
	legacyQueryServer := govkeeper.NewLegacyQueryServer(&am.keeper.Keeper)
	govtypesv1beta1.RegisterQueryServer(cfg.QueryServer(), legacyQueryServer)

	m := keeper.NewMigrator(am.keeper, am.legacySubspace)
	err := cfg.RegisterMigration(govtypes.ModuleName, 1, m.Migrate1to2)
	if err != nil {
		panic(err)
	}

	err = cfg.RegisterMigration(govtypes.ModuleName, 2, m.Migrate2to3)
	if err != nil {
		panic(err)
	}

	err = cfg.RegisterMigration(govtypes.ModuleName, 3, m.Migrate3to4)
	if err != nil {
		panic(err)
	}

	err = cfg.RegisterMigration(govtypes.ModuleName, 4, m.Migrate4to5)
	if err != nil {
		panic(err)
	}
}

// InitGenesis performs genesis initialization for the governance module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState typesv1.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, am.accountKeeper, am.bankKeeper, &genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the governance module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (am AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

// EndBlock implements the Cosmos SDK EndBlock module function.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}

//____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of this module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// RegisterStoreDecoder registers a decoder for gov module.
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[govtypes.StoreKey] = govsim.NewDecodeStore(am.cdc)
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return nil
}
