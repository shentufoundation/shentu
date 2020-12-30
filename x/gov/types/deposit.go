package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

// NewDeposit creates a new Deposit instance.
func NewDeposit(proposalID uint64, depositor sdk.AccAddress, amount sdk.Coins, txhash string) Deposit {
	deposit := types.NewDeposit(proposalID, depositor, amount)
	return Deposit{
		&deposit,
		txhash,
	}
}

// Deposits is a collection of Deposit objects.
type Deposits []Deposit
