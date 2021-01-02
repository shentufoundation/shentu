package keeper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/engine"

	"github.com/certikfoundation/shentu/x/cvm/types"
)

func TestState_NewState(t *testing.T) {
	testInput := CreateTestInput(t)
	ctx := testInput.Ctx
	cvmk := testInput.CvmKeeper
	state := cvmk.NewState(ctx)

	callframe := engine.NewCallFrame(state, acmstate.Named("TxCache"))
	cache := callframe.Cache
	fmt.Println(cache)
	addr, err := crypto.AddressFromBytes(Addrs[0].Bytes())
	require.Nil(t, err)
	err = state.SetAddressMeta(addr, nil)
	require.Nil(t, err)
}

func TestState_UpdateAccount(t *testing.T) {
	testInput := CreateTestInput(t)
	ctx := testInput.Ctx
	cvmk := testInput.CvmKeeper
	ak := testInput.AccountKeeper
	state := cvmk.NewState(ctx)

	addr, err := crypto.AddressFromBytes(Addrs[0].Bytes())
	require.Nil(t, err)
	acc, err := state.GetAccount(addr)
	require.Nil(t, err)
	acc.Balance = 123
	err = state.UpdateAccount(acc)
	require.Nil(t, err)

	sdkAcc := ak.GetAccount(ctx, Addrs[0])
	err = sdkAcc.SetCoins(sdk.Coins{sdk.NewInt64Coin("uctk", 1234)})
	require.Nil(t, err)
	ak.SetAccount(ctx, sdkAcc)
	sdkAcc = ak.GetAccount(ctx, Addrs[0])
	acc, err = state.GetAccount(addr)
	sdkCoins := sdkAcc.GetCoins().AmountOf("uctk").Uint64()
	accAddressHex, err := sdk.AccAddressFromHex(addr.String())
	require.Nil(t, err)
	require.Equal(t, Addrs[0], accAddressHex)
	require.Equal(t, sdkCoins, acc.Balance)
	require.Less(t, len(acc.EVMCode), 1)
	require.Nil(t, acc.ContractMeta)

	var nilAcc *acm.Account
	fmt.Println(nilAcc)
	err = state.UpdateAccount(nilAcc)
	require.NotNil(t, err)

	acc.Address[0] = 0x00
	err = state.UpdateAccount(acc)
	require.Nil(t, err)
	accAddressHex, err = sdk.AccAddressFromHex(acc.Address.String())
	newsdkAcc := ak.GetAccount(ctx, accAddressHex)
	sdkCoins = newsdkAcc.GetCoins().AmountOf("uctk").Uint64()
	require.Equal(t, sdkCoins, acc.Balance)
}

func TestState_RemoveAccount(t *testing.T) {
	testInput := CreateTestInput(t)
	ctx := testInput.Ctx
	cvmk := testInput.CvmKeeper
	state := cvmk.NewState(ctx)

	addr, err := crypto.AddressFromBytes(Addrs[0].Bytes())
	require.Nil(t, err)
	acc, err := state.GetAccount(addr)
	require.Nil(t, err)
	acc.Balance = 123
	err = state.UpdateAccount(acc)
	require.Nil(t, err)

	err = state.RemoveAccount(acc.Address)
	require.Nil(t, err)

	require.Nil(t, state.store.Get(types.CodeStoreKey(acc.Address)))
	require.Nil(t, state.store.Get(types.AbiStoreKey(acc.Address)))
	require.Nil(t, state.store.Get(types.AddressMetaStoreKey(acc.Address)))

	nilAddr := append([]byte{0x00}, acc.Address[1:]...)
	addr, err = crypto.AddressFromBytes(nilAddr)
	require.Nil(t, err)
	err = state.RemoveAccount(addr)
	require.NotNil(t, err)
}
