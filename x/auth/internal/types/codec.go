package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

var ModuleCdc = codec.New()

func RegisterCodec(cdc *codec.Codec) {
	// Copied from Cosmos
	cdc.RegisterInterface((*exported.GenesisAccount)(nil), nil)
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "cosmos-sdk/BaseAccount", nil)
	cdc.RegisterConcrete(auth.StdTx{}, "cosmos-sdk/StdTx", nil)

	// Custom types
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
