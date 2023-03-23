package types_test

import (
	"reflect"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/test-go/testify/require"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func Test_TaskParams(t *testing.T) {
	p1 := types.DefaultTaskParams()
	p2 := types.DefaultTaskParams()
	p3 := types.NewTaskParams(time.Duration(24)*time.Hour, int64(40), sdk.NewInt(40), sdk.NewInt(40), sdk.NewInt(2), sdk.NewInt(200), sdk.NewDecWithPrec(50, 2))

	require.True(t, reflect.DeepEqual(p1, p2))
	require.False(t, reflect.DeepEqual(p1, p3))
}

func Test_LockedPoolParams(t *testing.T) {
	p1 := types.DefaultLockedPoolParams()
	p2 := types.DefaultLockedPoolParams()
	p3 := types.NewLockedPoolParams(20, 4000)

	require.True(t, reflect.DeepEqual(p1, p2))
	require.False(t, reflect.DeepEqual(p1, p3))
}
