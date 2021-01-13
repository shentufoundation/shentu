package simulation_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/oracle/simulation"
	"github.com/certikfoundation/shentu/x/oracle/types"
)

func TestDecodeStore(t *testing.T) {
	cdc, _ := simapp.MakeCodecs()

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
		ClosingBlock:  rand.Int63n(10000),
		WaitingBlocks: rand.Int63n(1000),
		Status:        types.TaskStatus(rand.Intn(4)),
	}

	taskIDs := []types.TaskID{
		{
			Contract: task.Contract,
			Function: task.Function,
		},
	}

	operatorAddr, err := sdk.AccAddressFromBech32(operator.Address)
	require.NoError(t, err)
	withdrawAddr, err := sdk.AccAddressFromBech32(withdraw.Address)
	require.NoError(t, err)
	KVPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.OperatorStoreKey(operatorAddr), Value: cdc.MustMarshalBinaryLengthPrefixed(&operator)},
			{Key: types.WithdrawStoreKey(withdrawAddr, withdraw.DueBlock), Value: cdc.MustMarshalBinaryLengthPrefixed(&withdraw)},
			{Key: types.TotalCollateralKey(), Value: cdc.MustMarshalBinaryLengthPrefixed(&types.CoinsProto{Coins: totalCollateral})},
			{Key: types.TaskStoreKey(task.Contract, task.Function), Value: cdc.MustMarshalBinaryLengthPrefixed(&task)},
			{Key: types.ClosingTaskIDsStoreKey(task.ClosingBlock), Value: cdc.MustMarshalBinaryLengthPrefixed(&types.TaskIDs{TaskIds: taskIDs})},
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
	privkeySeed := make([]byte, 15)
	rand.Read(privkeySeed)

	privKey := secp256k1.GenPrivKeySecp256k1(privkeySeed)
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
	coins, err := sdk.ParseCoins(strconv.Itoa(amount) + denom)
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
