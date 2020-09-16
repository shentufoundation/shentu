package vesting

import (
	"github.com/cosmos/cosmos-sdk/codec"
	vtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec) {
	// Register Cosmos types
	vtypes.RegisterCodec(cdc)

	// Register custom types
	cdc.RegisterConcrete(&TriggeredVestingAccount{}, "auth/TriggeredVestingAccount", nil)
	cdc.RegisterConcrete(&ManualVestingAccount{}, "auth/ManualVestingAccount", nil)
}
