package migrate_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/cert/legacy/migrate"
	"github.com/certikfoundation/shentu/x/cert/types"
)

func hasListStoreKeys(store sdk.KVStore, prefix []byte) bool {
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	return iterator.Valid()
}

func TestMigrate(t *testing.T) {
	t.Run("Testing Cert-NFT Migration", func(t *testing.T) {
		app := simapp.Setup(false)
		ctx := app.BaseApp.NewContext(false, tmproto.Header{})
		store := ctx.KVStore(app.GetKey(types.StoreKey))

		// Generate certifiers with aliases
		certifiers := simapp.AddTestAddrs(app, ctx, 3, sdk.NewInt(10000))
		app.CertKeeper.SetCertifier(ctx,
			types.NewCertifier(certifiers[0], "auditing", certifiers[0], ""))
		app.CertKeeper.SetCertifier(ctx,
			types.NewCertifier(certifiers[1], "identity", certifiers[1], ""))
		app.CertKeeper.SetCertifier(ctx,
			types.NewCertifier(certifiers[2], "general", certifiers[2], ""))

		// Issue five legacy auditing certificates
		for i := 0; i < 5; i++ {
			cert, err := types.NewCertificate("auditing", "AuditingContent",
				"", "", "Audited by CertiK", certifiers[0])
			require.NoError(t, err, "Error defining an auditing certificate")

			_, err = app.CertKeeper.IssueCertificate(ctx, cert)
			require.NoError(t, err, "Cannot issue an auditing certificate")
		}

		// Issue four legacy identity certificates
		for i := 0; i < 4; i++ {
			cert, err := types.NewCertificate("identity", "IdentityContent",
				"", "", "Identity Certified by CertiK", certifiers[1])
			require.NoError(t, err, "Error defining an identity certificate")

			_, err = app.CertKeeper.IssueCertificate(ctx, cert)
			require.NoError(t, err, "Cannot issue an identity certificate")
		}

		// Issue three legacy general certificates
		for i := 0; i < 3; i++ {
			cert, err := types.NewCertificate("general", "GeneralContent",
				"", "", "Certified by CertiK", certifiers[2])
			require.NoError(t, err, "Error defining a general certificate.")

			_, err = app.CertKeeper.IssueCertificate(ctx, cert)
			require.NoError(t, err, "Cannot issue a general certificate")
		}

		require.True(t, hasListStoreKeys(store, types.CertifiersStoreKey()),
			"Legacy Cert module should store certifiers")
		require.True(t, hasListStoreKeys(store, types.CertificatesStoreKey()),
			"Legacy Cert module should store certificates")
		require.True(t, store.Has(types.NextCertificateIDStoreKey()),
			"Legacy Cert module should store next certificate ID")

		// Run migration
		migrator := migrate.NewMigrator(app.CertKeeper, app.NFTKeeper)
		require.NoError(t, migrator.MigrateCertToNFT(ctx, app.GetKey(types.StoreKey)))

		// Check kept keys
		require.True(t, hasListStoreKeys(store, types.CertifiersStoreKey()),
			"New Cert module should still store certifiers")

		// Check deleted keys
		require.False(t, hasListStoreKeys(store, types.CertificatesStoreKey()),
			"New Cert module should not store any certificates")
		require.False(t, store.Has(types.NextCertificateIDStoreKey()),
			"New Cert module should not store next certificate ID")

		require.Len(t, app.NFTKeeper.GetNFTs(ctx, "certikauditing"), 5,
			"NFT module should contain five auditing cert NFTs")
		require.Len(t, app.NFTKeeper.GetNFTs(ctx, "certikidentity"), 4,
			"NFT module should contain four identity cert NFTs")
		require.Len(t, app.NFTKeeper.GetNFTs(ctx, "certikgeneral"), 3,
			"NFT module should contain three general cert NFTs")
	})
}
