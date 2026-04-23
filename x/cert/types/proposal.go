package types

import (
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// Deprecated: CertifierUpdateProposal is the legacy v1beta1 governance proposal.
// Retained only so existing on-chain proposals can be deserialized/queried.
// New certifier updates use MsgUpdateCertifier.

var _ v1beta1.Content = &CertifierUpdateProposal{}

func (cup CertifierUpdateProposal) GetTitle() string       { return cup.Title }
func (cup CertifierUpdateProposal) GetDescription() string { return cup.Description }
func (cup CertifierUpdateProposal) ProposalRoute() string  { return ModuleName }
func (cup CertifierUpdateProposal) ProposalType() string   { return "CertifierUpdate" }
func (cup CertifierUpdateProposal) ValidateBasic() error {
	return v1beta1.ValidateAbstract(&cup)
}
