package types

import (
	"crypto/sha256"
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func concat(bytes ...[]byte) []byte {
	a := []byte{}
	for _, b := range bytes {
		a = append(a, b...)
	}
	return a
}

const (
	// CertifierStoreKeyPrefix is the prefix byte for certifier kv-store keys.
	CertifierStoreKeyPrefix byte = 0x0
	// CertificateStoreKeyPrefix is the prefix byte for certificate kv-store keys.
	CertificateStoreKeyPrefix byte = 0x5
	// NextCertificateIDKeyPrefix is the prefix byte for the next certificate ID.
	NextCertificateIDKeyPrefix byte = 0x8
)

var (
	// certifierStoreKeyPrefix is the prefix of certifier kv-store keys.
	certifierStoreKeyPrefix = []byte{CertifierStoreKeyPrefix}

	// certificateStoreKeyPrefix is the prefix of certificate kv-store keys.
	certificateStoreKeyPrefix = []byte{CertificateStoreKeyPrefix}

	// nextCertificateIDKeyPrefix is the prefix of the next certificate ID to assign.
	nextCertificateIDKeyPrefix = []byte{NextCertificateIDKeyPrefix}

	// Secondary certificate index prefixes. All ≥ 0x10 to avoid conflicts.

	// certifierIndexKeyPrefix indexes certificates by certifier address.
	certifierIndexKeyPrefix = []byte{0x10}

	// typeIndexKeyPrefix indexes certificates by certificate type.
	typeIndexKeyPrefix = []byte{0x11}

	// certifierTypeIndexKeyPrefix indexes certificates by certifier + certificate type.
	certifierTypeIndexKeyPrefix = []byte{0x12}

	// contentIndexKeyPrefix indexes certificates by SHA-256 hash of content string.
	contentIndexKeyPrefix = []byte{0x13}

	// typeContentIndexKeyPrefix indexes certificates by certificate type + SHA-256 hash of content.
	typeContentIndexKeyPrefix = []byte{0x14}
)

// CertifierStoreKey returns the kv-store key for the certifier registration.
func CertifierStoreKey(certifier sdk.AccAddress) []byte {
	return concat(certifierStoreKeyPrefix, certifier.Bytes())
}

// CertifiersStoreKey returns the kv-store key for accessing all current certifiers in the security council.
func CertifiersStoreKey() []byte {
	return certifierStoreKeyPrefix
}

// CertificateStoreKey returns the kv-store key for accessing a given certificate (ID).
func CertificateStoreKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	return concat(certificateStoreKeyPrefix, bz)
}

// CertificatesStoreKey returns the kv-store key for accessing all certificates.
func CertificatesStoreKey() []byte {
	return certificateStoreKeyPrefix
}

// NextCertificateIDStoreKey returns the kv-store key for next certificate ID to assign.
func NextCertificateIDStoreKey() []byte {
	return nextCertificateIDKeyPrefix
}

// contentHash returns the SHA-256 hash of a content string as a byte slice.
func contentHash(content string) []byte {
	h := sha256.Sum256([]byte(content))
	return h[:]
}

// certIDBigEndian encodes a certificate ID as big-endian 8 bytes for use in index keys.
// Big-endian ensures lexicographic iteration order matches ID order.
func certIDBigEndian(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// CertIDFromIndexKey extracts the certificate ID from an index key.
// All index keys store the cert ID as big-endian uint64 in the last 8 bytes.
func CertIDFromIndexKey(key []byte) uint64 {
	return binary.BigEndian.Uint64(key[len(key)-8:])
}

// CertifierIndexPrefix returns the prefix for all index entries for a given certifier.
// Format: [0x10][20B addr]
func CertifierIndexPrefix(certifier sdk.AccAddress) []byte {
	return concat(certifierIndexKeyPrefix, certifier.Bytes())
}

// CertifierIndexKey returns the full index key for a certifier+certID pair.
// Format: [0x10][20B addr][8B id_BE]
func CertifierIndexKey(certifier sdk.AccAddress, certID uint64) []byte {
	return concat(CertifierIndexPrefix(certifier), certIDBigEndian(certID))
}

// TypeIndexPrefix returns the prefix for all index entries for a given certificate type.
// Format: [0x11][1B type]
func TypeIndexPrefix(certType CertificateType) []byte {
	return concat(typeIndexKeyPrefix, []byte{byte(certType)})
}

// TypeIndexKey returns the full index key for a type+certID pair.
// Format: [0x11][1B type][8B id_BE]
func TypeIndexKey(certType CertificateType, certID uint64) []byte {
	return concat(TypeIndexPrefix(certType), certIDBigEndian(certID))
}

// CertifierTypeIndexPrefix returns the prefix for the certifier+type composite index.
// Format: [0x12][20B addr][1B type]
func CertifierTypeIndexPrefix(certifier sdk.AccAddress, certType CertificateType) []byte {
	return concat(certifierTypeIndexKeyPrefix, certifier.Bytes(), []byte{byte(certType)})
}

// CertifierTypeIndexKey returns the full index key for a certifier+type+certID triple.
// Format: [0x12][20B addr][1B type][8B id_BE]
func CertifierTypeIndexKey(certifier sdk.AccAddress, certType CertificateType, certID uint64) []byte {
	return concat(CertifierTypeIndexPrefix(certifier, certType), certIDBigEndian(certID))
}

// ContentIndexPrefix returns the prefix for the content hash index.
// Format: [0x13][32B sha256(content)]
func ContentIndexPrefix(content string) []byte {
	return concat(contentIndexKeyPrefix, contentHash(content))
}

// ContentIndexKey returns the full index key for a content+certID pair.
// Format: [0x13][32B sha256(content)][8B id_BE]
func ContentIndexKey(content string, certID uint64) []byte {
	return concat(ContentIndexPrefix(content), certIDBigEndian(certID))
}

// TypeContentIndexPrefix returns the prefix for the type+content composite index.
// Format: [0x14][1B type][32B sha256(content)]
func TypeContentIndexPrefix(certType CertificateType, content string) []byte {
	return concat(typeContentIndexKeyPrefix, []byte{byte(certType)}, contentHash(content))
}

// TypeContentIndexKey returns the full index key for a type+content+certID triple.
// Format: [0x14][1B type][32B sha256(content)][8B id_BE]
func TypeContentIndexKey(certType CertificateType, content string, certID uint64) []byte {
	return concat(TypeContentIndexPrefix(certType, content), certIDBigEndian(certID))
}
