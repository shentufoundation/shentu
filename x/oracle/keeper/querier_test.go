package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	. "github.com/shentufoundation/shentu/v2/x/oracle/keeper"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func TestQueryOperators(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, sdk.NewInt(80000*1e6))
	ok1 := app.OracleKeeper

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := NewQuerier(ok1, app.LegacyAmino())

	params := ok1.GetLockedPoolParams(ctx)
	err := ok1.CreateOperator(ctx, addrs[1], sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}, addrs[0], "operator1")
	require.Nil(t, err)

	path := []string{"operator", addrs[1].String()}
	bz, err := querier(ctx, path, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	ok2 := app.OracleKeeper
	querier = NewQuerier(ok2, app.LegacyAmino())

	params = ok2.GetLockedPoolParams(ctx)
	err = ok2.CreateOperator(ctx, addrs[2], sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}, addrs[0], "operator2")
	require.Nil(t, err)

	path = []string{"operators"}
	bz, err = querier(ctx, path, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	err = ok2.RemoveOperator(ctx, addrs[2].String(), addrs[2].String())
	require.Nil(t, err)

	path = []string{"withdraws"}
	bz, err = querier(ctx, path, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	path = []string{"operator", addrs[1].String()}
	bz, err = querier(ctx, path, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	path = []string{"operator", addrs[2].String()}
	_, err = querier(ctx, path, query)
	require.Error(t, err)
}

func TestQueryTask(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	legacyQuerierCdc := app.LegacyAmino()
	querier := NewQuerier(ok, legacyQuerierCdc)

	params := ok.GetLockedPoolParams(ctx)
	err := ok.CreateOperator(ctx, addrs[0], sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}, addrs[0], "operator")
	require.Nil(t, err)

	contract := "0x1234567890abcdef"
	function := "func"
	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 5000)}
	description := "testing"
	expiration := time.Now().Add(time.Hour).UTC()
	waitingBlocks := int64(50)

	err = ok.CreateTask(ctx, contract, function, bounty, description, expiration, addrs[0], waitingBlocks)
	require.Nil(t, err)

	taskParams := types.QueryTaskParams{
		Contract: contract,
		Function: function,
	}
	data := legacyQuerierCdc.MustMarshalJSON(taskParams)
	query := abci.RequestQuery{
		Path: "",
		Data: data,
	}

	path := []string{"task"}
	bz, err := querier(ctx, path, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	err = ok.RespondToTask(ctx, contract, function, 20, addrs[0])
	require.Nil(t, err)

	responseParams := types.QueryResponseParams{
		Contract: contract,
		Function: function,
		Operator: addrs[0],
	}
	data = legacyQuerierCdc.MustMarshalJSON(responseParams)
	query = abci.RequestQuery{
		Path: "",
		Data: data,
	}

	path = []string{"response"}
	bz, err = querier(ctx, path, query)
	require.NoError(t, err)
	require.NotNil(t, bz)
}
