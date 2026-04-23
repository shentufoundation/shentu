package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewCertifier returns a new certifier.
func NewCertifier(address sdk.AccAddress, description string) Certifier {
	return Certifier{
		Address:     address.String(),
		Description: description,
	}
}

// Certifiers is a collection of certifier objects.
type Certifiers []Certifier

func (c Certifiers) String() (out string) {
	for _, val := range c {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}
