package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers concrete types on the LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgCreatePool{}, "shield/MsgCreatePool", nil)
	cdc.RegisterConcrete(MsgUpdatePool{}, "shield/MsgUpdatePool", nil)
	cdc.RegisterConcrete(MsgPausePool{}, "shield/MsgPausePool", nil)
	cdc.RegisterConcrete(MsgResumePool{}, "shield/MsgResumePool", nil)
	cdc.RegisterConcrete(MsgDepositCollateral{}, "shield/MsgDepositCollateral", nil)
	cdc.RegisterConcrete(MsgWithdrawCollateral{}, "shield/MsgWithdrawCollateral", nil)
	cdc.RegisterConcrete(MsgWithdrawRewards{}, "shield/MsgWithdrawRewards", nil)
	cdc.RegisterConcrete(MsgWithdrawForeignRewards{}, "shield/MsgWithdrawForeignRewards", nil)
	cdc.RegisterConcrete(ShieldClaimProposal{}, "shield/ShieldClaimProposal", nil)
	cdc.RegisterConcrete(MsgPurchase{}, "shield/MsgPurchaseShield", nil)
	cdc.RegisterConcrete(MsgUpdateSponsor{}, "shield/MsgUpdateSponsor", nil)
	cdc.RegisterConcrete(MsgUnstake{}, "shield/MsgUnstakeFromShield", nil)
}

// RegisterInterfaces registers the x/shield interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePool{},
		&MsgUpdatePool{},
		&MsgPausePool{},
		&MsgResumePool{},
		&MsgDepositCollateral{},
		&MsgWithdrawCollateral{},
		&MsgWithdrawRewards{},
		&MsgWithdrawForeignRewards{},
		&MsgPurchase{},
		&MsgUpdateSponsor{},
		&MsgUnstake{},
	)
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&ShieldClaimProposal{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/shield module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/shield and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
