package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers the account types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreatePool{}, "shield/MsgCreatePool", nil)
	cdc.RegisterConcrete(MsgUpdatePool{}, "shield/MsgUpdatePool", nil)
	cdc.RegisterConcrete(MsgPausePool{}, "shield/MsgPausePool", nil)
	cdc.RegisterConcrete(MsgResumePool{}, "shield/MsgResumePool", nil)
	cdc.RegisterConcrete(MsgDepositCollateral{}, "shield/MsgDepositCollateral", nil)
	cdc.RegisterConcrete(ShieldClaimProposal{}, "shield/ShieldClaimProposal", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	ModuleCdc = cdc.Seal()
}
