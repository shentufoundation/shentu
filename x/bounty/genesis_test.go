package bounty_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/bounty"
)

func TestExportGenesis(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	k := app.BountyKeeper

	exported := bounty.ExportGenesis(ctx, k)

	app2 := shentuapp.Setup(false)
	ctx2 := app2.BaseApp.NewContext(false, tmproto.Header{})
	k2 := app2.BountyKeeper

	bounty.InitGenesis(ctx2, k2, *exported)
	exported2 := bounty.ExportGenesis(ctx, k)
	require.True(t, reflect.DeepEqual(exported, exported2))
}
