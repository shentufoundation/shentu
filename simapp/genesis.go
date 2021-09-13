package simapp

import "github.com/cosmos/cosmos-sdk/simapp"

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() simapp.GenesisState {
	encCfg := MakeTestEncodingConfig()
	return ModuleBasics.DefaultGenesis(encCfg.Codec)
}
