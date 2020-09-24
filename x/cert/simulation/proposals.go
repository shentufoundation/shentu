package simulation

import (
	"math/rand"

	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	certTypes "github.com/certikfoundation/shentu/x/cert/internal/types"
)

// OpWeightSubmitTextProposal app params key for text proposal
const OpWeightSubmitTextProposal = "op_weight_submit_text_proposal"

// ProposalContents defines the module weighted proposals' contents
func ProposalContents() []simulation.WeightedProposalContent {
	return []simulation.WeightedProposalContent{
		{
			AppParamsKey:       OpWeightSubmitTextProposal,
			DefaultWeight:      simappparams.DefaultWeightTextProposal,
			ContentSimulatorFn: SimulateSoftwareUpgradeProposalContent,
		},
	}
}

// SimulateTextProposalContent returns a random software upgrade proposal content.
func SimulateSoftwareUpgradeProposalContent(r *rand.Rand, _ sdk.Context, _ []simulation.Account) types.Content {
	var addorremove certTypes.AddOrRemove
	switch r.Intn(2) {
	case 0:
		addorremove = certTypes.Add
	case 1:
		addorremove = certTypes.Remove
	}
	return certTypes.NewCertifierUpdateProposal(
		simulation.RandStringOfLength(r, 140),
		simulation.RandStringOfLength(r, 5000),
		simulation.RandomAccounts(r, 1)[0].Address,
		simulation.RandStringOfLength(r, 30),
		simulation.RandomAccounts(r, 1)[0].Address,
		addorremove,
	)
}
