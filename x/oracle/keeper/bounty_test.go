package keeper_test

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func TestGetAllLeftBounties(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	leftBounties := ok.GetAllLeftBounties(ctx)
	require.Equal(t, 0, len(leftBounties))

	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	ok.AddLeftBounty(ctx, addrs[0], bounty)
	ok.AddLeftBounty(ctx, addrs[1], bounty)
	ok.AddLeftBounty(ctx, addrs[2], bounty)
	ok.AddLeftBounty(ctx, addrs[3], bounty)

	leftBounties = ok.GetAllLeftBounties(ctx)
	require.Equal(t, 4, len(leftBounties))
	require.Equal(t, bounty, leftBounties[0].Amount)
}

func TestCreatorLeftBounties(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	leftBounty := types.LeftBounty{
		Address: addrs[0].String(),
		Amount:  bounty,
	}

	ok.SetCreatorLeftBounty(ctx, leftBounty)
	savedLeftBounty, err := ok.GetCreatorLeftBounty(ctx, addrs[0])
	require.NoError(t, err)
	require.Equal(t, savedLeftBounty.Amount, leftBounty.Amount)
}
