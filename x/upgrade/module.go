package upgrade

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	"github.com/certikfoundation/shentu/x/upgrade/simulation"
)

// AppModule implements the sdk.AppModule interface
type AppModule struct {
	upgrade.AppModule
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper upgrade.Keeper) AppModule {
	return AppModule{
		upgrade.NewAppModule(keeper),
	}
}

// ProposalContents returns functions that generate gov proposals for the module
func (AppModule) ProposalContents(_ module.SimulationState) []sim.WeightedProposalContent {
	return simulation.ProposalContents()
}
