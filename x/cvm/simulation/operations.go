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
	}
}

// SimulateMsgCall creates a message operation of MsgCall with randomized field values.
func SimulateMsgCall(k keeper.Keeper, contractAddr sdk.AccAddress, value uint64, data string) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)
		data, err := hex.DecodeString(data)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgCall(caller.Address, contractAddr, value, data)

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
			fmt.Printf("<<<<<<<<< call error: %s\n", err)
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// call contract and check ret

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgDeployHello55 creates a massage operation of MsgDeploy with randomized field values.
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
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		// check pure/view function ret
		data, err := hex.DecodeString(Hello55SayHiData)
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

		fmt.Printf("<<<<<<<<< hello55 ret: %d\n", value)
		if value != 55 {
			panic("return value incorrect")
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgDeploySimple creates a massage operation of MsgDeploy with randomized field values.
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
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		fmt.Printf("<<<<<<<<< simple addr: %s\n", sdk.AccAddress(res.Data))

		futureOperations := []simulation.FutureOperation{
			{
				BlockHeight: int(ctx.BlockHeight()) + 1,
				Op:          SimulateMsgCall(k, res.Data, 0, SimpleSetData),
			},
		}

		return simulation.NewOperationMsg(msg, true, ""), futureOperations, nil
	}
}
