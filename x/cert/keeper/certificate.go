package keeper

import (
	"context"
	"encoding/binary"
	"fmt"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	qtypes "github.com/cosmos/cosmos-sdk/types/query"

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
	has, err := k.HasCertificateByID(ctx, certificate.CertificateId)
	if err != nil {
		return err
	}
	if !has {
		return types.ErrCertificateNotExists
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
	if certificateData == nil {
		return types.Certificate{}, types.ErrCertificateNotExists
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

// GetCertificatesFiltered gets certificates filtered by certifier, type, and/or content.
// It chooses the narrowest available index prefix for the requested filters and
// applies standard Cosmos SDK pagination semantics.
func (k Keeper) GetCertificatesFiltered(ctx context.Context, params types.QueryCertificatesParams) ([]types.Certificate, *qtypes.PageResponse, error) {
	if !types.IsValidCertificateType(params.CertificateType) {
		return nil, nil, fmt.Errorf("invalid certificate type: %d", params.CertificateType)
	}

	store, loadFromIndex := k.selectCertificateQueryStore(ctx, params)
	pageReq := params.Pagination
	if pageReq == nil {
		pageReq = &qtypes.PageRequest{}
	}

	offset := pageReq.Offset
	key := pageReq.Key
	limit := pageReq.Limit
	countTotal := pageReq.CountTotal
	reverse := pageReq.Reverse

	if offset > 0 && key != nil {
		return nil, nil, fmt.Errorf("invalid request, either offset or key is expected, got both")
	}

	if limit == 0 {
		limit = qtypes.DefaultLimit
		countTotal = true
	}

	certs := make([]types.Certificate, 0)
	appendMatch := func(key, value []byte) (bool, error) {
		cert, err := k.loadCertificateForQuery(ctx, key, value, loadFromIndex)
		if err != nil {
			return false, err
		}
		if !matchesCertificateQuery(cert, params) {
			return false, nil
		}
		certs = append(certs, cert)
		return true, nil
	}

	if len(key) != 0 {
		iterator := certificateQueryIterator(store, key, reverse)
		defer iterator.Close()

		var returned uint64
		var nextKey []byte

		for ; iterator.Valid(); iterator.Next() {
			if iterator.Error() != nil {
				return nil, nil, iterator.Error()
			}
			if returned == limit {
				nextKey = iterator.Key()
				break
			}
			matched, err := appendMatch(iterator.Key(), iterator.Value())
			if err != nil {
				return nil, nil, err
			}
			if matched {
				returned++
			}
		}

		return certs, &qtypes.PageResponse{NextKey: nextKey}, nil
	}

	iterator := certificateQueryIterator(store, nil, reverse)
	defer iterator.Close()

	end := offset + limit
	var matchedCount uint64
	var nextKey []byte

	for ; iterator.Valid(); iterator.Next() {
		if iterator.Error() != nil {
			return nil, nil, iterator.Error()
		}

		cert, err := k.loadCertificateForQuery(ctx, iterator.Key(), iterator.Value(), loadFromIndex)
		if err != nil {
			return nil, nil, err
		}
		if !matchesCertificateQuery(cert, params) {
			continue
		}

		matchedCount++
		if matchedCount <= offset {
			continue
		}
		if matchedCount <= end {
			certs = append(certs, cert)
			continue
		}
		if matchedCount == end+1 {
			nextKey = iterator.Key()
			if !countTotal {
				break
			}
		}
	}

	pageRes := &qtypes.PageResponse{NextKey: nextKey}
	if countTotal {
		pageRes.Total = matchedCount
	}
	return certs, pageRes, nil
}

func (k Keeper) selectCertificateQueryStore(ctx context.Context, params types.QueryCertificatesParams) (storetypes.KVStore, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	hasCertifier := len(params.Certifier) > 0
	hasType := params.CertificateType != types.CertificateTypeNil
	hasContent := params.Content != ""

	switch {
	case hasType && hasContent:
		return prefix.NewStore(store, types.TypeContentIndexPrefix(params.CertificateType, params.Content)), true
	case hasCertifier && hasType:
		return prefix.NewStore(store, types.CertifierTypeIndexPrefix(params.Certifier, params.CertificateType)), true
	case hasContent:
		return prefix.NewStore(store, types.ContentIndexPrefix(params.Content)), true
	case hasCertifier:
		return prefix.NewStore(store, types.CertifierIndexPrefix(params.Certifier)), true
	case hasType:
		return prefix.NewStore(store, types.TypeIndexPrefix(params.CertificateType)), true
	default:
		return prefix.NewStore(store, types.CertificatesStoreKey()), false
	}
}

func (k Keeper) loadCertificateForQuery(ctx context.Context, key, value []byte, loadFromIndex bool) (types.Certificate, error) {
	if !loadFromIndex {
		var cert types.Certificate
		k.cdc.MustUnmarshal(value, &cert)
		return cert, nil
	}

	id := types.CertIDFromIndexKey(key)
	return k.GetCertificateByID(ctx, id)
}

func matchesCertificateQuery(cert types.Certificate, params types.QueryCertificatesParams) bool {
	if len(params.Certifier) > 0 && !cert.GetCertifier().Equals(params.Certifier) {
		return false
	}
	if params.CertificateType != types.CertificateTypeNil && types.TranslateCertificateType(cert) != params.CertificateType {
		return false
	}
	if params.Content != "" && cert.GetContentString() != params.Content {
		return false
	}
	return true
}

func certificateQueryIterator(store storetypes.KVStore, start []byte, reverse bool) storetypes.Iterator {
	if reverse {
		var end []byte
		if start != nil {
			iterator := store.Iterator(start, nil)
			defer iterator.Close()
			if iterator.Valid() {
				iterator.Next()
				end = iterator.Key()
			}
		}
		return store.ReverseIterator(nil, end)
	}
	return store.Iterator(start, nil)
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
