// Package staking defines the staking module.
package staking

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingSim "github.com/cosmos/cosmos-sdk/x/staking/simulation"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/staking/client/rest"
	"github.com/certikfoundation/shentu/x/staking/internal/keeper"
	"github.com/certikfoundation/shentu/x/staking/internal/types"
	"github.com/certikfoundation/shentu/x/staking/simulation"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the staking module.
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the staking module's name.
func (AppModuleBasic) Name() string {
	return staking.AppModuleBasic{}.Name()
}

// RegisterCodec registers the staking module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	staking.RegisterCodec(cdc)
	// *staking.ModuleCdc = *ModuleCdc
}

// DefaultGenesis returns default genesis state as raw bytes for the staking module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	defaultGenesisState := staking.DefaultGenesisState()
	defaultGenesisState.Params.BondDenom = common.MicroCTKDenom
	defaultGenesisState.Params.UnbondingTime = time.Hour * 24 * 7 * 3
	defaultGenesisState.Params.MaxValidators = 125

	return staking.ModuleCdc.MustMarshalJSON(defaultGenesisState)
}

// ValidateGenesis performs genesis state validation for the staking module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return staking.AppModuleBasic{}.ValidateGenesis(bz)
}

// RegisterRESTRoutes registers the REST routes for the staking module.
func (AppModuleBasic) RegisterRESTRoutes(cliCtx context.CLIContext, route *mux.Router) {
	rest.RegisterRoutes(cliCtx, route)
	staking.AppModuleBasic{}.RegisterRESTRoutes(cliCtx, route)
}

// GetTxCmd returns the root tx command for the staking module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return staking.AppModuleBasic{}.GetTxCmd(cdc)
}

// GetQueryCmd returns the root query command for the staking module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return staking.AppModuleBasic{}.GetQueryCmd(cdc)
}

// CreateValidatorMsgHelpers is used for gen-tx.
func (AppModuleBasic) CreateValidatorMsgHelpers(ipDefault string) (
	fs *pflag.FlagSet, nodeIDFlag, pubkeyFlag, amountFlag, defaultsDesc string) {
	return staking.AppModuleBasic{}.CreateValidatorMsgHelpers(ipDefault)
}

// PrepareFlagsForTxCreateValidator is used for gen-tx.
func (AppModuleBasic) PrepareFlagsForTxCreateValidator(cfg *config.Config, nodeID,
	chainID string, valPubKey crypto.PubKey) {
	staking.AppModuleBasic{}.PrepareFlagsForTxCreateValidator(cfg, nodeID, chainID, valPubKey)
}

// BuildCreateValidatorMsg is used for gen-tx.
func (AppModuleBasic) BuildCreateValidatorMsg(cliCtx context.CLIContext,
	txBldr authtypes.TxBuilder) (authtypes.TxBuilder, sdk.Msg, error) {
	return staking.AppModuleBasic{}.BuildCreateValidatorMsg(cliCtx, txBldr)
}

// AppModule implements an application module for the staking module.
type AppModule struct {
	AppModuleBasic
	cosmosAppModule staking.AppModule
	authKeeper      stakingTypes.AccountKeeper
	certKeeper      types.CertKeeper
	keeper          keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	stakingKeeper keeper.Keeper,
	accountKeeper stakingTypes.AccountKeeper,
	supplyKeeper stakingTypes.SupplyKeeper,
	certKeeper types.CertKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic:  AppModuleBasic{},
		cosmosAppModule: staking.NewAppModule(stakingKeeper.Keeper, accountKeeper, supplyKeeper),
		authKeeper:      accountKeeper,
		certKeeper:      certKeeper,
		keeper:          stakingKeeper,
	}
}

// Name returns the staking module's name.
func (am AppModule) Name() string {
	return am.cosmosAppModule.Name()
}

// RegisterInvariants registers the staking module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	am.cosmosAppModule.RegisterInvariants(ir)
}

// Route returns the message routing key for the staking module.
func (am AppModule) Route() string {
	return am.cosmosAppModule.Route()
}

// NewHandler returns an sdk.Handler for the staking module.
func (am AppModule) NewHandler() sdk.Handler {
	return staking.NewHandler(am.keeper.Keeper)
}

// QuerierRoute returns the staking module's querier route name.
func (am AppModule) QuerierRoute() string { return am.cosmosAppModule.QuerierRoute() }

// NewQuerierHandler returns the staking module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier { return am.cosmosAppModule.NewQuerierHandler() }

// InitGenesis performs genesis initialization for the staking module.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return am.cosmosAppModule.InitGenesis(ctx, data)
}

// ExportGenesis returns the exported genesis state as raw bytes for the staking module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return am.cosmosAppModule.ExportGenesis(ctx)
}

// BeginBlock implements the Cosmos SDK BeginBlock module function.
func (am AppModule) BeginBlock(ctx sdk.Context, rbb abci.RequestBeginBlock) {
	am.cosmosAppModule.BeginBlock(ctx, rbb)
}

// EndBlock implements the Cosmos SDK EndBlock module function.
func (am AppModule) EndBlock(ctx sdk.Context, rbb abci.RequestEndBlock) []abci.ValidatorUpdate {
	return am.cosmosAppModule.EndBlock(ctx, rbb)
}

// GenerateGenesisState creates a randomized GenState of this module.
func (AppModuleBasic) GenerateGenesisState(simState *module.SimulationState) {
	stakingSim.RandomizedGenState(simState)
}

// RegisterStoreDecoder registers a decoder for staking module.
func (AppModuleBasic) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[StoreKey] = stakingSim.DecodeStore
}

// WeightedOperations returns staking operations for use in simulations.
func (am AppModule) WeightedOperations(simState module.SimulationState) []sim.WeightedOperation {
	return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.authKeeper, am.certKeeper, am.keeper.Keeper)
}

// ProposalContents returns functions that generate proposals for the module.
func (AppModule) ProposalContents(_ module.SimulationState) []sim.WeightedProposalContent {
	return nil
}

// RandomizedParams returns functions that generate params for the module.
func (AppModuleBasic) RandomizedParams(r *rand.Rand) []sim.ParamChange {
	return []sim.ParamChange{}
	//return stakingSim.ParamChanges(r)
}
