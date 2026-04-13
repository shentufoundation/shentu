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
// It performs the following tasks:
//  1. Rebuilds all secondary certificate indexes from existing primary certificate records.
//  2. Deletes obsolete validator store entries (prefix 0x1).
//  3. Deletes obsolete platform store entries (prefix 0x2).
//  4. Deletes obsolete library store entries (prefix 0x6).
//  5. Deletes obsolete certifier alias index entries (prefix 0x7).
//  6. Converts the primary certifier and certificate stores plus the
//     NextCertificateID counter to the cosmossdk.io/collections layout.
//     Wire-format changes:
//     - Certifier values: length-prefixed proto → non-prefixed proto
//       (key bytes also change: raw addr concat → sdk.AccAddressKey length-prefixed).
//     - Certificate primary keys: [0x05][8B id_LE] → [0x05][8B id_BE]
//       (collections.Uint64Key is big-endian).
//     - NextCertificateID: [0x08] 8B LE → [0x08] 8B BE (collections.Sequence).
//
// Secondary certificate indexes (0x10–0x14) retain their on-disk format and
// are still written/read via the hand-rolled key helpers in x/cert/types.
// All deletions tolerate already-empty prefixes.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	// 1. Rebuild secondary indexes for all existing certificates.
	// NOTE: runs before the collections conversion so it can still read the
	// pre-collections primary key format (LE).
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

	// 6. Convert primary stores to collections layout.
	if err := m.migrateCertifiersToCollections(ctx); err != nil {
		return err
	}
	if err := m.migrateCertificatesToCollections(ctx); err != nil {
		return err
	}
	return m.migrateNextCertificateID(ctx)
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

// migrateCertifiersToCollections rewrites every entry under the certifier
// prefix, replacing the length-prefixed proto value with a non-prefixed one.
// Key bytes also change: raw address concat → length-prefixed (sdk.AccAddressKey).
func (m Migrator) migrateCertifiersToCollections(ctx sdk.Context) error {
	store := runtime.KVStoreAdapter(m.keeper.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(store, types.CertifiersStoreKey())

	type entry struct {
		rawKey    []byte
		certifier types.Certifier
	}
	var entries []entry

	for ; iter.Valid(); iter.Next() {
		var certifier types.Certifier
		m.keeper.cdc.MustUnmarshalLengthPrefixed(iter.Value(), &certifier)

		cp := make([]byte, len(iter.Key()))
		copy(cp, iter.Key())
		entries = append(entries, entry{rawKey: cp, certifier: certifier})
	}
	iter.Close()

	for _, e := range entries {
		store.Delete(e.rawKey)
		if err := m.keeper.SetCertifier(ctx, e.certifier); err != nil {
			return err
		}
	}
	return nil
}

// certificateMigrationBatchSize bounds the number of records held in memory
// per migration batch. With small batches the migration's peak memory stays
// constant even on chains with very large certificate histories.
const certificateMigrationBatchSize = 256

// migrateCertificatesToCollections rewrites every certificate primary entry so
// the ID portion of the key is big-endian (as collections.Uint64Key expects).
// The proto value format is already non-prefixed — only the key bytes change.
//
// To keep peak memory bounded, the work is done in batches of
// certificateMigrationBatchSize records: each pass collects up to that many
// old-format keys, deletes them, and writes the new-format entries via
// collections, then resumes iteration just past the last processed old key.
//
// Old- vs new-format detection: both encodings live under the same `0x05`
// prefix byte, so we can't tell them apart by key prefix. Instead, after
// unmarshaling the proto value we compare its CertificateId field against
// big-endian and little-endian decodings of the key tail. A new (collections)
// entry satisfies BE(tail) == CertificateId and is skipped; an old entry
// satisfies LE(tail) == CertificateId and is migrated. This allows the loop
// to safely walk past collections entries written during earlier batches.
func (m Migrator) migrateCertificatesToCollections(ctx sdk.Context) error {
	store := runtime.KVStoreAdapter(m.keeper.storeService.OpenKVStore(ctx))
	prefix := types.CertificatesStoreKey()
	end := storetypes.PrefixEndBytes(prefix)
	start := prefix

	type entry struct {
		rawKey []byte
		cert   types.Certificate
	}

	for {
		iter := store.Iterator(start, end)
		batch := make([]entry, 0, certificateMigrationBatchSize)
		var lastOldKey []byte

		for ; iter.Valid() && len(batch) < certificateMigrationBatchSize; iter.Next() {
			key := iter.Key()
			if len(key) < 8 {
				continue
			}
			var cert types.Certificate
			m.keeper.cdc.MustUnmarshal(iter.Value(), &cert)

			tail := key[len(key)-8:]
			// Already in collections (big-endian) format — skip.
			if binary.BigEndian.Uint64(tail) == cert.CertificateId {
				continue
			}
			// Old (little-endian) format — migrate. CertificateId is already
			// set from the unmarshaled value, which matches LE(tail).

			cp := make([]byte, len(key))
			copy(cp, key)
			batch = append(batch, entry{rawKey: cp, cert: cert})
			lastOldKey = cp
		}
		iter.Close()

		if len(batch) == 0 {
			return nil
		}

		for _, e := range batch {
			store.Delete(e.rawKey)
			// Write via collections; secondary indexes are already in place so
			// we don't call writeCertificateIndexes again.
			if err := m.keeper.Certificates.Set(ctx, e.cert.CertificateId, e.cert); err != nil {
				return err
			}
		}

		// Resume strictly after the last processed old key. Any collections
		// entries written for this batch use BE(id) for the key, which sorts
		// far from old LE keys, so they don't interfere with forward progress.
		start = append(append(make([]byte, 0, len(lastOldKey)+1), lastOldKey...), 0x00)
	}
}

// migrateNextCertificateID rewrites the counter from 8B LE to 8B BE via the
// collections.Sequence wrapper. The key byte (0x08) is unchanged.
func (m Migrator) migrateNextCertificateID(ctx sdk.Context) error {
	store := runtime.KVStoreAdapter(m.keeper.storeService.OpenKVStore(ctx))
	bz := store.Get(types.NextCertificateIDStoreKey())
	if bz == nil {
		return nil
	}
	id := binary.LittleEndian.Uint64(bz)
	store.Delete(types.NextCertificateIDStoreKey())
	return m.keeper.NextCertificateID.Set(ctx, id)
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
