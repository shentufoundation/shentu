package simulation

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

func TestDecodeStore(t *testing.T) {
	cdc := makeTestCodec()
	rand.Seed(time.Now().UnixNano())

	certifier := types.Certifier{
		Address:     RandomAccount().Address,
		Proposer:    RandomAccount().Address,
		Description: "this is a test case.",
	}

	validator := types.Validator{
		PubKey:    RandomAccount().PubKey,
		Certifier: RandomAccount().Address,
	}

	platform := types.Platform{
		Address:     sdk.GetConsAddress(RandomAccount().PubKey),
		Description: "This is a test case.",
	}

	library := types.Library{
		Address:   sdk.AccAddress("f23908hf932"),
		Publisher: sdk.AccAddress("0092uf32"),
	}

	aliasCertifier := types.Certifier{
		Address:     RandomAccount().Address,
		Alias:       "Alice",
		Proposer:    RandomAccount().Address,
		Description: "this is a test case.",
	}

	KVPairs := kv.Pairs{
		kv.Pair{Key: types.CertifierStoreKey(certifier.Address), Value: cdc.MustMarshalBinaryLengthPrefixed(&certifier)},
		kv.Pair{Key: types.ValidatorStoreKey(validator.PubKey), Value: cdc.MustMarshalBinaryLengthPrefixed(&validator)},
		kv.Pair{Key: types.PlatformStoreKey(platform.Address), Value: []byte(platform.Description)},
		kv.Pair{Key: types.LibraryStoreKey(library.Address), Value: cdc.MustMarshalBinaryLengthPrefixed(&library)},
		kv.Pair{Key: types.CertifierAliasStoreKey(aliasCertifier.Alias), Value: cdc.MustMarshalBinaryLengthPrefixed(&aliasCertifier)},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Certifier", fmt.Sprintf("%v\n%v", certifier, certifier)},
		{"Validator", fmt.Sprintf("%v\n%v", validator, validator)},
		{"Platform", fmt.Sprintf("%s\n%s", platform.Description, platform.Description)},
		{"Library", fmt.Sprintf("%v\n%v", library, library)},
		{"Alias certifier", fmt.Sprintf("%v\n%v", aliasCertifier, aliasCertifier)},
		{"other", ""},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if i == len(tests)-1 { // nolint
				require.Panics(t, func() { DecodeStore(cdc, KVPairs[i], KVPairs[i]) }, tt.name) // nolint
			} else {
				require.Equal(t, tt.expectedLog, DecodeStore(cdc, KVPairs[i], KVPairs[i]), tt.name) // nolint
			}
		})
	}
}

// RandomAccount generates a random Account object.
func RandomAccount() simulation.Account {
	privkeySeed := make([]byte, 15)
	rand.Read(privkeySeed)

	privKey := secp256k1.GenPrivKeySecp256k1(privkeySeed)
	pubKey := privKey.PubKey()
	address := sdk.AccAddress(pubKey.Address())

	return simulation.Account{
		PrivKey: privKey,
		PubKey:  pubKey,
		Address: address,
	}
}
