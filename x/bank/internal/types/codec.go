package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/x/bank/exported"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*exported.SupplyI)(nil), nil)
	cdc.RegisterConcrete(&bankTypes.Supply{}, "bank/Supply", nil)
	cdc.RegisterConcrete(&bankTypes.MsgSend{}, "bank/MsgSend", nil)
	cdc.RegisterConcrete(&bankTypes.MsgMultiSend{}, "bank/MsgMultiSend", nil)
	cdc.RegisterConcrete(&MsgLockedSend{}, "bank/MsgLockedSend", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&bankTypes.MsgSend{},
		&bankTypes.MsgMultiSend{},
		&MsgLockedSend{},
	)

	registry.RegisterInterface(
		"cosmos.bank.v1beta1.SupplyI",
		(*exported.SupplyI)(nil),
		&bankTypes.Supply{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/bank module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
