package simulation

import (
	"github.com/certikfoundation/shentu/common"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

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
			weightMsgCreateOperator = simappparams.DefaultWeightMsgMultiSend
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
func SimulateMsgCreateOperator(k keeper.Keeper, ak types.AuthKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		operator, _ := simulation.RandomAcc(r, accs)
		proposer, _ := simulation.RandomAcc(r, accs)

		operatorAcc := ak.GetAccount(ctx, operator.Address)
		collateral, err := simulation.RandomFees(r, ctx, operatorAcc.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		if collateral.AmountOf(common.MicroCTKDenom).Int64() < k.GetLockedPoolParams(ctx).MinimumCollateral {
			return simulation.NewOperationMsgBasic(types.ModuleName, "NoOp: randomized collateral not enough, skip this tx", "", false, nil), nil, nil
		}

		proposerAcc := ak.GetAccount(ctx, proposer.Address)
		fees, err := simulation.RandomFees(r, ctx, proposerAcc.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgCreateOperator(operator.Address, collateral, proposer.Address, "an operator")

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{proposerAcc.GetAccountNumber()},
			[]uint64{proposerAcc.GetSequence()},
			proposer.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		futureOperations := []simulation.FutureOperation{
			{
				BlockHeight: int(ctx.BlockHeight()) + 3,
				Op:          SimulateMsgAddCollateral(k, ak, operator),
			},
		}

		return simulation.NewOperationMsg(msg, true, ""), futureOperations, nil
	}
}

// SimulateMsgAddCollateral generates a MsgAddCollateral object with all of its fields randomized.
func SimulateMsgAddCollateral(k keeper.Keeper, ak types.AuthKeeper, operator simulation.Account) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		_, err := k.GetOperator(ctx, operator.Address)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		operatorAcc := ak.GetAccount(ctx, operator.Address)
		collateralIncrement, err := simulation.RandomFees(r, ctx, operatorAcc.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

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
			operator.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
