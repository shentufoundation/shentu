package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// RegisterCodec registers concrete types on codec.
func RegisterCodec(cdc *codec.Codec) {
	// Cosmos types
	cdc.RegisterConcrete(bank.MsgSend{}, "bank/MsgSend", nil)
	cdc.RegisterConcrete(bank.MsgMultiSend{}, "bank/MsgMultiSend", nil)

	// Custom types
	cdc.RegisterConcrete(MsgSendLock{}, "bank/MsgSendLock", nil)
}

// ModuleCdc defines module codec.
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
