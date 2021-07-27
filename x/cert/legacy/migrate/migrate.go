package migrate

import (
	"github.com/google/uuid"

	sdk "github.com/cosmos/cosmos-sdk/types"

	certkeeper "github.com/certikfoundation/shentu/x/cert/legacy/keeper"
	certtypes "github.com/certikfoundation/shentu/x/cert/legacy/types"
	nftkeeper "github.com/certikfoundation/shentu/x/nft/keeper"
	nfttypes "github.com/certikfoundation/shentu/x/nft/types"
)

type Migrator struct {
	keeper    certkeeper.Keeper
	nftKeeper nftkeeper.Keeper
}

func NewMigrator(keeper certkeeper.Keeper, nftkeeper nftkeeper.Keeper) Migrator {
	return Migrator{
		keeper:    keeper,
		nftKeeper: nftkeeper,
	}
}

func (m Migrator) MigrateCertToNFT(ctx sdk.Context, storeKey sdk.StoreKey) error {
	store := ctx.KVStore(storeKey)

	var err error
	// Migrate certificates to certificate NFTs
	m.keeper.IterateAllCertificate(ctx, func(legacyCertificate certtypes.Certificate) bool {
		// Set token parameters based on certificate type
		var denomID, tokenNm string
		switch certtypes.TranslateCertificateType(legacyCertificate) {
		case certtypes.CertificateTypeAuditing:
			denomID = "certikauditing"
			tokenNm = "Auditing"
		case certtypes.CertificateTypeIdentity:
			denomID = "certikidentity"
			tokenNm = "Identity"
		default:
			denomID = "certikgeneral"
			tokenNm = "General"
		}
		tokenID := "certik" + uuid.NewString()

		// Set appropriate fields for certificate data
		certificate := nfttypes.Certificate{
			Content:     legacyCertificate.GetContentString(),
			Description: legacyCertificate.Description,
			Certifier:   legacyCertificate.Certifier,
		}

		// Issue certificate NFT
		if err = m.nftKeeper.IssueCertificate(ctx, denomID, tokenID, tokenNm, "", certificate); err != nil {
			return true
		}

		return false
	})

	if err != nil {
		return err
	}

	// Delete certificate stores
	store.Delete(certtypes.CertificatesStoreKey())
	store.Delete(certtypes.NextCertificateIDStoreKey())

	// Delete unused stores
	store.Delete(certtypes.PlatformsStoreKey())
	store.Delete(certtypes.LibrariesStoreKey())

	return nil
}
