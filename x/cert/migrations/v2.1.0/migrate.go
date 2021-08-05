package migrate

import (
	"github.com/google/uuid"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/keeper"
	"github.com/certikfoundation/shentu/x/cert/types"
	nftkeeper "github.com/certikfoundation/shentu/x/nft/keeper"
	nfttypes "github.com/certikfoundation/shentu/x/nft/types"
)

type Migrator struct {
	keeper    keeper.Keeper
	nftKeeper nftkeeper.Keeper
}

func NewMigrator(keeper keeper.Keeper, nftkeeper nftkeeper.Keeper) Migrator {
	return Migrator{
		keeper:    keeper,
		nftKeeper: nftkeeper,
	}
}

func deleteListStoreKeys(store sdk.KVStore, prefix []byte) {
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

func (m Migrator) MigrateCertToNFT(ctx sdk.Context, storeKey sdk.StoreKey) error {
	store := ctx.KVStore(storeKey)

	var err error
	// Migrate certificates to certificate NFTs
	m.keeper.IterateAllCertificate(ctx, func(legacyCertificate types.Certificate) bool {
		// Set token parameters based on certificate type
		var denomID, tokenNm string
		switch types.TranslateCertificateType(legacyCertificate) {
		case types.CertificateTypeAuditing:
			denomID = "certikauditing"
			tokenNm = "Auditing"
		case types.CertificateTypeIdentity:
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
	deleteListStoreKeys(store, types.CertificatesStoreKey())
	store.Delete(types.NextCertificateIDStoreKey())

	// Delete unused stores
	deleteListStoreKeys(store, types.PlatformsStoreKey())
	deleteListStoreKeys(store, types.LibrariesStoreKey())

	return nil
}
