package types

import (
	"encoding/json"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeCertifierUpdate defines the type for a CertifierUpdateProposal
	ProposalTypeCertifierUpdate = "CertifierUpdate"
)

// Assert CertifierUpdateProposal implements govtypes.Content at compile-time
var _ govtypes.Content = &CertifierUpdateProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeCertifierUpdate)
	govtypes.RegisterProposalTypeCodec(CertifierUpdateProposal{}, "cosmos-sdk/CertifierUpdateProposal")
}

// NewCertifierUpdateProposal creates a new certifier update proposal.
func NewCertifierUpdateProposal(title,
	description string,
	certifier sdk.AccAddress,
	alias string,
	proposer sdk.AccAddress,
	addorremove AddOrRemove,
) *CertifierUpdateProposal {
	return &CertifierUpdateProposal{
		Title:       title,
		Proposer:    proposer.String(),
		Alias:       alias,
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
	err := govtypes.ValidateAbstract(&cup)
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

type AddOrRemove bool

const (
	Add    AddOrRemove = false
	Remove AddOrRemove = true

	StringAdd    = "add"
	StringRemove = "remove"
)

func (aor AddOrRemove) String() string {
	switch aor {
	case Add:
		return StringAdd
	case Remove:
		return StringRemove
	default:
		panic("Invalid AddOrRemove value")
	}
}

func AddOrRemoveFromString(str string) (AddOrRemove, error) {
	switch strings.ToLower(str) {
	case StringAdd:
		return Add, nil
	case StringRemove:
		return Remove, nil
	default:
		return AddOrRemove(false), ErrAddOrRemove
	}
}

// MarshalJSON marshals to JSON using string
func (aor AddOrRemove) MarshalJSON() ([]byte, error) {
	return json.Marshal(aor.String())
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding
func (aor *AddOrRemove) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := AddOrRemoveFromString(s)
	if err != nil {
		return err
	}

	*aor = bz2
	return nil
}
