package keeper_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
)

func TestPrecogTaskBasic(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	scoringWaitTime := uint64(1800)
	businessHash := sha256.Sum256([]byte("hello"))
	hash := hex.EncodeToString(businessHash[:])

	expiration1 := time.Now().Add(time.Hour).UTC()
	require.NoError(t, ok.CreatePrecogTask(ctx, addrs[0].String(), "1", bounty, scoringWaitTime, expiration1, hash))

	task1, err := ok.GetPrecogTask(ctx, hash)
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), task1.Creator)
	require.Equal(t, hash, task1.BusinessTxHash)
	require.Equal(t, expiration1, task1.UsageExpirationTime)

	businessHash2 := sha256.Sum256([]byte("hello world"))
	hash2 := hex.EncodeToString(businessHash2[:])
	expiration2 := time.Now().Add(time.Hour * 2).UTC()
	require.NoError(t, ok.CreatePrecogTask(ctx, addrs[0].String(), "1", bounty, scoringWaitTime, expiration2, hash2))

	task2, err := ok.GetPrecogTask(ctx, hash2)
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), task2.Creator)
	require.Equal(t, hash2, task2.BusinessTxHash)
	require.Equal(t, expiration2, task2.UsageExpirationTime)
}
