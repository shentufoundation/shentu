package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc defines the cert codec.
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgProposeCertifier{}, "cert/ProposeCertifier", nil)
	cdc.RegisterConcrete(MsgCertifyValidator{}, "cert/CertifyValidator", nil)
	cdc.RegisterConcrete(MsgDecertifyValidator{}, "cert/DecertifyValidator", nil)
	cdc.RegisterConcrete(MsgCertifyPlatform{}, "cert/CertifyPlatform", nil)
	cdc.RegisterConcrete(MsgCertifyGeneral{}, "cert/CertifyGeneral", nil)
	cdc.RegisterConcrete(MsgCertifyCompilation{}, "cert/CertifyCompilation", nil)
	cdc.RegisterConcrete(CertifierUpdateProposal{}, "cert/CertifierUpdateProposal", nil)
	cdc.RegisterConcrete(MsgRevokeCertificate{}, "cert/RevokeCertificate", nil)
	cdc.RegisterInterface((*Certificate)(nil), nil)
	cdc.RegisterConcrete(&GeneralCertificate{}, "cert/GeneralCertificate", nil)
	cdc.RegisterConcrete(&CompilationCertificate{}, "cert/CompilationCertificate", nil)
}
