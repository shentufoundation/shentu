package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Withdraw struct {
	Address  sdk.AccAddress `json:"address"`
	Amount   sdk.Coins      `json:"amount"`
	DueBlock int64          `json:"due_block"`
}

// NewWithdraw returns a Withdraw object.
func NewWithdraw(address sdk.AccAddress, amount sdk.Coins, dueBlock int64) Withdraw {
	return Withdraw{
		Address:  address,
		Amount:   amount,
		DueBlock: dueBlock,
	}
}

// String returns a human readable string representation of a withdraw object.
func (w Withdraw) String() string {
	return fmt.Sprintf(`Withdraw
  Address: %s
  Amount: %s
  DueBlock: %d`,
		w.Address, w.Amount, w.DueBlock)
}

type Withdraws []Withdraw

func (withdraws Withdraws) String() (out string) {
	for _, w := range withdraws {
		out += w.String() + "\n"
	}
	return strings.TrimSpace(out)
}
