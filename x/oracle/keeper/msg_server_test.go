package keeper_test

import (
	"crypto/sha256"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/oracle"
	"github.com/shentufoundation/shentu/v2/x/oracle/keeper"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

// TestMsgServer_CreateTxTask This function demonstrates the ability to create and delete a transaction task, and to have collateral accounts approve or reject said task.
func TestMsgServer_CreateTxTask(t *testing.T) {
	ctx, ok, msgServer, addrs := DoInit(t)
	DepositCollateral(ctx, ok, addrs[0])
	DepositCollateral(ctx, ok, addrs[1])
	DepositCollateral(ctx, ok, addrs[2])

	taskParsms := ok.GetTaskParams(ctx)
	taskParsms.ShortcutQuorum = sdk.NewInt(1).ToDec()
	ok.SetTaskParams(ctx, taskParsms)

	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 1)}
	txBytes, txHash := GetBytesHash("ethereum transaction")

	CreateTxTask(ctx, addrs[2], txBytes, bounty, 3600, true)
	CheckTask(ctx, txHash, types.TaskStatusPending, -1)

	CreateTxTask(ctx, addrs[2], txBytes, bounty, 3600, false)

	RespondToTxTask(ctx, txHash, 90, addrs[0], true)
	ctx = PassBlocks(ctx, ok, t, 10, 0)
	RespondToTxTask(ctx, txHash, 80, addrs[1], true)
	//Duplicate submission
	RespondToTxTask(ctx, txHash, 8, addrs[1], false)

	ctx = PassBlocks(ctx, ok, t, 709, 0)
	ctx = PassBlocks(ctx, ok, t, 1, 1)                    //after one hour, the txtask should become invalid
	CheckTask(ctx, txHash, types.TaskStatusSucceeded, 85) // 85=(80+90)/2

	ctx = PassBlocks(ctx, ok, t, 22, 0)
	require.Len(t, ok.GetAllTasks(ctx), 1)

	msgDel := types.NewMsgDeleteTxTask(txHash, addrs[1])
	_, err := msgServer.DeleteTxTask(sdk.WrapSDKContext(ctx), msgDel)
	require.Error(t, err) // should fail because it's not the creator

	msgDel = types.NewMsgDeleteTxTask(txHash, addrs[2])
	_, err = msgServer.DeleteTxTask(sdk.WrapSDKContext(ctx), msgDel)
	require.NoError(t, err)
	require.Len(t, ok.GetAllTasks(ctx), 0)
}

func TestMsgServer_pending(t *testing.T) {
	ctx, ok, _, addrs := DoInit(t)
	DepositCollateral(ctx, ok, addrs[0])
	DepositCollateral(ctx, ok, addrs[1])

	taskParsms := ok.GetTaskParams(ctx)
	taskParsms.ShortcutQuorum = sdk.NewInt(1).ToDec()
	ok.SetTaskParams(ctx, taskParsms)

	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	txBytes, txHash := GetBytesHash("ethereum transaction")

	RespondToTxTask(ctx, txHash, 78, addrs[0], true)

	ctx = PassBlocks(ctx, ok, t, 2, 0)

	CreateTxTask(ctx, addrs[0], txBytes, bounty, 30, true)
	ctx = PassBlocks(ctx, ok, t, 3, 0)
	CheckTask(ctx, txHash, types.TaskStatusPending, -1)
	ctx = PassBlocks(ctx, ok, t, 3, 1)
	CheckTask(ctx, txHash, types.TaskStatusSucceeded, 78)

	//be noted the expiration time is counted from TxTaskResponse
	//the expirationDuration is one day. i.e. 86400 seconds
	//17280=86400/5
	ctx = PassBlocks(ctx, ok, t, 17280-9, 0)
	require.Len(t, ok.GetAllTasks(ctx), 1)
	//after passing one day, the txtask should be deleted
	ctx = PassBlocks(ctx, ok, t, 1, 0)
	require.Len(t, ok.GetAllTasks(ctx), 0)
}

func TestMsgServer_shortcut(t *testing.T) {
	ctx, ok, _, addrs := DoInit(t)
	DepositCollateral(ctx, ok, addrs[0])
	DepositCollateral(ctx, ok, addrs[1])
	DepositCollateral(ctx, ok, addrs[2])
	DepositCollateral(ctx, ok, addrs[3])

	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	txBytes, txHash := GetBytesHash("ethereum transaction")

	RespondToTxTask(ctx, txHash, 78, addrs[0], true)

	ctx = PassBlocks(ctx, ok, t, 2, 0)
	CheckTask(ctx, txHash, types.TaskStatusNil, -1)
	CreateTxTask(ctx, addrs[0], txBytes, bounty, 30, true)

	ctx = PassBlocks(ctx, ok, t, 1, 0)
	CheckTask(ctx, txHash, types.TaskStatusPending, -1)

	RespondToTxTask(ctx, txHash, 42, addrs[1], true)

	ctx = PassBlocks(ctx, ok, t, 0, 1)
	CheckTask(ctx, txHash, types.TaskStatusSucceeded, 60) // (78+42)/2

	ctx = PassBlocks(ctx, ok, t, 1, 0)
	RespondToTxTask(ctx, txHash, 42, addrs[2], false)
	ctx = PassBlocks(ctx, ok, t, 0, 0)
	ctx = PassBlocks(ctx, ok, t, 1, 0)
	RespondToTxTask(ctx, txHash, 78, addrs[0], false) //task closed

	txBytes, txHash = GetBytesHash("the second tx")

	RespondToTxTask(ctx, txHash, 78, addrs[0], true)
	RespondToTxTask(ctx, txHash, 78, addrs[0], false) //same operator
	RespondToTxTask(ctx, txHash, 78, addrs[1], true)
	RespondToTxTask(ctx, txHash, 78, addrs[2], true)
	ctx = PassBlocks(ctx, ok, t, 0, 0)

	CreateTxTask(ctx, addrs[3], txBytes, bounty, 30, true)

	ctx = PassBlocks(ctx, ok, t, 0, 1)
	CheckTask(ctx, txHash, types.TaskStatusSucceeded, 78)
	_ = PassBlocks(ctx, ok, t, 6, 0) //the two task should already be removed from closingTaskIDs
}

func TestMsgServer_WithdrawBounty(t *testing.T) {
	ctx, ok, _, addrs := DoInit(t)
	msgServer := ctx.Value("msgServer").(types.MsgServer)
	withdrawBounty := types.NewMsgWithdrawBounty(addrs[0])
	_, err := msgServer.WithdrawBounty(sdk.WrapSDKContext(ctx), withdrawBounty)
	require.Error(t, err)

	leftCoins := sdk.Coins{sdk.NewInt64Coin("uctk", 10000)}
	ok.AddRemainingBounty(ctx, addrs[0], leftCoins)
	_, err = msgServer.WithdrawBounty(sdk.WrapSDKContext(ctx), withdrawBounty)
	require.Error(t, err)

	CreateTxTask(ctx, addrs[0], []byte("test left"), leftCoins, 30, true)
	remainingBounty, err := ok.GetRemainingBounty(ctx, addrs[0])
	require.NoError(t, err)
	require.Equal(t, leftCoins, remainingBounty.Amount)
	_, err = msgServer.WithdrawBounty(sdk.WrapSDKContext(ctx), withdrawBounty)
	require.NoError(t, err)
	_, err = ok.GetRemainingBounty(ctx, addrs[0])
	require.Error(t, err)
}

func PassBlocks(ctx sdk.Context, ok keeper.Keeper, t require.TestingT, n int64, m int) sdk.Context {
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + n).WithBlockTime(ctx.BlockTime().Add(time.Second * 5 * time.Duration(n)))
	require.Len(t, append(ok.GetInvalidTaskIDs(ctx), ok.GetShortcutTasks(ctx)...), m)
	oracle.EndBlocker(ctx, ok)
	return ctx
}

func DoInit(t *testing.T) (sdk.Context, keeper.Keeper, types.MsgServer, []sdk.AccAddress) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 5, sdk.NewInt(80000*1e6))
	ok := app.OracleKeeper
	msgServer := keeper.NewMsgServerImpl(ok)
	ctx = ctx.WithValue("msgServer", msgServer).WithValue("t", t).WithValue("ok", ok)
	return ctx, ok, msgServer, addrs
}

func DepositCollateral(ctx sdk.Context, ok keeper.Keeper, addr sdk.AccAddress) {
	t := ctx.Value("t").(*testing.T)
	params := ok.GetLockedPoolParams(ctx)
	collateral := sdk.Coins{sdk.NewInt64Coin("uctk", params.MinimumCollateral)}
	require.NoError(t, ok.CreateOperator(ctx, addr, collateral, addr, "operator0"))
}

func RespondToTxTask(ctx sdk.Context, txHash []byte, score int64, addr sdk.AccAddress, success bool) {
	msgResp := types.NewMsgTxTaskResponse(txHash, score, addr)
	msgServer := ctx.Value("msgServer").(types.MsgServer)
	t := ctx.Value("t").(*testing.T)
	_, err := msgServer.TxTaskResponse(sdk.WrapSDKContext(ctx), msgResp)
	if success {
		require.NoError(t, err)
	} else {
		require.Error(t, err)
	}
}

func CheckTask(ctx sdk.Context, txHash []byte, expectedStatus types.TaskStatus, expectedScore int) {
	t := ctx.Value("t").(*testing.T)
	ok := ctx.Value("ok").(keeper.Keeper)
	txTaskRes, err := ok.GetTask(ctx, txHash)
	require.Nil(t, err)
	require.Equal(t, expectedStatus, txTaskRes.GetStatus())
	if expectedScore >= 0 {
		require.Equal(t, int64(expectedScore), txTaskRes.GetScore())
	}
}

func CreateTxTask(ctx sdk.Context, creator sdk.AccAddress, txBytes []byte, bounty sdk.Coins, deltaSecond int, succeed bool) {
	t := ctx.Value("t").(*testing.T)
	msgServer := ctx.Value("msgServer").(types.MsgServer)
	msgCreate := types.NewMsgCreateTxTask(creator, "1", txBytes, bounty, ctx.BlockTime().Add(time.Second*time.Duration(deltaSecond)))
	res, err := msgServer.CreateTxTask(sdk.WrapSDKContext(ctx), msgCreate)
	if succeed {
		require.NoError(t, err)
		txHash := sha256.Sum256(txBytes)
		require.Equal(t, res.AtxHash, txHash[:])
	} else {
		require.Error(t, err)
	}
}

func GetBytesHash(tx string) ([]byte, []byte) {
	txBytes := []byte(tx)
	txHash := sha256.Sum256(txBytes)
	return txBytes, txHash[:]
}
