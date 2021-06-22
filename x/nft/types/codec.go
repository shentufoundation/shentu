package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	"github.com/irisnet/irismod/modules/nft/exported"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

// RegisterLegacyAminoCodec concrete types on codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	nfttypes.RegisterLegacyAminoCodec(cdc)
	cdc.RegisterConcrete(&MsgCreateAdmin{}, "nft/MsgCreateAdmin", nil)
	cdc.RegisterConcrete(&MsgRevokeAdmin{}, "nft/MsgRevokeAdmin", nil)
	cdc.RegisterConcrete(MsgIssueCertificate{}, "nft/IssueCertificate", nil)
	cdc.RegisterConcrete(MsgRevokeCertificate{}, "nft/RevokeCertificate", nil)
	cdc.RegisterConcrete(&Compilation{}, "nft/Compilation", nil)
	cdc.RegisterConcrete(&Auditing{}, "nft/Auditing", nil)
	cdc.RegisterConcrete(&Proof{}, "nft/Proof", nil)
	cdc.RegisterConcrete(&OracleOperator{}, "nft/OracleOperator", nil)
	cdc.RegisterConcrete(&ShieldPoolCreator{}, "nft/ShieldPoolCreator", nil)
	cdc.RegisterConcrete(&Identity{}, "nft/Identity", nil)
	cdc.RegisterConcrete(&General{}, "nft/General", nil)

	cdc.RegisterInterface((*Content)(nil), nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&nfttypes.MsgIssueDenom{},
		&nfttypes.MsgTransferNFT{},
		&nfttypes.MsgEditNFT{},
		&nfttypes.MsgMintNFT{},
		&nfttypes.MsgBurnNFT{},
		&MsgCreateAdmin{},
		&MsgRevokeAdmin{},
		&MsgIssueCertificate{},
		&MsgRevokeCertificate{},
	)

	registry.RegisterInterface(
		"shentu.nft.v1alpha1.Content",
		(*Content)(nil),
		&Compilation{},
		&Auditing{},
		&Proof{},
		&OracleOperator{},
		&ShieldPoolCreator{},
		&Identity{},
		&General{},
	)

	registry.RegisterImplementations((*exported.NFT)(nil),
		&nfttypes.BaseNFT{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
