package keeper

import (
	"context"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// SetCertifier sets a certifier.
func (k Keeper) SetCertifier(ctx context.Context, certifier types.Certifier) error {
	store := k.storeService.OpenKVStore(ctx)
	certifierAddr, err := sdk.AccAddressFromBech32(certifier.Address)
	if err != nil {
		return err
	}
	if err := store.Set(types.CertifierStoreKey(certifierAddr), k.cdc.MustMarshalLengthPrefixed(&certifier)); err != nil {
		return err
	}
	if certifier.Alias != "" {
		if err := store.Set(types.CertifierAliasStoreKey(certifier.Alias), k.cdc.MustMarshalLengthPrefixed(&certifier)); err != nil {
			return err
		}
	}
	return nil
}

// deleteCertifier deletes a certifier.
func (k Keeper) deleteCertifier(ctx context.Context, certifierAddress sdk.AccAddress) error {
	store := k.storeService.OpenKVStore(ctx)
	certifier, err := k.GetCertifier(ctx, certifierAddress)
	if err != nil {
		return err
	}
	alias := certifier.Alias
	if err := store.Delete(types.CertifierAliasStoreKey(alias)); err != nil {
		return err
	}
	if err := store.Delete(types.CertifierStoreKey(certifierAddress)); err != nil {
		return err
	}
	return nil
}

// IsCertifier checks if an address is a certifier.
func (k Keeper) IsCertifier(ctx context.Context, address sdk.AccAddress) (bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	return store.Has(types.CertifierStoreKey(address))
}

// HasCertifierAlias checks if the alias of a certifier exists.
func (k Keeper) HasCertifierAlias(ctx context.Context, alias string) (bool, error) {
	if alias == "" {
		return true, nil
	}
	store := k.storeService.OpenKVStore(ctx)
	return store.Has(types.CertifierAliasStoreKey(alias))
}

// GetCertifier returns the certification information for a certifier.
func (k Keeper) GetCertifier(ctx context.Context, certifierAddress sdk.AccAddress) (types.Certifier, error) {
	store := k.storeService.OpenKVStore(ctx)
	certifierData, err := store.Get(types.CertifierStoreKey(certifierAddress))
	if err != nil {
		return types.Certifier{}, err
	}
	var certifier types.Certifier
	k.cdc.MustUnmarshalLengthPrefixed(certifierData, &certifier)
	return certifier, nil
}

// GetCertifierByAlias returns the certification information for a certifier by its alias.
func (k Keeper) GetCertifierByAlias(ctx context.Context, alias string) (types.Certifier, error) {
	if alias == "" {
		return types.Certifier{}, types.ErrInvalidCertifierAlias
	}
	store := k.storeService.OpenKVStore(ctx)
	certifierData, err := store.Get(types.CertifierAliasStoreKey(alias))
	if err != nil {
		return types.Certifier{}, err
	}

	var certifier types.Certifier
	k.cdc.MustUnmarshalLengthPrefixed(certifierData, &certifier)
	return certifier, nil
}

// IterateAllCertifiers iterates over the all the stored certifiers and performs a callback function.
func (k Keeper) IterateAllCertifiers(ctx context.Context, callback func(certifier types.Certifier) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.CertifiersStoreKey())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var certifier types.Certifier
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &certifier)

		if callback(certifier) {
			break
		}
	}
}

// GetAllCertifiers gets all certifiers.
func (k Keeper) GetAllCertifiers(ctx context.Context) types.Certifiers {
	certifiers := types.Certifiers{}
	k.IterateAllCertifiers(ctx, func(certifier types.Certifier) bool {
		certifiers = append(certifiers, certifier)
		return false
	})
	return certifiers
}
