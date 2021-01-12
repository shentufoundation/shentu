package cli

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/types"
)

type (
	// CertifierUpdateProposalJSON defines a CertifierUpdateProposal with a deposit
	CertifierUpdateProposalJSON struct {
		Title       string            `json:"title" yaml:"title"`
		Description string            `json:"description" yaml:"description"`
		Certifier   sdk.AccAddress    `json:"certifier" yaml:"certifier"`
		Alias       string            `json:"alias" yaml:"alias"`
		AddOrRemove types.AddOrRemove `json:"add_or_remove" yaml:"add_or_remove"`
		Deposit     sdk.Coins         `json:"deposit" yaml:"deposit"`
	}
)

// ParseCertifierUpdateProposalJSON reads and parses a CertifierUpdateProposalJSON from a file.
func ParseCertifierUpdateProposalJSON(cdc *codec.Codec, proposalFile string) (CertifierUpdateProposalJSON, error) {
	proposal := CertifierUpdateProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
