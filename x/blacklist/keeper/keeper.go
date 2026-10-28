package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	store "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/blacklist/types"
)

// Keeper maintains the set of blacklisted addresses.
type Keeper struct {
	cdc           codec.Codec
	storeService  store.KVStoreService
	authority     string
	accountKeeper types.AccountKeeper

	// Entries maps a bech32-encoded address to its blacklist entry.
	Entries collections.Map[string, types.Entry]

	schema collections.Schema
}

// NewKeeper constructs a new blacklist Keeper.
func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	accountKeeper types.AccountKeeper,
	authority string,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:           cdc,
		storeService:  storeService,
		accountKeeper: accountKeeper,
		authority:     authority,
		Entries: collections.NewMap(
			sb,
			collections.NewPrefix([]byte{types.EntriesKeyPrefix}),
			"entries",
			collections.StringKey,
			codec.CollValue[types.Entry](cdc),
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.schema = schema

	return k
}

// GetAuthority returns the module's authority address (governance).
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Cdc returns the keeper's codec, used by the ante decorator to unpack
// authz MsgExec inner messages and resolve message signers.
func (k Keeper) Cdc() codec.Codec {
	return k.cdc
}

// canonicalize normalizes any valid bech32 spelling (e.g. all-uppercase) to
// the canonical lowercase form, so that storage keys and lookups always use
// the same representation.
func canonicalize(address string) (string, error) {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

// IsBlacklisted returns the blacklist entry for the given bech32-encoded
// address and reports whether it is blacklisted. Releases are handled by
// governance via MsgUnbanAddress.
func (k Keeper) IsBlacklisted(ctx context.Context, address string) (types.Entry, bool, error) {
	address, err := canonicalize(address)
	if err != nil {
		return types.Entry{}, false, err
	}
	entry, err := k.Entries.Get(ctx, address)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.Entry{}, false, nil
		}
		return types.Entry{}, false, err
	}
	return entry, true, nil
}

// GetEntry returns a single blacklist entry by address.
func (k Keeper) GetEntry(ctx context.Context, address string) (types.Entry, error) {
	address, err := canonicalize(address)
	if err != nil {
		return types.Entry{}, err
	}
	return k.Entries.Get(ctx, address)
}

// GetAllEntries returns every blacklist entry.
func (k Keeper) GetAllEntries(ctx context.Context) ([]types.Entry, error) {
	var entries []types.Entry
	err := k.Entries.Walk(ctx, nil, func(_ string, entry types.Entry) (bool, error) {
		entries = append(entries, entry)
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return entries, nil
}

// SetEntry stores a blacklist entry under the canonical (lowercase) address
// spelling.
func (k Keeper) SetEntry(ctx context.Context, entry types.Entry) error {
	address, err := canonicalize(entry.Address)
	if err != nil {
		return err
	}
	entry.Address = address
	return k.Entries.Set(ctx, address, entry)
}

// RemoveEntry deletes a blacklist entry. It is a no-op if the address is not
// blacklisted.
func (k Keeper) RemoveEntry(ctx context.Context, address string) error {
	address, err := canonicalize(address)
	if err != nil {
		return err
	}
	return k.Entries.Remove(ctx, address)
}
