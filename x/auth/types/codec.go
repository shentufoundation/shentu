package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
)

// RegisterLegacyAminoCodec registers the necessary x/auth interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUnlock{}, "auth/MsgUnlock", nil)
	cdc.RegisterConcrete(&ManualVestingAccount{}, "auth/ManualVestingAccount", nil)
}

// RegisterInterfaces registers the x/auth interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUnlock{},
	)

	registry.RegisterInterface(
		"cosmos.vesting.v1beta1.VestingAccount",
		(*exported.VestingAccount)(nil),
		&ManualVestingAccount{},
	)

	registry.RegisterImplementations(
		(*authtypes.AccountI)(nil),
		&ManualVestingAccount{},
	)

	registry.RegisterImplementations(
		(*authtypes.GenesisAccount)(nil),
		&ManualVestingAccount{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/auth module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/auth and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
