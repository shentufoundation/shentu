package ante_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/blacklist/ante"
	"github.com/shentufoundation/shentu/v2/x/blacklist/types"
)

// testSignedTx is a minimal SigVerifiableTx for decorator tests.
type testSignedTx struct {
	authsigning.SigVerifiableTx
	signers [][]byte
	msgs    []sdk.Msg
}

func (tx testSignedTx) GetSigners() ([][]byte, error) { return tx.signers, nil }
func (tx testSignedTx) GetMsgs() []sdk.Msg            { return tx.msgs }

// testUnsignedTx has no signers.
type testUnsignedTx struct{ sdk.Tx }

func (testUnsignedTx) GetMsgs() []sdk.Msg { return nil }

// testFeeTx only carries a fee granter.
type testFeeTx struct {
	sdk.FeeTx
	granter []byte
}

func (tx testFeeTx) FeeGranter() []byte { return tx.granter }
func (tx testFeeTx) GetMsgs() []sdk.Msg { return nil }

func setupDecorator(t *testing.T) *decoratorEnv {
	t.Helper()
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, math.NewInt(10000))
	return &decoratorEnv{
		app:    app,
		ctx:    ctx,
		victim: addrs[0],
		clean:  addrs[1],
		clean2: addrs[2],
		dec:    ante.NewBlacklistDecorator(app.BlacklistKeeper),
	}
}

type decoratorEnv struct {
	app    *shentuapp.ShentuApp
	ctx    sdk.Context
	victim sdk.AccAddress
	clean  sdk.AccAddress
	clean2 sdk.AccAddress
	dec    ante.BlacklistDecorator
}

func (e *decoratorEnv) ban(t *testing.T, addr sdk.AccAddress) {
	t.Helper()
	require.NoError(t, e.app.BlacklistKeeper.SetEntry(e.ctx, types.Entry{
		Address:  addr.String(),
		Category: types.CATEGORY_COMPROMISED,
	}))
}

func nextCalledTracker() (*bool, sdk.AnteHandler) {
	called := false
	return &called, func(ctx sdk.Context, tx sdk.Tx, _ bool) (sdk.Context, error) {
		called = true
		return ctx, nil
	}
}

func TestBlacklistedSignerRejected(t *testing.T) {
	e := setupDecorator(t)
	e.ban(t, e.victim)

	called, next := nextCalledTracker()
	tx := testSignedTx{signers: [][]byte{e.victim.Bytes()}}
	_, err := e.dec.AnteHandle(e.ctx, tx, false, next)

	require.ErrorIs(t, err, types.ErrAddressBlacklisted)
	require.False(t, *called, "next handler must not be called for blacklisted txs")
}

func TestCleanSignerPasses(t *testing.T) {
	e := setupDecorator(t)
	e.ban(t, e.victim)

	called, next := nextCalledTracker()
	tx := testSignedTx{signers: [][]byte{e.clean.Bytes()}}
	_, err := e.dec.AnteHandle(e.ctx, tx, false, next)

	require.NoError(t, err)
	require.True(t, *called)
}

func TestFeeGranterRejected(t *testing.T) {
	e := setupDecorator(t)
	e.ban(t, e.victim)

	// the tx signer is clean but the fee is paid by the blacklisted address
	tx := testFeeTx{granter: e.victim.Bytes()}
	called, next := nextCalledTracker()
	_, err := e.dec.AnteHandle(e.ctx, tx, false, next)

	require.ErrorIs(t, err, types.ErrAddressBlacklisted)
	require.False(t, *called)
}

func TestAuthzExecRejected(t *testing.T) {
	e := setupDecorator(t)
	e.ban(t, e.victim)

	// clean grantee executes a MsgSend drawn from the blacklisted account
	inner := &banktypes.MsgSend{
		FromAddress: e.victim.String(),
		ToAddress:   e.clean.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin("uctk", math.NewInt(1))),
	}
	anyInner, err := codectypes.NewAnyWithValue(inner)
	require.NoError(t, err)
	exec := &authz.MsgExec{Grantee: e.clean.String(), Msgs: []*codectypes.Any{anyInner}}

	tx := testSignedTx{signers: [][]byte{e.clean.Bytes()}, msgs: []sdk.Msg{exec}}
	called, next := nextCalledTracker()
	_, err = e.dec.AnteHandle(e.ctx, tx, false, next)

	require.ErrorIs(t, err, types.ErrAddressBlacklisted)
	require.False(t, *called)
}

func TestNestedAuthzExecRejected(t *testing.T) {
	e := setupDecorator(t)
	e.ban(t, e.victim)

	inner := &banktypes.MsgSend{
		FromAddress: e.victim.String(),
		ToAddress:   e.clean.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin("uctk", math.NewInt(1))),
	}
	anyInner, err := codectypes.NewAnyWithValue(inner)
	require.NoError(t, err)

	// clean -> clean -> victim: two levels of authz nesting
	nested := &authz.MsgExec{Grantee: e.clean2.String(), Msgs: []*codectypes.Any{anyInner}}
	anyNested, err := codectypes.NewAnyWithValue(nested)
	require.NoError(t, err)
	outer := &authz.MsgExec{Grantee: e.clean.String(), Msgs: []*codectypes.Any{anyNested}}

	tx := testSignedTx{signers: [][]byte{e.clean.Bytes()}, msgs: []sdk.Msg{outer}}
	called, next := nextCalledTracker()
	_, err = e.dec.AnteHandle(e.ctx, tx, false, next)

	require.ErrorIs(t, err, types.ErrAddressBlacklisted)
	require.False(t, *called)
}
