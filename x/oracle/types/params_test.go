package types_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/test-go/testify/require"

	"github.com/certikfoundation/shentu/x/oracle/types"
)

func Test_TaskParams(t *testing.T) {
	p1 := types.DefaultTaskParams()
	p2 := types.DefaultTaskParams()
	ok := equalsForTastParams(&p1, &p2)
	require.True(t, ok)

	p2 = types.NewTaskParams(time.Duration(24)*time.Hour, int64(40), sdk.NewInt(40), sdk.NewInt(40), sdk.NewInt(2), sdk.NewInt(200))

	ok = equalsForTastParams(&p1, &p2)
	require.False(t, ok)
}

func Test_LockedPoolParams(t *testing.T) {
	p1 := types.DefaultLockedPoolParams()
	p2 := types.DefaultLockedPoolParams()

	ok := equalsForLockedPoolParams(&p1, &p2)
	require.True(t, ok)

	p2 = types.NewLockedPoolParams(20, 4000)
	ok = equalsForLockedPoolParams(&p1, &p2)
	require.False(t, ok)
}

func equalsForTastParams(t1 *types.TaskParams, t2 *types.TaskParams) bool {
	if t2 == nil {
		return t1 == nil
	} else if t1 == nil {
		return false
	}
	if t1.ExpirationDuration != t2.ExpirationDuration {
		return false
	}
	if t1.AggregationWindow != t2.AggregationWindow {
		return false
	}
	if t1.AggregationResult != t2.AggregationResult {
		return false
	}
	if t1.ThresholdScore != t2.ThresholdScore {
		return false
	}
	if t1.Epsilon1 != t2.Epsilon1 {
		return false
	}
	if !t1.Epsilon2.Equal(t2.Epsilon2) {
		return false
	}
	return true
}

func equalsForLockedPoolParams(t1 *types.LockedPoolParams, t2 *types.LockedPoolParams) bool {
	if t2 == nil {
		return t1 == nil
	} else if t1 == nil {
		return false
	}
	if t1.LockedInBlocks != t2.LockedInBlocks {
		return false
	}
	if t1.MinimumCollateral != t2.MinimumCollateral {
		return false
	}
	return true
}
