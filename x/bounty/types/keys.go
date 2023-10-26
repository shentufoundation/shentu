package types

const (
	ModuleName = "bounty"

	RouterKey = ModuleName

	StoreKey = ModuleName

	QuerierRoute = ModuleName
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
func GetProgramKey(id string) []byte {
	return append(ProgramsKey, []byte(id)...)
}

// GetFindingKey creates the key for a program
func GetFindingKey(id string) []byte {
	return append(FindingKey, []byte(id)...)
}

func GetProgramIDFindingListKey(id string) []byte {
	return append(ProgramIDFindingListKey, []byte(id)...)
}
