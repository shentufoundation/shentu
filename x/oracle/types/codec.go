package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
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
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/oracle module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/oracle and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
