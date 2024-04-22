package keeper_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/engine"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/cvm/types"
)

func TestState_NewState(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))
	cvmk := app.CVMKeeper
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
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	bondDenom := app.StakingKeeper.BondDenom(ctx)

	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))
	cvmk := app.CVMKeeper
	ak := app.AccountKeeper
	state := cvmk.NewState(ctx)

	addr, err := crypto.AddressFromBytes(addrs[0].Bytes())
	require.Nil(t, err)
	acc, err := state.GetAccount(addr)
	require.Nil(t, err)
	acc.Balance = 123
	err = state.UpdateAccount(acc)
	require.Nil(t, err)

	sdkAcc := ak.GetAccount(ctx, addrs[0])
	err = testutil.FundAccount(app.BankKeeper, ctx, addrs[0], sdk.Coins{sdk.NewInt64Coin("uctk", 1234)})
	require.Nil(t, err)
	ak.SetAccount(ctx, sdkAcc)

	acc, err = state.GetAccount(addr)
	sdkCoins := app.BankKeeper.GetAllBalances(ctx, addr.Bytes()).AmountOf(bondDenom).Uint64()
	accAddressHex, err := sdk.AccAddressFromHexUnsafe(addr.String())
	require.Nil(t, err)
	require.Equal(t, addrs[0], accAddressHex)
	require.Equal(t, sdkCoins, acc.Balance)
	require.Less(t, len(acc.EVMCode), 1)
	require.Len(t, acc.ContractMeta, 0)

	var nilAcc *acm.Account
	err = state.UpdateAccount(nilAcc)
	require.NotNil(t, err)

	acc.Address[0] = 0x00
	err = state.UpdateAccount(acc)
	require.Nil(t, err)
	accAddressHex, err = sdk.AccAddressFromHexUnsafe(acc.Address.String())
	sdkCoins = app.BankKeeper.GetAllBalances(ctx, accAddressHex.Bytes()).AmountOf("uctk").Uint64()
	require.Equal(t, sdkCoins, acc.Balance)
}

func TestState_RemoveAccount(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))
	cvmk := app.CVMKeeper
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

	require.Nil(t, cvmk.GetAbi(ctx, acc.Address))
	addrMetas, _ := state.GetAddressMeta(acc.Address)
	require.Len(t, addrMetas, 0)

	nilAddr := append([]byte{0x00}, acc.Address[1:]...)
	addr, err = crypto.AddressFromBytes(nilAddr)
	require.Nil(t, err)
	err = state.RemoveAccount(addr)
	require.NotNil(t, err)
}
