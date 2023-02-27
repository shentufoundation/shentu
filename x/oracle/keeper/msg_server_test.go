package keeper_test

import (
	"crypto/sha256"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/oracle/keeper"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func TestMsgServer_CreateTxTask(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(80000*1e6))

	msgServer := keeper.NewMsgServerImpl(app.OracleKeeper)

	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	expiration1 := time.Now().Add(time.Hour).UTC()
	businessTransaction := []byte("ethereum transaction")

	msgCreateTxTask := types.NewMsgCreateTxTask(addrs[0].String(), "1", bounty, expiration1, businessTransaction)

	res, err := msgServer.CreateTxTask(sdk.WrapSDKContext(ctx), msgCreateTxTask)
	require.NoError(t, err)

	businessHash1 := sha256.Sum256(businessTransaction)
	require.Equal(t, res.TxHash, businessHash1)
}
