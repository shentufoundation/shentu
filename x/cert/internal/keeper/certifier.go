package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

// SetCertifier sets a certifier.
func (k Keeper) SetCertifier(ctx sdk.Context, certifier types.Certifier) {
	store := ctx.KVStore(k.storeKey)
	certifierAddr, err := sdk.AccAddressFromBech32(certifier.Address)
	if err != nil {
		panic(err)
	}
	store.Set(types.CertifierStoreKey(certifierAddr), k.cdc.MustMarshalBinaryLengthPrefixed(&certifier))
	if certifier.Alias != "" {
		store.Set(types.CertifierAliasStoreKey(certifier.Alias), k.cdc.MustMarshalBinaryLengthPrefixed(&certifier))
	}
}

// deleteCertifier deletes a certifier.
func (k Keeper) deleteCertifier(ctx sdk.Context, certifierAddress sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)

	certifier, err := k.GetCertifier(ctx, certifierAddress)
	if err != nil {
		return err
	}
	alias := certifier.Alias
	store.Delete(types.CertifierAliasStoreKey(alias))
	store.Delete(types.CertifierStoreKey(certifierAddress))

	return nil
}

// IsCertifier checks if an address is a certifier.
func (k Keeper) IsCertifier(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.CertifierStoreKey(address))
}

// HasCertifierAlias checks if the alias of a certifier exists.
func (k Keeper) HasCertifierAlias(ctx sdk.Context, alias string) bool {
	if alias == "" {
		return true
	}
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.CertifierAliasStoreKey(alias))
}

// GetCertifier returns the certification information for a certifier.
func (k Keeper) GetCertifier(ctx sdk.Context, certifierAddress sdk.AccAddress) (types.Certifier, error) {
	if certifierData := ctx.KVStore(k.storeKey).Get(types.CertifierStoreKey(certifierAddress)); certifierData != nil {
		var certifier types.Certifier
		k.cdc.MustUnmarshalBinaryLengthPrefixed(certifierData, &certifier)
		return certifier, nil
	}
	return types.Certifier{}, types.ErrCertifierNotExists
}

// GetCertifierByAlias returns the certification information for a certifier by its alias.
func (k Keeper) GetCertifierByAlias(ctx sdk.Context, alias string) (types.Certifier, error) {
	if alias == "" {
		return types.Certifier{}, types.ErrInvalidCertifierAlias
	}
	if certifierData := ctx.KVStore(k.storeKey).Get(types.CertifierAliasStoreKey(alias)); certifierData != nil {
		var certifier types.Certifier
		k.cdc.MustUnmarshalBinaryLengthPrefixed(certifierData, &certifier)
		return certifier, nil
	}
	return types.Certifier{}, types.ErrCertifierNotExists
}

// IterateAllCertifiers iterates over the all the stored certifiers and performs a callback function.
func (k Keeper) IterateAllCertifiers(ctx sdk.Context, callback func(certifier types.Certifier) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.CertifiersStoreKey())

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var certifier types.Certifier
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &certifier)

		if callback(certifier) {
			break
		}
	}
}

// GetAllCertifiers gets all certifiers.
func (k Keeper) GetAllCertifiers(ctx sdk.Context) types.Certifiers {
	certifiers := types.Certifiers{}
	k.IterateAllCertifiers(ctx, func(certifier types.Certifier) bool {
		certifiers = append(certifiers, certifier)
		return false
	})
	return certifiers
}
