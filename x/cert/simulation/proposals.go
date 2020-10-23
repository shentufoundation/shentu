package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/app/params"
	"github.com/certikfoundation/shentu/x/cert/internal/keeper"
	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

// OpWeightSubmitCertifierUpdateProposal app params key for certifier update proposal
const OpWeightSubmitCertifierUpdateProposal = "op_weight_submit_certifier_update_proposal"

// ProposalContents defines the module weighted proposals' contents
func ProposalContents(k keeper.Keeper) []simulation.WeightedProposalContent {
	return []simulation.WeightedProposalContent{
		{
			AppParamsKey:       OpWeightSubmitCertifierUpdateProposal,
			DefaultWeight:      params.DefaultWeightCertifierUpdateProposal,
			ContentSimulatorFn: SimulateCertifierUpdateProposalContent(k),
		},
	}
}

// SimulateCertifierUpdateProposalContent generates random certifier update proposal content
// nolint: funlen
func SimulateCertifierUpdateProposalContent(k keeper.Keeper) simulation.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simulation.Account) govtypes.Content {
		certifiers := k.GetAllCertifiers(ctx)
		proposer_index := r.Intn(len(certifiers))

		var addorremove types.AddOrRemove
		var certifier sdk.AccAddress
		switch r.Intn(2) {
		case 0:
			addorremove = types.Add
			for _, acc := range accs {
				if k.IsCertifier(ctx, acc.Address) {
					continue
				} else {
					certifier = acc.Address
					break
				}
			}
		case 1:
			addorremove = types.Remove
			certifier_index := r.Intn(len(certifiers))
			certifier = certifiers[certifier_index].Address
		}

		return types.NewCertifierUpdateProposal(
			simulation.RandStringOfLength(r, 140),
			simulation.RandStringOfLength(r, 5000),
			certifier,
			simulation.RandStringOfLength(r, 30),
			certifiers[proposer_index].Address,
			addorremove,
		)
	}
}
