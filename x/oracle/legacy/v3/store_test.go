package v3_test

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/common"
	v280 "github.com/shentufoundation/shentu/v2/x/oracle/legacy/v3"
	oracletypes "github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func Test_MigrateAllTaskStore(t *testing.T) {
	app := shentuapp.Setup(false)
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(common.Bech32PrefixAccAddr, common.Bech32PrefixAccPub)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	cdc := shentuapp.MakeEncodingConfig().Marshaler

	store := ctx.KVStore(app.GetKey(oracletypes.StoreKey))
	operator, _ := common.PrefixToCertik(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes()).String())
	// mock old data

	response := oracletypes.Response{
		Operator: operator,
		Weight:   sdk.NewInt(50),
		Reward:   nil,
	}

	tasks := make(map[string]oracletypes.TaskI)
	for i := 0; i < 10; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		beginBlock := r.Int63n(100)
		waitingBlocks := r.Int63n(10) + 1
		expireHeight := beginBlock + waitingBlocks
		status := r.Intn(4)
		creator, _ := common.PrefixToCertik(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes()).String())

		task := oracletypes.Task{
			Contract:      simtypes.RandStringOfLength(r, 10),
			Function:      simtypes.RandStringOfLength(r, 5),
			BeginBlock:    beginBlock,
			Bounty:        nil,
			Description:   simtypes.RandStringOfLength(r, 5),
			Expiration:    time.Time{},
			Creator:       creator,
			Responses:     nil,
			Result:        simtypes.RandomAmount(r, sdk.NewInt(100)),
			ExpireHeight:  expireHeight,
			WaitingBlocks: waitingBlocks,
			Status:        oracletypes.TaskStatus(status),
		}
		task.AddResponse(response)
		tasks[string(oracletypes.NewTaskID(task.Contract, task.Function))] = &task

		bz, err := cdc.MarshalInterface(&task)
		if err != nil {
			panic(err)
		}
		store.Set(oracletypes.TaskStoreKey(task.GetID()), bz)

		txTask := oracletypes.TxTask{
			AtxHash:    []byte(hex.EncodeToString([]byte(simtypes.RandStringOfLength(r, 10)))),
			Bounty:     nil,
			ValidTime:  time.Time{},
			Expiration: time.Time{},
			Creator:    creator,
			Responses:  nil,
			Status:     oracletypes.TaskStatus(status),
			Score:      r.Int63n(100),
		}
		txTask.AddResponse(response)
		tasks[string(txTask.GetID())] = &txTask

		bz, err = cdc.MarshalInterface(&txTask)
		if err != nil {
			panic(err)
		}
		store.Set(oracletypes.TaskStoreKey(txTask.GetID()), bz)
	}

	err := v280.MigrateAllTaskStore(ctx, app.GetKey(oracletypes.StoreKey), cdc)
	require.Nil(t, err)

	app.OracleKeeper.IteratorAllTasks(ctx, func(task oracletypes.TaskI) bool {
		creator := task.GetCreator()
		_, _, err := bech32.DecodeAndConvert(creator)
		require.NoError(t, err)

		tk, ok := tasks[string(task.GetID())]
		require.True(t, ok)
		shentuAddr, err := common.PrefixToShentu(tk.GetCreator())
		require.NoError(t, err)
		require.Equal(t, shentuAddr, task.GetCreator())

		require.Equal(t, tk.GetBounty(), task.GetBounty())
		require.Equal(t, tk.GetScore(), task.GetScore())
		return false
	})

}

func Test_MigrateOperator(t *testing.T) {
	app := shentuapp.Setup(false)
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(common.Bech32PrefixAccAddr, common.Bech32PrefixAccPub)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	cdc := shentuapp.MakeEncodingConfig().Marshaler

	store := ctx.KVStore(app.GetKey(oracletypes.StoreKey))
	// mock old data
	for i := 0; i < 10; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		operatorAddr, _ := common.PrefixToCertik(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes()).String())
		proposerAddr, _ := common.PrefixToCertik(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes()).String())

		operator := oracletypes.Operator{
			Address:            operatorAddr,
			Proposer:           proposerAddr,
			Collateral:         nil,
			AccumulatedRewards: nil,
			Name:               simtypes.RandStringOfLength(r, 10),
		}

		bz := cdc.MustMarshalLengthPrefixed(&operator)
		_, addrBz, err := bech32.DecodeAndConvert(operatorAddr)
		if err != nil {
			panic(err)
		}
		addr := sdk.AccAddress(addrBz)
		store.Set(oracletypes.OperatorStoreKey(addr), bz)
	}

	err := v280.MigrateOperatorStore(ctx, app.GetKey(oracletypes.StoreKey), cdc)
	require.Nil(t, err)

	app.OracleKeeper.IterateAllOperators(ctx, func(operator oracletypes.Operator) bool {
		hrp, _, err := bech32.DecodeAndConvert(operator.Address)
		require.NoError(t, err)
		require.Equal(t, hrp, "shentu")

		hrp, _, err = bech32.DecodeAndConvert(operator.Proposer)
		require.NoError(t, err)
		require.Equal(t, hrp, "shentu")
		return false
	})
}

func Test_MigrateWithdraw(t *testing.T) {
	app := shentuapp.Setup(false)
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(common.Bech32PrefixAccAddr, common.Bech32PrefixAccPub)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	cdc := shentuapp.MakeEncodingConfig().Marshaler

	store := ctx.KVStore(app.GetKey(oracletypes.StoreKey))
	// mock old data
	for i := 0; i < 10; i++ {
		withdrawsAddr, _ := common.PrefixToCertik(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes()).String())

		withdraw := oracletypes.Withdraw{
			Address:  withdrawsAddr,
			Amount:   nil,
			DueBlock: rand.Int63(),
		}

		bz := cdc.MustMarshalLengthPrefixed(&withdraw)
		withdrawAcc, err := sdk.AccAddressFromBech32(withdrawsAddr)
		if err != nil {
			panic(err)
		}
		store.Set(oracletypes.WithdrawStoreKey(withdrawAcc, withdraw.DueBlock), bz)
	}

	err := v280.MigrateWithdrawStore(ctx, app.GetKey(oracletypes.StoreKey), cdc)
	require.Nil(t, err)

	app.OracleKeeper.IterateAllWithdraws(ctx, func(withdraw oracletypes.Withdraw) bool {
		hrp, _, err := bech32.DecodeAndConvert(withdraw.Address)
		require.NoError(t, err)
		require.Equal(t, hrp, "shentu")
		return false
	})
}
