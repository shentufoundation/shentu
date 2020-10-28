package simulation

import (
	"encoding/hex"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/x/cvm/internal/keeper"
	"github.com/certikfoundation/shentu/x/cvm/internal/types"
)

const (
	Hello55Code  = "6080604052348015600f57600080fd5b5060888061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80630c49c36c14602d575b600080fd5b60336049565b6040518082815260200191505060405180910390f35b6000603790509056fea2646970667358221220be801fb2205f223dcf6751ff8f3d1996fa2aa8bd72fec7015b75e7c4826e09a264736f6c634300060a0033"
	Hello55Abi   = "[{\"inputs\":[],\"name\":\"sayHi\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]"
	Hello55SayHi = "0c49c36c"

	SimpleCode      = "608060405234801561001057600080fd5b5060c78061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806360fe47b11460375780636d4ce63c146062575b600080fd5b606060048036036020811015604b57600080fd5b8101908080359060200190929190505050607e565b005b60686088565b6040518082815260200191505060405180910390f35b8060008190555050565b6000805490509056fea264697066735822122009c206964bec9aab615f7e1679af3050ffceeb7925dc31b5aae347317d0a74c164736f6c634300060a0033"
	SimpleAbi       = "[{\"inputs\":[],\"name\":\"get\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"set\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	SimpleSetPrefix = "60fe47b1"
	SimpleSetSample = "60fe47b10000000000000000000000000000000000000000000000000000000000000429" // hex for 1065
	SimpleGet       = "6d4ce63c"

	SimpleeventCode      = "608060405234801561001057600080fd5b50610275806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806360fe47b11461003b5780636d4ce63c14610069575b600080fd5b6100676004803603602081101561005157600080fd5b8101908080359060200190929190505050610087565b005b610071610214565b6040518082815260200191505060405180910390f35b806000819055507f6c2b4666ba8da5a95717621d879a77de725f3d816709b9cbe9f059b8f875e2846001600054016040518082815260200191505060405180910390a17f6c96f7d2523b7c8c2eaeba423e4234f7443e80284664c7aa5e285aac4d2613bd60005460020260405180828152602001806020018281038252600e8152602001807f7365636f6e644576656e742121210000000000000000000000000000000000008152506020019250505060405180910390a161014761021d565b6040518060600160405280600160ff168152602001600260ff168152602001600360ff1681525090507f2e71aa6814b90353bea2dc6b23acde4924508f6646936e13a137601134a73d9560036000540182306040518084815260200183600360200280838360005b838110156101ca5780820151818401526020810190506101af565b505050509050018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001935050505060405180910390a15050565b60008054905090565b604051806060016040528060039060208202803683378082019150509050509056fea2646970667358221220b22a748b3370f1eee88fe8672cc53af0fed90a6552bd0e861665ae873bd967fe64736f6c634300060a0033"
	SimpleeventAbi       = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"myNumber\",\"type\":\"uint256\"}],\"name\":\"MyEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"mySecondNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"s\",\"type\":\"string\"}],\"name\":\"MySecondEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"myThirdNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8[3]\",\"name\":\"arr\",\"type\":\"uint8[3]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"a\",\"type\":\"address\"}],\"name\":\"MyThirdEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"get\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"set\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	SimpleeventSetPrefix = "60fe47b1"
	SimpleeventSetSample = "60fe47b100000000000000000000000000000000000000000000000000000000000003e8" // hex for 1000
	SimpleeventGet       = "6d4ce63c"

	StorageCode        = "608060405234801561001057600080fd5b50610121806100206000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80632e64cec11460415780636057361d14605d5780638f3eff7b146088575b600080fd5b604760d0565b6040518082815260200191505060405180910390f35b608660048036036020811015607157600080fd5b810190808035906020019092919050505060d9565b005b608e60e3565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b60008054905090565b8060008190555050565b60003390509056fea2646970667358221220870d6416d9143f210b232ec3a5a9df94095b8ee3803922bab202cf85f2e6a46264736f6c634300060a0033"
	StorageAbi         = "[{\"inputs\":[],\"name\":\"retrieve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sayMyAddres\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	StorageStorePrefix = "6057361d"
	StorageStoreSample = "6057361d0000000000000000000000000000000000000000000000000000000000000225" // hex for 549
	StorageRetrieve    = "2e64cec1"
	StorageSayMyAddres = "8f3eff7b"
)

// DeployContract delivers a deploy tx and returns msg, contract address and error.
func DeployContract(caller simulation.Account, contractCode string, contractAbi string, k keeper.Keeper, r *rand.Rand,
	ctx sdk.Context, chainID string, app *baseapp.BaseApp) (msg types.MsgDeploy, contractAddr sdk.AccAddress, err error) {
	code, err := hex.DecodeString(contractCode)
	if err != nil {
		return msg, nil, err
	}

	msg = types.NewMsgDeploy(caller.Address, uint64(0), code, contractAbi, nil, false, false)

	account := k.AuthKeeper().GetAccount(ctx, caller.Address)
	fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
	if err != nil {
		return msg, nil, err
	}

	tx := helpers.GenTx(
		[]sdk.Msg{msg},
		fees,
		helpers.DefaultGenTxGas,
		chainID,
		[]uint64{account.GetAccountNumber()},
		[]uint64{account.GetSequence()},
		caller.PrivKey,
	)

	_, res, err := app.Deliver(tx)
	if err != nil {
		return msg, nil, err
	}

	return msg, res.Data, nil
}
