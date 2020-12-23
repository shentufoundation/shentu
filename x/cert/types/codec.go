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
	cdc.RegisterConcrete(MsgProposeCertifier{}, "cert/ProposeCertifier", nil)
	cdc.RegisterConcrete(MsgCertifyValidator{}, "cert/CertifyValidator", nil)
	cdc.RegisterConcrete(MsgDecertifyValidator{}, "cert/DecertifyValidator", nil)
	cdc.RegisterConcrete(MsgCertifyPlatform{}, "cert/CertifyPlatform", nil)
	cdc.RegisterConcrete(MsgCertifyGeneral{}, "cert/CertifyGeneral", nil)
	cdc.RegisterConcrete(MsgCertifyCompilation{}, "cert/CertifyCompilation", nil)
	cdc.RegisterConcrete(CertifierUpdateProposal{}, "cert/CertifierUpdateProposal", nil)
	cdc.RegisterConcrete(MsgRevokeCertificate{}, "cert/RevokeCertificate", nil)
	cdc.RegisterConcrete(&GeneralCertificate{}, "cert/GeneralCertificate", nil)
	cdc.RegisterConcrete(&CompilationCertificate{}, "cert/CompilationCertificate", nil)

	cdc.RegisterInterface((*Certificate)(nil), nil)
}

// RegisterInterfaces registers the x/oracle interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgProposeCertifier{},
		&MsgCertifyValidator{},
		&MsgDecertifyValidator{},
		&MsgCertifyPlatform{},
		&MsgCertifyGeneral{},
		&MsgCertifyCompilation{},
		&CertifierUpdateProposal{},
		&MsgRevokeCertificate{},
		&GeneralCertificate{},
		&CompilationCertificate{},
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

// // ModuleCdc defines the cert codec.
// var ModuleCdc *codec.Codec

// func init() {
// 	ModuleCdc = codec.New()
// 	RegisterCodec(ModuleCdc)
// 	codec.RegisterCrypto(ModuleCdc)
// 	ModuleCdc.Seal()
// }

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
