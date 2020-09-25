package simulation

import (
	// "fmt"
	"math/rand"
	"time"

	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

// OpWeightSubmitCommunitySpendProposal app params key for community spend proposal
const OpWeightSubmitCommunitySpendProposal = "op_weight_submit_community_spend_proposal"

// ProposalContents defines the module weighted proposals' contents
func ProposalContents() []simulation.WeightedProposalContent {
	return []simulation.WeightedProposalContent{
		{
			AppParamsKey:       OpWeightSubmitCommunitySpendProposal,
			DefaultWeight:      simappparams.DefaultWeightCommunitySpendProposal,
			ContentSimulatorFn: SimulateSoftwareUpgradeProposalContent(),
		},
	}
}

// SimulateSoftwareUpgradeProposalContent generates random software upgrade proposal content
// nolint: funlen
func SimulateSoftwareUpgradeProposalContent() simulation.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simulation.Account) govtypes.Content {
		plan := upgrade.Plan{
			Name:   simulation.RandStringOfLength(r, 140),
			Time:   time.Now(),
			Height: ctx.BlockHeight() + 50,
			Info:   simulation.RandStringOfLength(r, 140),
		}
		return upgrade.NewSoftwareUpgradeProposal(
			simulation.RandStringOfLength(r, 140),
			simulation.RandStringOfLength(r, 5000),
			plan,
		)
	}
}
