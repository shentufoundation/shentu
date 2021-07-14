package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/engine"
	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/native"
	"github.com/hyperledger/burrow/permission"

	"github.com/certikfoundation/shentu/x/cvm/types"
)

type CertificateCallable struct {
	ctx        sdk.Context
	certKeeper types.CertKeeper
}

const (
	// TODO: consolidate native contract gas consumption
	GasBase int64 = 1000
)

// registerCVMNative registers precompile contracts in CVM.
func registerCVMNative(cc CertificateCallable, nonce []byte) engine.Options {
	return engine.Options{
		Natives: native.MustDefaultNatives().
			MustFunction("General", leftPadAddress(101), permission.None, cc.checkGeneral).
			MustFunction("Proof", leftPadAddress(102), permission.None, cc.checkProof).
			MustFunction("Compilation", leftPadAddress(103), permission.None, cc.checkCompilation),
		Nonce: nonce,
	}
}

// checkGeneral checks if certificates for a given content exists.
func (cc CertificateCallable) checkGeneral(ctx native.Context) (output []byte, err error) {
	input := string(ctx.Input)
	if cc.certKeeper.IsContentCertified(cc.ctx, input) {
		return []byte{0x01}, nil
	}
	return []byte{0x00}, nil
}

// checkProof checks if proof certificates for a given content exists.
func (cc CertificateCallable) checkProof(ctx native.Context) (output []byte, err error) {
	return cc.checkFunc(ctx, "proof")
}

// checkFunc checks if a certain type of certificates for a given content exists.
func (cc CertificateCallable) checkFunc(ctx native.Context, certType string) ([]byte, error) {
	gasRequired := big.NewInt(GasBase)
	if ctx.Gas.Cmp(gasRequired) == -1 {
		return nil, errors.Codes.InsufficientGas
	} else {
		*ctx.Gas = *ctx.Gas.Sub(ctx.Gas, gasRequired)
	}
	input := string(ctx.Input)
	addr, err := sdk.AccAddressFromBech32(input)
	if err != nil {
		return []byte{0x00}, err
	}
	if cc.certKeeper.IsCertified(cc.ctx, addr.String(), certType) {
		return []byte{0x01}, nil
	}
	return []byte{0x00}, nil
}

// checkCompilation checks if compilation certificates for a given content exists.
func (cc CertificateCallable) checkCompilation(ctx native.Context) (output []byte, err error) {
	gasRequired := big.NewInt(GasBase)
	if ctx.Gas.Cmp(gasRequired) == -1 {
		return nil, errors.Codes.InsufficientGas
	} else {
		*ctx.Gas = *ctx.Gas.Sub(ctx.Gas, gasRequired)
	}
	input := string(ctx.Input)
	if cc.certKeeper.IsCertified(cc.ctx, input, "compilation") {
		return []byte{0x01}, nil
	}
	return []byte{0x00}, nil
}

func leftPadAddress(bs ...byte) crypto.Address {
	return crypto.AddressFromWord256(binary.LeftPadWord256(bs))
}
