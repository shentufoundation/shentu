package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func TestQueries(t *testing.T) {
	var (
		programPath  = []string{types.QueryProgram}
		programsPath = []string{types.QueryPrograms}
	)

	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	legacyQuerierCdc := app.LegacyAmino()
	querier := keeper.NewQuerier(app.BountyKeeper, legacyQuerierCdc)

	// default
	req := abci.RequestQuery{
		Data: []byte{},
		Path: "",
	}
	bz, err := querier(ctx, []string{"other"}, req)
	require.Error(t, err)
	require.Nil(t, bz)

	// program
	req = abci.RequestQuery{
		Data: []byte{},
		Path: fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryProgram),
	}
	bz, err = querier(ctx, programPath, req)
	require.Error(t, err)
	require.Nil(t, bz)

	// programs
	req = abci.RequestQuery{
		Data: []byte{},
		Path: fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPrograms),
	}

	bz, err = querier(ctx, programsPath, req)
	require.Error(t, err)
	require.Nil(t, bz)
}
