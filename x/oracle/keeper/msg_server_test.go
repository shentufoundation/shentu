package keeper_test

import (
	"crypto/sha256"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	certkeeper "github.com/shentufoundation/shentu/v2/x/cert/keeper"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	"github.com/shentufoundation/shentu/v2/x/oracle"
	"github.com/shentufoundation/shentu/v2/x/oracle/keeper"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func TestMsgServer_CreateTxTask(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper

	err := AddCertificate(ctx, app.CertKeeper, addrs)
	require.NoError(t, err)
	params := ok.GetLockedPoolParams(ctx)
	collateral := sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}
	require.NoError(t, ok.CreateOperator(ctx, addrs[0], collateral, addrs[0], "operator1"))
	require.NoError(t, ok.CreateOperator(ctx, addrs[1], collateral, addrs[1], "operator2"))

	msgServer := keeper.NewMsgServerImpl(app.OracleKeeper)

	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	validTime := ctx.BlockTime().Add(time.Hour).UTC()
	businessTransaction := []byte("ethereum transaction")

	msgCreateTxTask := types.NewMsgCreateTxTask(addrs[2], "1", businessTransaction, bounty, validTime)
	res, err := msgServer.CreateTxTask(sdk.WrapSDKContext(ctx), msgCreateTxTask)
	require.NoError(t, err)
	businessHash1 := sha256.Sum256(businessTransaction)
	require.Equal(t, res.TxHash, businessHash1[:])
	txTaskRes, err := ok.GetTask(ctx, businessHash1[:])
	require.NoError(t, err)
	require.Equal(t, types.TaskStatusPending, txTaskRes.GetStatus())

	_, err = msgServer.CreateTxTask(sdk.WrapSDKContext(ctx), msgCreateTxTask)
	require.Error(t, err)

	msgResp := types.NewMsgTxTaskResponse(businessHash1[:], 90, addrs[0])
	_, err = msgServer.TxTaskResponse(sdk.WrapSDKContext(ctx), msgResp)
	require.NoError(t, err)

	ctx = PassBlocks(ctx, ok, t, 10, 0)

	msgResp = types.NewMsgTxTaskResponse(businessHash1[:], 80, addrs[1])
	_, err = msgServer.TxTaskResponse(sdk.WrapSDKContext(ctx), msgResp)
	require.NoError(t, err)

	ctx = PassBlocks(ctx, ok, t, 709, 0)
	ctx = PassBlocks(ctx, ok, t, 1, 1) //after one hour, the txtask should become invalid
	txTaskRes, err = ok.GetTask(ctx, businessHash1[:])
	require.Nil(t, err)
	require.Equal(t, types.TaskStatusSucceeded, txTaskRes.GetStatus())
	require.Equal(t, int64(85), txTaskRes.GetScore()) // 85=(80+90)/2

	ctx = PassBlocks(ctx, ok, t, 22, 0)
	require.Len(t, ok.GetAllTasks(ctx), 1)

	msgDel := types.NewMsgDeleteTxTask(businessHash1[:], addrs[1])
	_, err = msgServer.DeleteTxTask(sdk.WrapSDKContext(ctx), msgDel)
	require.Error(t, err) // should fail because it's not the creator

	msgDel = types.NewMsgDeleteTxTask(businessHash1[:], addrs[2])
	_, err = msgServer.DeleteTxTask(sdk.WrapSDKContext(ctx), msgDel)
	require.NoError(t, err)
	require.Len(t, ok.GetAllTasks(ctx), 0)
}

func TestMsgServer_pending(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper
	msgServer := keeper.NewMsgServerImpl(ok)

	err := AddCertificate(ctx, app.CertKeeper, addrs)
	require.NoError(t, err)
	params := ok.GetLockedPoolParams(ctx)
	collateral := sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}
	require.NoError(t, ok.CreateOperator(ctx, addrs[0], collateral, addrs[0], "operator1"))

	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	txBytes := []byte("ethereum transaction")
	txHash := sha256.Sum256(txBytes)

	msgResp := types.NewMsgTxTaskResponse(txHash[:], 78, addrs[0])
	_, err = msgServer.TxTaskResponse(sdk.WrapSDKContext(ctx), msgResp)
	require.NoError(t, err)

	ctx = PassBlocks(ctx, ok, t, 2, 0)

	msgCreate := types.NewMsgCreateTxTask(addrs[0], "1", txBytes, bounty, ctx.BlockTime().Add(time.Second*30))
	res, err := msgServer.CreateTxTask(sdk.WrapSDKContext(ctx), msgCreate)
	require.NoError(t, err)
	require.Equal(t, res.TxHash, txHash[:])
	ctx = PassBlocks(ctx, ok, t, 3, 0)
	txTaskRes, err := ok.GetTask(ctx, txHash[:])
	require.Nil(t, err)
	require.Equal(t, types.TaskStatusPending, txTaskRes.GetStatus())
	ctx = PassBlocks(ctx, ok, t, 3, 1)
	txTaskRes, err = ok.GetTask(ctx, txHash[:])
	require.Nil(t, err)
	require.Equal(t, types.TaskStatusSucceeded, txTaskRes.GetStatus())
	require.Equal(t, int64(78), txTaskRes.GetScore())

	//be noted the expiration time is counted from TxTaskResponse
	//the expirationDuration is one day. i.e. 86400 seconds
	//17280=86400/5
	ctx = PassBlocks(ctx, ok, t, 17280-9, 0)
	require.Len(t, ok.GetAllTasks(ctx), 1)
	//after passing one day, the txtask should be deleted
	ctx = PassBlocks(ctx, ok, t, 1, 0)
	require.Len(t, ok.GetAllTasks(ctx), 0)
}

func PassBlocks(ctx sdk.Context, ok keeper.Keeper, t require.TestingT, n int64, m int) sdk.Context {
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + n).WithBlockTime(ctx.BlockTime().Add(time.Second * 5 * time.Duration(n)))
	require.Len(t, ok.GetInvalidTaskIDs(ctx), m)
	oracle.EndBlocker(ctx, ok)
	return ctx
}

func AddCertificate(ctx sdk.Context, ck certkeeper.Keeper, addrs []sdk.AccAddress) error {
	for _, addr := range addrs {
		ck.SetCertifier(ctx, certtypes.NewCertifier(addr, "", addr, ""))

		certificate, err := certtypes.NewCertificate("ORACLEOPERATOR", addr.String(), "", "", "", addr)
		if err != nil {
			return err
		}
		id := ck.GetNextCertificateID(ctx)
		certificate.CertificateId = id
		ck.SetNextCertificateID(ctx, id+1)
		ck.SetCertificate(ctx, certificate)
	}
	return nil
}
