package simulation_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types/kv"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	. "github.com/shentufoundation/shentu/v2/x/cert/simulation"
	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

func TestDecodeStore(t *testing.T) {
	app := shentuapp.Setup(false)
	cdc := app.Codec()

	rand.Seed(time.Now().UnixNano())

	certifier := types.Certifier{
		Address:     RandomAccount().Address.String(),
		Proposer:    RandomAccount().Address.String(),
		Description: "this is a test case.",
	}

	validatorPubKey := RandomAccount().PubKey
	var pkAny *codectypes.Any
	if validatorPubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(validatorPubKey); err != nil {
			panic(err)
		}
	}

	platformPubKey := RandomAccount().PubKey
	if validatorPubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(validatorPubKey); err != nil {
			panic(err)
		}
	}
	platform := types.Platform{
		ValidatorPubkey: pkAny,
		Description:     "This is a test case.",
	}

	libraryAddr := sdk.AccAddress("f23908hf932")
	library := types.Library{
		Address:   libraryAddr.String(),
		Publisher: sdk.AccAddress("0092uf32").String(),
	}

	aliasCertifier := types.Certifier{
		Address:     RandomAccount().Address.String(),
		Alias:       "Alice",
		Proposer:    RandomAccount().Address.String(),
		Description: "this is a test case.",
	}

	certifierAddr, err := sdk.AccAddressFromBech32(certifier.Address)
	require.NoError(t, err)

	kvPairs := []kv.Pair{
		{Key: types.CertifierStoreKey(certifierAddr), Value: cdc.MustMarshalLengthPrefixed(&certifier)},
		{Key: types.PlatformStoreKey(platformPubKey), Value: cdc.MustMarshalLengthPrefixed(&platform)},
		{Key: types.LibraryStoreKey(libraryAddr), Value: cdc.MustMarshalLengthPrefixed(&library)},
		{Key: types.CertifierAliasStoreKey(aliasCertifier.Alias), Value: cdc.MustMarshalLengthPrefixed(&aliasCertifier)},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Certifier", fmt.Sprintf("%v\n%v", certifier, certifier)},
		{"Platform", fmt.Sprintf("%v\n%v", platform, platform)},
		{"Library", fmt.Sprintf("%v\n%v", library, library)},
		{"Alias certifier", fmt.Sprintf("%v\n%v", aliasCertifier, aliasCertifier)},
		{"other", ""},
	}

	decoder := NewDecodeStore(cdc)

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if i == len(tests)-1 {
				require.Panics(t, func() { decoder(kvPairs[i], kvPairs[i]) }, tt.name)
			} else {
				require.Equal(t, tt.expectedLog, decoder(kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}

// RandomAccount generates a random Account object.
func RandomAccount() simtypes.Account {
	privkeySeed := make([]byte, 15)
	rand.Read(privkeySeed)

	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	address := sdk.AccAddress(pubKey.Address())

	return simtypes.Account{
		PrivKey: privKey,
		PubKey:  pubKey,
		Address: address,
	}
}
