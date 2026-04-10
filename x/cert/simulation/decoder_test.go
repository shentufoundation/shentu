package simulation_test

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/shentufoundation/shentu/v2/x/cert/simulation"
	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

func TestDecodeStore_Certifier(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig()
	types.RegisterInterfaces(encCfg.InterfaceRegistry)
	cdc := encCfg.Codec

	certifier := types.Certifier{
		Address:     "shentu1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5z5tpw",
		Proposer:    "shentu1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5z5tpw",
		Description: "test certifier",
	}

	kvA := kv.Pair{
		Key:   append(types.CertifiersStoreKey(), []byte("addr1")...),
		Value: cdc.MustMarshalLengthPrefixed(&certifier),
	}
	kvB := kv.Pair{
		Key:   append(types.CertifiersStoreKey(), []byte("addr2")...),
		Value: cdc.MustMarshalLengthPrefixed(&certifier),
	}

	decoder := simulation.NewDecodeStore(cdc)
	result := decoder(kvA, kvB)
	require.Contains(t, result, certifier.Address)
}

func TestDecodeStore_Certificate(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig()
	types.RegisterInterfaces(encCfg.InterfaceRegistry)
	cdc := encCfg.Codec

	cert := types.Certificate{
		CertificateId: 1,
		Certifier:     "shentu1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5z5tpw",
		Description:   "test cert",
	}

	kvA := kv.Pair{
		Key:   append(types.CertificatesStoreKey(), []byte{0, 0, 0, 0, 0, 0, 0, 1}...),
		Value: cdc.MustMarshal(&cert),
	}
	kvB := kv.Pair{
		Key:   append(types.CertificatesStoreKey(), []byte{0, 0, 0, 0, 0, 0, 0, 2}...),
		Value: cdc.MustMarshal(&cert),
	}

	decoder := simulation.NewDecodeStore(cdc)
	result := decoder(kvA, kvB)
	require.Contains(t, result, fmt.Sprintf("%v", cert))
}

func TestDecodeStore_NextCertificateID(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig()
	cdc := encCfg.Codec.(codec.BinaryCodec)

	bzA := make([]byte, 8)
	binary.LittleEndian.PutUint64(bzA, 42)
	bzB := make([]byte, 8)
	binary.LittleEndian.PutUint64(bzB, 99)

	kvA := kv.Pair{Key: types.NextCertificateIDStoreKey(), Value: bzA}
	kvB := kv.Pair{Key: types.NextCertificateIDStoreKey(), Value: bzB}

	decoder := simulation.NewDecodeStore(cdc)
	result := decoder(kvA, kvB)
	require.Contains(t, result, "42")
	require.Contains(t, result, "99")
}

func TestDecodeStore_InvalidPrefix(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig()
	cdc := encCfg.Codec.(codec.BinaryCodec)

	kvA := kv.Pair{Key: []byte{0xFF, 0x01}, Value: []byte("data")}
	kvB := kv.Pair{Key: []byte{0xFF, 0x02}, Value: []byte("data")}

	decoder := simulation.NewDecodeStore(cdc)
	require.Panics(t, func() { decoder(kvA, kvB) })
}
