package cli

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

const (
	keyFile = "./dec-key.json"
)

func TestAnyToBytes(t *testing.T) {
	decKey, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	desc := "test"
	encryptedDesc, err := ecies.Encrypt(rand.Reader, &decKey.PublicKey, []byte(desc), nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	var descAny *codectypes.Any
	encDesc := types.EciesEncryptedDesc{
		EncryptedDesc: encryptedDesc,
	}
	if descAny, err = codectypes.NewAnyWithValue(&encDesc); err != nil {
		t.Fatal(err)
	}

	descBytes := descAny.GetValue()[2:]
	descDecrypt, err := decKey.Decrypt(descBytes, nil, nil)

	if string(descDecrypt) != desc {
		t.Fatal("error")
	}
}

func TestSaveLoadKey(t *testing.T) {
	decKey, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	// TODO: avoid overwriting silently
	SaveKey(decKey, "./")

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
