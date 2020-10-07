package common

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	BlocksPerMinute = uint64(12)
	BlocksPerHour   = BlocksPerMinute * 60
	BlocksPerDay    = BlocksPerHour * 24
	BlocksPerWeek   = BlocksPerDay * 7
	BlocksPerMonth  = BlocksPerDay * 30
	BlocksPerYear   = BlocksPerDay * 365

	BlocksPerEpoch = BlocksPerWeek
)

var (
	BlocksPerSecondDec = sdk.NewDec(int64(BlocksPerMinute)).Quo(sdk.NewDec(60))
)
