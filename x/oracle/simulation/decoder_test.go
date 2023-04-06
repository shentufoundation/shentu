package simulation_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/oracle/simulation"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func TestDecodeStore(t *testing.T) {
	cdc, _ := shentuapp.MakeCodecs()

	dec := simulation.NewDecodeStore(cdc)
	rand.Seed(time.Now().UnixNano())

	operator := types.Operator{
		Address:            RandomAccount().Address.String(),
		Proposer:           RandomAccount().Address.String(),
		Collateral:         RandomCoins(100000),
		AccumulatedRewards: RandomCoins(100000),
		Name:               RandomString(10),
	}

	withdraw := types.Withdraw{
		Address:  RandomAccount().Address.String(),
		Amount:   RandomCoins(100000),
		DueBlock: rand.Int63n(1000) + 1,
	}

	totalCollateral := RandomCoins(1000000)

	task := types.Task{
		Contract:      RandomString(30),
		Function:      RandomString(15),
		Bounty:        RandomCoins(100000),
		Description:   RandomString(10),
		Expiration:    time.Time{},
		Creator:       RandomAccount().Address.String(),
		Responses:     []types.Response{RandomResponse()},
		Result:        sdk.NewInt(rand.Int63n(256)),
		ExpireHeight:  rand.Int63n(10000),
		WaitingBlocks: rand.Int63n(1000),
		Status:        types.TaskStatus(rand.Intn(4)),
	}

	taskIDs := []types.TaskID{
		{
			Tid: types.NewTaskID(task.Contract, task.Function),
		},
	}

	leftBounty := types.LeftBounty{
		Address: RandomAccount().Address.String(),
		Amount:  RandomCoins(100000),
	}
	operatorAddr, err := sdk.AccAddressFromBech32(operator.Address)
	require.NoError(t, err)
	withdrawAddr, err := sdk.AccAddressFromBech32(withdraw.Address)
	require.NoError(t, err)
	leftBountyAddr, err := sdk.AccAddressFromBech32(leftBounty.Address)
	require.NoError(t, err)

	KVPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.OperatorStoreKey(operatorAddr), Value: cdc.MustMarshalLengthPrefixed(&operator)},
			{Key: types.WithdrawStoreKey(withdrawAddr, withdraw.DueBlock), Value: cdc.MustMarshalLengthPrefixed(&withdraw)},
			{Key: types.TotalCollateralKey(), Value: cdc.MustMarshalLengthPrefixed(&types.CoinsProto{Coins: totalCollateral})},
			{Key: types.TaskStoreKey(types.NewTaskID(task.Contract, task.Function)), Value: cdc.MustMarshalLengthPrefixed(&task)},
			{Key: types.ClosingTaskIDsStoreKey(task.ExpireHeight), Value: cdc.MustMarshalLengthPrefixed(&types.TaskIDs{TaskIds: taskIDs})},
			{Key: types.LeftBountyStoreKey(leftBountyAddr), Value: cdc.MustMarshalLengthPrefixed(&leftBounty)},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Operator", fmt.Sprintf("%v\n%v", operator, operator)},
		{"Withdraw", fmt.Sprintf("%v\n%v", withdraw, withdraw)},
		{"TotalCollateral", fmt.Sprintf("%s\n%s", totalCollateral, totalCollateral)},
		{"Task", fmt.Sprintf("%v\n%v", task, task)},
		{"TaskIDs", fmt.Sprintf("%v\n%v", taskIDs, taskIDs)},
		{"other", ""},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if i == len(tests)-1 { // nolint
				require.Panics(t, func() { dec(KVPairs.Pairs[i], KVPairs.Pairs[i]) }, tt.name) // nolint
			} else {
				require.Equal(t, tt.expectedLog, dec(KVPairs.Pairs[i], KVPairs.Pairs[i]), tt.name) // nolint
			}
		})
	}
}

func RandomAccount() simtypes.Account {
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	address := sdk.AccAddress(pubKey.Address())

	return simtypes.Account{
		PrivKey: privKey,
		PubKey:  pubKey,
		Address: address,
	}
}

func RandomCoins(n int) sdk.Coins {
	amount := rand.Intn(n)
	denom := sdk.DefaultBondDenom
	coins, err := sdk.ParseCoinsNormalized(strconv.Itoa(amount) + denom)
	if err != nil {
		panic(err)
	}
	return coins
}

func RandomString(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return string(bytes)
}

func RandomResponse() types.Response {
	return types.Response{
		Score:    sdk.NewInt(rand.Int63n(256)),
		Operator: RandomAccount().Address.String(),
		Weight:   sdk.NewInt(rand.Int63n(256)),
		Reward:   RandomCoins(100000),
	}
}
