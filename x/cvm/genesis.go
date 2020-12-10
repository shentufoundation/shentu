package cvm

import (
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/execution/engine"
	"github.com/hyperledger/burrow/permission"

	"github.com/certikfoundation/shentu/x/cvm/internal/keeper"
	"github.com/certikfoundation/shentu/x/cvm/internal/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	k.SetGasRate(ctx, data.GasRate)
	state := k.NewState(ctx)

	callframe := engine.NewCallFrame(state, acmstate.Named("TxCache"))
	cache := callframe.Cache

	for _, contract := range data.Contracts {
		cryptoAddr, err := crypto.AddressFromBytes(contract.Address)
		if err != nil {
			panic(err)
		}
		if contract.Abi != nil {
			k.SetAbi(ctx, cryptoAddr, contract.Abi)
		}

		for _, kv := range contract.Storage {
			key := binary.LeftPadWord256(kv.Key)
			if err := state.SetStorage(cryptoAddr, key, kv.Value); err != nil {
				panic(err)
			}
		}

		// Address Metadata is stored separately.
		var addrMetas []*acm.ContractMeta
		for _, addrMeta := range contract.Meta {
			newMeta := acm.ContractMeta{
				CodeHash:     addrMeta.CodeHash,
				MetadataHash: addrMeta.MetadataHash,
			}
			addrMetas = append(addrMetas, &newMeta)
		}

		if len(addrMetas) > 0 {
			if err := state.SetAddressMeta(cryptoAddr, addrMetas); err != nil {
				panic(err)
			}
		}

		// Register contract account. Since account can already exist from Account InitGenesis,
		// we need to import those first.
		account, err := state.GetAccount(cryptoAddr)
		if err != nil {
			panic(err)
		}
		var balance uint64
		if account == nil {
			balance = 0
		} else {
			balance = account.Balance
		}
		var evmCode, wasmCode acm.Bytecode
		if contract.Code.CodeType == types.CVMCodeTypeEVMCode {
			evmCode = contract.Code.Code
		} else {
			wasmCode = contract.Code.Code
		}
		newAccount := acm.Account{
			Address:  cryptoAddr,
			Balance:  balance,
			EVMCode:  evmCode,
			WASMCode: wasmCode,
			Permissions: permission.AccountPermissions{
				Base: permission.BasePermissions{
					Perms: permission.Call | permission.CreateContract,
				},
			},
			ContractMeta: addrMetas,
		}

		if err := state.UpdateAccount(&newAccount); err != nil {
			panic(err)
		}

	}
	if err := cache.Sync(state); err != nil {
		panic(err)
	}
	for _, metadata := range data.Metadatas {
		if len(metadata.Hash) != 32 {
			panic("metadata hash is not 256 bits")
		}
		var metahash acmstate.MetadataHash
		copy(metahash[:], metadata.Hash[:32])
		if err := state.SetMetadata(metahash, metadata.Metadata); err != nil {
			panic(err)
		}
	}

	keeper.RegisterGlobalPermissionAcc(ctx, k)
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	gasRate := k.GetGasRate(ctx)
	contracts := k.GetAllContracts(ctx)
	metadatas := k.GetAllMetas(ctx)

	return GenesisState{
		GasRate:   gasRate,
		Contracts: contracts,
		Metadatas: metadatas,
	}
}
