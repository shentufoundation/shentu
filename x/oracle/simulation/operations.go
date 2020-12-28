package simulation

import (
	"math/rand"
	"time"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/x/oracle/keeper"
	"github.com/certikfoundation/shentu/x/oracle/types"
)

const (
	OpWeightMsgCreateOperator = "op_weight_msg_create_operator"
	OpWeightMsgCreateTask     = "op_weight_msg_create_task"
)

// WeightedOperations returns all the operations from the module with their respective weights.
func WeightedOperations(appParams simtypes.AppParams, cdc codec.JSONMarshaler, k keeper.Keeper, ak types.AccountKeeper) simulation.WeightedOperations {
	var weightMsgCreateOperator int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreateOperator, &weightMsgCreateOperator, nil,
		func(_ *rand.Rand) {
			weightMsgCreateOperator = simappparams.DefaultWeightMsgSend
		},
	)

	var weightMsgCreateTask int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreateTask, &weightMsgCreateTask, nil,
		func(_ *rand.Rand) {
			weightMsgCreateTask = simappparams.DefaultWeightMsgSend
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateOperator,
			SimulateMsgCreateOperator(k, ak),
		),

		simulation.NewWeightedOperation(
			weightMsgCreateTask,
			SimulateMsgCreateTask(ak, k),
		),
	}
}

// SimulateMsgCreateOperator generates a MsgCreateOperator object with all of its fields randomized.
// This operation leads a series of future operations.
func SimulateMsgCreateOperator(k keeper.Keeper, ak types.AuthKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		operator, _ := simulation.RandomAcc(r, accs)

		if k.IsOperator(ctx, operator.Address) {
			return simulation.NewOperationMsgBasic(types.ModuleName,
				"NoOp: operator already exists, skip this tx", "", false, nil), nil, nil
		}

		operatorAcc := ak.GetAccount(ctx, operator.Address)
		collateral := simulation.RandSubsetCoins(r, operatorAcc.SpendableCoins(ctx.BlockTime()))
		if collateral.AmountOf(sdk.DefaultBondDenom).Int64() < k.GetLockedPoolParams(ctx).MinimumCollateral {
			return simulation.NewOperationMsgBasic(types.ModuleName,
				"NoOp: randomized collateral not enough, skip this tx", "", false, nil), nil, nil
		}

		fees, err := simulation.RandomFees(r, ctx, operatorAcc.SpendableCoins(ctx.BlockTime()).Sub(collateral))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgCreateOperator(operator.Address, collateral, operator.Address, "an operator")

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{operatorAcc.GetAccountNumber()},
			[]uint64{operatorAcc.GetSequence()},
			operator.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		stdOperator := types.NewOperator(operator.Address, operator.Address, collateral, nil, "an operator")
		futureOperations := []simulation.FutureOperation{
			{
				BlockHeight: int(ctx.BlockHeight()) + simulation.RandIntBetween(r, 0, 20),
				Op:          SimulateMsgAddCollateral(k, ak, &stdOperator, operator.PrivKey),
			},
			{
				BlockHeight: int(ctx.BlockHeight()) + simulation.RandIntBetween(r, 0, 20),
				Op:          SimulateMsgReduceCollateral(k, ak, &stdOperator, operator.PrivKey),
			},
			{
				BlockHeight: int(ctx.BlockHeight()) + simulation.RandIntBetween(r, 0, 20),
				Op:          SimulateMsgWithdrawReward(k, ak, &stdOperator, operator.PrivKey),
			},
			{
				BlockHeight: int(ctx.BlockHeight()) + simulation.RandIntBetween(r, 20, 25),
				Op:          SimulateMsgRemoveOperator(k, ak, &stdOperator, operator.PrivKey),
			},
		}

		return simulation.NewOperationMsg(msg, true, ""), futureOperations, nil
	}
}

// SimulateMsgAddCollateral generates a MsgAddCollateral object with all of its fields randomized.
func SimulateMsgAddCollateral(k keeper.Keeper, ak types.AuthKeeper, stdOperator *types.Operator,
	operatorPrivKey crypto.PrivKey) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		operator, err := k.GetOperator(ctx, stdOperator.Address)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		if err := checkConsistency(operator, *stdOperator); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		operatorAcc := ak.GetAccount(ctx, operator.Address)
		collateralIncrement := simulation.RandSubsetCoins(r, operatorAcc.SpendableCoins(ctx.BlockTime()))
		stdOperator.Collateral = stdOperator.Collateral.Add(collateralIncrement...)

		fees, err := simulation.RandomFees(r, ctx, operatorAcc.SpendableCoins(ctx.BlockTime()).Sub(collateralIncrement))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgAddCollateral(operator.Address, collateralIncrement)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{operatorAcc.GetAccountNumber()},
			[]uint64{operatorAcc.GetSequence()},
			operatorPrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgReduceCollateral generates a MsgReduceCollateral object with all of its fields randomized.
func SimulateMsgReduceCollateral(k keeper.Keeper, ak types.AuthKeeper, stdOperator *types.Operator,
	operatorPrivKey crypto.PrivKey) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		operator, err := k.GetOperator(ctx, stdOperator.Address)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		if err := checkConsistency(operator, *stdOperator); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		collateralDecrement := simulation.RandSubsetCoins(r, operator.Collateral)
		newCollateral := operator.Collateral.Sub(collateralDecrement)
		if newCollateral.AmountOf(sdk.DefaultBondDenom).Int64() < k.GetLockedPoolParams(ctx).MinimumCollateral {
			return simulation.NewOperationMsgBasic(types.ModuleName,
				"NoOp: randomized collateral not enough, skip this tx", "", false, nil), nil, nil
		}
		stdOperator.Collateral = newCollateral

		operatorAcc := ak.GetAccount(ctx, operator.Address)
		fees, err := simulation.RandomFees(r, ctx, operatorAcc.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgReduceCollateral(operator.Address, collateralDecrement)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{operatorAcc.GetAccountNumber()},
			[]uint64{operatorAcc.GetSequence()},
			operatorPrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgRemoveOperator generates a MsgRemoveOperator object with all of its fields randomized.
func SimulateMsgRemoveOperator(k keeper.Keeper, ak types.AuthKeeper, stdOperator *types.Operator,
	operatorPrivKey crypto.PrivKey) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		operator, err := k.GetOperator(ctx, stdOperator.Address)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		if err := checkConsistency(operator, *stdOperator); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		operatorAcc := ak.GetAccount(ctx, operator.Address)
		fees, err := simulation.RandomFees(r, ctx, operatorAcc.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgRemoveOperator(operator.Address, operator.Address)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{operatorAcc.GetAccountNumber()},
			[]uint64{operatorAcc.GetSequence()},
			operatorPrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgWithdrawReward generates a MsgWithdrawReward object with all of its fields randomized.
func SimulateMsgWithdrawReward(k keeper.Keeper, ak types.AuthKeeper, stdOperator *types.Operator,
	operatorPrivKey crypto.PrivKey) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		operator, err := k.GetOperator(ctx, stdOperator.Address)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		if err := checkConsistency(operator, *stdOperator); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgWithdrawReward(operator.Address)

		operatorAcc := ak.GetAccount(ctx, operator.Address)
		fees, err := simulation.RandomFees(r, ctx, operatorAcc.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{operatorAcc.GetAccountNumber()},
			[]uint64{operatorAcc.GetSequence()},
			operatorPrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func checkConsistency(operator1, operator2 types.Operator) error {
	if !operator1.Address.Equals(operator2.Address) ||
		!operator1.Proposer.Equals(operator2.Proposer) ||
		!operator1.Collateral.IsEqual(operator2.Collateral) ||
		operator1.Name != operator2.Name {
		return types.ErrInconsistentOperators
	}
	return nil
}

// SimulateMsgCreateTask generates a MsgCreateTask object with all of its fields randomized.
// This operation leads a series of future operations.
func SimulateMsgCreateTask(ak types.AuthKeeper, k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		contract := simulation.RandStringOfLength(r, 10)
		function := simulation.RandStringOfLength(r, 10)
		description := simulation.RandStringOfLength(r, 20)
		creator, _ := simulation.RandomAcc(r, accs)
		creatorAcc := ak.GetAccount(ctx, creator.Address)
		bounty := simulation.RandSubsetCoins(r, creatorAcc.SpendableCoins(ctx.BlockTime()))
		wait := simulation.RandIntBetween(r, 5, 20)

		msg := types.NewMsgCreateTask(contract, function, bounty, description, creator.Address, int64(wait), time.Duration(0))

		fees, err := simulation.RandomFees(r, ctx, creatorAcc.SpendableCoins(ctx.BlockTime()).Sub(bounty))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{creatorAcc.GetAccountNumber()},
			[]uint64{creatorAcc.GetSequence()},
			creator.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		futureOperations := []simulation.FutureOperation{
			{
				BlockHeight: int(ctx.BlockHeight()) + simulation.RandIntBetween(r, 0, 20),
				Op:          SimulateMsgInquiryTask(ak, contract, function),
			},
			{
				BlockHeight: int(ctx.BlockHeight()) + simulation.RandIntBetween(r, 20, 25),
				Op:          SimulateMsgDeleteTask(ak, contract, function, creator),
			},
		}

		for _, acc := range accs {
			if k.IsOperator(ctx, acc.Address) && simulation.RandIntBetween(r, 0, 100) < 10 {
				futureOperations = append(futureOperations, simulation.FutureOperation{
					BlockHeight: int(ctx.BlockHeight()) + simulation.RandIntBetween(r, 0, wait),
					Op:          SimulateMsgTaskResponse(ak, k, contract, function, acc),
				})
			}
		}

		return simulation.NewOperationMsg(msg, true, ""), futureOperations, nil
	}
}

// SimulateMsgInquiryTask generates a MsgInquiryTask object with all of its fields randomized.
func SimulateMsgInquiryTask(ak types.AuthKeeper, contract, function string) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		txHash := simulation.RandStringOfLength(r, 20)
		inquirer, _ := simulation.RandomAcc(r, accs)

		msg := types.NewMsgInquiryTask(contract, function, txHash, inquirer.Address)

		inquirerAcc := ak.GetAccount(ctx, inquirer.Address)
		fees, err := simulation.RandomFees(r, ctx, inquirerAcc.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{inquirerAcc.GetAccountNumber()},
			[]uint64{inquirerAcc.GetSequence()},
			inquirer.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgTaskResponse generates a MsgTaskResponse object with all of its fields randomized.
func SimulateMsgTaskResponse(ak types.AuthKeeper, k keeper.Keeper, contract, function string,
	simAcc simulation.Account) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		if !k.IsOperator(ctx, simAcc.Address) {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		score := r.Int63n(100) + 1

		msg := types.NewMsgTaskResponse(contract, function, score, simAcc.Address)

		operatorAcc := ak.GetAccount(ctx, simAcc.Address)
		fees, err := simulation.RandomFees(r, ctx, operatorAcc.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{operatorAcc.GetAccountNumber()},
			[]uint64{operatorAcc.GetSequence()},
			simAcc.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgDeleteTask generates a MsgDeleteTask object with all of its fields randomized.
func SimulateMsgDeleteTask(ak types.AuthKeeper, contract, function string, creator simulation.Account) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		msg := types.NewMsgDeleteTask(contract, function, true, creator.Address)

		creatorAcc := ak.GetAccount(ctx, creator.Address)
		fees, err := simulation.RandomFees(r, ctx, creatorAcc.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{creatorAcc.GetAccountNumber()},
			[]uint64{creatorAcc.GetSequence()},
			creator.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
