package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/certikfoundation/shentu/simapp"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/engine"

	. "github.com/certikfoundation/shentu/x/cvm/keeper"
)

func TestState_NewState(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	cvmk := app.CVMKeeper
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
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	cvmk := app.CVMKeeper
	ak := app.AccountKeeper
	state := cvmk.NewState(ctx)

	addr, err := crypto.AddressFromBytes(Addrs[0].Bytes())
	require.Nil(t, err)
	acc, err := state.GetAccount(addr)
	require.Nil(t, err)
	acc.Balance = 123
	err = state.UpdateAccount(acc)
	require.Nil(t, err)

	sdkAcc := ak.GetAccount(ctx, Addrs[0])
	err = app.BankKeeper.SetBalances(ctx, Addrs[0], sdk.Coins{sdk.NewInt64Coin("uctk", 1234)})
	require.Nil(t, err)
	ak.SetAccount(ctx, sdkAcc)
	sdkAcc = ak.GetAccount(ctx, Addrs[0])
	acc, err = state.GetAccount(addr)
	sdkCoins := app.BankKeeper.GetAllBalances(ctx, addr.Bytes()).AmountOf("uctk").Uint64()
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
	sdkCoins = app.BankKeeper.GetAllBalances(ctx, accAddressHex.Bytes()).AmountOf("uctk").Uint64()
	require.Equal(t, sdkCoins, acc.Balance)
}

func TestState_RemoveAccount(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	cvmk := app.CVMKeeper
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

	require.Nil(t, cvmk.GetAbi(ctx, acc.Address))
	addrMetas, _ := state.GetAddressMeta(acc.Address)
	require.Nil(t, addrMetas)

	nilAddr := append([]byte{0x00}, acc.Address[1:]...)
	addr, err = crypto.AddressFromBytes(nilAddr)
	require.Nil(t, err)
	err = state.RemoveAccount(addr)
	require.NotNil(t, err)
}
