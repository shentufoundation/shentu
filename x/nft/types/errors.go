package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	nfttypes "github.com/irisnet/irismod/modules/nft/types"
)

var (
	ErrAdminNotFound = sdkerrors.Register(nfttypes.ModuleName, 13, "nft admin not found")
)
