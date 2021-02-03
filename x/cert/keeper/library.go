package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/types"
)

// SetLibrary sets a new Certificate library registry.
func (k Keeper) SetLibrary(ctx sdk.Context, library sdk.AccAddress, publisher sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	libraryData := types.Library{Address: library, Publisher: publisher}
	store.Set(types.LibraryStoreKey(library), k.cdc.MustMarshalBinaryLengthPrefixed(libraryData))
}

// deleteLibrary deletes a Certificate library registry.
func (k Keeper) deleteLibrary(ctx sdk.Context, library sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.LibraryStoreKey(library))
}

// IsLibrary checks if an address is a Certificate library.
func (k Keeper) IsLibrary(ctx sdk.Context, library sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.LibraryStoreKey(library))
}

// getLibraryPublisher gets the library publisher.
func (k Keeper) getLibraryPublisher(ctx sdk.Context, library sdk.AccAddress) (sdk.AccAddress, error) {
	store := ctx.KVStore(k.storeKey)
	if bPublisher := store.Get(types.LibraryStoreKey(library)); bPublisher != nil {
		var libraryData types.Library
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bPublisher, &libraryData)
		return libraryData.Publisher, nil
	}
	return nil, types.ErrLibraryNotExists
}

// PublishLibrary publishes a new Certificate library.
func (k Keeper) PublishLibrary(ctx sdk.Context, library sdk.AccAddress, publisher sdk.AccAddress) error {
	if k.IsLibrary(ctx, library) {
		return types.ErrLibraryAlreadyExists
	}
	k.SetLibrary(ctx, library, publisher)
	return nil
}

// InvalidateLibrary invalidate a Certificate library.
func (k Keeper) InvalidateLibrary(ctx sdk.Context, library sdk.AccAddress, invalidator sdk.AccAddress) error {
	if !k.IsCertifier(ctx, invalidator) {
		return types.ErrUnqualifiedCertifier
	}
	publisher, err := k.getLibraryPublisher(ctx, library)
	if err != nil {
		return err
	}
	// Can only be invalidate if the invalidator is the original publisher, or that the original publisher is no longer certifier.
	if !(publisher.Equals(invalidator) || !k.IsCertifier(ctx, publisher)) {
		return types.ErrUnqualifiedCertifier
	}
	k.deleteLibrary(ctx, library)
	return nil
}

// IterateAllLibraries iterates over the all the stored libraries and performs a callback function.
func (k Keeper) IterateAllLibraries(ctx sdk.Context, callback func(library types.Library) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LibrariesStoreKey())

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var library types.Library
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &library)

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
		libraryAddresses = append(libraryAddresses, library.Address)
		return false
	})
	return
}
