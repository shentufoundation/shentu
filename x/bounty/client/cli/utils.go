package cli

import (
	"crypto/rand"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

	"github.com/tendermint/tendermint/libs/tempfile"
)

// SaveKey saves the given key to a file as json and panics on error.
func SaveKey(privKey *ecies.PrivateKey, dirPath string) {
	if dirPath == "" {
		panic("cannot save private key: filePath not set")
	}

	decKeyBz := crypto.FromECDSA(privKey.ExportECDSA())
	if err := tempfile.WriteFileAtomic(dirPath+"/dec-key.json", decKeyBz, 0666); err != nil {
		panic(err)
	}
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

func GenerateKey() (*ecies.PrivateKey, *ecies.PublicKey, error) {
	decKey, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	if err != nil {
		return nil, nil, err
	}
	encKey := &decKey.PublicKey
	return decKey, encKey, nil
}
