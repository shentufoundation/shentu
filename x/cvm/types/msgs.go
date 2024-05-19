package types

import (
	"encoding/json"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/txs/payload"
)

// Governance message types and routes
const (
	TypeMsgDeploy = "deploy"
	TypeMsgCall   = "call"
)

var _ sdk.Msg = &MsgCall{}
var _ sdk.Msg = &MsgDeploy{}

// NewMsgCall returns a new CVM call message.
func NewMsgCall(caller, callee string, value uint64, data []byte) MsgCall {
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
	_, err := sdk.AccAddressFromBech32(m.Callee)
	_, err2 := sdk.AccAddressFromBech32(m.Caller)
	if m.Caller == "" || m.Callee == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Caller)
	}
	if err != nil || err2 != nil {
		if err != nil {
			return err
		}
		return err2
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
	addr, _ := sdk.AccAddressFromBech32(m.Caller)
	return []sdk.AccAddress{addr}
}

// NewMsgDeploy returns a new CVM deploy message.
func NewMsgDeploy(caller string, value uint64, code acm.Bytecode, abi string, meta []*payload.ContractMeta, isEWASM, isRuntime bool) MsgDeploy {
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
	_, err := sdk.AccAddressFromBech32(m.Caller)
	if m.Caller == "" || err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Caller)
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
	addr, _ := sdk.AccAddressFromBech32(m.Caller)
	return []sdk.AccAddress{addr}
}
