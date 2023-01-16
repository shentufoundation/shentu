package types

import (
	"strconv"
)

const (
	ModuleName = "bounty"

	RouterKey = ModuleName

	StoreKey = ModuleName
)

var (
	ProgramsKey      = []byte{0x01}
	NextProgramIDKey = []byte{0x02}
	FindingKey       = []byte{0x03}
	NextFindingIDKey = []byte{0x04}

	ProgramIDFindingListKey = []byte{0x10}
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

// GetFindingKey creates the key for a program
func GetFindingKey(id uint64) []byte {
	return append(FindingKey, []byte(strconv.FormatUint(id, 10))...)
}

// GetNextFindingIDKey creates the key for the validator with address
func GetNextFindingIDKey() []byte {
	return NextFindingIDKey
}

func GetProgramIDFindingListKey(id uint64) []byte {
	return append(ProgramIDFindingListKey, []byte(strconv.FormatUint(id, 10))...)
}
