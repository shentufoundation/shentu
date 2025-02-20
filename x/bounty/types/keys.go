package types

import (
	"cosmossdk.io/collections"
)

const (
	ModuleName = "bounty"

	RouterKey = ModuleName

	StoreKey = ModuleName

	QuerierRoute = ModuleName
)

var (
	ProgramKey = collections.NewPrefix(1)
	FindingKey = collections.NewPrefix(2)

	ProgramIDFindingListKey = collections.NewPrefix(10)

	// TheoremIDKey stores the sequence representing the next theorem ID.
	TheoremIDKey = collections.NewPrefix(11)
	// TheoremsKeyProofPKeyPrefix stores the theorem raw bytes.
	TheoremsKeyProofPKeyPrefix = collections.NewPrefix(12)
	// ActiveTheoremQueuePrefix stores the active theorems.
	ActiveTheoremQueuePrefix = collections.NewPrefix(13)
	// InactiveTheoremQueuePrefix stores the inactive theorems.
	InactiveTheoremQueuePrefix = collections.NewPrefix(14)
	// ProofPeriodProposalKeyPrefix stores which theorem are on proof period.
	ProofPeriodProposalKeyPrefix = collections.NewPrefix(15)
	// ProofsKeyPrefix stores the proof raw bytes.
	ProofsKeyPrefix = collections.NewPrefix(16)
	// ActiveProofQueuePrefix stores the active proofs.
	ActiveProofQueuePrefix = collections.NewPrefix(17)
	// GrantsKeyPrefix stores grants.
	GrantsKeyPrefix = collections.NewPrefix(18)

	ProofByTheoremIndexKey = collections.NewPrefix(31) // key for proofs by a theorem
	GrantByTheoremIndexKey = collections.NewPrefix(32) // key for grants by a theorem
	GrantByAddressIndexKey = collections.NewPrefix(33) // key for grants by an address
	// ParamsKey stores the module's params.
	ParamsKey = collections.NewPrefix(41)
)

// GetProgramKey creates the key for a program
// VALUE: staking/Validator
func GetProgramKey(id string) []byte {
	return append(ProgramKey, []byte(id)...)
}

// GetFindingKey creates the key for a program
func GetFindingKey(id string) []byte {
	return append(FindingKey, []byte(id)...)
}

func GetProgramIDFindingListKey(id string) []byte {
	return append(ProgramIDFindingListKey, []byte(id)...)
}
