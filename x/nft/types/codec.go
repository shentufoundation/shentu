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

	registry.RegisterImplementations((*exported.NFT)(nil),
		&nfttypes.BaseNFT{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
