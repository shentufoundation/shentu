package simulation_test

import (
	gobin "encoding/binary"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	. "github.com/shentufoundation/shentu/v2/x/cvm/simulation"
	"github.com/shentufoundation/shentu/v2/x/cvm/types"
)

func TestDecodeStore(t *testing.T) {
	cdc := shentuapp.MakeEncodingConfig()
	dec := NewDecodeStore(cdc.Marshaler)

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

	meta1 := acm.ContractMeta{
		CodeHash:     bytes1,
		MetadataHash: bytes2,
	}
	meta2 := acm.ContractMeta{
		CodeHash:     bytes1,
		MetadataHash: bytes3,
	}
	metadata := types.ContractMetas{
		Metas: []*acm.ContractMeta{
			&meta1,
			&meta2,
		},
	}

	KVPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.StorageStoreKey(address, key), Value: value1},
			{Key: types.BlockHashStoreKey(int64(height)), Value: value2},
			{Key: types.CodeStoreKey(address), Value: value3},
			{Key: types.AbiStoreKey(address), Value: value4},
			{Key: types.MetaHashStoreKey(metahash), Value: []byte(str)},
			{Key: types.AddressMetaStoreKey(address), Value: cdc.Marshaler.MustMarshal(&metadata)},
		},
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
				require.Panics(t, func() { dec(KVPairs.Pairs[i], KVPairs.Pairs[i]) }, tt.name) // nolint
			} else {
				require.Equal(t, tt.expectedLog, dec(KVPairs.Pairs[i], KVPairs.Pairs[i]), tt.name) // nolint
			}
		})
	}
}
