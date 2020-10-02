package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/txs/payload"
)

// MsgCall is the CVM call message.
type MsgCall struct {
	// Caller is the sender of the CVM-message.
	Caller sdk.AccAddress

	// Callee is the recipient of the CVM-message.
	Callee sdk.AccAddress

	// Value is the amount of CTK transferred with the call.
	Value uint64

	// Data is the binary call data.
	Data acm.Bytecode
}

// NewMsgCall returns a new CVM call message.
func NewMsgCall(caller, callee sdk.AccAddress, value uint64, data []byte) MsgCall {
	return MsgCall{
		Caller: caller,
		Callee: callee,
		Value:  value,
		Data:   data,
	}
}

// Route returns the module name.
func (m MsgCall) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgCall) Type() string { return "call" }

// ValidateBasic runs stateless checks on the message.
func (m MsgCall) ValidateBasic() error {
	if m.Caller.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Caller.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgCall) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgCall) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Caller}
}

// MsgDeploy is the CVM deploy message.
type MsgDeploy struct {
	// Caller is the sender of the CVM-message.
	Caller sdk.AccAddress

	// Value is the amount of CTK transferred with the call.
	Value uint64

	// Code is the contract byte code.
	Code acm.Bytecode

	// Abi is the Solidity ABI bytes for the contract code.
	Abi string

	// Meta is the metadata for the contract.
	Meta []*payload.ContractMeta

	// IsEWASM is true if the code is EWASM code.
	IsEWASM bool

	// IsRuntime is true if the code is runtime code.
	IsRuntime bool
}

// NewMsgDeploy returns a new CVM deploy message.
func NewMsgDeploy(caller sdk.AccAddress, value uint64, code acm.Bytecode, abi string, meta []*payload.ContractMeta, isEWASM, isRuntime bool) MsgDeploy {
	return MsgDeploy{
		Caller:    caller,
		Value:     value,
		Code:      code,
		Abi:       abi,
		Meta:      meta,
		IsEWASM:   isEWASM,
		IsRuntime: isRuntime,
	}
}

// Route returns the module name.
func (m MsgDeploy) Route() string { return ModuleName }

// Type returns the action name.
func (m MsgDeploy) Type() string { return "deploy" }

// ValidateBasic runs stateless checks on the message.
func (m MsgDeploy) ValidateBasic() error {
	if m.Caller.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Caller.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgDeploy) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required.
func (m MsgDeploy) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Caller}
}
