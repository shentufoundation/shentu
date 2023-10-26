package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/bounty interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgCreateProgram{}, "bounty/CreateProgram", nil)
	cdc.RegisterConcrete(MsgEditProgram{}, "bounty/EditProgram", nil)
	cdc.RegisterConcrete(MsgCloseProgram{}, "bounty/CloseProgram", nil)
	cdc.RegisterConcrete(MsgSubmitFinding{}, "bounty/SubmitFinding", nil)
	cdc.RegisterConcrete(MsgAcceptFinding{}, "bounty/AcceptFinding", nil)
	cdc.RegisterConcrete(MsgRejectFinding{}, "bounty/RejectFinding", nil)
	cdc.RegisterConcrete(MsgCloseFinding{}, "bounty/CloseFinding", nil)
	cdc.RegisterConcrete(MsgReleaseFinding{}, "bounty/ReleaseFinding", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateProgram{},
		&MsgEditProgram{},
		&MsgCloseProgram{},
		&MsgSubmitFinding{},
		&MsgAcceptFinding{},
		&MsgRejectFinding{},
		&MsgCloseFinding{},
		&MsgReleaseFinding{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/bounty module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/bounty and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
