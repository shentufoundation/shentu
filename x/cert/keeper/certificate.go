package keeper

import (
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/types"
)

// SetCertificate stores a certificate using its ID field.
func (k Keeper) SetCertificate(ctx sdk.Context, certificate types.Certificate) {
	store := ctx.KVStore(k.storeKey)
	bz := k.MustMarshalCertificate(certificate)
	store.Set(types.CertificateStoreKey(certificate.ID()), bz)
}

// MustMarshalCertificate attempts to encode a Certificate object and returns the
// raw encoded bytes. It panics on error.
func (k Keeper) MustMarshalCertificate(certificate types.Certificate) []byte {
	// marshals a Certificate interface. If the given type implements
	// the Marshaler interface, it is treated as a Proto-defined message and
	// serialized that way. Otherwise, it falls back on the internal Amino codec.
	bz, err := k.cdc.MarshalInterface(certificate)
	if err != nil {
		panic(fmt.Errorf("failed to encode certificate: %w", err))
	}
	return bz
}

// DeleteCertificate deletes a certificate using its ID field.
func (k Keeper) DeleteCertificate(ctx sdk.Context, certificate types.Certificate) error {
	if !k.HasCertificateByID(ctx, certificate.ID()) {
		return types.ErrCertificateNotExists
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.CertificateStoreKey(certificate.ID()))
	return nil
}

// HasCertificateByID checks if a certificate exists given an ID.
func (k Keeper) HasCertificateByID(ctx sdk.Context, id uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.CertificateStoreKey(id))
}

// GetCertificateByID retrieves a certificate given an ID.
func (k Keeper) GetCertificateByID(ctx sdk.Context, id uint64) (types.Certificate, error) {
	store := ctx.KVStore(k.storeKey)
	certificateData := store.Get(types.CertificateStoreKey(id))
	if certificateData == nil {
		return nil, types.ErrCertificateNotExists
	}

	var cert types.Certificate
	err := k.cdc.UnmarshalInterface(certificateData, &cert)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// GetCertificateType gets type of a certificate by certificate ID.
func (k Keeper) GetCertificateType(ctx sdk.Context, id uint64) (types.CertificateType, error) {
	certificate, err := k.GetCertificateByID(ctx, id)
	if err != nil {
		return types.CertificateTypeNil, err
	}
	return types.TranslateCertificateType(certificate), nil
}

// IsCertified checks if a certificate of given type and content exists.
func (k Keeper) IsCertified(ctx sdk.Context, requestContentType string, content string, certType string) bool {
	certificateType := types.CertificateTypeFromString(certType)
	certificates := k.GetCertificatesByTypeAndContent(ctx, certificateType,
		types.RequestContentTypeFromString(requestContentType), content)
	return len(certificates) > 0
}

// IsContentCertified checks if a certificate of given content exists.
func (k Keeper) IsContentCertified(ctx sdk.Context, requestContent string) bool {
	for _, requestContentType := range types.RequestContentTypes {
		if len(k.GetCertificatesByContent(ctx, requestContentType, requestContent)) > 0 {
			return true
		}
	}
	return false
}

// IssueCertificate issues a certificate.
func (k Keeper) IssueCertificate(ctx sdk.Context, c types.Certificate) (uint64, error) {
	if !k.IsCertifier(ctx, c.Certifier()) {
		return 0, types.ErrUnqualifiedCertifier
	}

	certificateID := k.GetNextCertificateID(ctx)
	c.SetCertificateID(certificateID)

	k.SetNextCertificateID(ctx, certificateID+1)
	k.SetCertificate(ctx, c)

	return c.ID(), nil
}

// IterateAllCertificate iterates over the all the stored certificates and performs a callback function.
func (k Keeper) IterateAllCertificate(ctx sdk.Context, callback func(certificate types.Certificate) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.CertificatesStoreKey())

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var cert types.Certificate
		err := k.cdc.UnmarshalInterface(iterator.Value(), &cert)
		if err != nil {
			panic(err)
		}

		if callback(cert) {
			break
		}
	}
}

// GetAllCertificates gets all certificates.
func (k Keeper) GetAllCertificates(ctx sdk.Context) (certificates []types.Certificate) {
	k.IterateAllCertificate(ctx, func(certificate types.Certificate) bool {
		certificates = append(certificates, certificate)
		return false
	})
	return certificates
}

// GetCertificatesByCertifier gets certificates certified by a given certifier.
func (k Keeper) GetCertificatesByCertifier(ctx sdk.Context, certifier sdk.AccAddress) []types.Certificate {
	certificates := []types.Certificate{}
	k.IterateAllCertificate(ctx, func(certificate types.Certificate) bool {
		if certificate.Certifier().Equals(certifier) {
			certificates = append(certificates, certificate)
		}
		return false
	})
	return certificates
}

// GetCertificatesByContent retrieves all certificates with given content.
func (k Keeper) GetCertificatesByContent(ctx sdk.Context, contentType types.RequestContentType, content string) []types.Certificate {
	certificates := []types.Certificate{}
	k.IterateAllCertificate(ctx, func(certificate types.Certificate) bool {
		if certificate.Content().GetType() == contentType && certificate.Content().GetContent() == content {
			certificates = append(certificates, certificate)
		}
		return false
	})
	return certificates
}

// GetCertificatesByTypeAndContent retrieves all certificates with given certificate type and content.
func (k Keeper) GetCertificatesByTypeAndContent(ctx sdk.Context, certType types.CertificateType,
	contentType types.RequestContentType, content string) []types.Certificate {
	certificates := []types.Certificate{}
	k.IterateAllCertificate(ctx, func(certificate types.Certificate) bool {
		if certificate.Content().GetType() == contentType && certificate.Content().GetContent() == content &&
			types.TranslateCertificateType(certificate) == certType {
			certificates = append(certificates, certificate)
		}
		return false
	})
	return certificates
}

// GetCertificatesFiltered gets certificates filtered.
func (k Keeper) GetCertificatesFiltered(ctx sdk.Context, params types.QueryCertificatesParams) (uint64, []types.Certificate, error) {
	filteredCertificates := []types.Certificate{}
	k.IterateAllCertificate(ctx, func(certificate types.Certificate) bool {
		if (params.Certifier != nil && !certificate.Certifier().Equals(params.Certifier)) ||
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
func (k Keeper) RevokeCertificate(ctx sdk.Context, certificate types.Certificate, revoker sdk.AccAddress) error {
	if !k.IsCertifier(ctx, revoker) {
		return types.ErrUnqualifiedRevoker
	}
	return k.DeleteCertificate(ctx, certificate)
}

// GetCertifiedIdentities returns a list of addresses certified as identities.
func (k Keeper) GetCertifiedIdentities(ctx sdk.Context) []sdk.AccAddress {
	identities := []sdk.AccAddress{}
	k.IterateAllCertificate(ctx, func(certificate types.Certificate) (stop bool) {
		if types.TranslateCertificateType(certificate) == types.CertificateTypeIdentity &&
			certificate.Content().GetType() == types.RequestContentTypeAddress {
			addr, _ := sdk.AccAddressFromBech32(certificate.Content().GetContent())
			identities = append(identities, addr)
		}
		return false
	})
	return identities
}

// SetNextCertificateID sets the next certificate ID to store.
func (k Keeper) SetNextCertificateID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	store.Set(types.NextCertificateIDStoreKey(), bz)
}

// GetNextCertificateID gets the next certificate ID from store.
func (k Keeper) GetNextCertificateID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	opBz := store.Get(types.NextCertificateIDStoreKey())
	return binary.LittleEndian.Uint64(opBz)
}
