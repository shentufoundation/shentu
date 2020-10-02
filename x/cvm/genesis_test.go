package cvm

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hyperledger/burrow/txs/payload"

	"github.com/certikfoundation/shentu/x/cvm/internal/keeper"
)

func TestExportGenesis(t *testing.T) {
	testInput := keeper.CreateTestInput(t)
	ctx := testInput.Ctx
	k := testInput.CvmKeeper

	code, err := hex.DecodeString(keeper.BasicTestsBytecodeString)
	require.Nil(t, err)

	_, _ = k.Call(ctx, keeper.Addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
	exported := ExportGenesis(ctx, k)

	testInput2 := keeper.CreateTestInput(t)
	ctx2 := testInput2.Ctx
	k2 := testInput2.CvmKeeper
	InitGenesis(ctx2, k2, exported)
	exported2 := ExportGenesis(ctx, k)

	fmt.Println(exported)
	fmt.Println(exported2)

	require.True(t, reflect.DeepEqual(exported, exported2))
}
