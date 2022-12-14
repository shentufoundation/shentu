package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/tendermint/tendermint/libs/tempfile"
	"os"
)

// SaveKey saves the given key to a file as json and panics on error.
func SaveKeys(pubKey cryptotypes.PubKey, privKey cryptotypes.PrivKey, dirPath string, cdc codec.Codec) {
	if dirPath == "" {
		panic("cannot save PrivValidator key: filePath not set")
	}

	decKeyBz := cdc.MustMarshalJSON(privKey)
	if err := tempfile.WriteFileAtomic(dirPath + "/dec-key.json", decKeyBz, 0666); err != nil {
		panic(err)
	}
}

// LoadKey loads the key at the given location by loading the stored private key and getting the public key part.
func LoadPubKey(filePath string, cdc codec.Codec) cryptotypes.PubKey {
	keyJSONBytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var privKey secp256k1.PrivKey
	cdc.MustUnmarshalJSON(keyJSONBytes, &privKey)
	if err != nil {
		panic(err)
	}

	return privKey.PubKey()
}