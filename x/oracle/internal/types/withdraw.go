package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewWithdraw returns a Withdraw object.
func NewWithdraw(address sdk.AccAddress, amount sdk.Coins, dueBlock int64) Withdraw {
	return Withdraw{
		Address:  address.String(),
		Amount:   amount,
		DueBlock: dueBlock,
	}
}

type Withdraws []Withdraw

func (withdraws Withdraws) String() (out string) {
	for _, w := range withdraws {
		out += w.String() + "\n"
	}
	return strings.TrimSpace(out)
}
