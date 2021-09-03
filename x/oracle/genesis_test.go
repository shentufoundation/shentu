package oracle_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/certikfoundation/shentu/v2/simapp"
	"github.com/certikfoundation/shentu/v2/x/oracle"
)

func TestExportGenesis(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	k := app.OracleKeeper

	exported := oracle.ExportGenesis(ctx, k)

	app2 := simapp.Setup(false)
	ctx2 := app2.BaseApp.NewContext(false, tmproto.Header{})
	k2 := app2.OracleKeeper

	oracle.InitGenesis(ctx2, k2, exported)
	exported2 := oracle.ExportGenesis(ctx, k)
	require.True(t, reflect.DeepEqual(exported, exported2))
}
