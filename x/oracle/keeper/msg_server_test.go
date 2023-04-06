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

// TestMsgServer_CreateAtxTask This function demonstrates the ability to create and delete a transaction task, and to have collateral accounts approve or reject said task.
func TestMsgServer_CreateAtxTask(t *testing.T) {
	ctx, ok, msgServer, addrs := DoInit(t)
	DepositCollateral(ctx, ok, addrs[0])
	DepositCollateral(ctx, ok, addrs[1])
	DepositCollateral(ctx, ok, addrs[2])

	taskParsms := ok.GetTaskParams(ctx)
	taskParsms.ShortcutQuorum = sdk.NewInt(1).ToDec()
	ok.SetTaskParams(ctx, taskParsms)

	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 1)}
	atxBytes, atxHash := GetBytesHash("ethereum transaction")

	CreateAtxTask(ctx, addrs[2], atxBytes, bounty, 3600, true)
	CheckTask(ctx, atxHash, types.TaskStatusPending, -1)

	CreateAtxTask(ctx, addrs[2], atxBytes, bounty, 3600, false)

	RespondToAtxTask(ctx, atxHash, 90, addrs[0], true)
	ctx = PassBlocks(ctx, ok, t, 10, 0)
	RespondToAtxTask(ctx, atxHash, 80, addrs[1], true)
	//Duplicate submission
	RespondToAtxTask(ctx, atxHash, 8, addrs[1], false)

	ctx = PassBlocks(ctx, ok, t, 709, 0)
	ctx = PassBlocks(ctx, ok, t, 1, 1)                     //after one hour, the txtask should become invalid
	CheckTask(ctx, atxHash, types.TaskStatusSucceeded, 85) // 85=(80+90)/2

	ctx = PassBlocks(ctx, ok, t, 22, 0)
	require.Len(t, ok.GetAllTasks(ctx), 1)

	msgDel := types.NewMsgDeleteAtxTask(atxHash, addrs[1])
	_, err := msgServer.DeleteAtxTask(sdk.WrapSDKContext(ctx), msgDel)
	require.Error(t, err) // should fail because it's not the creator

	msgDel = types.NewMsgDeleteAtxTask(atxHash, addrs[2])
	_, err = msgServer.DeleteAtxTask(sdk.WrapSDKContext(ctx), msgDel)
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
	atxBytes, atxHash := GetBytesHash("ethereum transaction")

	RespondToAtxTask(ctx, atxHash, 78, addrs[0], true)

	ctx = PassBlocks(ctx, ok, t, 2, 0)

	CreateAtxTask(ctx, addrs[0], atxBytes, bounty, 30, true)
	ctx = PassBlocks(ctx, ok, t, 3, 0)
	CheckTask(ctx, atxHash, types.TaskStatusPending, -1)
	ctx = PassBlocks(ctx, ok, t, 3, 1)
	CheckTask(ctx, atxHash, types.TaskStatusSucceeded, 78)

	//be noted the expiration time is counted from AtxTaskResponse
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
	atxBytes, atxHash := GetBytesHash("ethereum transaction")

	RespondToAtxTask(ctx, atxHash, 78, addrs[0], true)

	ctx = PassBlocks(ctx, ok, t, 2, 0)
	CheckTask(ctx, atxHash, types.TaskStatusNil, -1)
	CreateAtxTask(ctx, addrs[0], atxBytes, bounty, 30, true)

	ctx = PassBlocks(ctx, ok, t, 1, 0)
	CheckTask(ctx, atxHash, types.TaskStatusPending, -1)

	RespondToAtxTask(ctx, atxHash, 42, addrs[1], true)

	ctx = PassBlocks(ctx, ok, t, 0, 1)
	CheckTask(ctx, atxHash, types.TaskStatusSucceeded, 60) // (78+42)/2

	ctx = PassBlocks(ctx, ok, t, 1, 0)
	RespondToAtxTask(ctx, atxHash, 42, addrs[2], false)
	ctx = PassBlocks(ctx, ok, t, 0, 0)
	ctx = PassBlocks(ctx, ok, t, 1, 0)
	RespondToAtxTask(ctx, atxHash, 78, addrs[0], false) //task closed

	atxBytes, atxHash = GetBytesHash("the second tx")

	RespondToAtxTask(ctx, atxHash, 78, addrs[0], true)
	RespondToAtxTask(ctx, atxHash, 78, addrs[0], false) //same operator
	RespondToAtxTask(ctx, atxHash, 78, addrs[1], true)
	RespondToAtxTask(ctx, atxHash, 78, addrs[2], true)
	ctx = PassBlocks(ctx, ok, t, 0, 0)

	CreateAtxTask(ctx, addrs[3], atxBytes, bounty, 30, true)

	ctx = PassBlocks(ctx, ok, t, 0, 1)
	CheckTask(ctx, atxHash, types.TaskStatusSucceeded, 78)
	_ = PassBlocks(ctx, ok, t, 6, 0) //the two task should already be removed from closingTaskIDs
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

func RespondToAtxTask(ctx sdk.Context, atxHash []byte, score int64, addr sdk.AccAddress, success bool) {
	msgResp := types.NewMsgAtxTaskResponse(atxHash, score, addr)
	msgServer := ctx.Value("msgServer").(types.MsgServer)
	t := ctx.Value("t").(*testing.T)
	_, err := msgServer.AtxTaskResponse(sdk.WrapSDKContext(ctx), msgResp)
	if success {
		require.NoError(t, err)
	} else {
		require.Error(t, err)
	}
}

func CheckTask(ctx sdk.Context, atxHash []byte, expectedStatus types.TaskStatus, expectedScore int) {
	t := ctx.Value("t").(*testing.T)
	ok := ctx.Value("ok").(keeper.Keeper)
	txTaskRes, err := ok.GetTask(ctx, atxHash)
	require.Nil(t, err)
	require.Equal(t, expectedStatus, txTaskRes.GetStatus())
	if expectedScore >= 0 {
		require.Equal(t, int64(expectedScore), txTaskRes.GetScore())
	}
}

func CreateAtxTask(ctx sdk.Context, creator sdk.AccAddress, atxBytes []byte, bounty sdk.Coins, deltaSecond int, succeed bool) {
	t := ctx.Value("t").(*testing.T)
	msgServer := ctx.Value("msgServer").(types.MsgServer)
	msgCreate := types.NewMsgCreateAtxTask(creator, "1", atxBytes, bounty, ctx.BlockTime().Add(time.Second*time.Duration(deltaSecond)))
	res, err := msgServer.CreateAtxTask(sdk.WrapSDKContext(ctx), msgCreate)
	if succeed {
		require.NoError(t, err)
		atxHash := sha256.Sum256(atxBytes)
		require.Equal(t, res.AtxHash, atxHash[:])
	} else {
		require.Error(t, err)
	}
}

func GetBytesHash(tx string) ([]byte, []byte) {
	atxBytes := []byte(tx)
	atxHash := sha256.Sum256(atxBytes)
	return atxBytes, atxHash[:]
}
