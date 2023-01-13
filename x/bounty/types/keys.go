package types

import (
	"strconv"
)

const (
	ModuleName = "bounty"

	RouterKey = ModuleName

	StoreKey = ModuleName

	QuerierRoute = ModuleName
)

var (
	ProgramsKey      = []byte{0x01}
	NextProgramIDKey = []byte{0x02}
)

// GetProgramKey creates the key for a program
// VALUE: staking/Validator
func GetProgramKey(id uint64) []byte {
	return append(ProgramsKey, []byte(strconv.FormatUint(id, 10))...)
}

// GetNextProgramIDKey creates the key for the validator with address
// VALUE: staking/Validator
func GetNextProgramIDKey() []byte {
	return NextProgramIDKey
}
