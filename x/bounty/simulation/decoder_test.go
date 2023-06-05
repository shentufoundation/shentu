package simulation_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/bounty/simulation"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func TestDecodeStore(t *testing.T) {
	cdc, _ := shentuapp.MakeCodecs()

	dec := simulation.NewDecodeStore(cdc)
	rand.Seed(time.Now().UnixNano())

	encKey := types.EciesPubKey{
		EncryptionKey: []byte{4, 160, 29, 82, 27, 80},
	}
	encKeyAny, _ := codectypes.NewAnyWithValue(&encKey)
	fmt.Printf("key: %v\n", encKeyAny)
	coins, _ := sdk.ParseCoinsNormalized(strconv.Itoa(9738548213) + sdk.DefaultBondDenom)
	endTime, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", "7386-05-27 10:44:51 +0000 UTC")

	program := types.Program{
		ProgramId:         100,
		CreatorAddress:    "cosmos1aaffxj9kdecpdkzm6909w2watzxh4p5afmdcf6",
		SubmissionEndTime: endTime,
		Description:       "simulation desc 4716372193818942079",
		EncryptionKey:     encKeyAny,
		Deposit:           coins,
		CommissionRate:    sdk.NewDec(6),
		Active:            true,
	}

	programPair := kv.Pair{
		Key:   types.ProgramsKey,
		Value: cdc.MustMarshal(&program),
	}

	require.Equal(t, fmt.Sprintf("%v\n%v", program, program), dec(programPair, programPair))
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

func RandomString(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return string(bytes)
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
