package cli

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/gogo/protobuf/proto"

	"github.com/tendermint/tendermint/libs/tempfile"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
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

func KeyAnyToPubKey(keyAny *codectypes.Any) (*ecies.PublicKey, error) {
	if keyAny == nil {
		return nil, fmt.Errorf("empty public key")
	}
	var encryptionKey types.EciesPubKey
	err := proto.Unmarshal(keyAny.GetValue(), &encryptionKey)
	if err != nil {
		return nil, err
	}

	pubEcdsa, err := crypto.UnmarshalPubkey(encryptionKey.EncryptionKey)
	if err != nil {
		return nil, err
	}
	eciesEncKey := ecies.ImportECDSAPublic(pubEcdsa)

	return eciesEncKey, nil
}

func GetRandBytes() ([]byte, *bytes.Reader) {
	randBytes := make([]byte, RandBytesLen)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic("could not read from random source: " + err.Error())
	}
	return randBytes, bytes.NewReader(randBytes)
}
