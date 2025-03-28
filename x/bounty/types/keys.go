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
	// Program related keys
	ProgramKeyPrefix      = collections.NewPrefix(1)
	FindingKeyPrefix      = collections.NewPrefix(2)
	ProgramFindingListKey = collections.NewPrefix(10)

	// Theorem related keys
	TheoremIDKey          = collections.NewPrefix(21)
	TheoremKeyPrefix      = collections.NewPrefix(22)
	ActiveTheoremQueueKey = collections.NewPrefix(23)

	// Proof related keys
	ProofKeyPrefix      = collections.NewPrefix(31)
	ActiveProofQueueKey = collections.NewPrefix(32)

	// Grant and deposit related keys
	GrantKeyPrefix   = collections.NewPrefix(41)
	DepositKeyPrefix = collections.NewPrefix(42)
	RewardKeyPrefix  = collections.NewPrefix(43)

	// Relationship keys
	TheoremProofPrefix   = collections.NewPrefix(51)
	ProofByTheoremPrefix = collections.NewPrefix(52)
	GrantByTheoremPrefix = collections.NewPrefix(53)

	// Parameter key
	ParamsKey = collections.NewPrefix(61)
)
