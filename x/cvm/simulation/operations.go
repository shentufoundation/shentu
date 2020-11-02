package simulation

import (
	"encoding/hex"
	"math/rand"
	"strconv"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/x/cvm/internal/keeper"
	"github.com/certikfoundation/shentu/x/cvm/internal/types"
)

const (
	OpWeightMsgDeploy = "op_weight_msg_deploy"
)

// WeightedOperations creates an operation with a weight for each type of message generators.
func WeightedOperations(appParams simulation.AppParams, cdc *codec.Codec, k keeper.Keeper) simulation.WeightedOperations {
	var weightMsgDeploy int
	appParams.GetOrGenerate(cdc, OpWeightMsgDeploy, &weightMsgDeploy, nil,
		func(_ *rand.Rand) {
			weightMsgDeploy = simappparams.DefaultWeightMsgSend
		})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgDeploy, SimulateMsgDeployHello55(k)),
		simulation.NewWeightedOperation(weightMsgDeploy, SimulateMsgDeploySimple(k)),
		simulation.NewWeightedOperation(weightMsgDeploy, SimulateMsgDeploySimpleEvent(k)),
		simulation.NewWeightedOperation(weightMsgDeploy, SimulateMsgDeployStorage(k)),
		simulation.NewWeightedOperation(weightMsgDeploy, SimulateMsgDeployStringTest(k)),
	}
}

// SimulateMsgDeployHello55 creates a massage deploying /tests/hello55.sol contract.
func SimulateMsgDeployHello55(k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)

		// deploy hello55.sol
		msg, contractAddr, err := DeployContract(caller, Hello55Code, Hello55Abi, k, r, ctx, chainID, app)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check sayHi() ret
		data, err := hex.DecodeString(Hello55SayHi)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		ret, err := k.Call(ctx, caller.Address, contractAddr, 0, data, nil, true, false, false)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		value, err := strconv.ParseInt(hex.EncodeToString(ret), 16, 32)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		if value != 55 {
			panic("return value incorrect")
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgDeploySimple creates a massage deploying /tests/simple.sol contract.
func SimulateMsgDeploySimple(k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)

		// deploy simple.sol
		msg, contractAddr, err := DeployContract(caller, SimpleCode, SimpleAbi, k, r, ctx, chainID, app)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check get() ret
		data, err := hex.DecodeString(SimpleGet)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		ret, err := k.Call(ctx, caller.Address, contractAddr, 0, data, nil, true, false, false)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		value, err := strconv.ParseInt(hex.EncodeToString(ret), 16, 64)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		if value != 0 {
			panic("return value incorrect")
		}

		futureOperations := []simulation.FutureOperation{
			{
				BlockHeight: int(ctx.BlockHeight()) + r.Intn(10),
				Op:          SimulateMsgCallSimpleSet(k, contractAddr, int(r.Uint32())),
			},
		}

		return simulation.NewOperationMsg(msg, true, ""), futureOperations, nil
	}
}

// SimulateMsgCallSimpleSet creates a message calling set() in /tests/simple.sol contract.
func SimulateMsgCallSimpleSet(k keeper.Keeper, contractAddr sdk.AccAddress, varValue int) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)

		hexStr := strconv.FormatInt(int64(varValue), 16)
		length := len(hexStr)
		for i := 0; i < 64-length; i++ {
			hexStr = "0" + hexStr
		}

		// call set()
		msg, _, err := CallFunction(caller, SimpleSetPrefix, hexStr, contractAddr, k, ctx, r, chainID, app)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check get() ret
		data, err := hex.DecodeString(SimpleGet)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		ret, err := k.Call(ctx, caller.Address, contractAddr, 0, data, nil, true, false, false)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		value, err := strconv.ParseInt(hex.EncodeToString(ret), 16, 64)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		if value != int64(varValue) {
			panic("return value incorrect")
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgDeploySimpleEvent creates a massage deploying /tests/simpleevent.sol contract.
func SimulateMsgDeploySimpleEvent(k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)

		// deploy simpleevent.sol
		msg, contractAddr, err := DeployContract(caller, SimpleeventCode, SimpleeventAbi, k, r, ctx, chainID, app)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check get() ret
		data, err := hex.DecodeString(SimpleeventGet)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		ret, err := k.Call(ctx, caller.Address, contractAddr, 0, data, nil, true, false, false)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		value, err := strconv.ParseInt(hex.EncodeToString(ret), 16, 64)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		if value != 0 {
			panic("return value incorrect")
		}

		futureOperations := []simulation.FutureOperation{
			{
				BlockHeight: int(ctx.BlockHeight()) + r.Intn(10),
				Op:          SimulateMsgCallSimpleEventSet(k, contractAddr, int(r.Uint32())),
			},
		}

		return simulation.NewOperationMsg(msg, true, ""), futureOperations, nil
	}
}

// SimulateMsgCallSimpleEventSet creates a message calling set() in /tests/simpleevent.sol contract.
func SimulateMsgCallSimpleEventSet(k keeper.Keeper, contractAddr sdk.AccAddress, varValue int) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)

		hexStr := strconv.FormatInt(int64(varValue), 16)
		length := len(hexStr)
		for i := 0; i < 64-length; i++ {
			hexStr = "0" + hexStr
		}

		// call set()
		msg, _, err := CallFunction(caller, SimpleeventSetPrefix, hexStr, contractAddr, k, ctx, r, chainID, app)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check get() ret
		data, err := hex.DecodeString(SimpleeventGet)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		ret, err := k.Call(ctx, caller.Address, contractAddr, 0, data, nil, true, false, false)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		value, err := strconv.ParseInt(hex.EncodeToString(ret), 16, 64)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		if value != int64(varValue) {
			panic("return value incorrect")
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgDeployStorage creates a massage deploying /tests/storage.sol contract.
func SimulateMsgDeployStorage(k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)

		// deploy storage.sol
		msg, contractAddr, err := DeployContract(caller, StorageCode, StorageAbi, k, r, ctx, chainID, app)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check retrieve() ret
		data, err := hex.DecodeString(StorageRetrieve)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		ret, err := k.Call(ctx, caller.Address, contractAddr, 0, data, nil, true, false, false)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		value, err := strconv.ParseInt(hex.EncodeToString(ret), 16, 64)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		if value != 0 {
			panic("return value incorrect")
		}

		// check sayMyAddres() ret
		data, err = hex.DecodeString(StorageSayMyAddres)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		ret, err = k.Call(ctx, caller.Address, contractAddr, 0, data, nil, true, false, false)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		sender := sdk.AccAddress(ret[12:])
		if !sender.Equals(caller.Address) {
			panic("return value incorrect")
		}

		futureOperations := []simulation.FutureOperation{
			{
				BlockHeight: int(ctx.BlockHeight()) + r.Intn(10),
				Op:          SimulateMsgCallStorageStore(k, contractAddr, int(r.Uint32())),
			},
		}

		return simulation.NewOperationMsg(msg, true, ""), futureOperations, nil
	}
}

// SimulateMsgCallStorageStore creates a message calling store() in /tests/storage.sol contract.
func SimulateMsgCallStorageStore(k keeper.Keeper, contractAddr sdk.AccAddress, varValue int) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)

		hexStr := strconv.FormatInt(int64(varValue), 16)
		length := len(hexStr)
		for i := 0; i < 64-length; i++ {
			hexStr = "0" + hexStr
		}

		// call store()
		msg, _, err := CallFunction(caller, StorageStorePrefix, hexStr, contractAddr, k, ctx, r, chainID, app)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check retrieve() ret
		data, err := hex.DecodeString(StorageRetrieve)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		ret, err := k.Call(ctx, caller.Address, contractAddr, 0, data, nil, true, false, false)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		value, err := strconv.ParseInt(hex.EncodeToString(ret), 16, 64)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		if value != int64(varValue) {
			panic("return value incorrect")
		}

		// check sayMyAddres() ret
		data, err = hex.DecodeString(StorageSayMyAddres)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		ret, err = k.Call(ctx, caller.Address, contractAddr, 0, data, nil, true, false, false)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		sender := sdk.AccAddress(ret[12:])
		if !sender.Equals(caller.Address) {
			panic("return value incorrect")
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgDeployStringTest creates a massage deploying /tests/stringtest.sol contract.
func SimulateMsgDeployStringTest(k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)

		// deploy stringtest.sol
		msg, contractAddr, err := DeployContract(caller, StringtestCode, StringtestAbi, k, r, ctx, chainID, app)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		var ref string // hex str shared among future operations for checking purpose

		futureOperations := []simulation.FutureOperation{
			{
				BlockHeight: int(ctx.BlockHeight()) + 1,
				Op:          SimulateMsgCallStringTestGetl(k, contractAddr, &ref),
			},
			{
				BlockHeight: int(ctx.BlockHeight()) + 1,
				Op:          SimulateMsgCallStringTestGets(k, contractAddr, &ref),
			},
			{
				BlockHeight: int(ctx.BlockHeight()) + 2,
				Op:          SimulateMsgCallStringTestChangeString(k, contractAddr, &ref),
			},
			{
				BlockHeight: int(ctx.BlockHeight()) + 3,
				Op:          SimulateMsgCallStringTestGets(k, contractAddr, &ref),
			},
			{
				BlockHeight: int(ctx.BlockHeight()) + 3,
				Op:          SimulateMsgCallStringTestGetl(k, contractAddr, &ref),
			},
			{
				BlockHeight: int(ctx.BlockHeight()) + r.Intn(10),
				Op:          SimulateMsgCallStringTestChangeGiven(k, contractAddr),
			},
			{
				BlockHeight: int(ctx.BlockHeight()) + r.Intn(10),
				Op:          SimulateMsgCallStringTestTestStuff(k, contractAddr),
			},
		}

		return simulation.NewOperationMsg(msg, true, ""), futureOperations, nil
	}
}

// SimulateMsgCallStringTestChangeString creates a message calling changeString() in /tests/stringtest.sol contract.
func SimulateMsgCallStringTestChangeString(k keeper.Keeper, contractAddr sdk.AccAddress, ref *string) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)

		// turn length into a hex string of length 64
		length := r.Intn(32) + 1
		hexLen := strconv.FormatInt(int64(length), 16)
		l := len(hexLen)
		for i := 0; i < 64-l; i++ {
			hexLen = "0" + hexLen
		}

		// turn string into a hex string of length 64
		*ref = simulation.RandStringOfLength(r, length)
		hexStr := hex.EncodeToString([]byte(*ref))
		l = len(hexStr)
		for i := 0; i < 64-l; i++ {
			hexStr = hexStr + "0"
		}

		// call changeString()
		msg, ret, err := CallFunction(caller, StringtestChangeStringPrefix, hexLen+hexStr, contractAddr, k, ctx, r, chainID, app)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check ret and update ref
		*ref = StringtestChangeStringPrefix[8:] + hexLen + hexStr
		if hex.EncodeToString(ret) != *ref {
			panic("return value incorrect")
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgCallStringTestChangeGiven creates a message calling changeGiven() in /tests/stringtest.sol contract.
func SimulateMsgCallStringTestChangeGiven(k keeper.Keeper, contractAddr sdk.AccAddress) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)

		// assemble length into a hex string of length 64
		length := r.Intn(30) + 3
		hexLen := strconv.FormatInt(int64(length), 16)
		l := len(hexLen)
		for i := 0; i < 64-l; i++ {
			hexLen = "0" + hexLen
		}

		// assemble string into a hex string of length 64
		str := simulation.RandStringOfLength(r, length)
		hexStr := hex.EncodeToString([]byte(str))
		l = len(hexStr)
		for i := 0; i < 64-l; i++ {
			hexStr = hexStr + "0"
		}

		// call changeGiven()
		msg, ret, err := CallFunction(caller, StringtestChangeGivenPrefix, hexLen+hexStr, contractAddr, k, ctx, r, chainID, app)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check ret
		ref := StringtestChangeStringPrefix[8:] + hexLen + hex.EncodeToString([]byte("Abc")) + hexStr[6:]
		if hex.EncodeToString(ret) != ref {
			panic("return value incorrect")
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgCallStringTestGets creates a message calling gets() in /tests/stringtest.sol contract.
func SimulateMsgCallStringTestGets(k keeper.Keeper, contractAddr sdk.AccAddress, ref *string) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)

		// call gets()
		msg, ret, err := CallFunction(caller, StringtestGets, "", contractAddr, k, ctx, r, chainID, app)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check ret
		if *ref == "" && len(ret) != 64 {
			panic("return value incorrect")
		}
		if *ref != "" && hex.EncodeToString(ret) != *ref {
			panic("return value incorrect")
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgCallStringTestGetl creates a message calling getl() in /tests/stringtest.sol contract.
func SimulateMsgCallStringTestGetl(k keeper.Keeper, contractAddr sdk.AccAddress, ref *string) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)

		// call getl()
		msg, ret, err := CallFunction(caller, StringtestGetl, "", contractAddr, k, ctx, r, chainID, app)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check ret
		length, err := strconv.ParseInt(hex.EncodeToString(ret), 16, 32)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		if *ref == "" && length != 0 {
			panic("return value incorrect")
		}
		str := *ref
		if str != "" && str[64:128] != hex.EncodeToString(ret) {
			panic("return value incorrect")
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgCallStringTestTestStuff creates a message calling testStuff() in /tests/stringtest.sol contract.
func SimulateMsgCallStringTestTestStuff(k keeper.Keeper, contractAddr sdk.AccAddress) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)

		// call testStuff()
		msg, ret, err := CallFunction(caller, StringtestTestStuff, "", contractAddr, k, ctx, r, chainID, app)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check ret
		value, err := strconv.ParseInt(hex.EncodeToString(ret), 16, 32)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		if value != 123123 {
			panic("return value incorrect")
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
