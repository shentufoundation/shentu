package simulation

import (
	"encoding/hex"
	"fmt"
	"math/rand"

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
	}
}

// SimulateMsgCall creates a message operation of MsgCall with randomized field values.
func SimulateMsgCall(k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)
		callee, _ := simulation.RandomAcc(r, accs)
		value := uint64(0)
		var data []byte

		msg := types.NewMsgCall(caller.Address, callee.Address, value, data)

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

		contractAddr, err := sdk.AccAddressFromBech32(string(res.Events.ToABCIEvents()[2].Attributes[0].Value))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		fmt.Printf("<<<<<<<<<<<<<< Contract addr: %s\n", contractAddr)

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
