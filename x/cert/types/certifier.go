package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Certifier is a type for certifier.
type Certifier struct {
	Address     sdk.AccAddress `json:"certifier"`
	Alias       string         `json:"alias"`
	Proposer    sdk.AccAddress `json:"proposer"`
	Description string         `json:"description"`
}

// NewCertifier returns a new certifier.
func NewCertifier(
	address sdk.AccAddress,
	alias string,
	proposer sdk.AccAddress,
	description string,
) Certifier {
	return Certifier{
		Address:     address,
		Alias:       alias,
		Proposer:    proposer,
		Description: description,
	}
}

// String returns a human readable string representation of a validator.
func (c Certifier) String() string {
	return fmt.Sprintf(`Certifier
  Address: %s
  Proposer: %s
  Description: %s`,
		c.Address, c.Proposer, c.Description)
}

// Certifiers is a collection of certifier objects.
type Certifiers []Certifier

func (c Certifiers) String() (out string) {
	for _, val := range c {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}
