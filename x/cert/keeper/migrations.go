package keeper

import (
	"encoding/binary"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate1to2 migrates from version 1 to 2 (no-op; address prefix migration was already applied).
func (m Migrator) Migrate1to2(_ sdk.Context) error {
	return nil
}

// Migrate2to3 migrates from version 2 to 3.
//
// It performs four tasks:
//  1. Rebuilds all secondary certificate indexes from existing primary certificate records.
//  2. Deletes obsolete validator store entries (prefix 0x1).
//  3. Deletes obsolete platform store entries (prefix 0x2).
//  4. Deletes obsolete library store entries (prefix 0x6).
//  5. Deletes obsolete certifier alias index entries (prefix 0x7).
//
// All deletions tolerate already-empty prefixes.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	// 1. Rebuild secondary indexes for all existing certificates.
	if err := m.rebuildCertificateIndexes(ctx); err != nil {
		return err
	}

	// 2–5. Delete obsolete store prefixes.
	obsoletePrefixes := [][]byte{
		{0x1}, // validator certifications
		{0x2}, // platform certifications
		{0x6}, // library registrations
		{0x7}, // certifier alias index
	}
	for _, prefix := range obsoletePrefixes {
		if err := deleteStorePrefix(ctx, m.keeper, prefix); err != nil {
			return err
		}
	}
	return nil
}

// rebuildCertificateIndexes iterates every certificate in the primary store and
// writes all five secondary index entries for each one.
func (m Migrator) rebuildCertificateIndexes(ctx sdk.Context) error {
	store := runtime.KVStoreAdapter(m.keeper.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.CertificatesStoreKey())
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var cert types.Certificate
		m.keeper.cdc.MustUnmarshal(iterator.Value(), &cert)

		// Recover the certificate ID from the primary key.
		// Primary key format: [0x5][8B id_LE]
		key := iterator.Key()
		cert.CertificateId = binary.LittleEndian.Uint64(key[len(key)-8:])

		if err := m.keeper.writeCertificateIndexes(ctx, cert); err != nil {
			return err
		}
	}
	return nil
}

// deleteStorePrefix deletes all KV entries whose key starts with the given prefix.
// It is safe to call on an already-empty prefix.
func deleteStorePrefix(ctx sdk.Context, k Keeper, prefix []byte) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(store, prefix)

	// Collect keys before deleting to avoid iterator invalidation.
	var keys [][]byte
	for ; iter.Valid(); iter.Next() {
		cp := make([]byte, len(iter.Key()))
		copy(cp, iter.Key())
		keys = append(keys, cp)
	}
	iter.Close()

	for _, key := range keys {
		store.Delete(key)
	}
	return nil
}
