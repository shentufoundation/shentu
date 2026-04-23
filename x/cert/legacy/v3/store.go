package v3

import (
	"encoding/binary"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// certificateMigrationBatchSize bounds the number of records held in memory
// per migration batch. With small batches the migration's peak memory stays
// constant even on chains with very large certificate histories.
const certificateMigrationBatchSize = 256

// StoreWriter bundles the codec, the raw KV store, and the collections-aware
// write callbacks that the v2→v3 migration needs. Using closures keeps this
// package free of any dependency on x/cert/keeper.
type StoreWriter struct {
	Cdc            codec.BinaryCodec
	Store          storetypes.KVStore
	WriteIndexes   func(cert types.Certificate) error
	SetCertifier   func(types.Certifier) error
	SetCertificate func(id uint64, cert types.Certificate) error
	SetNextCertID  func(id uint64) error
}

// MigrateStore performs the v2→v3 migration:
//  1. Rebuilds all secondary certificate indexes from existing primary records.
//  2. Deletes obsolete validator store entries (prefix 0x1).
//  3. Deletes obsolete platform store entries (prefix 0x2).
//  4. Deletes obsolete library store entries (prefix 0x6).
//  5. Deletes obsolete certifier alias index entries (prefix 0x7).
//  6. Converts the primary certifier and certificate stores plus the
//     NextCertificateID counter to the cosmossdk.io/collections layout.
//     Wire-format changes:
//     - Certifier values: length-prefixed proto → non-prefixed proto
//     (key bytes also change: raw addr concat → sdk.AccAddressKey length-prefixed).
//     - Certificate primary keys: [0x05][8B id_LE] → [0x05][8B id_BE]
//     (collections.Uint64Key is big-endian).
//     - NextCertificateID: [0x08] 8B LE → [0x08] 8B BE (collections.Sequence).
//
// Secondary certificate indexes (0x10–0x14) retain their on-disk format and
// are still written/read via the hand-rolled key helpers in x/cert/types.
// All deletions tolerate already-empty prefixes.
func MigrateStore(w StoreWriter) error {
	// 1. Rebuild secondary indexes for all existing certificates.
	// NOTE: runs before the collections conversion so it can still read the
	// pre-collections primary key format (LE).
	if err := rebuildCertificateIndexes(w); err != nil {
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
		deleteStorePrefix(w.Store, prefix)
	}

	// 6. Convert primary stores to collections layout.
	if err := migrateCertifiersToCollections(w); err != nil {
		return err
	}
	if err := migrateCertificatesToCollections(w); err != nil {
		return err
	}
	return migrateNextCertificateID(w)
}

// rebuildCertificateIndexes iterates every certificate in the primary store and
// writes all five secondary index entries for each one.
func rebuildCertificateIndexes(w StoreWriter) error {
	iterator := storetypes.KVStorePrefixIterator(w.Store, types.CertificatesStoreKey())
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var cert types.Certificate
		w.Cdc.MustUnmarshal(iterator.Value(), &cert)

		// Recover the certificate ID from the primary key.
		// Primary key format: [0x5][8B id_LE]
		key := iterator.Key()
		cert.CertificateId = binary.LittleEndian.Uint64(key[len(key)-8:])

		if err := w.WriteIndexes(cert); err != nil {
			return err
		}
	}
	return nil
}

// migrateCertifiersToCollections rewrites every entry under the certifier
// prefix, replacing the length-prefixed proto value with a non-prefixed one.
// Key bytes also change: raw address concat → length-prefixed (sdk.AccAddressKey).
func migrateCertifiersToCollections(w StoreWriter) error {
	iter := storetypes.KVStorePrefixIterator(w.Store, types.CertifiersStoreKey())

	type entry struct {
		rawKey    []byte
		certifier types.Certifier
	}
	var entries []entry

	for ; iter.Valid(); iter.Next() {
		var certifier types.Certifier
		w.Cdc.MustUnmarshalLengthPrefixed(iter.Value(), &certifier)

		cp := make([]byte, len(iter.Key()))
		copy(cp, iter.Key())
		entries = append(entries, entry{rawKey: cp, certifier: certifier})
	}
	iter.Close()

	for _, e := range entries {
		w.Store.Delete(e.rawKey)
		if err := w.SetCertifier(e.certifier); err != nil {
			return err
		}
	}
	return nil
}

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
func migrateCertificatesToCollections(w StoreWriter) error {
	prefix := types.CertificatesStoreKey()
	end := storetypes.PrefixEndBytes(prefix)
	start := prefix

	type entry struct {
		rawKey []byte
		cert   types.Certificate
	}

	for {
		iter := w.Store.Iterator(start, end)
		batch := make([]entry, 0, certificateMigrationBatchSize)
		var lastOldKey []byte

		for ; iter.Valid() && len(batch) < certificateMigrationBatchSize; iter.Next() {
			key := iter.Key()
			if len(key) < 8 {
				continue
			}
			var cert types.Certificate
			w.Cdc.MustUnmarshal(iter.Value(), &cert)

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
			w.Store.Delete(e.rawKey)
			// Write via collections; secondary indexes are already in place so
			// we don't call WriteIndexes again.
			if err := w.SetCertificate(e.cert.CertificateId, e.cert); err != nil {
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
func migrateNextCertificateID(w StoreWriter) error {
	bz := w.Store.Get(types.NextCertificateIDStoreKey())
	if bz == nil {
		return nil
	}
	id := binary.LittleEndian.Uint64(bz)
	w.Store.Delete(types.NextCertificateIDStoreKey())
	return w.SetNextCertID(id)
}

// deleteStorePrefix deletes all KV entries whose key starts with the given prefix.
// It is safe to call on an already-empty prefix.
func deleteStorePrefix(store storetypes.KVStore, prefix []byte) {
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
}
