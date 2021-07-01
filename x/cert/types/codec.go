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
	cdc.RegisterConcrete(MsgCertifyPlatform{}, "cert/CertifyPlatform", nil)
	cdc.RegisterConcrete(MsgIssueCertificate{}, "cert/IssueCertificate", nil)
	cdc.RegisterConcrete(CertifierUpdateProposal{}, "cert/CertifierUpdateProposal", nil)
	cdc.RegisterConcrete(MsgRevokeCertificate{}, "cert/RevokeCertificate", nil)
	cdc.RegisterConcrete(&Compilation{}, "cert/Compilation", nil)
	cdc.RegisterConcrete(&Auditing{}, "cert/Auditing", nil)
	cdc.RegisterConcrete(&Proof{}, "cert/Proof", nil)
	cdc.RegisterConcrete(&OracleOperator{}, "cert/OracleOperator", nil)
	cdc.RegisterConcrete(&ShieldPoolCreator{}, "cert/ShieldPoolCreator", nil)
	cdc.RegisterConcrete(&Identity{}, "cert/Identity", nil)
	cdc.RegisterConcrete(&General{}, "cert/General", nil)

	cdc.RegisterInterface((*Content)(nil), nil)
}

// RegisterInterfaces registers the x/cert interfaces types with the interface registry
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgProposeCertifier{},
		&MsgCertifyPlatform{},
		&MsgIssueCertificate{},
		&MsgRevokeCertificate{},
	)

	registry.RegisterImplementations((*govtypes.Content)(nil),
		&CertifierUpdateProposal{},
	)

	registry.RegisterInterface(
		"shentu.cert.v1alpha1.Content",
		(*Content)(nil),
		&Compilation{},
		&Auditing{},
		&Proof{},
		&OracleOperator{},
		&ShieldPoolCreator{},
		&Identity{},
		&General{},
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
