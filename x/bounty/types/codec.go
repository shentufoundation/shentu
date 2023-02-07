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
	cdc.RegisterConcrete(MsgSubmitFinding{}, "bounty/SubmitFinding", nil)
	cdc.RegisterConcrete(MsgHostAcceptFinding{}, "bounty/HostAcceptFinding", nil)
	cdc.RegisterConcrete(MsgHostRejectFinding{}, "bounty/HostRejectFinding", nil)
	cdc.RegisterConcrete(MsgReleaseFinding{}, "bounty/MsgReleaseFinding", nil)

	cdc.RegisterConcrete(&EciesPubKey{}, "bounty/EciesPubKey", nil)

	cdc.RegisterConcrete(&EciesEncryptedDesc{}, "bounty/EciesEncryptedDesc", nil)
	cdc.RegisterConcrete(&EciesEncryptedPoc{}, "bounty/EciesEncryptedPoc", nil)
	cdc.RegisterConcrete(&EciesEncryptedComment{}, "bounty/EciesEncryptedComment", nil)
	cdc.RegisterConcrete(&PlainTextDesc{}, "bounty/PlainTextDesc", nil)
	cdc.RegisterConcrete(&PlainTextPoc{}, "bounty/PlainTextPoc", nil)
	cdc.RegisterConcrete(&PlainTextComment{}, "bounty/PlainTextComment", nil)

	cdc.RegisterInterface((*EncryptionKey)(nil), nil)
	cdc.RegisterInterface((*FindingDesc)(nil), nil)
	cdc.RegisterInterface((*FindingPoc)(nil), nil)
	cdc.RegisterInterface((*FindingComment)(nil), nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateProgram{},
		&MsgSubmitFinding{},
		&MsgHostAcceptFinding{},
		&MsgHostRejectFinding{},
		&MsgReleaseFinding{},
	)

	registry.RegisterInterface(
		"shentu.bounty.v1.EncryptionKey",
		(*EncryptionKey)(nil),
		&EciesPubKey{},
	)

	registry.RegisterInterface(
		"shentu.bounty.v1.FindingDesc",
		(*FindingDesc)(nil),
		&EciesEncryptedDesc{},
		&PlainTextDesc{},
	)

	registry.RegisterInterface(
		"shentu.bounty.v1.FindingPoc",
		(*FindingPoc)(nil),
		&EciesEncryptedPoc{},
		&PlainTextPoc{},
	)

	registry.RegisterInterface(
		"shentu.bounty.v1.FindingComment",
		(*FindingComment)(nil),
		&EciesEncryptedComment{},
		&PlainTextComment{},
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
