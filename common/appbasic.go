// Package common provides common modules, constants and functions for the application.
package common

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

// AppModuleBasic defines common app module basics object.
type AppModuleBasic struct {
	moduleName          string
	regCodec            func(cdc *codec.Codec)
	modCdc              *codec.Codec
	defaultGenesisState interface{}
	validateGenesis     func(data json.RawMessage) error
	storeKey            string
	registerRoutes      func(cliCtx client.CLIContext, r *mux.Router)
	getQueryCmd         func(storeKey string, cdc *codec.Codec) *cobra.Command
	getTxCmd            func(cdc *codec.Codec) *cobra.Command
}

// NewAppModuleBasic create a new common AppModuleBasic object.
func NewAppModuleBasic(
	moduleName string,
	regCodec func(cdc *codec.Codec),
	modCdc *codec.Codec,
	defaultGenesisState interface{},
	validateGenesis func(data json.RawMessage) error,
	storeKey string,
	registerRoutes func(cliCtx client.CLIContext, r *mux.Router),
	getQueryCmd func(storeKey string, cdc *codec.Codec) *cobra.Command,
	getTxCmd func(cdc *codec.Codec) *cobra.Command,
) AppModuleBasic {
	amb := AppModuleBasic{}
	amb.moduleName = moduleName
	amb.regCodec = regCodec
	amb.modCdc = modCdc
	amb.defaultGenesisState = defaultGenesisState
	amb.validateGenesis = validateGenesis
	amb.storeKey = storeKey
	amb.registerRoutes = registerRoutes
	amb.getQueryCmd = getQueryCmd
	amb.getTxCmd = getTxCmd
	return amb
}

// Name returns the module name.
func (amb AppModuleBasic) Name() string {
	return amb.moduleName
}

// RegisterCodec registers module codec.
func (amb AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	amb.regCodec(cdc)
}

// DefaultGenesis returns the default genesis state.
func (amb AppModuleBasic) DefaultGenesis() json.RawMessage {
	return amb.modCdc.MustMarshalJSON(amb.defaultGenesisState)
}

// RegisterRESTRoutes registers REST routes for the module.
func (amb AppModuleBasic) RegisterRESTRoutes(ctx client.CLIContext, rtr *mux.Router) {
	amb.registerRoutes(ctx, rtr)
}

// GetQueryCmd gets the root query command of this module.
func (amb AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return amb.getQueryCmd(amb.storeKey, cdc)
}

// GetTxCmd gets the root tx command of this module.
func (amb AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return amb.getTxCmd(cdc)
}

// ValidateGenesis validates the module's genesis.
func (amb AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return amb.validateGenesis(bz)
}
