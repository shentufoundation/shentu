package simulation

import (
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	govsim "github.com/cosmos/cosmos-sdk/x/gov/simulation"
)

// ProposalMsgs returns the module-owned weighted proposal messages
// used by the simulation framework when it decides what the next
// governance proposal should contain. Shentu does not surface any new
// gov-owned proposal messages on top of upstream's text proposal, so
// this is a thin passthrough — kept here to give AppModule a single
// simulation entry point and to leave room to register shentu-specific
// messages later without touching module.go.
func ProposalMsgs() []simtypes.WeightedProposalMsg {
	return govsim.ProposalMsgs()
}
