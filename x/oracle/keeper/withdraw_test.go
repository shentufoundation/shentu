package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func TestWithdraw(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	params := ok.GetLockedPoolParams(ctx)

	amount := sdk.Coins{sdk.NewInt64Coin("uctk", 1000)}
	require.NoError(t, ok.CreateWithdraw(ctx, addrs[0], amount))

	ctx = ctx.WithBlockHeight(1)
	require.NoError(t, ok.CreateWithdraw(ctx, addrs[1], amount))

	ctx = ctx.WithBlockHeight(2)
	require.NoError(t, ok.CreateWithdraw(ctx, addrs[0], amount))
	require.NoError(t, ok.CreateWithdraw(ctx, addrs[1], amount))

	require.NoError(t, ok.DeleteWithdraw(ctx, addrs[1], params.LockedInBlocks+1))

	withdraws := ok.GetAllWithdraws(ctx)
	require.Len(t, withdraws, 3)
	require.Equal(t, params.LockedInBlocks, withdraws[0].DueBlock)
	require.Equal(t, params.LockedInBlocks+2, withdraws[1].DueBlock)
	require.Equal(t, params.LockedInBlocks+2, withdraws[2].DueBlock)

	withdraws = ok.GetAllWithdrawsForExport(ctx)
	require.Len(t, withdraws, 3)
	require.Equal(t, params.LockedInBlocks-2, withdraws[0].DueBlock)
	require.Equal(t, params.LockedInBlocks, withdraws[1].DueBlock)
	require.Equal(t, params.LockedInBlocks, withdraws[2].DueBlock)
}

// Test set withdraw
func TestSetWithdraw(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	params := ok.GetLockedPoolParams(ctx)
	amount := sdk.Coins{sdk.NewInt64Coin("uctk", 1000)}
	dueBlock := ctx.BlockHeight() + params.LockedInBlocks
	withdraw := types.NewWithdraw(addrs[0], amount, dueBlock)
	ok.SetWithdraw(ctx, withdraw)
	withdraws := ok.GetAllWithdraws(ctx)
	require.Len(t, withdraws, 1)
	require.Equal(t, params.LockedInBlocks, withdraws[0].DueBlock)
	withdraws = ok.GetAllWithdrawsForExport(ctx)
	require.Len(t, withdraws, 1)
	require.Equal(t, params.LockedInBlocks, withdraws[0].DueBlock)
}
