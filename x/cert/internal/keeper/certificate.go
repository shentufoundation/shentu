package keeper

import (
	"encoding/binary"
	"encoding/hex"
	"strings"

	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

//
// certID -> certificate
//

// SetCertificate stores a certificate using its ID field.
func (k Keeper) SetCertificate(ctx sdk.Context, certificate types.Certificate) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CertificateStoreKey(certificate.ID()), certificate.Bytes(k.cdc))
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
	var certificate types.Certificate
	k.cdc.MustUnmarshalBinaryLengthPrefixed(certificateData, &certificate)
	return certificate, nil
}

// GetNextCertificateID gets the next unused certificate ID.
func (k Keeper) GetNextCertificateID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	opBz := store.Get(types.GetNextCertificateIDKey())
	return binary.LittleEndian.Uint64(opBz)
}

// SetNextCertificateID stores the latest certificate ID.
func (k Keeper) SetNextCertificateID(ctx sdk.Context, id uint64) {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetNextCertificateIDKey(), bz)
}

// GetCertificateType gets type of a certificate by certificate ID.
func (k Keeper) GetCertificateType(ctx sdk.Context, id uint64) (types.CertificateType, error) {
	certificate, err := k.GetCertificateByID(ctx, id)
	if err != nil {
		return types.CertificateTypeNil, err
	}
	return certificate.Type(), nil
}

//
// certifier -> []CertID
//

// SetCertifierCertIDs stores the list of certificate IDs under the
// given certifier key.
func (k Keeper) SetCertifierCertIDs(ctx sdk.Context, certifier sdk.AccAddress, ids []uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(ids)
	store.Set(types.CertifierCertIDsKey(certifier), bz)
}

// GetCertifierCertIDs retrieves the list of certificate IDs under
// the given certifier key.
func (k Keeper) GetCertifierCertIDs(ctx sdk.Context, certifier sdk.AccAddress) []uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CertifierCertIDsKey(certifier))
	var ids []uint64
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &ids)
		return ids
	}
	return ids
}

// AddCertIDToCertifier adds an ID to the list of certificate IDs
// certified by a given certifier.
func (k Keeper) AddCertIDToCertifier(ctx sdk.Context, certifier sdk.AccAddress, id uint64) {
	ids := k.GetCertifierCertIDs(ctx, certifier)
	ids = append(ids, id)
	k.SetCertifierCertIDs(ctx, certifier, ids)
}

// DeleteCertIDFromCertifier deletes an ID from the list of
// certificate IDs issued by a given certifier.
func (k Keeper) DeleteCertIDFromCertifier(ctx sdk.Context, certifier sdk.AccAddress, id uint64) error {
	ids := k.GetCertifierCertIDs(ctx, certifier)
	if len(ids) == 0 {
		return types.ErrCertificateNotExists
	}
	if len(ids) > 1 {
		for i := range ids {
			if ids[i] == id {
				ids = append(ids[:i], ids[i+1:]...)
				k.SetCertifierCertIDs(ctx, certifier, ids)
				return nil
			}
		}
	} else {
		store := ctx.KVStore(k.storeKey)
		store.Delete(types.CertifierCertIDsKey(certifier))
		return nil
	}
	return types.ErrCertificateNotExists
}

// GetCertificatesByCertifier gets certificates certified by a given certifier.
func (k Keeper) GetCertificatesByCertifier(ctx sdk.Context, certifier sdk.AccAddress) []types.Certificate {
	ids := k.GetCertifierCertIDs(ctx, certifier)

	certificates := []types.Certificate{}
	for _, id := range ids {
		certificate, err := k.GetCertificateByID(ctx, id)
		if err != nil {
			panic(err)
		}
		certificates = append(certificates, certificate)
	}
	return certificates
}

//
// cert_type | content -> CertID
//

// SetContentCertID stores the certificate ID corresponding to the given
// content.
func (k Keeper) SetContentCertID(ctx sdk.Context, certType types.CertificateType, content types.RequestContent, id uint64) {
	// TODO: Cannot assume unique content?
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ContentCertIDKey(certType, content.RequestContentType, content.RequestContent), bz)
}

// GetContentCertID retrieves the certificate ID corresponding to the
// given content.
func (k Keeper) GetContentCertID(ctx sdk.Context, certType types.CertificateType, content types.RequestContent) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ContentCertIDKey(certType, content.RequestContentType, content.RequestContent))
	if bz == nil {
		return 0, false
	}
	return binary.LittleEndian.Uint64(bz), true
}

// DeleteContentCertID deletes the content - certificate ID pair from
// the store.
func (k Keeper) DeleteContentCertID(ctx sdk.Context, certType types.CertificateType, content types.RequestContent) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ContentCertIDKey(certType, content.RequestContentType, content.RequestContent))
}

// GetCertificateByTypeAndContent retrieves the certificate with the
// given certificate type and content.
func (k Keeper) GetCertificateByTypeAndContent(ctx sdk.Context, certType types.CertificateType, requestContent types.RequestContent) (types.Certificate, bool) {
	var certificate types.Certificate
	id, found := k.GetContentCertID(ctx, certType, requestContent)
	if !found {
		return certificate, false
	}
	certificate, err := k.GetCertificateByID(ctx, id)
	if err != nil {
		panic(err)
	}
	return certificate, true
}

// GetCertificatesByContent retrieves all certificates with given content.
func (k Keeper) GetCertificatesByContent(ctx sdk.Context, requestContent types.RequestContent) []types.Certificate {
	var certificates []types.Certificate
	for _, certType := range types.CertificateTypes {
		if certificate, found := k.GetCertificateByTypeAndContent(ctx, certType, requestContent); found {
			certificates = append(certificates, certificate)
		}
	}
	return certificates
}

// IterateCertificatesByType iterates over certificates with identical given certificate type.
func (k Keeper) IterateCertificatesByType(ctx sdk.Context, certType types.CertificateType, callback func(id uint64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := append(types.ContentCertIDStoreKeyPrefix, certType.Bytes()...)
	iterator := sdk.KVStorePrefixIterator(store, prefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		if callback(binary.LittleEndian.Uint64(iterator.Value())) {
			break
		}
	}
}

// IsCertified checks if a certificate of given type and content exists.
func (k Keeper) IsCertified(ctx sdk.Context, contentType string, content string, certType string) bool {
	requestContent, err := types.NewRequestContent(contentType, content)
	if err != nil {
		return false
	}
	certificateType := types.CertificateTypeFromString(certType)

	_, found := k.GetCertificateByTypeAndContent(ctx, certificateType, requestContent)
	return found
}

// IsContentCertified checks if a certificate of given content exists.
func (k Keeper) IsContentCertified(ctx sdk.Context, requestContent string) bool {
	for _, certType := range types.CertificateTypes {
		for _, requestContentType := range types.RequestContentTypes {
			requestContent := types.RequestContent{RequestContentType: requestContentType, RequestContent: requestContent}
			if _, found := k.GetCertificateByTypeAndContent(ctx, certType, requestContent); found {
				return true
			}
		}
	}
	return false
}

// IssueCertificate issues a certificate.
func (k Keeper) IssueCertificate(ctx sdk.Context, c types.Certificate) (uint64, error) {
	if !k.IsCertifier(ctx, c.Certifier()) {
		return 0, types.ErrUnqualifiedCertifier
	}
	if k.IsCertified(ctx, c.RequestContent().RequestContentType.String(), c.RequestContent().RequestContent, c.Type().String()) {
		return 0, types.ErrDuplicateCertificate
	}

	c.SetCertificateID(k.GetNextCertificateID(ctx))
	c.SetTxHash(hex.EncodeToString(tmhash.Sum(ctx.TxBytes())))

	k.AddCertIDToCertifier(ctx, c.Certifier(), c.ID())
	k.SetContentCertID(ctx, c.Type(), c.RequestContent(), c.ID())
	k.SetCertificate(ctx, c)

	k.SetNextCertificateID(ctx, c.ID() + 1)
	return c.ID(), nil
}

// IterateAllCertificate iterates over the all the stored certificates and performs a callback function.
func (k Keeper) IterateAllCertificate(ctx sdk.Context, callback func(certificate types.Certificate) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.CertificatesStoreKey())

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var certificate types.Certificate
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &certificate)

		if callback(certificate) {
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

// GetCertificatesFiltered gets certificates filtered.
func (k Keeper) GetCertificatesFiltered(ctx sdk.Context, params types.QueryCertificatesParams) (uint64, []types.Certificate, error) {
	filteredCertificates := []types.Certificate{}
	callback := func(certificate types.Certificate) bool {
		if len(params.Certifier) != 0 && !certificate.Certifier().Equals(params.Certifier) {
			return false
		}
		if params.ContentType != "" &&
			(strings.ToUpper(params.ContentType) != strings.ToUpper(certificate.RequestContent().RequestContentType.String()) ||
				certificate.RequestContent().RequestContent != params.Content) {
			return false
		}
		filteredCertificates = append([]types.Certificate{certificate}, filteredCertificates...)
		return false
	}

	// Choose an efficient iteration mechanism.
	if len(params.Certifier) != 0 {
		if params.ContentType != "" && params.Content != "" {
			for _, certType := range types.CertificateTypes {
				requestContent, err := types.NewRequestContent(params.ContentType, params.Content)
				if err != nil {
					return 0, nil, err
				}
				certificate, found := k.GetCertificateByTypeAndContent(ctx, certType, requestContent)
				if !found {
					return 0, nil, types.ErrCertificateNotExists
				}
				return 1, []types.Certificate{certificate}, nil
			}
		} else {
			k.IterateAllCertificate(ctx, callback)
		}
	} else if params.ContentType != "" && params.Content != "" {
		requestContent, err := types.NewRequestContent(params.ContentType, params.Content)
		if err != nil {
			return 0, nil, err
		}
		filteredCertificates = k.GetCertificatesByContent(ctx, requestContent)
	} else {
		k.IterateAllCertificate(ctx, callback)
	}

	// Post-processing
	total := uint64(len(filteredCertificates))

	start, end := client.Paginate(len(filteredCertificates), params.Page, params.Limit, 100)
	if start < 0 || end < 0 {
		filteredCertificates = []types.Certificate{}
	} else {
		filteredCertificates = filteredCertificates[start:end]
	}

	return total, filteredCertificates, nil
}

// RevokeCertificate revokes a certificate.
func (k Keeper) RevokeCertificate(ctx sdk.Context, certificate types.Certificate, revoker sdk.AccAddress) error {
	if !k.IsCertifier(ctx, revoker) {
		return types.ErrUnqualifiedRevoker
	}

	if err := k.DeleteCertIDFromCertifier(ctx, certificate.Certifier(), certificate.ID()); err != nil {
		return err
	}
	k.DeleteContentCertID(ctx, certificate.Type(), certificate.RequestContent())
	if err := k.DeleteCertificate(ctx, certificate); err != nil {
		return err
	}
	return nil
}

// GetCertifiedIdentities returns a list of addresses certified as identities.
func (k Keeper) GetCertifiedIdentities(ctx sdk.Context) []sdk.AccAddress {
	var ids []uint64
	k.IterateCertificatesByType(ctx, types.CertificateTypeIdentity, func(id uint64) (stop bool) {
		ids = append(ids, id)
		return false
	})

	identities := []sdk.AccAddress{}
	for _, id := range ids {
		certificate, err := k.GetCertificateByID(ctx, id)
		if err != nil {
			panic(err)
		}
		addr, _ := sdk.AccAddressFromBech32(certificate.RequestContent().RequestContent)
		identities = append(identities, addr)	
	}
	return identities
}
