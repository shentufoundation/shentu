package wasmtest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	shentuapp "github.com/certikfoundation/shentu/v2/app"
)

func CreateTestInput() (*shentuapp.ShentuApp, sdk.Context) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "testid", Time: time.Now().UTC()})
	return app, ctx
}

func FundAccountSim(t *testing.T, ctx sdk.Context, app *shentuapp.ShentuApp, acct sdk.AccAddress) {
	err := simapp.FundAccount(app.BankKeeper, ctx, acct, sdk.NewCoins(
		sdk.NewCoin("uctk", sdk.NewInt(10000000000)),
	))
	require.NoError(t, err)
}

// we need to make this deterministic (same every test run), as content might affect gas costs
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func RandomAccountAddress() sdk.AccAddress {
	_, _, addr := keyPubAddr()
	return addr
}

func RandomBech32AccountAddress() string {
	return RandomAccountAddress().String()
}
