package simulation

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	return cdc
}

func TestDecodeStore(t *testing.T) {
	cdc := makeTestCodec()

	rand.Seed(time.Now().UnixNano())

	operator := types.Operator{
		Address:            RandomAccount().Address,
		Proposer:           RandomAccount().Address,
		Collateral:         RandomCoins(100000),
		AccumulatedRewards: RandomCoins(100000),
		Name:               RandomString(10),
	}

	withdraw := types.Withdraw{
		Address:  RandomAccount().Address,
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
		Creator:       RandomAccount().Address,
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

	KVPairs := kv.Pairs{
		kv.Pair{Key: types.OperatorStoreKey(operator.Address), Value: cdc.MustMarshalBinaryLengthPrefixed(&operator)},
		kv.Pair{Key: types.WithdrawStoreKey(withdraw.Address, withdraw.DueBlock), Value: cdc.MustMarshalBinaryLengthPrefixed(&withdraw)},
		kv.Pair{Key: types.TotalCollateralKey(), Value: cdc.MustMarshalBinaryLengthPrefixed(&totalCollateral)},
		kv.Pair{Key: types.TaskStoreKey(task.Contract, task.Function), Value: cdc.MustMarshalBinaryLengthPrefixed(&task)},
		kv.Pair{Key: types.ClosingTaskIDsStoreKey(task.ClosingBlock), Value: cdc.MustMarshalBinaryLengthPrefixed(&taskIDs)},
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
				require.Panics(t, func() { DecodeStore(cdc, KVPairs[i], KVPairs[i]) }, tt.name) // nolint
			} else {
				require.Equal(t, tt.expectedLog, DecodeStore(cdc, KVPairs[i], KVPairs[i]), tt.name) // nolint
			}
		})
	}
}

func RandomAccount() simulation.Account {
	privkeySeed := make([]byte, 15)
	rand.Read(privkeySeed)

	privKey := secp256k1.GenPrivKeySecp256k1(privkeySeed)
	pubKey := privKey.PubKey()
	address := sdk.AccAddress(pubKey.Address())

	return simulation.Account{
		PrivKey: privKey,
		PubKey:  pubKey,
		Address: address,
	}
}

func RandomCoins(n int) sdk.Coins {
	amount := rand.Intn(n)
	denom := common.MicroCTKDenom
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
		Contract: RandomString(30),
		Function: RandomString(15),
		Score:    sdk.NewInt(rand.Int63n(256)),
		Operator: RandomAccount().Address,
	}
}
