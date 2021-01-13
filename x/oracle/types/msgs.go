package types

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateOperator   = "create_operator"
	TypeMsgRemoveOperator   = "remove_operator"
	TypeMsgAddCollateral    = "add_collateral"
	TypeMsgReduceCollateral = "reduce_collateral"
	TypeMsgWithdrawReward   = "withdraw_reward"
	TypeMsgCreateTask       = "create_task"
	TypeMsgRespondToTask    = "respond_to_task"
	TypeMsgInquireTask      = "inquire_task"
	TypeMsgDeleteTask       = "delete_task"
)

// NewMsgCreateOperator returns the message for creating an operator.
func NewMsgCreateOperator(address sdk.AccAddress, collateral sdk.Coins, proposer sdk.AccAddress, name string) *MsgCreateOperator {
	return &MsgCreateOperator{
		Address:    address.String(),
		Collateral: collateral,
		Proposer:   proposer.String(),
		Name:       name,
	}
}

// Route returns the module name.
func (MsgCreateOperator) Route() string { return ModuleName }

// Type returns the action name.
func (MsgCreateOperator) Type() string { return TypeMsgCreateOperator }

// ValidateBasic runs stateless checks on the message.
func (m MsgCreateOperator) ValidateBasic() error {
	addr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		panic(err)
	}
	if addr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, string(addr.Bytes()))
	}
	if m.Collateral.IsAnyNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.Collateral.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgCreateOperator) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgCreateOperator) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgRemoveOperator returns the message for removing an operator.
func NewMsgRemoveOperator(address sdk.AccAddress, proposer sdk.AccAddress) *MsgRemoveOperator {
	return &MsgRemoveOperator{
		Address:  address.String(),
		Proposer: proposer.String(),
	}
}

// Route returns the module name.
func (MsgRemoveOperator) Route() string { return ModuleName }

// Type returns the action name.
func (MsgRemoveOperator) Type() string { return TypeMsgRemoveOperator }

// ValidateBasic runs stateless checks on the message.
func (m MsgRemoveOperator) ValidateBasic() error {
	addr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		panic(err)
	}
	if addr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, string(addr.Bytes()))
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgRemoveOperator) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgRemoveOperator) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Proposer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgAddCollateral returns the message for adding collateral.
func NewMsgAddCollateral(address sdk.AccAddress, increment sdk.Coins) *MsgAddCollateral {
	return &MsgAddCollateral{
		Address:             address.String(),
		CollateralIncrement: increment,
	}
}

// Route returns the module name.
func (MsgAddCollateral) Route() string { return ModuleName }

// Type returns the action name.
func (MsgAddCollateral) Type() string { return TypeMsgAddCollateral }

// ValidateBasic runs stateless checks on the message.
func (m MsgAddCollateral) ValidateBasic() error {
	addr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		panic(err)
	}
	if addr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, string(addr.Bytes()))
	}
	if m.CollateralIncrement.IsAnyNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.CollateralIncrement.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgAddCollateral) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgAddCollateral) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgReduceCollateral returns the message for reducing collateral.
func NewMsgReduceCollateral(address sdk.AccAddress, decrement sdk.Coins) *MsgReduceCollateral {
	return &MsgReduceCollateral{
		Address:             address.String(),
		CollateralDecrement: decrement,
	}
}

// Route returns the module name.
func (MsgReduceCollateral) Route() string { return ModuleName }

// Type returns the action name.
func (MsgReduceCollateral) Type() string { return TypeMsgReduceCollateral }

// ValidateBasic runs stateless checks on the message.
func (m MsgReduceCollateral) ValidateBasic() error {
	addr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		panic(err)
	}
	if addr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, string(addr.Bytes()))
	}
	if m.CollateralDecrement.IsAnyNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.CollateralDecrement.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgReduceCollateral) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgReduceCollateral) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgWithdrawReward returns the message for withdrawing reward.
func NewMsgWithdrawReward(address sdk.AccAddress) *MsgWithdrawReward {
	return &MsgWithdrawReward{
		Address: address.String(),
	}
}

// Route returns the module name.
func (MsgWithdrawReward) Route() string { return ModuleName }

// Type returns the action name.
func (MsgWithdrawReward) Type() string { return TypeMsgWithdrawReward }

// ValidateBasic runs stateless checks on the message.
func (m MsgWithdrawReward) ValidateBasic() error {
	addr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		panic(err)
	}
	if addr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, string(addr.Bytes()))
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgWithdrawReward) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgWithdrawReward) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgCreateTask returns a new message for creating a task.
func NewMsgCreateTask(contract, function string, bounty sdk.Coins, description string,
	creator sdk.AccAddress, wait int64, validDuration time.Duration) *MsgCreateTask {
	return &MsgCreateTask{
		Contract:      contract,
		Function:      function,
		Bounty:        bounty,
		Description:   description,
		Creator:       creator.String(),
		Wait:          wait,
		ValidDuration: validDuration,
	}
}

// Route returns the module name.
func (MsgCreateTask) Route() string { return ModuleName }

// Type returns the action name.
func (MsgCreateTask) Type() string { return TypeMsgCreateTask }

// ValidateBasic runs stateless checks on the message.
func (m MsgCreateTask) ValidateBasic() error {
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgCreateTask) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgCreateTask) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgTaskResponse returns a new message for responding to a task.
func NewMsgTaskResponse(contract, function string, score int64, operator sdk.AccAddress) *MsgTaskResponse {
	return &MsgTaskResponse{
		Contract: contract,
		Function: function,
		Score:    score,
		Operator: operator.String(),
	}
}

// Route returns the module name.
func (MsgTaskResponse) Route() string { return ModuleName }

// Type returns the action name.
func (MsgTaskResponse) Type() string { return TypeMsgRespondToTask }

// ValidateBasic runs stateless checks on the message.
func (m MsgTaskResponse) ValidateBasic() error {
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgTaskResponse) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgTaskResponse) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Operator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgInquiryTask returns a new MsgInquiryTask instance.
func NewMsgInquiryTask(contract, function, txhash string, inquirer sdk.AccAddress) *MsgInquiryTask {
	return &MsgInquiryTask{
		Contract: contract,
		Function: function,
		TxHash:   txhash,
		Inquirer: inquirer.String(),
	}
}

// Route returns the module name.
func (MsgInquiryTask) Route() string { return ModuleName }

// Type returns the action name.
func (MsgInquiryTask) Type() string { return TypeMsgInquireTask }

// ValidateBasic runs stateless checks on the message.
func (m MsgInquiryTask) ValidateBasic() error {
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgInquiryTask) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgInquiryTask) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Inquirer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgDeleteTask returns a new MsgDeleteTask instance.
func NewMsgDeleteTask(contract, function string, force bool, deleter sdk.AccAddress) *MsgDeleteTask {
	return &MsgDeleteTask{
		Contract: contract,
		Function: function,
		Force:    force,
		Deleter:  deleter.String(),
	}
}

// Route returns the module name.
func (MsgDeleteTask) Route() string { return ModuleName }

// Type returns the action name.
func (MsgDeleteTask) Type() string { return TypeMsgDeleteTask }

// ValidateBasic runs stateless checks on the message.
func (m MsgDeleteTask) ValidateBasic() error {
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgDeleteTask) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgDeleteTask) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Deleter)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
