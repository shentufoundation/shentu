package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/native"
	"github.com/hyperledger/burrow/permission"

	"github.com/certikfoundation/shentu/vm"
	certtypes "github.com/certikfoundation/shentu/x/cert/types"
	"github.com/certikfoundation/shentu/x/cvm/types"
)

type CertificateCallable struct {
	ctx        sdk.Context
	certKeeper types.CertKeeper
}

const (
	// TODO: consolidate native contract gas consumption
	GasBase uint64 = 1000
)

// registerCVMNative registers precompile contracts in CVM.
func registerCVMNative(options *vm.CVMOptions, cc CertificateCallable) {
	options.Natives = native.MustDefaultNatives().
		MustFunction("General", leftPadAddress(9), permission.None, cc.checkGeneral).
		MustFunction("Proof", leftPadAddress(10), permission.None, cc.checkProof).
		MustFunction("Compilation", leftPadAddress(11), permission.None, cc.checkCompilation).
		MustFunction("CertifyValidator", leftPadAddress(12), permission.None, cc.certifyValidator)
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
	gasRequired := GasBase
	if *ctx.Gas < gasRequired {
		return nil, errors.Codes.InsufficientGas
	} else {
		*ctx.Gas -= gasRequired
	}
	input := string(ctx.Input)
	addr, err := sdk.AccAddressFromBech32(input)
	if err != nil {
		return []byte{0x00}, err
	}
	if cc.certKeeper.IsCertified(cc.ctx, "address", addr.String(), certType) {
		return []byte{0x01}, nil
	}
	return []byte{0x00}, nil
}

// checkCompilation checks if compilation certificates for a given content exists.
func (cc CertificateCallable) checkCompilation(ctx native.Context) (output []byte, err error) {
	gasRequired := GasBase
	if *ctx.Gas < gasRequired {
		return nil, errors.Codes.InsufficientGas
	} else {
		*ctx.Gas -= gasRequired
	}
	input := string(ctx.Input)
	if cc.certKeeper.IsCertified(cc.ctx, "sourcecodehash", input, "compilation") {
		return []byte{0x01}, nil
	}
	return []byte{0x00}, nil
}

// certifyValidator certifies a validator.
func (cc CertificateCallable) certifyValidator(ctx native.Context) (output []byte, err error) {
	gasRequired := GasBase
	if *ctx.Gas < gasRequired {
		return nil, errors.Codes.InsufficientGas
	} else {
		*ctx.Gas -= gasRequired
	}
	if !cc.certKeeper.IsCertifier(cc.ctx, ctx.Origin.Bytes()) {
		return nil, certtypes.ErrUnqualifiedCertifier
	}
	input := string(ctx.Input)
	pubKey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, input)
	if err != nil {
		return []byte{0x00}, err
	}
	cc.certKeeper.SetValidator(cc.ctx, pubKey, ctx.Caller.Bytes())
	return []byte{0x01}, nil
}

func leftPadAddress(bs ...byte) crypto.Address {
	return crypto.AddressFromWord256(binary.LeftPadWord256(bs))
}
