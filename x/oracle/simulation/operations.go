package simulation

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	simutil "github.com/shentufoundation/shentu/v2/x/cvm/simulation"
	"github.com/shentufoundation/shentu/v2/x/oracle/keeper"
	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

const (
	OpWeightMsgCreateOperator = "op_weight_msg_create_operator"
	OpWeightMsgCreateTask     = "op_weight_msg_create_task"
	OpWeightMsgCreateTxTask   = "op_weight_msg_create_tx_task"
)

// WeightedOperations returns all the operations from the module with their respective weights.
func WeightedOperations(appParams simtypes.AppParams, cdc codec.JSONCodec, k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simulation.WeightedOperations {
	var (
		weightMsgCreateOperator int
		weightMsgCreateTask     int
		weightMsgCreateTxTask   int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgCreateOperator, &weightMsgCreateOperator, nil,
		func(_ *rand.Rand) {
			weightMsgCreateOperator = simappparams.DefaultWeightMsgSend
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgCreateTask, &weightMsgCreateTask, nil,
		func(_ *rand.Rand) {
			weightMsgCreateTask = simappparams.DefaultWeightMsgSend
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgCreateTxTask, &weightMsgCreateTxTask, nil,
		func(_ *rand.Rand) {
			weightMsgCreateTxTask = simappparams.DefaultWeightMsgSend
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgCreateOperator, SimulateMsgCreateOperator(k, ak, bk)),
		simulation.NewWeightedOperation(weightMsgCreateTask, SimulateMsgCreateTask(ak, k, bk)),
		simulation.NewWeightedOperation(weightMsgCreateTxTask, SimulateMsgCreateTxTask(ak, k, bk)),
	}
}

// SimulateMsgCreateOperator generates a MsgCreateOperator object with all of its fields randomized.
// This operation leads a series of future operations.
func SimulateMsgCreateOperator(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		operator, _ := simtypes.RandomAcc(r, accs)

		if k.IsOperator(ctx, operator.Address) {
			return simtypes.NewOperationMsgBasic(types.ModuleName, "NoOp: operator already exists, skip this tx", "", false, nil), nil, nil
		}

		operatorAcc := ak.GetAccount(ctx, operator.Address)
		collateral := simtypes.RandSubsetCoins(r, bk.SpendableCoins(ctx, operatorAcc.GetAddress()))
		if collateral.Empty() {
			return simtypes.NewOperationMsgBasic(types.ModuleName, "NoOp: empty collateral, skip this tx", "", false, nil), nil, nil
		}
		if collateral.AmountOf(sdk.DefaultBondDenom).Int64() < k.GetLockedPoolParams(ctx).MinimumCollateral {
			return simtypes.NewOperationMsgBasic(types.ModuleName,
				"NoOp: randomized collateral not enough, skip this tx", "", false, nil), nil, nil
		}

		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, operatorAcc.GetAddress()).Sub(collateral))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateOperator, err.Error()), nil, err
		}

		msg := types.NewMsgCreateOperator(operator.Address, collateral, operator.Address, "an operator")

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{operatorAcc.GetAccountNumber()},
			[]uint64{operatorAcc.GetSequence()},
			operator.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		stdOperator := types.NewOperator(operator.Address, operator.Address, collateral, nil, "an operator")
		futureOperations := []simtypes.FutureOperation{
			{
				BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 0, 20),
				Op:          SimulateMsgAddCollateral(k, ak, bk, &stdOperator, operator.PrivKey),
			},
			{
				BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 0, 20),
				Op:          SimulateMsgReduceCollateral(k, ak, bk, &stdOperator, operator.PrivKey),
			},
			{
				BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 0, 20),
				Op:          SimulateMsgWithdrawReward(k, ak, bk, &stdOperator, operator.PrivKey),
			},
			{
				BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 20, 25),
				Op:          SimulateMsgRemoveOperator(k, ak, bk, &stdOperator, operator.PrivKey),
			},
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), futureOperations, nil
	}
}

// SimulateMsgAddCollateral generates a MsgAddCollateral object with all of its fields randomized.
func SimulateMsgAddCollateral(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, stdOperator *types.Operator,
	operatorPrivKey cryptotypes.PrivKey) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		stdOperatorAddr, err := sdk.AccAddressFromBech32(stdOperator.Address)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRemoveOperator, err.Error()), nil, err
		}
		operator, err := k.GetOperator(ctx, stdOperatorAddr)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAddCollateral, err.Error()), nil, err
		}

		if err := checkConsistency(operator, *stdOperator); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAddCollateral, err.Error()), nil, err
		}

		operatorAddr, err := sdk.AccAddressFromBech32(operator.Address)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRemoveOperator, err.Error()), nil, err
		}
		operatorAcc := ak.GetAccount(ctx, operatorAddr)
		collateralIncrement := simtypes.RandSubsetCoins(r, bk.SpendableCoins(ctx, operatorAcc.GetAddress()))
		if collateralIncrement.Empty() {
			return simtypes.NewOperationMsgBasic(types.ModuleName, "NoOp: empty collateral increment, skip this tx", "", false, nil), nil, nil
		}
		stdOperator.Collateral = stdOperator.Collateral.Add(collateralIncrement...)

		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, operatorAcc.GetAddress()).Sub(collateralIncrement))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAddCollateral, err.Error()), nil, err
		}

		msg := types.NewMsgAddCollateral(operatorAddr, collateralIncrement)

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{operatorAcc.GetAccountNumber()},
			[]uint64{operatorAcc.GetSequence()},
			operatorPrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgReduceCollateral generates a MsgReduceCollateral object with all of its fields randomized.
func SimulateMsgReduceCollateral(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, stdOperator *types.Operator,
	operatorPrivKey cryptotypes.PrivKey) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		stdOperatorAddr, err := sdk.AccAddressFromBech32(stdOperator.Address)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRemoveOperator, err.Error()), nil, err
		}
		operator, err := k.GetOperator(ctx, stdOperatorAddr)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgReduceCollateral, err.Error()), nil, err
		}

		if err := checkConsistency(operator, *stdOperator); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgReduceCollateral, err.Error()), nil, err
		}

		collateralDecrement := simtypes.RandSubsetCoins(r, operator.Collateral)
		if collateralDecrement.Empty() {
			return simtypes.NewOperationMsgBasic(types.ModuleName, "NoOp: empty collateral increment, skip this tx", "", false, nil), nil, nil
		}
		newCollateral := operator.Collateral.Sub(collateralDecrement)
		if newCollateral.AmountOf(sdk.DefaultBondDenom).Int64() < k.GetLockedPoolParams(ctx).MinimumCollateral {
			return simtypes.NewOperationMsgBasic(types.ModuleName,
				"NoOp: randomized collateral not enough, skip this tx", "", false, nil), nil, nil
		}
		stdOperator.Collateral = newCollateral

		operatorAddr, err := sdk.AccAddressFromBech32(operator.Address)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRemoveOperator, err.Error()), nil, err
		}
		operatorAcc := ak.GetAccount(ctx, operatorAddr)
		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, operatorAcc.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgReduceCollateral, err.Error()), nil, err
		}

		msg := types.NewMsgReduceCollateral(operatorAddr, collateralDecrement)

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{operatorAcc.GetAccountNumber()},
			[]uint64{operatorAcc.GetSequence()},
			operatorPrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgRemoveOperator generates a MsgRemoveOperator object with all of its fields randomized.
func SimulateMsgRemoveOperator(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper,
	stdOperator *types.Operator, operatorPrivKey cryptotypes.PrivKey) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		stdOperatorAddr, err := sdk.AccAddressFromBech32(stdOperator.Address)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRemoveOperator, err.Error()), nil, err
		}
		operator, err := k.GetOperator(ctx, stdOperatorAddr)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRemoveOperator, err.Error()), nil, err
		}

		if err := checkConsistency(operator, *stdOperator); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRemoveOperator, err.Error()), nil, err
		}

		operatorAddr, err := sdk.AccAddressFromBech32(operator.Address)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRemoveOperator, err.Error()), nil, err
		}
		operatorAcc := ak.GetAccount(ctx, operatorAddr)
		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, operatorAcc.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRemoveOperator, err.Error()), nil, err
		}

		msg := types.NewMsgRemoveOperator(operatorAddr, operatorAddr)

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{operatorAcc.GetAccountNumber()},
			[]uint64{operatorAcc.GetSequence()},
			operatorPrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgWithdrawReward generates a MsgWithdrawReward object with all of its fields randomized.
func SimulateMsgWithdrawReward(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper,
	stdOperator *types.Operator, operatorPrivKey cryptotypes.PrivKey) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		stdOperatorAddr, err := sdk.AccAddressFromBech32(stdOperator.Address)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRemoveOperator, err.Error()), nil, err
		}
		operator, err := k.GetOperator(ctx, stdOperatorAddr)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawReward, err.Error()), nil, err
		}

		if err := checkConsistency(operator, *stdOperator); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawReward, err.Error()), nil, err
		}

		operatorAddr, err := sdk.AccAddressFromBech32(operator.Address)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRemoveOperator, err.Error()), nil, err
		}
		msg := types.NewMsgWithdrawReward(operatorAddr)

		operatorAcc := ak.GetAccount(ctx, operatorAddr)
		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, operatorAcc.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{operatorAcc.GetAccountNumber()},
			[]uint64{operatorAcc.GetSequence()},
			operatorPrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

func checkConsistency(operator1, operator2 types.Operator) error {
	if operator1.Address != operator2.Address || operator1.Proposer != operator2.Proposer ||
		!operator1.Collateral.IsEqual(operator2.Collateral) || operator1.Name != operator2.Name {
		return types.ErrInconsistentOperators
	}
	return nil
}

// SimulateMsgCreateTask generates a MsgCreateTask object with all of its fields randomized.
// This operation leads a series of future operations.
func SimulateMsgCreateTask(ak types.AccountKeeper, k keeper.Keeper, bk types.BankKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		contract := simtypes.RandStringOfLength(r, 10)
		function := simtypes.RandStringOfLength(r, 10)
		description := simtypes.RandStringOfLength(r, 20)
		creator, _ := simtypes.RandomAcc(r, accs)
		creatorAcc := ak.GetAccount(ctx, creator.Address)
		bounty := simtypes.RandSubsetCoins(r, bk.SpendableCoins(ctx, creatorAcc.GetAddress()))
		wait := simtypes.RandIntBetween(r, 5, 20)

		msg := types.NewMsgCreateTask(contract, function, bounty, description, creator.Address, int64(wait), time.Duration(0))

		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, creatorAcc.GetAddress()).Sub(bounty))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{creatorAcc.GetAccountNumber()},
			[]uint64{creatorAcc.GetSequence()},
			creator.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		futureOperations := []simtypes.FutureOperation{
			{
				BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 20, 25),
				Op:          SimulateMsgDeleteTask(ak, bk, contract, function, creator),
			},
		}

		for _, acc := range accs {
			if k.IsOperator(ctx, acc.Address) && simtypes.RandIntBetween(r, 0, 100) < 10 {
				futureOperations = append(futureOperations, simtypes.FutureOperation{
					BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 0, wait),
					Op:          SimulateMsgTaskResponse(ak, k, bk, contract, function, acc),
				})
			}
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), futureOperations, nil
	}
}

// SimulateMsgTaskResponse generates a MsgTaskResponse object with all of its fields randomized.
func SimulateMsgTaskResponse(ak types.AccountKeeper, k keeper.Keeper, bk types.BankKeeper, contract, function string,
	simAcc simtypes.Account) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if !k.IsOperator(ctx, simAcc.Address) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRespondToTask, "not an operator"), nil, nil
		}

		score := r.Int63n(100) + 1

		msg := types.NewMsgTaskResponse(contract, function, score, simAcc.Address)

		operatorAcc := ak.GetAccount(ctx, simAcc.Address)
		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, operatorAcc.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRespondToTask, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{operatorAcc.GetAccountNumber()},
			[]uint64{operatorAcc.GetSequence()},
			simAcc.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgDeleteTask generates a MsgDeleteTask object with all of its fields randomized.
func SimulateMsgDeleteTask(ak types.AccountKeeper, bk types.BankKeeper, contract, function string, creator simtypes.Account) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := types.NewMsgDeleteTask(contract, function, true, creator.Address)

		creatorAcc := ak.GetAccount(ctx, creator.Address)
		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, creatorAcc.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDeleteTask, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{creatorAcc.GetAccountNumber()},
			[]uint64{creatorAcc.GetSequence()},
			creator.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgCreateTxTask generates a MsgCreateTxTask object with all of its fields randomized.
// This operation leads a series of future operations.
func SimulateMsgCreateTxTask(ak types.AccountKeeper, k keeper.Keeper, bk types.BankKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		creator, _ := simtypes.RandomAcc(r, accs)
		creatorAcc := ak.GetAccount(ctx, creator.Address)
		chainId := fmt.Sprintf("%d", simtypes.RandIntBetween(r, 1, 2000))
		bounty := simtypes.RandSubsetCoins(r, bk.SpendableCoins(ctx, creatorAcc.GetAddress()))
		validTime := ctx.BlockTime().Add(10 * time.Second).UTC()
		// mock business
		businessTx := []byte(simtypes.RandStringOfLength(r, 500))
		businessTxHash := sha256.Sum256(businessTx)

		blockWait := simtypes.RandIntBetween(r, 1, 20)

		msg := types.NewMsgCreateTxTask(creator.Address, chainId, businessTx, bounty, validTime)

		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, creatorAcc.GetAddress()).Sub(bounty))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{creatorAcc.GetAccountNumber()},
			[]uint64{creatorAcc.GetSequence()},
			creator.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		futureOperations := []simtypes.FutureOperation{
			{
				BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 20, 25),
				Op:          SimulateMsgDeleteTxTask(ak, bk, businessTxHash[:], creator),
			},
		}

		for _, acc := range accs {
			if k.IsOperator(ctx, acc.Address) && simtypes.RandIntBetween(r, 0, 100) < 10 {
				futureOperations = append(futureOperations, simtypes.FutureOperation{
					BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 0, blockWait),
					Op:          SimulateMsgTxTaskResponse(ak, k, bk, businessTxHash[:], acc),
				})
			}
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgTxTaskResponse generates a MsgTxTaskResponse object with all of its fields randomized.
func SimulateMsgTxTaskResponse(ak types.AccountKeeper, k keeper.Keeper, bk types.BankKeeper, txHash []byte,
	simAcc simtypes.Account) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if !k.IsOperator(ctx, simAcc.Address) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRespondToTxTask, "not an operator"), nil, nil
		}
		score := r.Int63n(100)

		msg := types.NewMsgTxTaskResponse(txHash, score, simAcc.Address)
		operatorAcc := ak.GetAccount(ctx, simAcc.Address)
		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, operatorAcc.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRespondToTask, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{operatorAcc.GetAccountNumber()},
			[]uint64{operatorAcc.GetSequence()},
			simAcc.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgDeleteTxTask generates a MsgDeleteTxTask object with all of its fields randomized.
func SimulateMsgDeleteTxTask(ak types.AccountKeeper, bk types.BankKeeper, txHash []byte, creator simtypes.Account) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := types.NewMsgDeleteTxTask(txHash, creator.Address)

		creatorAcc := ak.GetAccount(ctx, creator.Address)
		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, creatorAcc.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDeleteTxTask, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{creatorAcc.GetAccountNumber()},
			[]uint64{creatorAcc.GetSequence()},
			creator.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}

}
