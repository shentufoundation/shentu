package simulation

import (
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	"github.com/certikfoundation/shentu/app/params"
)

// OpWeightSubmitSoftwareUpgradeProposal app params key for software upgrade proposal
const OpWeightSubmitSoftwareUpgradeProposal = "op_weight_submit_software_upgrade_proposal"

// ProposalContents defines the module weighted proposals' contents
func ProposalContents() []simulation.WeightedProposalContent {
	return []simulation.WeightedProposalContent{
		{
			AppParamsKey:       OpWeightSubmitSoftwareUpgradeProposal,
			DefaultWeight:      params.DefaultWeightSoftwareUpgradeProposal,
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
