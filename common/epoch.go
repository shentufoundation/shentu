package common

const (
	BlocksPerMinute = uint64(12)
	BlocksPerHour   = BlocksPerMinute * 60
	BlocksPerDay    = BlocksPerHour * 24
	BlocksPerWeek   = BlocksPerDay * 7
	BlocksPerMonth  = BlocksPerDay * 30
	BlocksPerYear   = BlocksPerDay * 365

	BlocksPerEpoch = BlocksPerWeek

	SecondsPerBlock = uint64(5)
)
