package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/engine"

	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/cvm/internal/types"
)

func TestState_NewState(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))
	cvmk := app.CvmKeeper
	state := cvmk.NewState(ctx)

	callframe := engine.NewCallFrame(state, acmstate.Named("TxCache"))
	cache := callframe.Cache
	fmt.Println(cache)
	addr, err := crypto.AddressFromBytes(addrs[0].Bytes())
	require.Nil(t, err)
	err = state.SetAddressMeta(addr, nil)
	require.Nil(t, err)
}

func TestState_UpdateAccount(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))
	ak, cvmk := app.AccountKeeper, app.CvmKeeper
	state := cvmk.NewState(ctx)

	addr, err := crypto.AddressFromBytes(addrs[0].Bytes())
	require.Nil(t, err)
	acc, err := state.GetAccount(addr)
	require.Nil(t, err)
	acc.Balance = 123
	err = state.UpdateAccount(acc)
	require.Nil(t, err)

	sdkAcc := ak.GetAccount(ctx, addrs[0])
	err = sdkAcc.SetCoins(sdk.Coins{sdk.NewInt64Coin("uctk", 1234)})
	require.Nil(t, err)
	ak.SetAccount(ctx, sdkAcc)
	sdkAcc = ak.GetAccount(ctx, addrs[0])
	acc, err = state.GetAccount(addr)
	sdkCoins := sdkAcc.GetCoins().AmountOf("uctk").Uint64()
	accAddressHex, err := sdk.AccAddressFromHex(addr.String())
	require.Nil(t, err)
	require.Equal(t, addrs[0], accAddressHex)
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
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{Time: time.Now().UTC()}).WithGasMeter(NewGasMeter(10000000000000))
	addrs := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))
	cvmk := app.CvmKeeper
	state := cvmk.NewState(ctx)

	addr, err := crypto.AddressFromBytes(addrs[0].Bytes())
	require.Nil(t, err)
	acc, err := state.GetAccount(addr)
	require.Nil(t, err)
	acc.Balance = 123
	err = state.UpdateAccount(acc)
	require.Nil(t, err)

	require.Nil(t, getAbi(ctx, app.GetKey(types.StoreKey), addr))
	require.NotNil(t, getCode(ctx, app.GetKey(types.StoreKey), addr))
	require.NotNil(t, getAddressMeta(ctx, app.GetKey(types.StoreKey), addr))

	err = state.RemoveAccount(acc.Address)
	require.Nil(t, err)

	acc, err = state.GetAccount(acc.Address)
	require.Nil(t, err)

	require.Nil(t, getAbi(ctx, app.GetKey(types.StoreKey), addr))
	require.Nil(t, getCode(ctx, app.GetKey(types.StoreKey), addr))
	require.Nil(t, getAddressMeta(ctx, app.GetKey(types.StoreKey), addr))

	nilAddr := append([]byte{0x00}, acc.Address[1:]...)
	addr, err = crypto.AddressFromBytes(nilAddr)
	require.Nil(t, err)
	err = state.RemoveAccount(addr)
	require.NotNil(t, err)
}
