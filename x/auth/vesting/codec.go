package vesting

import (
	"github.com/cosmos/cosmos-sdk/codec"
	vtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

// RegisterLegacyAminoCodec registers the vesting interfaces and concrete types on the
// provided LegacyAmino codec. These types are used for Amino JSON serialization
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	vtypes.RegisterLegacyAminoCodec(cdc)
	cdc.RegisterConcrete(&ManualVestingAccount{}, "auth/ManualVestingAccount", nil)
}

// TODO RegisterInterface?

var amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
