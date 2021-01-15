package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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

// RegisterInterfaces registers the x/cert interfaces types with the interface registry
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {	
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgProposeCertifier{},
		&MsgCertifyValidator{},
		&MsgDecertifyValidator{},
		&MsgCertifyPlatform{},
		&MsgCertifyGeneral{},
		&MsgCertifyCompilation{},
		&MsgRevokeCertificate{},
	)

	registry.RegisterImplementations((*govtypes.Content)(nil),
		&CertifierUpdateProposal{},
	)

	registry.RegisterImplementations((*Certificate)(nil),
		&GeneralCertificate{},
		&CompilationCertificate{},
	)

	registry.RegisterInterface("shentu.cert.v1alpha1.Certificate", (*Certificate)(nil),
		&GeneralCertificate{},
		&CompilationCertificate{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/cert module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/cert and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
