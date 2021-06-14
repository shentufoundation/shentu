package types_test

import (
	"reflect"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/test-go/testify/require"

	"github.com/certikfoundation/shentu/x/oracle/types"
)

func Test_TaskParams(t *testing.T) {
	p1 := types.DefaultTaskParams()
	p2 := types.DefaultTaskParams()
	require.True(t, reflect.DeepEqual(p1, p2))

	p2 = types.NewTaskParams(time.Duration(24)*time.Hour, int64(40), sdk.NewInt(40), sdk.NewInt(40), sdk.NewInt(2), sdk.NewInt(200))
	require.False(t, reflect.DeepEqual(p1, p2))
}

func Test_LockedPoolParams(t *testing.T) {
	p1 := types.DefaultLockedPoolParams()
	p2 := types.DefaultLockedPoolParams()
	require.True(t, reflect.DeepEqual(p1, p2))

	p2 = types.NewLockedPoolParams(20, 4000)
	require.False(t, reflect.DeepEqual(p1, p2))
}
