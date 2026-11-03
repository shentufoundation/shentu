package app

import (
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"

	legacywasm "github.com/shentufoundation/shentu/v2/app/legacy/wasm"
	"github.com/shentufoundation/shentu/v2/app/params"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	// Historical governance proposals still contain this Any after x/wasm
	// removal. Register its wire type without restoring a message handler.
	encodingConfig.InterfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &legacywasm.MsgUpdateParams{})
	return encodingConfig
}
