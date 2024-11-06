package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&banktypes.MsgSend{}, "bank/MsgSend", nil)
	cdc.RegisterConcrete(&banktypes.MsgMultiSend{}, "bank/MsgMultiSend", nil)
	cdc.RegisterConcrete(&MsgLockedSend{}, "bank/MsgLockedSend", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&banktypes.MsgSend{},
		&banktypes.MsgMultiSend{},
		&MsgLockedSend{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
