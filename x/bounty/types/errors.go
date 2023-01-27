package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	ErrorEmptyProgramIDFindingList = "empty finding id list"
)

// Finding
const (
	errFindingNotExists uint32 = iota + 201
)

// [2xx] Finding
var (
	ErrFindingNotExists = sdkerrors.Register(ModuleName, errFindingNotExists, "finding does not exist")
)
