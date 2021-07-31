package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/execution/engine"
	"github.com/hyperledger/burrow/execution/native"
)

type CertificateCallable struct {
	ctx sdk.Context
}

const (
	// TODO: consolidate native contract gas consumption
	GasBase int64 = 1000
)

// registerCVMNative registers precompile contracts in CVM.
func registerCVMNative(cc CertificateCallable, nonce []byte) engine.Options {
	return engine.Options{
		Natives: native.MustDefaultNatives(),
		Nonce:   nonce,
	}
}
