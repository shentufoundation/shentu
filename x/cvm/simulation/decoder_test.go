package simulation

import (
	gobin "encoding/binary"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"

	"github.com/certikfoundation/shentu/x/cvm/internal/types"
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	return cdc
}

func TestDecodeStore(t *testing.T) {
	cdc := makeTestCodec()

	rand.Seed(time.Now().UnixNano())

	bytes1 := make([]byte, binary.Word160Length)
	rand.Read(bytes1)
	address := crypto.Address{}
	copy(address[:], bytes1)

	bytes2 := make([]byte, binary.Word256Bytes)
	rand.Read(bytes2)
	key := binary.Word256{}
	copy(key[:], bytes2)

	height := rand.Intn(50) + 1

	bytes3 := make([]byte, 32)
	rand.Read(bytes3)
	metahash := acmstate.MetadataHash{}
	copy(metahash[:], bytes3)

	value1 := make([]byte, 1+rand.Intn(50))
	rand.Read(value1)

	value2 := make([]byte, 1+rand.Intn(50))
	rand.Read(value2)

	value3 := make([]byte, 1+rand.Intn(50))
	rand.Read(value3)

	value4 := make([]byte, 1+rand.Intn(50))
	rand.Read(value4)

	gasRate := 1 + rand.Uint64()
	gasRateBytes := make([]byte, 8)
	gobin.LittleEndian.PutUint64(gasRateBytes, gasRate)

	str := "odjfg0834u89f"

	metadata := []acm.ContractMeta{
		{
			CodeHash:     bytes1,
			MetadataHash: bytes2,
		},
		{
			CodeHash:     bytes1,
			MetadataHash: bytes3,
		},
	}

	KVPairs := kv.Pairs{
		kv.Pair{Key: types.StorageStoreKey(address, key), Value: value1},
		kv.Pair{Key: types.BlockHashStoreKey(int64(height)), Value: value2},
		kv.Pair{Key: types.CodeStoreKey(address), Value: value3},
		kv.Pair{Key: types.AbiStoreKey(address), Value: value4},
		kv.Pair{Key: types.MetaHashStoreKey(metahash), Value: []byte(str)},
		kv.Pair{Key: types.AddressMetaStoreKey(address), Value: cdc.MustMarshalBinaryLengthPrefixed(metadata)},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Storage", fmt.Sprintf("%b\n%b", value1, value1)},
		{"BlockHash", fmt.Sprintf("%b\n%b", value2, value2)},
		{"Code", fmt.Sprintf("%b\n%b", value3, value3)},
		{"Abi", fmt.Sprintf("%b\n%b", value4, value4)},
		{"MetaHash", fmt.Sprintf("%s\n%s", str, str)},
		{"AddressMetaHash", fmt.Sprintf("%v\n%v", metadata, metadata)},
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
