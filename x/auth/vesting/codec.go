package vesting

import (
	"github.com/cosmos/cosmos-sdk/codec"
	vexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec) {
	// Copied from Cosmos
	cdc.RegisterInterface((*vexported.VestingAccount)(nil), nil)
	cdc.RegisterConcrete(&vtypes.BaseVestingAccount{}, "auth/BaseVestingAccount", nil)
	cdc.RegisterConcrete(&vtypes.ContinuousVestingAccount{}, "auth/ContinuousVestingAccount", nil)
	cdc.RegisterConcrete(&vtypes.DelayedVestingAccount{}, "auth/DelayedVestingAccount", nil)
	cdc.RegisterConcrete(&vtypes.PeriodicVestingAccount{}, "auth/PeriodicVestingAccount", nil)

	// Custom
	cdc.RegisterConcrete(&TriggeredVestingAccount{}, "auth/TriggeredVestingAccount", nil)
}
