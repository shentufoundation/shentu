package cert_test

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	certkeeper "github.com/shentufoundation/shentu/v2/x/cert/keeper"
	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// Test_CertificateIndexWrites verifies that SetCertificate populates all five
// secondary indexes and that IsCertified / IsContentCertified use them correctly.
func Test_CertificateIndexWrites(t *testing.T) {
	t.Run("SetCertificate populates all secondary indexes", func(t *testing.T) {
		app := shentuapp.Setup(t, false)
		ctx := app.BaseApp.NewContext(false)
		addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
		app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], ""))

		contentStr := "shentu1fdyv6hpukqj6kqdtwc42qacq9lpxm0pnggk5vn"
		certTypeStr := "auditing"

		// Before issuing: IsCertified and IsContentCertified must be false.
		require.False(t, app.CertKeeper.IsCertified(ctx, contentStr, certTypeStr))
		require.False(t, app.CertKeeper.IsContentCertified(ctx, contentStr))

		cert, err := types.NewCertificate(certTypeStr, contentStr, "", "", "", addrs[0])
		require.NoError(t, err)
		_, err = app.CertKeeper.IssueCertificate(ctx, cert)
		require.NoError(t, err)

		// After issuing: all index-based lookups must find the certificate.
		require.True(t, app.CertKeeper.IsCertified(ctx, contentStr, certTypeStr))
		require.True(t, app.CertKeeper.IsContentCertified(ctx, contentStr))

		// Certifier index: GetCertificatesByCertifier must return the certificate.
		certs := app.CertKeeper.GetCertificatesByCertifier(ctx, addrs[0])
		require.Len(t, certs, 1)
		require.Equal(t, contentStr, certs[0].GetContentString())

		// Type index: GetCertificatesFiltered with type-only filter must find it.
		params := types.NewQueryCertificatesParams(1, 100, nil, types.CertificateTypeAuditing)
		filtered, pagination, err := app.CertKeeper.GetCertificatesFiltered(ctx, params)
		require.NoError(t, err)
		require.Equal(t, uint64(1), pagination.Total)
		require.Len(t, filtered, 1)

		// Certifier+type index: combination filter must also find it.
		params2 := types.NewQueryCertificatesParams(1, 100, addrs[0], types.CertificateTypeAuditing)
		filtered2, pagination2, err := app.CertKeeper.GetCertificatesFiltered(ctx, params2)
		require.NoError(t, err)
		require.Equal(t, uint64(1), pagination2.Total)
		require.Len(t, filtered2, 1)
	})
}

// Test_CertificateIndexDeletes verifies that DeleteCertificate removes all secondary
// index entries so that subsequent index-based lookups no longer find the certificate.
func Test_CertificateIndexDeletes(t *testing.T) {
	t.Run("DeleteCertificate removes all secondary indexes", func(t *testing.T) {
		app := shentuapp.Setup(t, false)
		ctx := app.BaseApp.NewContext(false)
		addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
		app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], ""))

		contentStr := "unique-content-for-delete-test"
		cert, err := types.NewCertificate("general", contentStr, "", "", "", addrs[0])
		require.NoError(t, err)
		id, err := app.CertKeeper.IssueCertificate(ctx, cert)
		require.NoError(t, err)

		// Confirm the cert is reachable.
		require.True(t, app.CertKeeper.IsContentCertified(ctx, contentStr))

		// Delete the certificate.
		issued, err := app.CertKeeper.GetCertificateByID(ctx, id)
		require.NoError(t, err)
		err = app.CertKeeper.DeleteCertificate(ctx, issued)
		require.NoError(t, err)

		// After delete: all index-based lookups must find nothing.
		require.False(t, app.CertKeeper.IsContentCertified(ctx, contentStr))
		require.False(t, app.CertKeeper.IsCertified(ctx, contentStr, "general"))

		certs := app.CertKeeper.GetCertificatesByCertifier(ctx, addrs[0])
		require.Empty(t, certs)

		params := types.NewQueryCertificatesParams(1, 100, addrs[0], types.CertificateTypeNil)
		filtered, pagination, err := app.CertKeeper.GetCertificatesFiltered(ctx, params)
		require.NoError(t, err)
		require.Equal(t, uint64(0), pagination.Total)
		require.Empty(t, filtered)
	})
}

// Test_IsContentCertified verifies that IsContentCertified returns true only after
// a certificate with the given content has been issued.
func Test_IsContentCertified(t *testing.T) {
	t.Run("IsContentCertified uses content hash index", func(t *testing.T) {
		app := shentuapp.Setup(t, false)
		ctx := app.BaseApp.NewContext(false)
		addrs := shentuapp.AddTestAddrs(app, ctx, 1, math.NewInt(10000))
		app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], ""))

		content := "some-unique-content-string"
		require.False(t, app.CertKeeper.IsContentCertified(ctx, content))

		cert, err := types.NewCertificate("general", content, "", "", "", addrs[0])
		require.NoError(t, err)
		_, err = app.CertKeeper.IssueCertificate(ctx, cert)
		require.NoError(t, err)

		require.True(t, app.CertKeeper.IsContentCertified(ctx, content))
		// Different content must still return false.
		require.False(t, app.CertKeeper.IsContentCertified(ctx, "other-content"))
	})
}

// Test_IsBountyAdmin verifies that IsBountyAdmin returns true only after
// a BountyAdmin certificate for the given address has been issued.
func Test_IsBountyAdmin(t *testing.T) {
	t.Run("IsBountyAdmin uses type+content index", func(t *testing.T) {
		app := shentuapp.Setup(t, false)
		ctx := app.BaseApp.NewContext(false)
		addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))
		app.CertKeeper.SetCertifier(ctx, types.NewCertifier(addrs[0], addrs[0], ""))

		// addrs[1] is not yet a bounty admin.
		require.False(t, app.CertKeeper.IsBountyAdmin(ctx, addrs[1]))

		cert, err := types.NewCertificate("bountyadmin", addrs[1].String(), "", "", "", addrs[0])
		require.NoError(t, err)
		_, err = app.CertKeeper.IssueCertificate(ctx, cert)
		require.NoError(t, err)

		require.True(t, app.CertKeeper.IsBountyAdmin(ctx, addrs[1]))
		// addrs[0] was not certified as bounty admin.
		require.False(t, app.CertKeeper.IsBountyAdmin(ctx, addrs[0]))
	})
}

// Test_Migrate2to3IndexRebuild verifies that the Migrate2to3 migration rebuilds
// certificate secondary indexes from existing primary certificate store entries.
// It writes certificates directly to the primary KV store (bypassing the keeper's
// SetCertificate which writes indexes) to simulate the real pre-upgrade state.
func Test_Migrate2to3IndexRebuild(t *testing.T) {
	t.Run("Migrate2to3 rebuilds secondary indexes from scratch", func(t *testing.T) {
		app := shentuapp.Setup(t, false)
		ctx := app.BaseApp.NewContext(false)
		addrs := shentuapp.AddTestAddrs(app, ctx, 2, math.NewInt(10000))

		// Access the raw KV store for the cert module.
		storeKey := app.GetKey(types.StoreKey)
		rawStore := ctx.KVStore(storeKey)
		cdc := app.AppCodec()

		// Simulate pre-v4 certifier entry: raw addr-concat key + length-prefixed value.
		preV4Certifier := types.NewCertifier(addrs[0], addrs[0], "")
		rawStore.Set(types.CertifierStoreKey(addrs[0]), cdc.MustMarshalLengthPrefixed(&preV4Certifier))

		// Build two certificates.
		cert1, err := types.NewCertificate("auditing", "content-alpha", "", "", "", addrs[0])
		require.NoError(t, err)
		cert1.CertificateId = 1

		cert2, err := types.NewCertificate("bountyadmin", addrs[1].String(), "", "", "", addrs[0])
		require.NoError(t, err)
		cert2.CertificateId = 2

		// Write certificates ONLY to the primary store (bypass indexes).
		rawStore.Set(types.CertificateStoreKey(cert1.CertificateId), cdc.MustMarshal(&cert1))
		rawStore.Set(types.CertificateStoreKey(cert2.CertificateId), cdc.MustMarshal(&cert2))
		// Simulate pre-v4 next-cert-ID entry: 8B little-endian.
		nextIDBz := make([]byte, 8)
		binary.LittleEndian.PutUint64(nextIDBz, 3)
		rawStore.Set(types.NextCertificateIDStoreKey(), nextIDBz)

		// Before migration: index-based lookups must fail (no index entries).
		require.False(t, app.CertKeeper.IsCertified(ctx, "content-alpha", "auditing"))
		require.False(t, app.CertKeeper.IsContentCertified(ctx, "content-alpha"))
		require.False(t, app.CertKeeper.IsBountyAdmin(ctx, addrs[1]))
		require.Empty(t, app.CertKeeper.GetCertificatesByCertifier(ctx, addrs[0]))

		// Run the migration (rebuilds indexes + converts primary stores to collections).
		migrator := certkeeper.NewMigrator(app.CertKeeper)
		require.NoError(t, migrator.Migrate2to3(ctx))

		// After migration: all index-based lookups must succeed.
		require.True(t, app.CertKeeper.IsCertified(ctx, "content-alpha", "auditing"))
		require.True(t, app.CertKeeper.IsContentCertified(ctx, "content-alpha"))
		require.True(t, app.CertKeeper.IsBountyAdmin(ctx, addrs[1]))

		certs := app.CertKeeper.GetCertificatesByCertifier(ctx, addrs[0])
		require.Len(t, certs, 2)

		// Filtered query must also work.
		params := types.NewQueryCertificatesParams(1, 100, addrs[0], types.CertificateTypeAuditing)
		filtered, pagination, err := app.CertKeeper.GetCertificatesFiltered(ctx, params)
		require.NoError(t, err)
		require.Equal(t, uint64(1), pagination.Total)
		require.Len(t, filtered, 1)
	})
}

// Test_Migrate2to3DeletesObsoletePrefixes verifies that the migration deletes all
// entries stored under the validator (0x1), platform (0x2), library (0x6), and
// certifier alias (0x7) key prefixes.
func Test_Migrate2to3DeletesObsoletePrefixes(t *testing.T) {
	t.Run("Migrate2to3 deletes obsolete store prefixes", func(t *testing.T) {
		app := shentuapp.Setup(t, false)
		ctx := app.BaseApp.NewContext(false)

		// Seed dummy entries under each obsolete prefix.
		storeKey := app.GetKey(types.StoreKey)
		rawStore := ctx.KVStore(storeKey)
		obsoletePrefixes := []byte{0x1, 0x2, 0x6, 0x7}
		for _, pfx := range obsoletePrefixes {
			rawStore.Set([]byte{pfx, 0xAA, 0xBB}, []byte("dummy"))
		}
		for _, pfx := range obsoletePrefixes {
			require.True(t, rawStore.Has([]byte{pfx, 0xAA, 0xBB}), "expected entry under prefix 0x%x before migration", pfx)
		}

		// Run the migration.
		migrator := certkeeper.NewMigrator(app.CertKeeper)
		require.NoError(t, migrator.Migrate2to3(ctx))

		// Confirm all obsolete entries are gone.
		for _, pfx := range obsoletePrefixes {
			require.False(t, rawStore.Has([]byte{pfx, 0xAA, 0xBB}), "expected no entry under prefix 0x%x after migration", pfx)
		}
	})
}
