package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgBeginRedelegate defines the attributes of a bonding transaction.
type MsgCreatePool struct {
	From     sdk.AccAddress `json:"creator_address" yaml:"creator_address"`
	Coverage sdk.Coins      `json:"amount" yaml:"amount"`
	Deposit  MixedCoins     `json:"deposit" yaml:"deposit"`
	Sponsor  string         `json:"sponsor" yaml:"sponsor"`
}

// NewMsgBeginRedelegate creates a new MsgBeginRedelegate instance.
func NewMsgCreatePool(accAddr sdk.AccAddress, coverage sdk.Coins, deposit MixedCoins, sponsor string) (MsgCreatePool, error) {
	return MsgCreatePool{
		From:     accAddr,
		Coverage: coverage,
		Deposit:  deposit,
		Sponsor:  sponsor,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgCreatePool) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg MsgCreatePool) Type() string { return EventTypeCreatePool }

// GetSigners implements the sdk.Msg interface
func (msg MsgCreatePool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgCreatePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreatePool) ValidateBasic() error {
	if msg.Sponsor == "" {
		return ErrEmptySponsor
	}
	if msg.Deposit.Native == nil && msg.Deposit.Foreign == nil {
		return ErrNoDeposit
	}
	if msg.Coverage == nil {
		return ErrNoCoverage
	}
	return nil
}
