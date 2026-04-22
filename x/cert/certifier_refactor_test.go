package cert_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	certkeeper "github.com/shentufoundation/shentu/v2/x/cert/keeper"
	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

func TestMsgUpdateCertifier(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 4, math.NewInt(10000))
	authority := app.GovKeeper.GetGovernanceAccount(ctx).GetAddress()
	msgServer := certkeeper.NewMsgServerImpl(app.CertKeeper)

	addFirst := types.NewMsgUpdateCertifier(authority, addrs[0], "first certifier", types.Add)
	_, err := msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), addFirst)
	require.NoError(t, err)

	certifier, err := app.CertKeeper.GetCertifier(ctx, addrs[0])
	require.NoError(t, err)
	require.Equal(t, types.NewCertifier(addrs[0], "first certifier"), certifier)

	_, err = msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), addFirst)
	require.True(t, errorsmod.IsOf(err, types.ErrCertifierAlreadyExists), "expected ErrCertifierAlreadyExists, got %v", err)

	addSecond := types.NewMsgUpdateCertifier(authority, addrs[2], "second certifier", types.Add)
	_, err = msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), addSecond)
	require.NoError(t, err)

	removeFirst := types.NewMsgUpdateCertifier(authority, addrs[0], "", types.Remove)
	_, err = msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), removeFirst)
	require.NoError(t, err)

	_, err = app.CertKeeper.GetCertifier(ctx, addrs[0])
	require.True(t, errorsmod.IsOf(err, types.ErrCertifierNotExists), "expected ErrCertifierNotExists, got %v", err)

	removeLast := types.NewMsgUpdateCertifier(authority, addrs[2], "", types.Remove)
	_, err = msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), removeLast)
	require.True(t, errorsmod.IsOf(err, types.ErrOnlyOneCertifier), "expected ErrOnlyOneCertifier, got %v", err)
}

func TestMsgUpdateCertifierRejectsUnauthorizedAuthority(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	msgServer := certkeeper.NewMsgServerImpl(app.CertKeeper)

	msg := types.NewMsgUpdateCertifier(addrs[0], addrs[1], "", types.Add)
	_, err := msgServer.UpdateCertifier(sdk.WrapSDKContext(ctx), msg)
	require.True(t, errorsmod.IsOf(err, sdkerrors.ErrUnauthorized), "expected unauthorized error, got %v", err)
}

func TestCertifierQueryByAddress(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
	querier := certkeeper.Querier{Keeper: app.CertKeeper}

	expected := types.NewCertifier(addrs[0], "query me")
	require.NoError(t, app.CertKeeper.SetCertifier(ctx, expected))

	resp, err := querier.Certifier(sdk.WrapSDKContext(ctx), &types.QueryCertifierRequest{Address: addrs[0].String()})
	require.NoError(t, err)
	require.Equal(t, expected, resp.Certifier)
}
