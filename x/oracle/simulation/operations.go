package simulation

import (
	"math/rand"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/oracle/internal/keeper"
	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

const (
	OpWeightMsgCreateOperator = "op_weight_msg_create_operator"
)

// WeightedOperations returns all the operations from the module with their respective weights.
func WeightedOperations(
	appParams simulation.AppParams, cdc *codec.Codec, k keeper.Keeper, ak types.AuthKeeper) simulation.WeightedOperations {
	var weightMsgCreateOperator int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreateOperator, &weightMsgCreateOperator, nil,
		func(_ *rand.Rand) {
			weightMsgCreateOperator = simappparams.DefaultWeightMsgSend
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateOperator,
			SimulateMsgCreateOperator(k, ak),
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
		collateral, err := simulation.RandomFees(r, ctx, operatorAcc.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		if collateral.AmountOf(common.MicroCTKDenom).Int64() < k.GetLockedPoolParams(ctx).MinimumCollateral {
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
				BlockHeight: int(ctx.BlockHeight()) + 1,
				Op:          SimulateMsgAddCollateral(k, ak, &stdOperator, operator.PrivKey),
			},
			{
				BlockHeight: int(ctx.BlockHeight()) + 2,
				Op:          SimulateMsgReduceCollateral(k, ak, &stdOperator, operator.PrivKey),
			},
			{
				BlockHeight: int(ctx.BlockHeight()) + 3,
				Op:          SimulateMsgRemoveOperator(k, ak, &stdOperator, operator.PrivKey),
			},
		}

		return simulation.NewOperationMsg(msg, true, ""), futureOperations, nil
	}
}

// SimulateMsgAddCollateral generates a MsgAddCollateral object with all of its fields randomized.
func SimulateMsgAddCollateral(k keeper.Keeper, ak types.AuthKeeper, stdOperator *types.Operator, operatorPrivKey crypto.PrivKey) simulation.Operation {
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
		collateralIncrement, err := simulation.RandomFees(r, ctx, operatorAcc.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
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
func SimulateMsgReduceCollateral(k keeper.Keeper, ak types.AuthKeeper, stdOperator *types.Operator, operatorPrivKey crypto.PrivKey) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		operator, err := k.GetOperator(ctx, stdOperator.Address)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		if err := checkConsistency(operator, *stdOperator); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		collateralDecrement, err := simulation.RandomFees(r, ctx, operator.Collateral)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		newCollateral := operator.Collateral.Sub(collateralDecrement)
		if newCollateral.AmountOf(common.MicroCTKDenom).Int64() < k.GetLockedPoolParams(ctx).MinimumCollateral {
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
func SimulateMsgRemoveOperator(k keeper.Keeper, ak types.AuthKeeper, stdOperator *types.Operator, operatorPrivKey crypto.PrivKey) simulation.Operation {
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

func checkConsistency(operator1, operator2 types.Operator) error {
	if !operator1.Address.Equals(operator2.Address) ||
		!operator1.Proposer.Equals(operator2.Proposer) ||
		!operator1.Collateral.IsEqual(operator2.Collateral) ||
		!operator1.AccumulatedRewards.IsEqual(operator2.AccumulatedRewards) ||
		operator1.Name != operator2.Name {
		return types.ErrInconsistentOperators
	}
	return nil
}
