package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

// Deposit wraps the deposit made by an account address to an active proposal and corresponding txhash.
type Deposit struct {
	types.Deposit
	TxHash string `json:"txhash" yaml:"txhash"`
}

// NewDeposit creates a new Deposit instance.
func NewDeposit(proposalID uint64, depositor sdk.AccAddress, amount sdk.Coins, txhash string) Deposit {
	return Deposit{
		types.NewDeposit(proposalID, depositor, amount),
		txhash,
	}
}

// Deposits is a collection of Deposit objects.
type Deposits []Deposit
