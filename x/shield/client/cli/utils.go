package cli

import (
	"encoding/json"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ShieldClaimProposalJSON defines a shield claim proposal.
type ShieldClaimProposalJSON struct {
	PoolID      uint64    `json:"pool_id" yaml:"pool_id"`
	Loss        sdk.Coins `json:"loss" yaml:"loss"`
	Evidence    string    `json:"evidence" yaml:"evidence"`
	PurchaseID  uint64    `json:"purchase_id" yaml:"purchase_id"`
	Description string    `json:"description" yaml:"description"`
	Deposit     sdk.Coins `json:"deposit" yaml:"deposit"`
}

// ParseShieldClaimProposalJSON reads and parses a ShieldClaimProposalJSON from a file.
func ParseShieldClaimProposalJSON(proposalFile string) (ShieldClaimProposalJSON, error) {
	proposal := ShieldClaimProposalJSON{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := json.Unmarshal(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
