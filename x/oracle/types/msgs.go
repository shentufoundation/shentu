package types

import (
	"encoding/json"
	"fmt"
	"time"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgCreateOperator   = "create_operator"
	TypeMsgRemoveOperator   = "remove_operator"
	TypeMsgAddCollateral    = "add_collateral"
	TypeMsgReduceCollateral = "reduce_collateral"
	TypeMsgWithdrawReward   = "withdraw_reward"
	TypeMsgCreateTask       = "create_task"
	TypeMsgRespondToTask    = "respond_to_task"
	TypeMsgDeleteTask       = "delete_task"
	TypeMsgCreateTxTask     = "create_tx_task"
	TypeMsgRespondToTxTask  = "respond_to_tx_task"
	TypeMsgDeleteTxTask     = "delete_tx_task"
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
	_, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return err
	}

	_, err = sdk.AccAddressFromBech32(m.Proposer)
	if err != nil {
		return err
	}

	if !m.Collateral.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.Collateral.String())
	}

	if !m.Collateral.IsAllPositive() {
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
	addr := sdk.MustAccAddressFromBech32(m.Address)
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
	_, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return err
	}

	_, err = sdk.AccAddressFromBech32(m.Proposer)
	if err != nil {
		return err
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
	addr := sdk.MustAccAddressFromBech32(m.Proposer)
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
	_, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return err
	}

	if !m.CollateralIncrement.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.CollateralIncrement.String())
	}

	if !m.CollateralIncrement.IsAllPositive() {
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
	addr := sdk.MustAccAddressFromBech32(m.Address)
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
	_, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return err
	}

	if !m.CollateralDecrement.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.CollateralDecrement.String())
	}

	if !m.CollateralDecrement.IsAllPositive() {
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
	addr := sdk.MustAccAddressFromBech32(m.Address)
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
	_, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return err
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
	addr := sdk.MustAccAddressFromBech32(m.Address)
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
	addr := sdk.MustAccAddressFromBech32(m.Creator)
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
	addr := sdk.MustAccAddressFromBech32(m.Operator)
	return []sdk.AccAddress{addr}
}

// NewMsgDeleteTask returns a new MsgDeleteTask instance.
func NewMsgDeleteTask(contract, function string, force bool, deleter sdk.AccAddress) *MsgDeleteTask {
	return &MsgDeleteTask{
		Contract: contract,
		Function: function,
		Force:    force,
		From:     deleter.String(),
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
	addr := sdk.MustAccAddressFromBech32(m.From)
	return []sdk.AccAddress{addr}
}

// NewMsgCreateTxTask returns a new MsgCreateTxTask instance.
func NewMsgCreateTxTask(creator sdk.AccAddress, chainID string, txBytes []byte,
	bounty sdk.Coins, validTime time.Time) *MsgCreateTxTask {
	return &MsgCreateTxTask{
		Creator:   creator.String(),
		ChainId:   chainID,
		AtxBytes:  txBytes,
		Bounty:    bounty,
		ValidTime: validTime,
	}
}

// Route returns the module name.
func (MsgCreateTxTask) Route() string { return ModuleName }

// Type returns the action name.
func (MsgCreateTxTask) Type() string { return TypeMsgCreateTxTask }

// ValidateBasic runs stateless checks on the message.
func (m MsgCreateTxTask) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return err
	}

	if len(m.AtxBytes) == 0 {
		return fmt.Errorf("tx bytes cannot be empty")
	}
	if len(m.ChainId) == 0 {
		return fmt.Errorf("chain_id cannot be empty")
	}
	if !m.Bounty.IsValid() {
		return fmt.Errorf("invalid bounty")
	}
	return nil
}

// GetSignBytes encodes the message for signing.
// LegacyMsg interface for Amino
func (m MsgCreateTxTask) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgCreateTxTask) GetSigners() []sdk.AccAddress {
	addr := sdk.MustAccAddressFromBech32(m.Creator)

	return []sdk.AccAddress{addr}
}

func NewMsgTxTaskResponse(txHash []byte, score int64, operator sdk.AccAddress) *MsgTxTaskResponse {
	return &MsgTxTaskResponse{
		AtxHash:  txHash,
		Score:    score,
		Operator: operator.String(),
	}
}

// LegacyMsg interface for Amino
func (MsgTxTaskResponse) Route() string { return ModuleName }

// LegacyMsg interface for Amino
func (MsgTxTaskResponse) Type() string { return TypeMsgRespondToTxTask }

// Msg interface
func (m MsgTxTaskResponse) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Operator)
	if err != nil {
		return err
	}
	if len(m.AtxHash) == 0 {
		return fmt.Errorf("atx_hash cannot be empty")
	}
	return nil
}

// LegacyMsg interface for Amino
func (m MsgTxTaskResponse) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Msg interface, return the account that should sign the tx
func (m MsgTxTaskResponse) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Operator)
	return []sdk.AccAddress{addr}
}

func NewMsgDeleteTxTask(txHash []byte, deleter sdk.AccAddress) *MsgDeleteTxTask {
	return &MsgDeleteTxTask{
		AtxHash: txHash,
		From:    deleter.String(),
	}
}

// LegacyMsg interface for Amino
func (MsgDeleteTxTask) Route() string { return ModuleName }

// LegacyMsg interface for Amino
func (MsgDeleteTxTask) Type() string { return TypeMsgDeleteTxTask }

// Msg interface
func (m MsgDeleteTxTask) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		return nil
	}
	if len(m.AtxHash) == 0 {
		return fmt.Errorf("atx_hash cannot be empty")
	}
	return nil
}

// LegacyMsg interface for Amino
func (m MsgDeleteTxTask) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Msg interface, return the account that should sign the tx
func (m MsgDeleteTxTask) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{addr}
}
