package keeper

import (
	"context"
	"encoding/binary"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// SetCertificate stores a certificate using its ID field and maintains all secondary indexes.
func (k Keeper) SetCertificate(ctx context.Context, certificate types.Certificate) error {
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&certificate)
	if err := store.Set(types.CertificateStoreKey(certificate.CertificateId), bz); err != nil {
		return err
	}
	return k.writeCertificateIndexes(ctx, certificate)
}

// DeleteCertificate deletes a certificate and removes all secondary index entries.
func (k Keeper) DeleteCertificate(ctx context.Context, certificate types.Certificate) error {
	_, err := k.HasCertificateByID(ctx, certificate.CertificateId)
	if err != nil {
		return err
	}
	store := k.storeService.OpenKVStore(ctx)
	if err := store.Delete(types.CertificateStoreKey(certificate.CertificateId)); err != nil {
		return err
	}
	return k.deleteCertificateIndexes(ctx, certificate)
}

// writeCertificateIndexes writes all secondary index entries for a certificate.
func (k Keeper) writeCertificateIndexes(ctx context.Context, cert types.Certificate) error {
	store := k.storeService.OpenKVStore(ctx)
	certType := types.TranslateCertificateType(cert)
	content := cert.GetContentString()
	certifier := cert.GetCertifier()
	id := cert.CertificateId

	if err := store.Set(types.CertifierIndexKey(certifier, id), []byte{}); err != nil {
		return err
	}
	if err := store.Set(types.TypeIndexKey(certType, id), []byte{}); err != nil {
		return err
	}
	if err := store.Set(types.CertifierTypeIndexKey(certifier, certType, id), []byte{}); err != nil {
		return err
	}
	if err := store.Set(types.ContentIndexKey(content, id), []byte{}); err != nil {
		return err
	}
	return store.Set(types.TypeContentIndexKey(certType, content, id), []byte{})
}

// deleteCertificateIndexes removes all secondary index entries for a certificate.
func (k Keeper) deleteCertificateIndexes(ctx context.Context, cert types.Certificate) error {
	store := k.storeService.OpenKVStore(ctx)
	certType := types.TranslateCertificateType(cert)
	content := cert.GetContentString()
	certifier := cert.GetCertifier()
	id := cert.CertificateId

	if err := store.Delete(types.CertifierIndexKey(certifier, id)); err != nil {
		return err
	}
	if err := store.Delete(types.TypeIndexKey(certType, id)); err != nil {
		return err
	}
	if err := store.Delete(types.CertifierTypeIndexKey(certifier, certType, id)); err != nil {
		return err
	}
	if err := store.Delete(types.ContentIndexKey(content, id)); err != nil {
		return err
	}
	return store.Delete(types.TypeContentIndexKey(certType, content, id))
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
// Uses the type+content index for efficient lookup.
func (k Keeper) IsCertified(ctx context.Context, content string, certType string) bool {
	ct := types.CertificateTypeFromString(certType)
	prefix := types.TypeContentIndexPrefix(ct, content)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	return iterator.Valid()
}

// IsContentCertified checks if a certificate of given content exists.
// Uses the content hash index for efficient lookup.
func (k Keeper) IsContentCertified(ctx context.Context, content string) bool {
	prefix := types.ContentIndexPrefix(content)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	return iterator.Valid()
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
// Uses the certifier secondary index for efficient prefix-based lookup.
func (k Keeper) GetCertificatesByCertifier(ctx context.Context, certifier sdk.AccAddress) []types.Certificate {
	return k.loadCertsFromIndexPrefix(ctx, types.CertifierIndexPrefix(certifier))
}

// GetCertificatesByContent retrieves all certificates with given content.
// Uses the content hash secondary index for efficient lookup.
func (k Keeper) GetCertificatesByContent(ctx context.Context, content string) []types.Certificate {
	return k.loadCertsFromIndexPrefix(ctx, types.ContentIndexPrefix(content))
}

// GetCertificatesByTypeAndContent retrieves all certificates with given certificate type and content.
// Uses the type+content secondary index for efficient lookup.
func (k Keeper) GetCertificatesByTypeAndContent(ctx context.Context, certType types.CertificateType, content string) []types.Certificate {
	return k.loadCertsFromIndexPrefix(ctx, types.TypeContentIndexPrefix(certType, content))
}

// loadCertsFromIndexPrefix loads all certificates whose IDs appear under the given index prefix.
func (k Keeper) loadCertsFromIndexPrefix(ctx context.Context, prefix []byte) []types.Certificate {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	var certs []types.Certificate
	for ; iterator.Valid(); iterator.Next() {
		id := types.CertIDFromIndexKey(iterator.Key())
		cert, err := k.GetCertificateByID(ctx, id)
		if err != nil {
			continue
		}
		certs = append(certs, cert)
	}
	return certs
}

// GetCertificatesFiltered gets certificates filtered by certifier and/or type.
// Chooses the narrowest index for the given filter combination and paginates directly
// on the index prefix iterator, avoiding full-store scans.
// Returns the true total count of matching certificates and the requested page.
func (k Keeper) GetCertificatesFiltered(ctx context.Context, params types.QueryCertificatesParams) (uint64, []types.Certificate, error) {
	limit := params.Limit
	if limit <= 0 {
		limit = 100
	}
	skip := 0
	if params.Page > 1 {
		skip = (params.Page - 1) * limit
	}

	hasCertifier := len(params.Certifier) > 0
	hasType := params.CertificateType != types.CertificateTypeNil

	// Choose the narrowest matching index prefix.
	if hasCertifier && hasType {
		return k.paginateIndex(ctx, types.CertifierTypeIndexPrefix(params.Certifier, params.CertificateType), skip, limit)
	}
	if hasCertifier {
		return k.paginateIndex(ctx, types.CertifierIndexPrefix(params.Certifier), skip, limit)
	}
	if hasType {
		return k.paginateIndex(ctx, types.TypeIndexPrefix(params.CertificateType), skip, limit)
	}

	// No filter: iterate the primary certificate store directly.
	return k.paginatePrimary(ctx, skip, limit)
}

// paginateIndex iterates an index prefix, counts total matching entries, and loads
// the page of certificate objects from the primary store.
func (k Keeper) paginateIndex(ctx context.Context, prefix []byte, skip, limit int) (uint64, []types.Certificate, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	var total uint64
	var certs []types.Certificate
	for ; iterator.Valid(); iterator.Next() {
		total++
		if int(total) > skip && len(certs) < limit {
			id := types.CertIDFromIndexKey(iterator.Key())
			cert, err := k.GetCertificateByID(ctx, id)
			if err != nil {
				continue
			}
			certs = append(certs, cert)
		}
	}
	return total, certs, nil
}

// paginatePrimary iterates the primary certificate store for the no-filter case.
func (k Keeper) paginatePrimary(ctx context.Context, skip, limit int) (uint64, []types.Certificate, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.CertificatesStoreKey())
	defer iterator.Close()

	var total uint64
	var certs []types.Certificate
	for ; iterator.Valid(); iterator.Next() {
		total++
		if int(total) > skip && len(certs) < limit {
			var cert types.Certificate
			k.cdc.MustUnmarshal(iterator.Value(), &cert)
			certs = append(certs, cert)
		}
	}
	return total, certs, nil
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
// Uses the type+content secondary index for efficient lookup.
func (k Keeper) IsBountyAdmin(ctx context.Context, address sdk.AccAddress) bool {
	prefix := types.TypeContentIndexPrefix(types.CertificateTypeBountyAdmin, address.String())
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	return iterator.Valid()
}
