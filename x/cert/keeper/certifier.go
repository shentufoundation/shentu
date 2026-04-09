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
	return store.Set(types.CertifierStoreKey(certifierAddr), k.cdc.MustMarshalLengthPrefixed(&certifier))
}

// deleteCertifier deletes a certifier.
func (k Keeper) deleteCertifier(ctx context.Context, certifierAddress sdk.AccAddress) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Delete(types.CertifierStoreKey(certifierAddress))
}

// IsCertifier checks if an address is a certifier.
func (k Keeper) IsCertifier(ctx context.Context, address sdk.AccAddress) (bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	return store.Has(types.CertifierStoreKey(address))
}

// GetCertifier returns the certification information for a certifier.
func (k Keeper) GetCertifier(ctx context.Context, certifierAddress sdk.AccAddress) (types.Certifier, error) {
	store := k.storeService.OpenKVStore(ctx)
	certifierData, err := store.Get(types.CertifierStoreKey(certifierAddress))
	if err != nil {
		return types.Certifier{}, err
	}
	if certifierData == nil {
		return types.Certifier{}, types.ErrCertifierNotExists
	}
	var certifier types.Certifier
	k.cdc.MustUnmarshalLengthPrefixed(certifierData, &certifier)
	return certifier, nil
}

// UpdateCertifier applies an add/remove certifier operation through one keeper path.
func (k Keeper) UpdateCertifier(ctx context.Context, operation types.AddOrRemove, certifier types.Certifier) error {
	certifierAddr, err := sdk.AccAddressFromBech32(certifier.Address)
	if err != nil {
		return err
	}

	switch operation {
	case types.Add:
		isCertifier, err := k.IsCertifier(ctx, certifierAddr)
		if err != nil {
			return err
		}
		if isCertifier {
			return types.ErrCertifierAlreadyExists
		}
		return k.SetCertifier(ctx, certifier)
	case types.Remove:
		isCertifier, err := k.IsCertifier(ctx, certifierAddr)
		if err != nil {
			return err
		}
		if !isCertifier {
			return types.ErrCertifierNotExists
		}
		certifiers := k.GetAllCertifiers(ctx)
		if len(certifiers) == 1 {
			return types.ErrOnlyOneCertifier
		}
		return k.deleteCertifier(ctx, certifierAddr)
	default:
		return types.ErrAddOrRemove
	}
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
