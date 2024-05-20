package oracle_test

import (
	"reflect"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/oracle"
)

func TestExportGenesis(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	k := app.OracleKeeper

	exported := oracle.ExportGenesis(ctx, k)

	app2 := shentuapp.Setup(t, false)
	ctx2 := app2.BaseApp.NewContext(false, tmproto.Header{})
	k2 := app2.OracleKeeper

	oracle.InitGenesis(ctx2, k2, exported)
	exported2 := oracle.ExportGenesis(ctx, k)
	require.True(t, reflect.DeepEqual(exported, exported2))
}
