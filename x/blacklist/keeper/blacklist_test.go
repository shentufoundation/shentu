package keeper_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	blacklistkeeper "github.com/shentufoundation/shentu/v2/x/blacklist/keeper"
	"github.com/shentufoundation/shentu/v2/x/blacklist/types"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
)

func TestBanUnbanRoundtrip(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))

	k := app.BlacklistKeeper
	authority := k.GetAuthority()

	msg := &types.MsgBanAddress{
		Authority: authority,
		Address:   addrs[0].String(),
		Category:  types.CATEGORY_COMPROMISED,
	}
	ms := blacklistkeeper.NewMsgServerImpl(k)

	// ban
	_, err := ms.BanAddress(ctx, msg)
	require.NoError(t, err)

	entry, blacklisted, err := k.IsBlacklisted(ctx, addrs[0].String())
	require.NoError(t, err)
	require.True(t, blacklisted)
	require.Equal(t, types.CATEGORY_COMPROMISED, entry.Category)

	// duplicate ban overwrites; still blacklisted
	_, err = ms.BanAddress(ctx, msg)
	require.NoError(t, err)
	_, blacklisted, err = k.IsBlacklisted(ctx, addrs[0].String())
	require.NoError(t, err)
	require.True(t, blacklisted)

	// clean address is not blacklisted
	_, blacklisted, err = k.IsBlacklisted(ctx, addrs[1].String())
	require.NoError(t, err)
	require.False(t, blacklisted)

	// unban
	_, err = ms.UnbanAddress(ctx, &types.MsgUnbanAddress{
		Authority: authority,
		Address:   addrs[0].String(),
	})
	require.NoError(t, err)
	_, blacklisted, err = k.IsBlacklisted(ctx, addrs[0].String())
	require.NoError(t, err)
	require.False(t, blacklisted)

	// unban of a non-blacklisted address fails
	_, err = ms.UnbanAddress(ctx, &types.MsgUnbanAddress{
		Authority: authority,
		Address:   addrs[0].String(),
	})
	require.ErrorIs(t, err, types.ErrEntryNotFound)
}

func TestBanAddressAuthorization(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))

	ms := blacklistkeeper.NewMsgServerImpl(app.BlacklistKeeper)

	// non-authority is rejected
	_, err := ms.BanAddress(ctx, &types.MsgBanAddress{
		Authority: addrs[0].String(),
		Address:   addrs[0].String(),
	})
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	_, err = ms.UnbanAddress(ctx, &types.MsgUnbanAddress{
		Authority: addrs[0].String(),
		Address:   addrs[0].String(),
	})
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
}

func TestBanModuleAccountRejected(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	ms := blacklistkeeper.NewMsgServerImpl(app.BlacklistKeeper)

	// materialize the distribution module account and try to ban it
	distAddr := app.AccountKeeper.GetModuleAccount(ctx, distrtypes.ModuleName).GetAddress()
	_, err := ms.BanAddress(ctx, &types.MsgBanAddress{
		Authority: app.BlacklistKeeper.GetAuthority(),
		Address:   distAddr.String(),
		Category:  types.CATEGORY_SUSPECTED,
	})
	require.ErrorIs(t, err, types.ErrModuleAccount)
}

func TestGenesisRoundtrip(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))

	k := app.BlacklistKeeper
	ms := blacklistkeeper.NewMsgServerImpl(k)
	_, err := ms.BanAddress(ctx, &types.MsgBanAddress{
		Authority: k.GetAuthority(),
		Address:   addrs[0].String(),
		Category:  types.CATEGORY_COMPROMISED,
	})
	require.NoError(t, err)
	_, err = ms.BanAddress(ctx, &types.MsgBanAddress{
		Authority: k.GetAuthority(),
		Address:   addrs[1].String(),
		Category:  types.CATEGORY_SUSPECTED,
	})
	require.NoError(t, err)

	gs := k.ExportGenesis(ctx)
	require.Len(t, gs.Entries, 2)

	require.NoError(t, types.ValidateGenesis(*gs))
	require.Error(t, types.ValidateGenesis(types.GenesisState{
		Entries: []types.Entry{{Address: "not-an-address", Category: types.CATEGORY_SUSPECTED}},
	}))
	require.Error(t, types.ValidateGenesis(types.GenesisState{
		Entries: []types.Entry{{Address: addrs[0].String()}},
	}))

	// re-import into a fresh context (same store is fine: entries overwrite)
	k.InitGenesis(ctx, *gs)
	entry, blacklisted, err := k.IsBlacklisted(ctx, addrs[0].String())
	require.NoError(t, err)
	require.True(t, blacklisted)
	require.Equal(t, types.CATEGORY_COMPROMISED, entry.Category)
}

func TestCanonicalAddressSpelling(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))

	ms := blacklistkeeper.NewMsgServerImpl(app.BlacklistKeeper)
	canonical := addrs[0].String()
	upper := strings.ToUpper(canonical)
	require.NotEqual(t, canonical, upper)

	// ban submitted with an all-uppercase bech32 spelling
	_, err := ms.BanAddress(ctx, &types.MsgBanAddress{
		Authority: app.BlacklistKeeper.GetAuthority(),
		Address:   upper,
		Category:  types.CATEGORY_COMPROMISED,
	})
	require.NoError(t, err)

	// the canonical (lowercase) spelling must be blocked by the ante check
	entry, blacklisted, err := app.BlacklistKeeper.IsBlacklisted(ctx, canonical)
	require.NoError(t, err)
	require.True(t, blacklisted, "ban stored under a non-canonical spelling was silently ineffective")
	require.Equal(t, canonical, entry.Address, "stored address must be canonicalized")

	// unban with yet another spelling still hits the same entry
	_, err = ms.UnbanAddress(ctx, &types.MsgUnbanAddress{
		Authority: app.BlacklistKeeper.GetAuthority(),
		Address:   upper,
	})
	require.NoError(t, err)
	_, blacklisted, err = app.BlacklistKeeper.IsBlacklisted(ctx, canonical)
	require.NoError(t, err)
	require.False(t, blacklisted)
}

// TestCertifierEmergencyLane verifies the emergency lane: a registered
// certifier can ban and unban directly, without a governance vote. A
// non-certifier, non-governance submitter is rejected.
func TestCertifierEmergencyLane(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, math.NewInt(10000))
	certifier, target, outsider := addrs[0], addrs[1], addrs[2]

	require.NoError(t, app.CertKeeper.SetCertifier(ctx, certtypes.NewCertifier(certifier, "")))

	ms := blacklistkeeper.NewMsgServerImpl(app.BlacklistKeeper)

	// outsider (neither gov nor certifier) is rejected
	_, err := ms.BanAddress(ctx, &types.MsgBanAddress{
		Authority: outsider.String(),
		Address:   target.String(),
		Category:  types.CATEGORY_COMPROMISED,
	})
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)

	// certifier bans directly, no governance vote needed
	_, err = ms.BanAddress(ctx, &types.MsgBanAddress{
		Authority: certifier.String(),
		Address:   target.String(),
		Category:  types.CATEGORY_COMPROMISED,
	})
	require.NoError(t, err)
	_, blacklisted, err := app.BlacklistKeeper.IsBlacklisted(ctx, target.String())
	require.NoError(t, err)
	require.True(t, blacklisted)

	// certifier unbans directly
	_, err = ms.UnbanAddress(ctx, &types.MsgUnbanAddress{
		Authority: certifier.String(),
		Address:   target.String(),
	})
	require.NoError(t, err)
	_, blacklisted, err = app.BlacklistKeeper.IsBlacklisted(ctx, target.String())
	require.NoError(t, err)
	require.False(t, blacklisted)

	// governance lane still works: gov authority bans directly too
	_, err = ms.BanAddress(ctx, &types.MsgBanAddress{
		Authority: app.BlacklistKeeper.GetAuthority(),
		Address:   target.String(),
		Category:  types.CATEGORY_SUSPECTED,
	})
	require.NoError(t, err)
	_, blacklisted, err = app.BlacklistKeeper.IsBlacklisted(ctx, target.String())
	require.NoError(t, err)
	require.True(t, blacklisted)
}
