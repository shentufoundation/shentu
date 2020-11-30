package simulation

import (
	"math/rand"

	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simTypes "github.com/cosmos/cosmos-sdk/types/simulation"
	cosmosBank "github.com/cosmos/cosmos-sdk/x/bank"
	sim "github.com/cosmos/cosmos-sdk/x/bank/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/x/auth/vesting"
	"github.com/certikfoundation/shentu/x/bank/internal/keeper"
	"github.com/certikfoundation/shentu/x/bank/internal/types"
)

const (
	OpWeightMsgLockedSend      = "op_weight_msg_locked_send"
	DefaultWeightMsgLockedSend = 10
)

func WeightedOperations(appParams simTypes.AppParams, cdc codec.JSONMarshaler, ak types.AccountKeeper, bk keeper.Keeper) simulation.WeightedOperations {
	cosmosOps := sim.WeightedOperations(appParams, cdc, ak, bk)

	var weightMsgLockedSend int
	appParams.GetOrGenerate(cdc, OpWeightMsgLockedSend, &weightMsgLockedSend, nil,
		func(_ *rand.Rand) {
			weightMsgLockedSend = DefaultWeightMsgLockedSend
		},
	)

	op := simulation.NewWeightedOperation(weightMsgLockedSend, SimulateMsgLockedSend(ak, bk))
	return append(cosmosOps, op)
}

// SimulateMsgLockedSend tests and runs a single msg send where both
// accounts already exist.
// nolint: funlen
func SimulateMsgLockedSend(ak types.AccountKeeper, bk keeper.Keeper) simTypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simTypes.Account, chainID string,
	) (simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		for _, acc := range accs {
			account := ak.GetAccount(ctx, acc.Address)
			mvacc, ok := account.(*vesting.ManualVestingAccount)
			if !ok || mvacc.OriginalVesting.IsEqual(mvacc.VestedCoins) {
				continue
			}

			from, _ := simTypes.RandomAcc(r, accs)
			fromAcc := ak.GetAccount(ctx, from.Address)
			spendableCoins := fromAcc.SpendableCoins(ctx.BlockTime())
			sendCoins := simTypes.RandSubsetCoins(r, spendableCoins)
			if sendCoins.Empty() {
				return simTypes.NoOpMsg(cosmosBank.ModuleName), nil, nil
			}
			spendableCoins = spendableCoins.Sub(sendCoins)

			fees, err := simTypes.RandomFees(r, ctx, spendableCoins)
			if err != nil {
				return simTypes.NoOpMsg(cosmosBank.ModuleName), nil, err
			}

			msg := types.NewMsgLockedSend(fromAcc.GetAddress(), mvacc.Address, sdk.AccAddress{}, sendCoins)

			tx, err := helpers.GenTx(
				[]sdk.Msg{msg},
				fees,
				helpers.DefaultGenTxGas,
				chainID,
				[]uint64{fromAcc.GetAccountNumber()},
				[]uint64{fromAcc.GetSequence()},
				from.PrivKey,
			)
			if err != nil {
				return simTypes.NoOpMsg(cosmosBank.ModuleName), nil, err
			}

			txGen := simappparams.MakeTestEncodingConfig().TxConfig
			_, _, err = app.Deliver(txGen.TxEncoder(), tx)
			if err != nil {
				return simTypes.NoOpMsg(cosmosBank.ModuleName), nil, err
			}

			return simTypes.NewOperationMsg(msg, true, ""), nil, nil
		}
		return simTypes.NewOperationMsgBasic(cosmosBank.ModuleName,
			"NoOp: no available manual-vesting account found, skip this tx", "", false, nil), nil, nil
	}
}
