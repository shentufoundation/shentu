package keeper_test

import (
	"crypto/sha256"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
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

	txTask := types.TxTask{
		Creator:   addrs[0].String(),
		TxHash:    businessHash[:],
		Bounty:    bounty,
		ValidTime: expiration1,
		Status:    types.TaskStatusNil,
	}
	require.NoError(t, ok.CreateTxTask(ctx, &txTask))

	task1, err := ok.GetTask(ctx, txTask.GetID())
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), task1.GetCreator())
	require.Equal(t, businessHash[:], task1.GetID())
	_, validTime1 := task1.GetValidTime()
	require.Equal(t, expiration1, validTime1)

	businessHash2 := sha256.Sum256([]byte("hello world"))
	expiration2 := time.Now().Add(time.Hour * 2).UTC()
	txTask2 := types.TxTask{
		Creator:   addrs[0].String(),
		TxHash:    businessHash2[:],
		Bounty:    bounty,
		ValidTime: expiration2,
		Status:    types.TaskStatusNil,
	}
	require.NoError(t, ok.CreateTxTask(ctx, &txTask2))

	task2, err := ok.GetTask(ctx, txTask2.GetID())
	require.Nil(t, err)
	require.Equal(t, addrs[0].String(), task2.GetCreator())
	require.Equal(t, businessHash2[:], task2.GetID())

	_ = ok.DeleteTask(ctx, task2)
	_, err = ok.GetTask(ctx, task2.GetID())
	require.Error(t, err)

}
