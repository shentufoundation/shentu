package simulation

import (
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/shentufoundation/shentu/v2/x/cert/keeper"
)

// ProposalMsgs returns the weighted proposal messages that the
// simulator can wrap in a MsgSubmitProposal.
//
// MsgUpdateCertifier is intentionally NOT surfaced here. The upstream
// gov simulator picks the submit-proposal proposer at random from all
// sim accounts, but shentu's gov keeper rejects cert-update proposals
// from non-certifiers (see Keeper.ValidateCertUpdateProposer). Feeding
// cert-update into this pipeline therefore fails the sim run almost
// every time. The guard itself is covered by keeper unit tests.
func ProposalMsgs(_ keeper.Keeper) []simtypes.WeightedProposalMsg {
	return nil
}
