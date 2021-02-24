package migrate

import (
	"fmt"

	cvmtypes "github.com/certikfoundation/shentu/x/cvm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
)

// CVMCodeType is the type for code in CVM.
type CVMCodeType byte

// CVM code types
const (
	CVMCodeTypeEVMCode CVMCodeType = iota
	CVMCodeTypeEWASMCode
)

// CVMCode defines the data structure of code in CVM.
type CVMCode struct {
	CodeType CVMCodeType
	Code     acm.Bytecode
}

type Contract struct {
	Address crypto.Address `json:"address"`
	Code    CVMCode        `json:"code"`
	Storage []Storage      `json:"storage"`
	Abi     []byte         `json:"abi"`
	Meta    []ContractMeta `json:"meta"`
}

type ContractMeta struct {
	CodeHash     []byte
	MetadataHash []byte
}

type Storage struct {
	Key   binary.Word256 `json:"key"`
	Value []byte         `json:"value"`
}

type Metadata struct {
	Hash     acmstate.MetadataHash `json:"hash"`
	Metadata string                `json:"metadata"`
}

// GenesisState is a cvm genesis state.
type CVMGenesisState struct {
	// GasRate defines the gas exchange rate between Cosmos gas and CVM gas.
	// CVM gas equals to Cosmos Gas * gasRate.
	GasRate   uint64     `json:"gasrate"`
	Contracts []Contract `json:"contracts"`
	Metadata  []Metadata `json:"metadata"`
}

func RegisterCVMLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(acm.Bytecode{}, "acm/Bytecode", nil)
	cdc.RegisterConcrete(binary.Word256{}, "binary/Word256", nil)
	cdc.RegisterConcrete([]acm.ContractMeta{}, "cvm/ContractMeta", nil)
}

func migrateCVM(oldGenState CVMGenesisState) *cvmtypes.GenesisState {
	newContracts := cvmtypes.Contracts{}
	for _, cont := range oldGenState.Contracts {
		var storages []cvmtypes.Storage
		for _, s := range cont.Storage {
			storages = append(storages, cvmtypes.Storage{
				Key:   s.Key,
				Value: s.Value})
		}
		var metas []cvmtypes.ContractMeta
		for _, s := range cont.Meta {
			metas = append(metas, cvmtypes.ContractMeta{
				CodeHash:     s.CodeHash,
				MetadataHash: s.MetadataHash})
		}
		addContract := cvmtypes.Contract{
			Address: cont.Address,
			Code: cvmtypes.CVMCode{
				CodeType: int64(cont.Code.CodeType),
				Code:     cont.Code.Code,
			},
			Storage: storages,
			Abi:     cont.Abi,
			Meta:    metas,
		}
		newContracts = append(newContracts, addContract)
	}

	newMetas := cvmtypes.Metadatas{}
	for _, meta := range oldGenState.Metadata {
		fmt.Println(meta.Hash)
		fmt.Println(crypto.Keccak256([]byte(meta.Metadata)))
		newMeta := cvmtypes.Metadata{
			Hash:     meta.Hash.Bytes(),
			Metadata: meta.Metadata,
		}
		fmt.Println(newMeta.Hash)
		fmt.Println(newMeta.Hash)
		fmt.Println(len(newMeta.Hash))
		fmt.Println(len(newMeta.Hash))
		newMetas = append(newMetas, newMeta)
	}
	return &cvmtypes.GenesisState{
		GasRate:   oldGenState.GasRate,
		Contracts: newContracts,
		Metadatas: newMetas,
	}
}
