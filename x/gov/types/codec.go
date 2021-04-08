package types

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&govtypes.MsgSubmitProposal{},
		&govtypes.MsgVote{},
		&govtypes.MsgDeposit{},
	)
	registry.RegisterInterface(
		"cosmos.gov.v1beta1.Content",
		(*govtypes.Content)(nil),
		&govtypes.TextProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
