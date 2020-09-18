package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

var ModuleCdc = codec.New()

func RegisterCodec(cdc *codec.Codec) {
	// Register Cosmos types
	types.RegisterCodec(cdc)

	// Register custom types
	cdc.RegisterConcrete(MsgTriggerVesting{}, "auth/MsgTriggerVesting", nil)
	cdc.RegisterConcrete(MsgManualVesting{}, "auth/MsgManualVesting", nil)
}

func RegisterAccountTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
}
