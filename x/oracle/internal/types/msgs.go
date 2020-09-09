package types

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgCreateOperator struct {
	Address    sdk.AccAddress
	Collateral sdk.Coins
	Proposer   sdk.AccAddress
	Name       string
}

// NewMsgCreateOperator returns the message for creating an operator.
func NewMsgCreateOperator(address sdk.AccAddress, collateral sdk.Coins, proposer sdk.AccAddress, name string) MsgCreateOperator {
	return MsgCreateOperator{
		Address:    address,
		Collateral: collateral,
		Proposer:   proposer,
		Name:       name,
	}
}

// Route returns the module name.
func (MsgCreateOperator) Route() string { return ModuleName }

// Type returns the action name.
func (MsgCreateOperator) Type() string { return "create_operator" }

// ValidateBasic runs stateless checks on the message.
func (m MsgCreateOperator) ValidateBasic() error {
	if m.Address == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, string(m.Address.Bytes()))
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
	return []sdk.AccAddress{m.Proposer}
}

type MsgRemoveOperator struct {
	Address  sdk.AccAddress
	Proposer sdk.AccAddress
}

// NewMsgRemoveOperator returns the message for removing an operator.
func NewMsgRemoveOperator(address sdk.AccAddress, proposer sdk.AccAddress) MsgRemoveOperator {
	return MsgRemoveOperator{
		Address:  address,
		Proposer: proposer,
	}
}

// Route returns the module name.
func (MsgRemoveOperator) Route() string { return ModuleName }

// Type returns the action name.
func (MsgRemoveOperator) Type() string { return "remove_operator" }

// ValidateBasic runs stateless checks on the message.
func (m MsgRemoveOperator) ValidateBasic() error {
	if m.Address == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, string(m.Address.Bytes()))
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
	return []sdk.AccAddress{m.Proposer}
}

type MsgAddCollateral struct {
	Address             sdk.AccAddress
	CollateralIncrement sdk.Coins
}

// NewMsgAddCollateral returns the message for adding collateral.
func NewMsgAddCollateral(address sdk.AccAddress, increment sdk.Coins) MsgAddCollateral {
	return MsgAddCollateral{
		Address:             address,
		CollateralIncrement: increment,
	}
}

// Route returns the module name.
func (MsgAddCollateral) Route() string { return ModuleName }

// Type returns the action name.
func (MsgAddCollateral) Type() string { return "add_collateral" }

// ValidateBasic runs stateless checks on the message.
func (m MsgAddCollateral) ValidateBasic() error {
	if m.Address == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, string(m.Address.Bytes()))
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
	return []sdk.AccAddress{m.Address}
}

type MsgReduceCollateral struct {
	Address             sdk.AccAddress
	CollateralDecrement sdk.Coins
}

// NewMsgReduceCollateral returns the message for reducing collateral.
func NewMsgReduceCollateral(address sdk.AccAddress, decrement sdk.Coins) MsgReduceCollateral {
	return MsgReduceCollateral{
		Address:             address,
		CollateralDecrement: decrement,
	}
}

// Route returns the module name.
func (MsgReduceCollateral) Route() string { return ModuleName }

// Type returns the action name.
func (MsgReduceCollateral) Type() string { return "reduce_collateral" }

// ValidateBasic runs stateless checks on the message.
func (m MsgReduceCollateral) ValidateBasic() error {
	if m.Address == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, string(m.Address.Bytes()))
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
	return []sdk.AccAddress{m.Address}
}

type MsgWithdrawReward struct {
	Address sdk.AccAddress
}

// NewMsgWithdrawReward returns the message for withdrawing reward.
func NewMsgWithdrawReward(address sdk.AccAddress) MsgWithdrawReward {
	return MsgWithdrawReward{
		Address: address,
	}
}

// Route returns the module name.
func (MsgWithdrawReward) Route() string { return ModuleName }

// Type returns the action name.
func (MsgWithdrawReward) Type() string { return "withdraw_reward" }

// ValidateBasic runs stateless checks on the message.
func (m MsgWithdrawReward) ValidateBasic() error {
	if m.Address == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, string(m.Address.Bytes()))
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
	return []sdk.AccAddress{m.Address}
}

// MsgCreateTask is the message for creating a task.
type MsgCreateTask struct {
	Contract      string
	Function      string
	Bounty        sdk.Coins
	Description   string
	Creator       sdk.AccAddress
	Wait          int64
	Now           time.Time
	ValidDuration time.Duration
}

// NewMsgCreateTask returns a new message for creating a task.
func NewMsgCreateTask(contract, function string, bounty sdk.Coins, description string,
	creator sdk.AccAddress, wait int64, now time.Time, validDuration time.Duration) MsgCreateTask {
	return MsgCreateTask{
		Contract:      contract,
		Function:      function,
		Bounty:        bounty,
		Description:   description,
		Creator:       creator,
		Wait:          wait,
		Now:           now,
		ValidDuration: validDuration,
	}
}

// Route returns the module name.
func (MsgCreateTask) Route() string { return ModuleName }

// Type returns the action name.
func (MsgCreateTask) Type() string { return "create_task" }

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
	return []sdk.AccAddress{m.Creator}
}

// MsgTaskResponse is the message for responding to a task.
type MsgTaskResponse struct {
	Contract string
	Function string
	Score    int64
	Operator sdk.AccAddress
}

// NewMsgTaskResponse returns a new message for responding to a task.
func NewMsgTaskResponse(contract, function string, score int64, operator sdk.AccAddress) MsgTaskResponse {
	return MsgTaskResponse{
		Contract: contract,
		Function: function,
		Score:    score,
		Operator: operator,
	}
}

// Route returns the module name.
func (MsgTaskResponse) Route() string { return ModuleName }

// Type returns the action name.
func (MsgTaskResponse) Type() string { return "respond_to_task" }

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
	return []sdk.AccAddress{m.Operator}
}

// MsgInquiryTask is the message for inquiry a task.
type MsgInquiryTask struct {
	Contract string
	Function string
	TxHash   string
	Inquirer sdk.AccAddress
}

// NewMsgInquiryTask returns a new MsgInquiryTask instance.
func NewMsgInquiryTask(contract, function, txhash string, inquirer sdk.AccAddress) MsgInquiryTask {
	return MsgInquiryTask{
		Contract: contract,
		Function: function,
		TxHash:   txhash,
		Inquirer: inquirer,
	}
}

// Route returns the module name.
func (MsgInquiryTask) Route() string { return ModuleName }

// Type returns the action name.
func (MsgInquiryTask) Type() string { return "inquiry a task" }

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
	return []sdk.AccAddress{m.Inquirer}
}

// MsgDeleteTask is the msg type for delete a task.
type MsgDeleteTask struct {
	Contract string
	Function string
	Force    bool
	Deleter  sdk.AccAddress
}

// NewMsgDeleteTask returns a new MsgDeleteTask instance.
func NewMsgDeleteTask(contract, function string, force bool, deleter sdk.AccAddress) MsgDeleteTask {
	return MsgDeleteTask{
		Contract: contract,
		Function: function,
		Force:    force,
		Deleter:  deleter,
	}
}

// Route returns the module name.
func (MsgDeleteTask) Route() string { return ModuleName }

// Type returns the action name.
func (MsgDeleteTask) Type() string { return "delete_task" }

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
	return []sdk.AccAddress{m.Deleter}
}
