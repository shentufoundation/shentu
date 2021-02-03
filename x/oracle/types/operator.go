package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Operator struct {
	Address            sdk.AccAddress `json:"address"`
	Proposer           sdk.AccAddress `json:"proposer"`
	Collateral         sdk.Coins      `json:"collateral"`
	AccumulatedRewards sdk.Coins      `json:"accumulated_rewards"`
	Name               string         `json:"name"`
}

// NewOperator returns an Operator object.
func NewOperator(address sdk.AccAddress, proposer sdk.AccAddress, collateral sdk.Coins,
	accumulatedRewards sdk.Coins, name string) Operator {
	return Operator{
		Address:            address,
		Proposer:           proposer,
		Collateral:         collateral,
		AccumulatedRewards: accumulatedRewards,
		Name:               name,
	}
}

// String returns a human readable string representation of an operator.
func (o Operator) String() string {
	return fmt.Sprintf(`Operator
  Address: %s
  Proposer: %s
  Collatetal: %s
  AccumulatedRewards: %s
  Name: %s`,
		o.Address, o.Proposer, o.Collateral.String(), o.AccumulatedRewards.String(), o.Name)
}

type Operators []Operator

func (operators Operators) String() (out string) {
	for _, o := range operators {
		out += o.String() + "\n"
	}
	return strings.TrimSpace(out)
}
