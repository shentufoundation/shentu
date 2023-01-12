package cli

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func TestSaveLoadKey(t *testing.T) {
	decKey, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	creator := "cosmos1xxkueklal9vejv9unqu80w9vptyepfa95pd53u"
	accAddress, _ := sdk.AccAddressFromBech32(creator)

	decKeyBz := crypto.FromECDSA(decKey.ExportECDSA())
	hasher := sha256.New()
	hasher.Write([]byte(creator))
	hasher.Write(decKeyBz)
	keyFile := fmt.Sprintf("./dec-key-%s-%s.json", creator[6:12], hex.EncodeToString(hasher.Sum(nil))[:6])

	SaveKey(decKey, "./", accAddress.String())
	defer func() {
		if _, err := os.Stat(keyFile); err == nil {
			os.Remove(keyFile)
		}
	}()

	pubKeyBytes := LoadPubKey(keyFile)

	pubEcdsa, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		t.Fatal(err)
	}
	pubEcies := ecies.ImportECDSAPublic(pubEcdsa)

	message := []byte("Hello, world.")
	ct, err := ecies.Encrypt(rand.Reader, pubEcies, message, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Check that decrypting works.
	pt, err := decKey.Decrypt(ct, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pt, message) {
		t.Fatal("ecies: plaintext doesn't match message")
	}

	prvKey := LoadPrvKey(keyFile)
	pt2, err := prvKey.Decrypt(ct, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pt2, message) {
		t.Fatal("ecies: plaintext doesn't match message")
	}
}
