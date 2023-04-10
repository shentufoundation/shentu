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

func TestGetAllRemainingBounties(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	remainingBounties := ok.GetAllRemainingBounties(ctx)
	require.Equal(t, 0, len(remainingBounties))

	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	ok.AddRemainingBounty(ctx, addrs[0], bounty)
	ok.AddRemainingBounty(ctx, addrs[1], bounty)
	ok.AddRemainingBounty(ctx, addrs[2], bounty)
	ok.AddRemainingBounty(ctx, addrs[3], bounty)

	remainingBounties = ok.GetAllRemainingBounties(ctx)
	require.Equal(t, 4, len(remainingBounties))
	require.Equal(t, bounty, remainingBounties[0].Amount)
}

func TestCreatorRemainingBounties(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 4, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	remainingBounty := types.RemainingBounty{
		Address: addrs[0].String(),
		Amount:  bounty,
	}

	ok.SetRemainingBounty(ctx, remainingBounty)
	savedRemainingBounty, err := ok.GetRemainingBounty(ctx, addrs[0])
	require.NoError(t, err)
	require.Equal(t, savedRemainingBounty.Amount, remainingBounty.Amount)
}
