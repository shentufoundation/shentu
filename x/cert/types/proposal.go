package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	// ProposalTypeCertifierUpdate defines the type for a CertifierUpdateProposal
	ProposalTypeCertifierUpdate = "CertifierUpdate"
)

// Assert CertifierUpdateProposal implements govtypes.Content at compile-time
var _ v1beta1.Content = &CertifierUpdateProposal{}

func init() {
	v1beta1.RegisterProposalType(ProposalTypeCertifierUpdate)
	// govtypes.RegisterProposalTypeCodec(CertifierUpdateProposal{}, "cosmos-sdk/CertifierUpdateProposal")
}

// NewCertifierUpdateProposal creates a new certifier update proposal.
func NewCertifierUpdateProposal(title,
	description string,
	certifier sdk.AccAddress,
	proposer sdk.AccAddress,
	addorremove AddOrRemove,
) *CertifierUpdateProposal {
	return &CertifierUpdateProposal{
		Title:       title,
		Proposer:    proposer.String(),
		Certifier:   certifier.String(),
		Description: description,
		AddOrRemove: addorremove,
	}
}

// GetTitle returns the title of a certifier update proposal.
func (cup CertifierUpdateProposal) GetTitle() string { return cup.Title }

// GetDescription returns the description of a certifier update proposal.
func (cup CertifierUpdateProposal) GetDescription() string { return cup.Description }

// ProposalRoute returns the routing key of a certifier update proposal.
func (cup CertifierUpdateProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a certifier update proposal.
func (cup CertifierUpdateProposal) ProposalType() string { return ProposalTypeCertifierUpdate }

// ValidateBasic runs basic stateless validity checks
func (cup CertifierUpdateProposal) ValidateBasic() error {
	err := v1beta1.ValidateAbstract(&cup)
	if err != nil {
		return err
	}

	certifierAddr, err := sdk.AccAddressFromBech32(cup.Certifier)
	if err != nil {
		panic(err)
	}
	if certifierAddr.Empty() {
		return ErrEmptyCertifier
	}

	return nil
}
