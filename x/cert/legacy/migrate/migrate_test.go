package migrate_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/simapp"
	// certtypes "github.com/certikfoundation/shentu/x/cert/types"
	"github.com/certikfoundation/shentu/x/cert/legacy/migrate"
	"github.com/certikfoundation/shentu/x/cert/legacy/types"
)

func TestMigrate(t *testing.T) {
	t.Run("Testing Cert-NFT Migration", func(t *testing.T) {
		app := simapp.Setup(false)
		ctx := app.BaseApp.NewContext(false, tmproto.Header{})
		store := ctx.KVStore(app.GetKey(types.StoreKey))

		// Make sure all keys exist before migration
		keptStoreKeys := [][]byte{
			types.CertifiersStoreKey(),
			types.CertifierAliasesStoreKey(),
			types.ValidatorsStoreKey(),
		}
		for _, key := range keptStoreKeys {
			require.True(t, store.Has(key), "Legacy Cert module must contain all store keys")
		}

		deletedStoreKeys := [][]byte{
			types.CertificatesStoreKey(),
			types.LibrariesStoreKey(),
			types.PlatformsStoreKey(),
			types.NextCertificateIDStoreKey(),
		}
		for _, key := range deletedStoreKeys {
			require.True(t, store.Has(key), "Legacy Cert module must contain all store keys")
		}

		// Generate certifiers with aliases
		certifiers := simapp.AddTestAddrs(app, ctx, 3, sdk.NewInt(10000))
		app.CertLegacyKeeper.SetCertifier(ctx,
			types.NewCertifier(certifiers[0], "auditing", certifiers[0], ""))
		app.CertLegacyKeeper.SetCertifier(ctx,
			types.NewCertifier(certifiers[1], "identity", certifiers[1], ""))
		app.CertLegacyKeeper.SetCertifier(ctx,
			types.NewCertifier(certifiers[2], "general", certifiers[2], ""))

		// Issue five legacy auditing certificates
		for i := 0; i < 5; i++ {
			cert, err := types.NewCertificate("auditing", "AuditingContent",
				"", "", "Audited by CertiK", certifiers[0])
			require.NoError(t, err, "Error defining an auditing certificate")

			_, err = app.CertLegacyKeeper.IssueCertificate(ctx, cert)
			require.NoError(t, err, "Cannot issue an auditing certificate")
		}

		// Issue four legacy identity certificates
		for i := 0; i < 4; i++ {
			cert, err := types.NewCertificate("identity", "IdentityContent",
				"", "", "Identity Certified by CertiK", certifiers[1])
			require.NoError(t, err, "Error defining an identity certificate")

			_, err = app.CertLegacyKeeper.IssueCertificate(ctx, cert)
			require.NoError(t, err, "Cannot issue an identity certificate")
		}

		// Issue three legacy general certificates
		for i := 0; i < 3; i++ {
			cert, err := types.NewCertificate("general", "GeneralContent",
				"", "", "Certified by CertiK", certifiers[2])
			require.NoError(t, err, "Error defining a general certificate.")

			_, err = app.CertLegacyKeeper.IssueCertificate(ctx, cert)
			require.NoError(t, err, "Cannot issue a general certificate")
		}

		// Run migration
		migrator := migrate.NewMigrator(app.CertLegacyKeeper, app.NFTKeeper)
		require.NoError(t, migrator.MigrateCertToNFT(ctx, app.GetKey(types.StoreKey)))

		for _, key := range keptStoreKeys {
			require.True(t, store.Has(key), "New Cert module must contain kept store keys")
		}

		for _, key := range deletedStoreKeys {
			require.False(t, store.Has(key), "New Cert module must not contain deleted store keys")
		}
	})
}
