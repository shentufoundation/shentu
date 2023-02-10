package cli

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

	"github.com/tendermint/tendermint/libs/tempfile"
)

// SaveKey saves the given key to a file as json and panics on error.
func SaveKey(privKey *ecies.PrivateKey, dirPath string, pid uint64) string {
	if dirPath == "" {
		panic("cannot save private key: filePath not set")
	}

	decKeyBz := crypto.FromECDSA(privKey.ExportECDSA())
	//to create a unique file name
	hashBytes := sha256.Sum256(decKeyBz)
	filename := fmt.Sprintf("dec-key-%d-%x.json", pid, hashBytes[:3])
	fullPath := filepath.Join(dirPath, filename)
	if err := tempfile.WriteFileAtomic(fullPath, decKeyBz, 0666); err != nil {
		panic(err)
	}
	return fullPath
}

// LoadPubKey loads the key at the given location by loading the stored private key and getting the public key part.
func LoadPubKey(filePath string) []byte {
	keyBytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	prvK, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		panic(err)
	}

	return crypto.FromECDSAPub(&prvK.PublicKey)
}

// LoadPrvKey loads the key at the given location by loading the stored private key.
func LoadPrvKey(filePath string) *ecies.PrivateKey {
	keyBytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	prvK, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		panic(err)
	}

	return ecies.ImportECDSA(prvK)
}
