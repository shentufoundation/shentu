package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers concrete types on the LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgCreateOperator{}, "oracle/CreateOperator", nil)
	cdc.RegisterConcrete(MsgRemoveOperator{}, "oracle/RemoveOperator", nil)
	cdc.RegisterConcrete(MsgAddCollateral{}, "oracle/AddCollateral", nil)
	cdc.RegisterConcrete(MsgReduceCollateral{}, "oracle/ReduceCollateral", nil)
	cdc.RegisterConcrete(MsgWithdrawReward{}, "oracle/WithdrawReward", nil)
	cdc.RegisterConcrete(MsgCreateTask{}, "oracle/CreateTask", nil)
	cdc.RegisterConcrete(MsgTaskResponse{}, "oracle/RespondToTask", nil)
	cdc.RegisterConcrete(MsgDeleteTask{}, "oracle/DeleteTask", nil)
	cdc.RegisterConcrete(MsgCreateTxTask{}, "oracle/CreateTxTask", nil)
	cdc.RegisterConcrete(MsgTxTaskResponse{}, "oracle/RespondToTxTask", nil)
	cdc.RegisterConcrete(MsgDeleteTxTask{}, "oracle/DeleteTxTask", nil)

	cdc.RegisterInterface((*TaskI)(nil), nil)
	cdc.RegisterConcrete(Task{}, "oracle/Task", nil)
	cdc.RegisterConcrete(TxTask{}, "oracle/TxTask", nil)
}

// RegisterInterfaces registers the x/oracle interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateOperator{},
		&MsgRemoveOperator{},
		&MsgAddCollateral{},
		&MsgReduceCollateral{},
		&MsgWithdrawReward{},
		&MsgCreateTask{},
		&MsgTaskResponse{},
		&MsgDeleteTask{},
		&MsgCreateTxTask{},
		&MsgTxTaskResponse{},
		&MsgDeleteTxTask{},
	)
	registry.RegisterInterface("shentu.oracle.v1alpha1.TaskI", (*TaskI)(nil), &Task{}, &TxTask{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
