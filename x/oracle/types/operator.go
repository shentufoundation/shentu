package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewOperator returns an Operator object.
func NewOperator(
	address sdk.AccAddress,
	proposer sdk.AccAddress,
	collateral sdk.Coins,
	accumulatedRewards sdk.Coins,
	name string,
) Operator {
	return Operator{
		Address:            address.String(),
		Proposer:           proposer.String(),
		Collateral:         collateral,
		AccumulatedRewards: accumulatedRewards,
		Name:               name,
	}
}

type Operators []Operator

func (operators Operators) String() (out string) {
	for _, o := range operators {
		out += o.String() + "\n"
	}
	return strings.TrimSpace(out)
}
