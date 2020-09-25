package cli

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ShieldClaimProposalJSON defines a shield claim proposal.
type ShieldClaimProposalJSON struct {
	PoolID         int64     `json:"pool_id" yaml:"pool_id"`
	Loss           sdk.Coins `json:"loss" yaml:"loss"`
	Evidence       string    `json:"evidence" yaml:"evidence"`
	PurchaseTxHash string    `json:"purchase_txhash" yaml:"purchase_txhash"`
	Description    string    `json:"description" yaml:"description"`
	Deposit        sdk.Coins `json:"deposit" yaml:"deposit"`
}

// ParseShieldClaimProposalJSON reads and parses a ShieldClaimProposalJSON from a file.
func ParseShieldClaimProposalJSON(cdc *codec.Codec, proposalFile string) (ShieldClaimProposalJSON, error) {
	proposal := ShieldClaimProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
