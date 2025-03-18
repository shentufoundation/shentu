package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/bounty interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgCreateProgram{}, "bounty/CreateProgram", nil)
	cdc.RegisterConcrete(MsgEditProgram{}, "bounty/EditProgram", nil)
	cdc.RegisterConcrete(MsgActivateProgram{}, "bounty/ActivateProgram", nil)
	cdc.RegisterConcrete(MsgCloseProgram{}, "bounty/CloseProgram", nil)
	cdc.RegisterConcrete(MsgSubmitFinding{}, "bounty/SubmitFinding", nil)
	cdc.RegisterConcrete(MsgEditFinding{}, "bounty/EditFinding", nil)
	cdc.RegisterConcrete(MsgConfirmFinding{}, "bounty/ConfirmFinding", nil)
	cdc.RegisterConcrete(MsgConfirmFindingPaid{}, "bounty/ConfirmFindingPaid", nil)
	cdc.RegisterConcrete(MsgActivateFinding{}, "bounty/ActivateFinding", nil)
	cdc.RegisterConcrete(MsgCloseFinding{}, "bounty/CloseFinding", nil)
	cdc.RegisterConcrete(MsgPublishFinding{}, "bounty/PublishFinding", nil)
	cdc.RegisterConcrete(MsgCreateTheorem{}, "bounty/CreateTheorem", nil)
	cdc.RegisterConcrete(MsgSubmitProofHash{}, "bounty/SubmitProofHash", nil)
	cdc.RegisterConcrete(MsgSubmitProofDetail{}, "bounty/SubmitProofDetail", nil)
	cdc.RegisterConcrete(MsgSubmitProofVerification{}, "bounty/SubmitProofVerification", nil)
	cdc.RegisterConcrete(MsgGrant{}, "bounty/Grant", nil)
	cdc.RegisterConcrete(MsgWithdrawReward{}, "bounty/WithdrawReward", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateProgram{},
		&MsgEditProgram{},
		&MsgActivateProgram{},
		&MsgCloseProgram{},
		&MsgSubmitFinding{},
		&MsgEditFinding{},
		&MsgConfirmFinding{},
		&MsgConfirmFindingPaid{},
		&MsgActivateFinding{},
		&MsgCloseFinding{},
		&MsgPublishFinding{},
		&MsgCreateTheorem{},
		&MsgSubmitProofHash{},
		&MsgSubmitProofDetail{},
		&MsgSubmitProofVerification{},
		&MsgGrant{},
		&MsgWithdrawReward{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
