package simulation

import (
	"encoding/hex"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/hyperledger/burrow/txs/payload"

	"github.com/certikfoundation/shentu/x/cvm/internal/keeper"
	"github.com/certikfoundation/shentu/x/cvm/internal/types"
)

const (
	OpWeightMsgCall   = "op_weight_msg_call"
	OpWeightMsgDeploy = "op_weight_msg_deploy"
)

// WeightedOperations creates an operation with a weight for each type of message generators.
func WeightedOperations(appParams simulation.AppParams, cdc *codec.Codec, k keeper.Keeper) simulation.WeightedOperations {
	var weightMsgCall int
	appParams.GetOrGenerate(cdc, OpWeightMsgCall, &weightMsgCall, nil,
		func(_ *rand.Rand) {
			weightMsgCall = simappparams.DefaultWeightMsgSend
		})

	var weightMsgDeploy int
	appParams.GetOrGenerate(cdc, OpWeightMsgDeploy, &weightMsgDeploy, nil,
		func(_ *rand.Rand) {
			weightMsgDeploy = simappparams.DefaultWeightMsgSend
		})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgCall, SimulateMsgCall(k)),
		simulation.NewWeightedOperation(weightMsgDeploy, SimulateMsgDeploy(k)),
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

// SimulateMsgDeploy creates a massage operation of MsgDeploy with randomized field values.
func SimulateMsgDeploy(k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		caller, _ := simulation.RandomAcc(r, accs)
		value := uint64(0)
		code, _ := hex.DecodeString(keeper.Hello55BytecodeString)
		abi := keeper.Hello55AbiJsonString
		var meta []*payload.ContractMeta

		msg := types.NewMsgDeploy(caller.Address, value, code, abi, meta, false, false)

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
