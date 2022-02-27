package v1alpha1

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
	cdc.RegisterConcrete(MsgCreatePool{}, "shentu/v1beta1/MsgCreatePool", nil)
	cdc.RegisterConcrete(MsgUpdatePool{}, "shentu/v1beta1/MsgUpdatePool", nil)
	cdc.RegisterConcrete(MsgPausePool{}, "shentu/v1beta1/MsgPausePool", nil)
	cdc.RegisterConcrete(MsgResumePool{}, "shentu/v1beta1/MsgResumePool", nil)
	cdc.RegisterConcrete(MsgDepositCollateral{}, "shentu/v1beta1/MsgDepositCollateral", nil)
	cdc.RegisterConcrete(MsgWithdrawCollateral{}, "shentu/v1beta1/MsgWithdrawCollateral", nil)
	cdc.RegisterConcrete(MsgWithdrawRewards{}, "shentu/v1beta1/MsgWithdrawRewards", nil)
	cdc.RegisterConcrete(MsgWithdrawForeignRewards{}, "shentu/v1beta1/MsgWithdrawForeignRewards", nil)
	cdc.RegisterConcrete(ShieldClaimProposal{}, "shentu/v1beta1/ShieldClaimProposal", nil)
	cdc.RegisterConcrete(MsgPurchaseShield{}, "shentu/v1beta1/MsgPurchaseShield", nil)
	cdc.RegisterConcrete(MsgWithdrawReimbursement{}, "shentu/v1beta1/MsgWithdrawReimbursement", nil)
	cdc.RegisterConcrete(MsgUpdateSponsor{}, "shentu/v1beta1/MsgUpdateSponsor", nil)
	cdc.RegisterConcrete(MsgStakeForShield{}, "shentu/v1beta1/MsgStakeForShield", nil)
	cdc.RegisterConcrete(MsgUnstakeFromShield{}, "shentu/v1beta1/MsgUnstakeFromShield", nil)
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
		&MsgPurchaseShield{},
		&MsgWithdrawReimbursement{},
		&MsgUpdateSponsor{},
		&MsgStakeForShield{},
		&MsgUnstakeFromShield{},
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
