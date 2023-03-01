package keeper_test

import (
	"crypto/sha256"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
)

func TestTxTaskBasic(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	businessHash := sha256.Sum256([]byte("hello"))

	expiration1 := time.Now().Add(time.Hour).UTC()
	require.NoError(t, ok.CreateTxTask(ctx, addrs[0].String(), bounty, expiration1, businessHash[:]))

	task1, err := ok.GetTxTask(ctx, businessHash[:])
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), task1.Creator)
	require.Equal(t, businessHash[:], task1.TxHash)
	require.Equal(t, expiration1, task1.Expiration)

	businessHash2 := sha256.Sum256([]byte("hello world"))
	expiration2 := time.Now().Add(time.Hour * 2).UTC()
	require.NoError(t, ok.CreateTxTask(ctx, addrs[0].String(), bounty, expiration2, businessHash2[:]))

	task2, err := ok.GetTxTask(ctx, businessHash2[:])
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), task2.Creator)
	require.Equal(t, businessHash2[:], task2.TxHash)
	require.Equal(t, expiration2, task2.Expiration)

	_ = ok.DeleteTxTask(ctx, businessHash2[:])
	_, err = ok.GetTxTask(ctx, businessHash2[:])
	require.Error(t, err)

}
