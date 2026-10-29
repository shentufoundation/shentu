package app

import (
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// TestMain configures the bech32 prefixes once for the whole test package.
// Mutating the global SDK configuration mid-suite (instead of here) poisons
// subsequently constructed apps: keepers built under one prefix are later
// validated against another, panicking with "hrp does not match bech32
// prefix".
func TestMain(m *testing.M) {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("shentu", "shentupub")
	cfg.SetBech32PrefixForValidator("shentavaloper", "shentavaloperpub")
	cfg.SetBech32PrefixForConsensusNode("shentavalcons", "shentavalconspub")

	code := m.Run()

	// do NOT restore the prefixes here: restoring would invalidate cached
	// address strings of apps constructed during the run.
	os.Exit(code)
}

// TestSeedV220Blacklist verifies that the v2.20.0 upgrade handler seeds every
// bootstrap blacklist entry into the store, where the ante decorator picks
// them up immediately.
func TestSeedV220Blacklist(t *testing.T) {
	app := Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	require.NotEmpty(t, v220BootstrapBlacklist)
	require.NoError(t, app.SeedV220Blacklist(ctx))

	for _, entry := range v220BootstrapBlacklist {
		got, blacklisted, err := app.BlacklistKeeper.IsBlacklisted(ctx, entry.Address)
		require.NoError(t, err, entry.Address)
		require.True(t, blacklisted, entry.Address)
		require.Equal(t, entry.Category, got.Category, entry.Address)
	}
}
