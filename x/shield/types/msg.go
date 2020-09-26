package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCreatePool defines the attributes of a create-pool transaction.
type MsgCreatePool struct {
	From    sdk.AccAddress `json:"from" yaml:"from"`
	Shield  sdk.Coins      `json:"shield" yaml:"shield"`
	Deposit MixedCoins     `json:"deposit" yaml:"deposit"`
	Sponsor string         `json:"sponsor" yaml:"sponsor"`
}

// NewMsgCreatePool creates a new MsgBeginRedelegate instance.
func NewMsgCreatePool(accAddr sdk.AccAddress, coverage sdk.Coins, deposit MixedCoins, sponsor string) (MsgCreatePool, error) {
	return MsgCreatePool{
		From:    accAddr,
		Shield:  coverage,
		Deposit: deposit,
		Sponsor: sponsor,
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
	if msg.Shield == nil {
		return ErrNoCoverage
	}
	return nil
}

// MsgUpdatePool defines the attributes of a shield pool update transaction.
type MsgUpdatePool struct {
	From    sdk.AccAddress `json:"from" yaml:"from"`
	Shield  sdk.Coins      `json:"Shield" yaml:"Shield"`
	Deposit MixedCoins     `json:"deposit" yaml:"deposit"`
	Sponsor string         `json:"sponsor" yaml:"sponsor"`
}

// NewMsgUpdatePool creates a new MsgUpdatePool instance.
func NewMsgUpdatePool(accAddr sdk.AccAddress, coverage sdk.Coins, deposit MixedCoins, sponsor string) (MsgCreatePool, error) {
	return MsgCreatePool{
		From:    accAddr,
		Shield:  coverage,
		Deposit: deposit,
		Sponsor: sponsor,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdatePool) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg MsgUpdatePool) Type() string { return EventTypeUpdatePool }

// GetSigners implements the sdk.Msg interface
func (msg MsgUpdatePool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdatePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdatePool) ValidateBasic() error {
	if msg.Sponsor == "" {
		return ErrEmptySponsor
	}
	if msg.Deposit.Native == nil && msg.Deposit.Foreign == nil {
		return ErrNoDeposit
	}
	if msg.Shield == nil {
		return ErrNoCoverage
	}
	return nil
}

// MsgPausePool defines the attributes of a pausing a shield pool.
type MsgPausePool struct {
	From     sdk.AccAddress `json:"from" yaml:"from"`
	Sponsor  string         `json:"sponsor" yaml:"sponsor"`
}

// NewMsgPausePool creates a new NewMsgPausePool instance.
func NewMsgPausePool(accAddr sdk.AccAddress, sponsor string) (MsgPausePool, error) {
	return MsgPausePool{
		From:     accAddr,
		Sponsor:  sponsor,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgPausePool) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg MsgPausePool) Type() string { return EventTypePausePool }

// GetSigners implements the sdk.Msg interface
func (msg MsgPausePool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgPausePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgPausePool) ValidateBasic() error {
	if msg.Sponsor == "" {
		return ErrEmptySponsor
	}
	return nil
}

// MsgResumePool defines the attributes of a resuming a shield pool.
type MsgResumePool struct {
	From     sdk.AccAddress `json:"from" yaml:"from"`
	Sponsor  string         `json:"sponsor" yaml:"sponsor"`
}

// NewMsgResumePool creates a new NewMsgResumePool instance.
func NewMsgResumePool(accAddr sdk.AccAddress, sponsor string) (MsgResumePool, error) {
	return MsgResumePool{
		From:     accAddr,
		Sponsor:  sponsor,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgResumePool) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg MsgResumePool) Type() string { return EventTypeResumePool }

// GetSigners implements the sdk.Msg interface
func (msg MsgResumePool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgResumePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgResumePool) ValidateBasic() error {
	if msg.Sponsor == "" {
		return ErrEmptySponsor
	}
	return nil
}