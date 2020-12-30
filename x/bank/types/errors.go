package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

var (
	ErrCodeExists = sdkerrors.Register(bankTypes.ModuleName, 6, "can't perform multisend to a contract address")
)
