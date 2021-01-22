package cvm_test

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hyperledger/burrow/txs/payload"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/cvm"
	"github.com/certikfoundation/shentu/x/cvm/keeper"
)

func TestExportGenesis(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(10000))
	k := app.CVMKeeper

	code, err := hex.DecodeString(keeper.BasicTestsBytecodeString)
	require.Nil(t, err)

	_, _ = k.Tx(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
	exported := cvm.ExportGenesis(ctx, k)

	app2 := simapp.Setup(false)
	ctx2 := app2.BaseApp.NewContext(false, tmproto.Header{})
	k2 := app2.CVMKeeper

	cvm.InitGenesis(ctx2, k2, *exported)
	exported2 := cvm.ExportGenesis(ctx, k)

	fmt.Println(exported)
	fmt.Println(exported2)

	require.True(t, reflect.DeepEqual(exported, exported2))
}
