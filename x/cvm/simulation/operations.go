package simulation

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
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
	}
}

// SimulateMsgDeployHello55 creates a massage deploying /tests/hello55.sol contract.
func SimulateMsgDeployHello55(k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)
		code, err := hex.DecodeString(Hello55Code)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgDeploy(caller.Address, uint64(0), code, Hello55Abi, nil, false, false)

		account := k.AuthKeeper().GetAccount(ctx, caller.Address)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
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
			fmt.Printf("<<<<<<<<< hello55 deploy error: %s\n", err)
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check pure/view function ret
		data, err := hex.DecodeString(Hello55SayHi)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		ret, err := k.Call(ctx, caller.Address, res.Data, 0, data, nil, true, false, false)
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
		code, err := hex.DecodeString(SimpleCode)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgDeploy(caller.Address, uint64(0), code, SimpleAbi, nil, false, false)

		account := k.AuthKeeper().GetAccount(ctx, caller.Address)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
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
			fmt.Printf("<<<<<<<<< simple deploy error: %s\n", err)
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check pure/view function ret
		data, err := hex.DecodeString(SimpleGet)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		ret, err := k.Call(ctx, caller.Address, res.Data, 0, data, nil, true, false, false)
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
				Op:          SimulateMsgCallSimpleSet(k, res.Data, int(r.Uint32())),
			},
		}

		return simulation.NewOperationMsg(msg, true, ""), futureOperations, nil
	}
}

// SimulateMsgCallSimpleSet creates a message calling set func in /tests/simple.sol contract.
func SimulateMsgCallSimpleSet(k keeper.Keeper, contractAddr sdk.AccAddress, varValue int) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)

		hexStr := strconv.FormatInt(int64(varValue), 16)
		length := len(hexStr)
		for i := 0; i < 64-length; i++ {
			hexStr = "0" + hexStr
		}
		data, err := hex.DecodeString(SimpleSetPrefix + hexStr)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgCall(caller.Address, contractAddr, 0, data)

		account := k.AuthKeeper().GetAccount(ctx, caller.Address)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
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

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check pure/view function ret
		data, err = hex.DecodeString(SimpleGet)
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
		code, err := hex.DecodeString(SimpleeventCode)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgDeploy(caller.Address, uint64(0), code, SimpleeventAbi, nil, false, false)

		account := k.AuthKeeper().GetAccount(ctx, caller.Address)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
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
			fmt.Printf("<<<<<<<<< simpleevent deploy error: %s\n", err)
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check pure/view function ret
		data, err := hex.DecodeString(SimpleeventGet)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		ret, err := k.Call(ctx, caller.Address, res.Data, 0, data, nil, true, false, false)
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
				Op:          SimulateMsgCallSimpleEventSet(k, res.Data, int(r.Uint32())),
			},
		}

		return simulation.NewOperationMsg(msg, true, ""), futureOperations, nil
	}
}

// SimulateMsgCallSimpleEventSet creates a message calling set func in /tests/simpleevent.sol contract.
func SimulateMsgCallSimpleEventSet(k keeper.Keeper, contractAddr sdk.AccAddress, varValue int) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)

		hexStr := strconv.FormatInt(int64(varValue), 16)
		length := len(hexStr)
		for i := 0; i < 64-length; i++ {
			hexStr = "0" + hexStr
		}
		data, err := hex.DecodeString(SimpleeventSetPrefix + hexStr)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgCall(caller.Address, contractAddr, 0, data)

		account := k.AuthKeeper().GetAccount(ctx, caller.Address)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
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

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check pure/view function ret
		data, err = hex.DecodeString(SimpleeventGet)
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
