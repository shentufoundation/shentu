package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/irisnet/irismod/modules/nft/types"
)

var (
	ErrAdminNotFound = sdkerrors.Register(types.ModuleName, 13, "nft admin not found")
)