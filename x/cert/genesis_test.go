package cert_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	cert "github.com/shentufoundation/shentu/v2/x/cert"
	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

func TestGenesisRoundTrip(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, math.NewInt(10000))

	// Seed certifiers and certificates.
	for _, addr := range addrs {
		require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addr, addr, "genesis test")))
	}
	for i, content := range []string{"content-a", "content-b", "content-c"} {
		c, err := types.NewCertificate("general", content, "", "", "", addrs[i%len(addrs)])
		require.NoError(t, err)
		_, err = app.CertKeeper.IssueCertificate(ctx, c)
		require.NoError(t, err)
	}

	// Export genesis.
	exported := cert.ExportGenesis(ctx, app.CertKeeper)
	require.Len(t, exported.Certifiers, 3)
	require.Len(t, exported.Certificates, 3)
	require.True(t, exported.NextCertificateId > 0)

	// Set up a fresh app and import.
	app2 := shentuapp.Setup(t, false)
	ctx2 := app2.BaseApp.NewContext(false)
	cert.InitGenesis(ctx2, app2.CertKeeper, *exported)

	// Re-export and compare.
	reExported := cert.ExportGenesis(ctx2, app2.CertKeeper)
	require.Equal(t, len(exported.Certifiers), len(reExported.Certifiers))
	require.Equal(t, len(exported.Certificates), len(reExported.Certificates))
	require.Equal(t, exported.NextCertificateId, reExported.NextCertificateId)

	// Verify state is queryable after import.
	for _, addr := range addrs {
		_, err := app2.CertKeeper.GetCertifier(ctx2, addr)
		require.NoError(t, err)
	}
	require.True(t, app2.CertKeeper.IsContentCertified(ctx2, "content-a"))
	require.True(t, app2.CertKeeper.IsContentCertified(ctx2, "content-b"))
	require.True(t, app2.CertKeeper.IsContentCertified(ctx2, "content-c"))
}

func TestGenesisRoundTrip_OpenMath(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	addrs := shentuapp.AddTestAddrs(app, ctx, 3, math.NewInt(10000))

	require.NoError(t, app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], "")))

	// Issue openmath certs for two provers.
	for _, prover := range []string{addrs[1].String(), addrs[2].String()} {
		c, err := types.NewCertificate("openmath", prover, "", "", "prover cert", addrs[0])
		require.NoError(t, err)
		_, err = app.CertKeeper.IssueCertificate(ctx, c)
		require.NoError(t, err)
	}

	require.True(t, app.CertKeeper.IsCertified(ctx, addrs[1].String(), "openmath"))

	// Export and re-import.
	exported := cert.ExportGenesis(ctx, app.CertKeeper)
	require.Len(t, exported.Certificates, 2)

	app2 := shentuapp.Setup(t, false)
	ctx2 := app2.BaseApp.NewContext(false)
	cert.InitGenesis(ctx2, app2.CertKeeper, *exported)

	// OpenMath certs survive the round trip.
	require.True(t, app2.CertKeeper.IsCertified(ctx2, addrs[1].String(), "openmath"))
	require.True(t, app2.CertKeeper.IsCertified(ctx2, addrs[2].String(), "openmath"))
	require.False(t, app2.CertKeeper.IsCertified(ctx2, addrs[0].String(), "openmath"))
}

func TestGenesisDefault(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	cert.InitDefaultGenesis(ctx, app.CertKeeper)

	exported := cert.ExportGenesis(ctx, app.CertKeeper)
	require.Empty(t, exported.Certifiers)
	require.Empty(t, exported.Certificates)
}
