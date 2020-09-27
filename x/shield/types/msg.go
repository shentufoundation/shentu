package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCreatePool defines the attributes of a create-pool transaction.
type MsgCreatePool struct {
	From             sdk.AccAddress `json:"from" yaml:"from"`
	Shield           sdk.Coins      `json:"shield" yaml:"shield"`
	Deposit          MixedCoins     `json:"deposit" yaml:"deposit"`
	Sponsor          string         `json:"sponsor" yaml:"sponsor"`
	TimeOfCoverage   int64          `json:"time_of_coverage" yaml:"time_of_coverage"`
	BlocksOfCoverage int64          `json:"blocks_of_coverage" yaml:"blocks_of_coverage"`
}

// NewMsgCreatePool creates a new MsgBeginRedelegate instance.
func NewMsgCreatePool(
	accAddr sdk.AccAddress, shield sdk.Coins, deposit MixedCoins, sponsor string, timeOfCoverage,
	blocksOfCoverage int64) (MsgCreatePool, error) {
	return MsgCreatePool{
		From:             accAddr,
		Shield:           shield,
		Deposit:          deposit,
		Sponsor:          sponsor,
		TimeOfCoverage:   timeOfCoverage,
		BlocksOfCoverage: blocksOfCoverage,
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
		return ErrNoShield
	}
	if msg.TimeOfCoverage <= 0 && msg.BlocksOfCoverage <= 0 {
		return ErrInvalidDuration
	}
	return nil
}

// MsgUpdatePool defines the attributes of a shield pool update transaction.
type MsgUpdatePool struct {
	From             sdk.AccAddress `json:"from" yaml:"from"`
	Shield           sdk.Coins      `json:"Shield" yaml:"Shield"`
	Deposit          MixedCoins     `json:"deposit" yaml:"deposit"`
	PoolID           uint64         `json:"pool_id" yaml:"pool_id"`
	AdditionalTime   int64          `json:"additional_period" yaml:"additional_period"`
	AdditionalBlocks int64          `json:"additional_blocks" yaml:"additional_blocks"`
}

// NewMsgUpdatePool creates a new MsgUpdatePool instance.
func NewMsgUpdatePool(
	accAddr sdk.AccAddress, shield sdk.Coins, deposit MixedCoins, id uint64, additionalTime, additionalBlocks int64) (MsgUpdatePool, error) {
	return MsgUpdatePool{
		From:             accAddr,
		Shield:           shield,
		Deposit:          deposit,
		PoolID:           id,
		AdditionalTime:   additionalTime,
		AdditionalBlocks: additionalBlocks,
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
	if msg.PoolID == 0 {
		return ErrInvalidPoolID
	}
	if msg.Deposit.Native == nil && msg.Deposit.Foreign == nil {
		return ErrNoDeposit
	}
	if msg.Shield == nil {
		return ErrNoShield
	}
	if msg.AdditionalTime <= 0 && msg.AdditionalBlocks <= 0 {
		return ErrInvalidDuration
	}
	return nil
}

// MsgPausePool defines the attributes of a pausing a shield pool.
type MsgPausePool struct {
	From   sdk.AccAddress `json:"from" yaml:"from"`
	PoolID uint64         `json:"pool_id" yaml:"pool_id"`
}

// NewMsgPausePool creates a new NewMsgPausePool instance.
func NewMsgPausePool(accAddr sdk.AccAddress, id uint64) (MsgPausePool, error) {
	return MsgPausePool{
		From:   accAddr,
		PoolID: id,
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
	if msg.PoolID == 0 {
		return ErrInvalidPoolID
	}
	return nil
}

// MsgResumePool defines the attributes of a resuming a shield pool.
type MsgResumePool struct {
	From   sdk.AccAddress `json:"from" yaml:"from"`
	PoolID uint64         `json:"pool_id" yaml:"pool_id"`
}

// NewMsgResumePool creates a new NewMsgResumePool instance.
func NewMsgResumePool(accAddr sdk.AccAddress, id uint64) (MsgResumePool, error) {
	return MsgResumePool{
		From:   accAddr,
		PoolID: id,
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
	if msg.PoolID == 0 {
		return ErrInvalidPoolID
	}
	return nil
}

type MsgDepositCollateral struct {
	From       sdk.AccAddress `json:"sender" yaml:"sender"`
	PoolID     uint64         `json:"pool_id" yaml:"pool_id"`
	Collateral sdk.Coins      `json:"collateral" yaml:"collateral"`
}

// NewMsgDepositCollateral creates a new MsgDepositCollateral instance.
func NewMsgDepositCollateral(sender sdk.AccAddress, id uint64, collateral sdk.Coins) (MsgDepositCollateral, error) {
	return MsgDepositCollateral{
		From:       sender,
		PoolID:     id,
		Collateral: collateral,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgDepositCollateral) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg MsgDepositCollateral) Type() string { return "deposit_collateral" }

// GetSigners implements the sdk.Msg interface
func (msg MsgDepositCollateral) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgDepositCollateral) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgDepositCollateral) ValidateBasic() error {
	if msg.PoolID == 0 {
		return ErrInvalidPoolID
	}
	return nil
}
