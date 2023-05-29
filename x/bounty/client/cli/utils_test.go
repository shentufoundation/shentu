package cli

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/gogo/protobuf/proto"

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
	desc := "test234567891011"
	encryptedDesc, err := ecies.Encrypt(rand.Reader, &decKey.PublicKey, []byte(desc), nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	var descAny *codectypes.Any
	encDesc := types.EciesEncryptedDesc{
		FindingDesc: encryptedDesc,
	}
	if descAny, err = codectypes.NewAnyWithValue(&encDesc); err != nil {
		t.Fatal(err)
	}

	var descProto types.EciesEncryptedDesc
	err = proto.Unmarshal(descAny.GetValue(), &descProto)
	if err != nil {
		t.Fatal(err)
	}

	descDecrypt, err := decKey.Decrypt(descProto.FindingDesc, nil, nil)

	if string(descDecrypt) != desc {
		t.Fatal(err)
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

	message := []byte("test234567891011")
	_, reader := GetRandBytes()

	ct, err := ecies.Encrypt(reader, pubEcies, message, nil, nil)
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

func TestSaveLoadKey2(t *testing.T) {
	decKey, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	SaveKey(decKey, "./")

	var encAny *codectypes.Any
	encKey := crypto.FromECDSAPub(&decKey.ExportECDSA().PublicKey)

	encKeyMsg := types.EciesPubKey{
		EncryptionKey: encKey,
	}
	encAny, err = codectypes.NewAnyWithValue(&encKeyMsg)
	if err != nil {
		t.Fatal(err.Error())
	}

	testText := "Use a specific height to query state at (this can error if the node is pruning state)"
	encryptionKey := encAny.GetValue()
	pubEcdsa, err := crypto.UnmarshalPubkey(encryptionKey[2:])
	if err != nil {
		t.Fatal(err.Error())
	}
	eciesEncKey := ecies.ImportECDSAPublic(pubEcdsa)

	_, reader := GetRandBytes()
	encryptedDescBytes, err := ecies.Encrypt(reader, eciesEncKey, []byte(testText), nil, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	encDesc := types.EciesEncryptedDesc{
		FindingDesc: encryptedDescBytes,
	}
	descAny, err := codectypes.NewAnyWithValue(&encDesc)
	if err != nil {
		t.Fatal(err.Error())

	}
	//end encrypt
	//start decrypt
	prvKey := LoadPrvKey(keyFile)

	var descProto types.EciesEncryptedDesc
	err = proto.Unmarshal(descAny.GetValue(), &descProto)
	if err != nil {
		t.Fatal(err)
	}

	descBytes, err := prvKey.Decrypt(descProto.FindingDesc, nil, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(descBytes) != testText {
		t.Fatal("error")
	}
}

func TestGetRandBytes(t *testing.T) {
	randBytes, reader := GetRandBytes()

	if reader == nil || randBytes == nil {
		t.Fatal("GetRandBytes error")
	}
}
