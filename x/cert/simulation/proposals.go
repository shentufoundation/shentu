package simulation

//
//import (
//	"math/rand"
//
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
//	"github.com/cosmos/cosmos-sdk/x/simulation"
//
//	"github.com/shentufoundation/shentu/v2/app/params"
//	"github.com/shentufoundation/shentu/v2/x/cert/keeper"
//	"github.com/shentufoundation/shentu/v2/x/cert/types"
//)
//
//// OpWeightSubmitCertifierUpdateProposal app params key for certifier update proposal
//const OpWeightSubmitCertifierUpdateProposal = "op_weight_submit_certifier_update_proposal"
//
//// ProposalContents defines the module weighted proposals' contents
//func ProposalContents(k keeper.Keeper) []simtypes.WeightedProposalContent {
//	return []simtypes.WeightedProposalContent{
//		simulation.NewWeightedProposalContent(
//			OpWeightSubmitCertifierUpdateProposal,
//			params.DefaultWeightCertifierUpdateProposal,
//			SimulateCertifierUpdateProposalContent(k),
//		),
//	}
//}
//
//// SimulateCertifierUpdateProposalContent generates random certifier update proposal content
//// nolint: funlen
//func SimulateCertifierUpdateProposalContent(k keeper.Keeper) simtypes.ContentSimulatorFn {
//	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content {
//		certifiers := k.GetAllCertifiers(ctx)
//		proposer_index := r.Intn(len(certifiers))
//		proposerAddr, err := sdk.AccAddressFromBech32(certifiers[proposer_index].Address)
//		if err != nil {
//			panic(err)
//		}
//
//		var addorremove types.AddOrRemove
//		var certifier sdk.AccAddress
//		switch r.Intn(2) {
//		case 0:
//			addorremove = types.Add
//			for _, acc := range accs {
//				if k.IsCertifier(ctx, acc.Address) {
//					continue
//				} else {
//					certifier = acc.Address
//					break
//				}
//			}
//		case 1:
//			addorremove = types.Remove
//			certifier_index := r.Intn(len(certifiers))
//
//			certifierAddr, err := sdk.AccAddressFromBech32(certifiers[certifier_index].Address)
//			if err != nil {
//				panic(err)
//			}
//
//			if len(certifiers) == 1 {
//				addorremove = types.Add
//				for _, acc := range accs {
//					if k.IsCertifier(ctx, acc.Address) {
//						continue
//					} else {
//						certifierAddr = acc.Address
//						break
//					}
//				}
//			}
//
//			certifier = certifierAddr
//		}
//
//		return types.NewCertifierUpdateProposal(
//			simtypes.RandStringOfLength(r, 140),
//			simtypes.RandStringOfLength(r, 5000),
//			certifier,
//			simtypes.RandStringOfLength(r, 30),
//			proposerAddr,
//			addorremove,
//		)
//	}
//}
