package types

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/binary"
)

// ModuleCdc defines the cvm codec.
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCall{}, "cvm/Call", nil)
	cdc.RegisterConcrete(MsgDeploy{}, "cvm/Deploy", nil)
	cdc.RegisterConcrete(acm.Bytecode{}, "acm/Bytecode", nil)
	cdc.RegisterConcrete(binary.Word256{}, "binary/Word256", nil)
	cdc.RegisterConcrete([]acm.ContractMeta{}, "cvm/ContractMeta", nil)
}
