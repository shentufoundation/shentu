package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"

	"github.com/certikfoundation/shentu/x/cvm/internal/types"
)

// RandomizedGenState creates a random genesis state for module simulation.
func RandomizedGenState(simState *module.SimulationState) {
	r := simState.Rand
	gs := types.GenesisState{}

	gs.GasRate = 1 + r.Uint64()

	numOfContracts := 1 + r.Intn(50)
	for i := 0; i < numOfContracts; i++ {
		gs.Contracts = append(gs.Contracts, GenerateAContract(r))
	}

	numOfMetadata := 1 + r.Intn(50)
	for i := 0; i < numOfMetadata; i++ {
		gs.Metadata = append(gs.Metadata, GenerateAMetadata(r))
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(gs)
}

// GenerateAContract returns a Contract object with all of its fields randomized.
func GenerateAContract(r *rand.Rand) types.Contract {
	contract := types.Contract{}

	bytes := make([]byte, binary.Word160Length)
	r.Read(bytes)
	address := crypto.Address{}
	copy(address[:], bytes)
	contract.Address = address

	contract.Code = types.NewCVMCode(types.CVMCodeTypeEVMCode, []byte(simulation.RandStringOfLength(r, 1+r.Intn(50))))

	contract.Abi = []byte(simulation.RandStringOfLength(r, 1+r.Intn(50)))

	numOfContractMetas := r.Intn(50)
	for i := 0; i < numOfContractMetas; i++ {
		contract.Meta = append(contract.Meta, GenerateAContractMeta(r))
	}

	return contract
}

// GenerateAMetadata returns a Metadata object with all of its fields randomized.
func GenerateAMetadata(r *rand.Rand) types.Metadata {
	bytes := make([]byte, 32)
	r.Read(bytes)
	hash := acmstate.MetadataHash{}
	copy(hash[:], bytes)

	return types.Metadata{
		Hash:     hash,
		Metadata: simulation.RandStringOfLength(r, 1+r.Intn(50)),
	}
}

// GenerateAContractMeta returns a ContractMeta object with all of its fields randomized.
func GenerateAContractMeta(r *rand.Rand) types.ContractMeta {
	return types.ContractMeta{
		CodeHash:     []byte(simulation.RandStringOfLength(r, 1+r.Intn(50))),
		MetadataHash: []byte(simulation.RandStringOfLength(r, 1+r.Intn(50))),
	}
}
