package keeper

import (
	"context"
	"encoding/binary"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// SetCertificate stores a certificate using its ID field.
func (k Keeper) SetCertificate(ctx context.Context, certificate types.Certificate) error {
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&certificate)
	return store.Set(types.CertificateStoreKey(certificate.CertificateId), bz)
}

// DeleteCertificate deletes a certificate using its ID field.
func (k Keeper) DeleteCertificate(ctx context.Context, certificate types.Certificate) error {
	_, err := k.HasCertificateByID(ctx, certificate.CertificateId)
	if err != nil {
		return err
	}
	store := k.storeService.OpenKVStore(ctx)
	return store.Delete(types.CertificateStoreKey(certificate.CertificateId))
}

// HasCertificateByID checks if a certificate exists given an ID.
func (k Keeper) HasCertificateByID(ctx context.Context, id uint64) (bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	return store.Has(types.CertificateStoreKey(id))
}

// GetCertificateByID retrieves a certificate given an ID.
func (k Keeper) GetCertificateByID(ctx context.Context, id uint64) (types.Certificate, error) {
	store := k.storeService.OpenKVStore(ctx)
	certificateData, err := store.Get(types.CertificateStoreKey(id))
	if err != nil {
		return types.Certificate{}, err
	}
	var cert types.Certificate
	k.cdc.MustUnmarshal(certificateData, &cert)
	return cert, nil
}

// GetCertificateType gets type of a certificate by certificate ID.
func (k Keeper) GetCertificateType(ctx context.Context, id uint64) (types.CertificateType, error) {
	certificate, err := k.GetCertificateByID(ctx, id)
	if err != nil {
		return types.CertificateTypeNil, err
	}
	return types.TranslateCertificateType(certificate), nil
}

// IsCertified checks if a certificate of given type and content exists.
func (k Keeper) IsCertified(ctx context.Context, content string, certType string) bool {
	certificateType := types.CertificateTypeFromString(certType)
	certificates := k.GetCertificatesByTypeAndContent(ctx, certificateType, content)
	return len(certificates) > 0
}

// IsContentCertified checks if a certificate of given content exists.
func (k Keeper) IsContentCertified(ctx context.Context, content string) bool {
	return len(k.GetCertificatesByContent(ctx, content)) > 0
}

// IssueCertificate issues a certificate.
func (k Keeper) IssueCertificate(ctx context.Context, c types.Certificate) (uint64, error) {
	isCertifier, err := k.IsCertifier(ctx, c.GetCertifier())
	if err != nil {
		return 0, err
	}
	if !isCertifier {
		return 0, types.ErrUnqualifiedCertifier
	}

	certificateID, err := k.GetNextCertificateID(ctx)
	if err != nil {
		return 0, err
	}
	c.CertificateId = certificateID

	if err := k.SetNextCertificateID(ctx, certificateID+1); err != nil {
		return 0, err
	}
	if err := k.SetCertificate(ctx, c); err != nil {
		return 0, err
	}

	return c.CertificateId, nil
}

// IterateAllCertificate iterates over the all the stored certificates and performs a callback function.
func (k Keeper) IterateAllCertificate(ctx context.Context, callback func(certificate types.Certificate) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.CertificatesStoreKey())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var cert types.Certificate
		k.cdc.MustUnmarshal(iterator.Value(), &cert)

		if callback(cert) {
			break
		}
	}
}

// GetAllCertificates gets all certificates.
func (k Keeper) GetAllCertificates(ctx context.Context) (certificates []types.Certificate) {
	k.IterateAllCertificate(ctx, func(certificate types.Certificate) bool {
		certificates = append(certificates, certificate)
		return false
	})
	return certificates
}

// GetCertificatesByCertifier gets certificates certified by a given certifier.
func (k Keeper) GetCertificatesByCertifier(ctx context.Context, certifier sdk.AccAddress) []types.Certificate {
	certificates := []types.Certificate{}
	k.IterateAllCertificate(ctx, func(certificate types.Certificate) bool {
		if certificate.GetCertifier().Equals(certifier) {
			certificates = append(certificates, certificate)
		}
		return false
	})
	return certificates
}

// GetCertificatesByContent retrieves all certificates with given content.
func (k Keeper) GetCertificatesByContent(ctx context.Context, content string) []types.Certificate {
	certificates := []types.Certificate{}
	k.IterateAllCertificate(ctx, func(certificate types.Certificate) bool {
		if certificate.GetContentString() == content {
			certificates = append(certificates, certificate)
		}
		return false
	})
	return certificates
}

// GetCertificatesByTypeAndContent retrieves all certificates with given certificate type and content.
func (k Keeper) GetCertificatesByTypeAndContent(ctx context.Context, certType types.CertificateType, content string) []types.Certificate {
	certificates := []types.Certificate{}
	k.IterateAllCertificate(ctx, func(certificate types.Certificate) bool {
		if certificate.GetContentString() == content &&
			types.TranslateCertificateType(certificate) == certType {
			certificates = append(certificates, certificate)
		}
		return false
	})
	return certificates
}

// GetCertificatesFiltered gets certificates filtered.
func (k Keeper) GetCertificatesFiltered(ctx context.Context, params types.QueryCertificatesParams) (uint64, []types.Certificate, error) {
	filteredCertificates := []types.Certificate{}
	k.IterateAllCertificate(ctx, func(certificate types.Certificate) bool {
		if (len(params.Certifier) != 0 && !certificate.GetCertifier().Equals(params.Certifier)) ||
			(params.CertificateType != types.CertificateTypeNil && types.TranslateCertificateType(certificate) != params.CertificateType) {
			return false
		}
		filteredCertificates = append(filteredCertificates, certificate)
		return false
	})

	// Post-processing
	start, end := client.Paginate(len(filteredCertificates), params.Page, params.Limit, 100)
	if start < 0 || end < 0 {
		filteredCertificates = []types.Certificate{}
	} else {
		filteredCertificates = filteredCertificates[start:end]
	}

	return uint64(len(filteredCertificates)), filteredCertificates, nil
}

// RevokeCertificate revokes a certificate.
func (k Keeper) RevokeCertificate(ctx context.Context, certificate types.Certificate, revoker sdk.AccAddress) error {
	isCertifier, err := k.IsCertifier(ctx, revoker)
	if err != nil {
		return err
	}
	if !isCertifier {
		return types.ErrUnqualifiedCertifier
	}

	return k.DeleteCertificate(ctx, certificate)
}

// GetCertifiedIdentities returns a list of addresses certified as identities.
func (k Keeper) GetCertifiedIdentities(ctx context.Context) []sdk.AccAddress {
	identities := []sdk.AccAddress{}
	k.IterateAllCertificate(ctx, func(certificate types.Certificate) (stop bool) {
		if types.TranslateCertificateType(certificate) == types.CertificateTypeIdentity {
			addr, _ := sdk.AccAddressFromBech32(certificate.GetContentString())
			identities = append(identities, addr)
		}
		return false
	})
	return identities
}

// SetNextCertificateID sets the next certificate ID to store.
func (k Keeper) SetNextCertificateID(ctx context.Context, id uint64) error {
	store := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	return store.Set(types.NextCertificateIDStoreKey(), bz)
}

// GetNextCertificateID gets the next certificate ID from store.
func (k Keeper) GetNextCertificateID(ctx context.Context) (uint64, error) {
	store := k.storeService.OpenKVStore(ctx)
	opBz, err := store.Get(types.NextCertificateIDStoreKey())
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(opBz), nil
}

// IsBountyAdmin checks if an address is a bounty admin.
func (k Keeper) IsBountyAdmin(ctx context.Context, address sdk.AccAddress) bool {
	certificates := k.GetCertificatesByTypeAndContent(ctx, types.CertificateTypeBountyAdmin, address.String())

	for _, certificate := range certificates {
		if certificate.GetContentString() == address.String() {
			return true
		}
	}

	return false
}
