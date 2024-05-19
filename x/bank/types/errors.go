package types

import (
	sdkerrors "cosmossdk.io/errors"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

var (
	ErrCodeExists = sdkerrors.Register(bankTypes.ModuleName, 8, "can't perform multisend to a contract address")
)
