package keeper

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// SetLibrary sets a new Certificate library registry.
func (k Keeper) SetLibrary(ctx sdk.Context, library sdk.AccAddress, publisher sdk.AccAddress) error {
	store := k.storeService.OpenKVStore(ctx)
	libraryData := types.Library{Address: library.String(), Publisher: publisher.String()}
	return store.Set(types.LibraryStoreKey(library), k.cdc.MustMarshalLengthPrefixed(&libraryData))
}

// deleteLibrary deletes a Certificate library registry.
func (k Keeper) deleteLibrary(ctx sdk.Context, library sdk.AccAddress) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Delete(types.LibraryStoreKey(library))
}

// IsLibrary checks if an address is a Certificate library.
func (k Keeper) IsLibrary(ctx sdk.Context, library sdk.AccAddress) (bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	return store.Has(types.LibraryStoreKey(library))
}

// getLibraryPublisher gets the library publisher.
func (k Keeper) getLibraryPublisher(ctx sdk.Context, library sdk.AccAddress) (sdk.AccAddress, error) {
	store := k.storeService.OpenKVStore(ctx)
	bPublisher, err := store.Get(types.LibraryStoreKey(library))
	if err != nil {
		return nil, err
	}
	if bPublisher != nil {
		var libraryData types.Library
		k.cdc.MustUnmarshalLengthPrefixed(bPublisher, &libraryData)

		publisherAddr, err := sdk.AccAddressFromBech32(libraryData.Publisher)
		if err != nil {
			panic(err)
		}
		return publisherAddr, nil
	}
	return nil, types.ErrLibraryNotExists
}

// PublishLibrary publishes a new Certificate library.
func (k Keeper) PublishLibrary(ctx sdk.Context, library sdk.AccAddress, publisher sdk.AccAddress) error {
	isLibrary, err := k.IsLibrary(ctx, library)
	if err != nil {
		return err
	}
	if isLibrary {
		return types.ErrLibraryAlreadyExists
	}
	return k.SetLibrary(ctx, library, publisher)
}

// InvalidateLibrary invalidate a Certificate library.
func (k Keeper) InvalidateLibrary(ctx sdk.Context, library sdk.AccAddress, invalidator sdk.AccAddress) error {
	isCertifier, err := k.IsCertifier(ctx, invalidator)
	if err != nil {
		return err
	}
	if !isCertifier {
		return types.ErrRejectedValidator
	}

	publisher, err := k.getLibraryPublisher(ctx, library)
	if err != nil {
		return err
	}
	// Can only be invalidate if the invalidator is the original publisher, or that the original publisher is no longer certifier.
	isCertifier, err = k.IsCertifier(ctx, invalidator)
	if err != nil {
		return err
	}
	if !isCertifier {
		return types.ErrRejectedValidator
	}

	if !(publisher.Equals(invalidator) || !isCertifier) {
		return types.ErrUnqualifiedCertifier
	}
	return k.deleteLibrary(ctx, library)
}

// IterateAllLibraries iterates over the all the stored libraries and performs a callback function.
func (k Keeper) IterateAllLibraries(ctx sdk.Context, callback func(library types.Library) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.LibrariesStoreKey())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var library types.Library
		k.cdc.MustUnmarshalLengthPrefixed(iterator.Value(), &library)

		if callback(library) {
			break
		}
	}
}

// GetAllLibraries gets all libraries.
func (k Keeper) GetAllLibraries(ctx sdk.Context) (libraries types.Libraries) {
	k.IterateAllLibraries(ctx, func(library types.Library) bool {
		libraries = append(libraries, library)
		return false
	})
	return
}

// GetAllLibraryAddresses gets all library addresses.
func (k Keeper) GetAllLibraryAddresses(ctx sdk.Context) (libraryAddresses []sdk.AccAddress) {
	k.IterateAllLibraries(ctx, func(library types.Library) bool {
		publisherAddr, err := sdk.AccAddressFromBech32(library.Publisher)
		if err != nil {
			panic(err)
		}
		libraryAddresses = append(libraryAddresses, publisherAddr)
		return false
	})
	return
}
