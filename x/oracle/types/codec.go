package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is a generic sealed codec to be used throughout this module.
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers concrete types on codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateOperator{}, "oracle/CreateOperator", nil)
	cdc.RegisterConcrete(MsgRemoveOperator{}, "oracle/RemoveOperator", nil)
	cdc.RegisterConcrete(MsgAddCollateral{}, "oracle/AddCollateral", nil)
	cdc.RegisterConcrete(MsgReduceCollateral{}, "oracle/ReduceCollateral", nil)
	cdc.RegisterConcrete(MsgWithdrawReward{}, "oracle/WithdrawReward", nil)
	cdc.RegisterConcrete(MsgCreateTask{}, "oracle/CreateTask", nil)
	cdc.RegisterConcrete(MsgTaskResponse{}, "oracle/RespondToTask", nil)
	cdc.RegisterConcrete(MsgInquiryTask{}, "oracle/InquiryTask", nil)
	cdc.RegisterConcrete(MsgDeleteTask{}, "oracle/DeleteTask", nil)
}
