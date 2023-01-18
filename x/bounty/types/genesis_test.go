package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultGenesisState(t *testing.T) {
	var startingProgramId uint64 = 1
	state1 := DefaultGenesisState()
	require.Equal(t, state1.StartingProgramId, startingProgramId)

	// TODO add GenesisState equal
}
