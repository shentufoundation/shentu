package types

import (
	//sdkerrors "cosmossdk.io/errors"
	sdkerrors "cosmossdk.io/errors"

	"github.com/hyperledger/burrow/execution/errors"
)

// BurrowErrorCodeStart is the default sdk code type.
const BurrowErrorCodeStart = 200

// ErrCodedError wraps execution CodedError into sdk Error.
func ErrCodedError(error errors.CodedError) *sdkerrors.Error {
	return sdkerrors.New(ModuleName, BurrowErrorCodeStart+error.ErrorCode().Number, error.ErrorCode().Name)
}
