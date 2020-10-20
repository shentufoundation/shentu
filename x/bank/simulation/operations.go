package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosBank "github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/certikfoundation/shentu/x/auth/vesting"
	"github.com/certikfoundation/shentu/x/bank/internal/keeper"
	"github.com/certikfoundation/shentu/x/bank/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sim "github.com/cosmos/cosmos-sdk/x/bank/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

const (
	OpWeightMsgLockedSend      = "op_weight_msg_locked_send"
	DefaultWeightMsgLockedSend = 10
)

func WeightedOperations(appParams simulation.AppParams, cdc *codec.Codec, ak types.AccountKeeper, bk keeper.Keeper) simulation.WeightedOperations {
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
func SimulateMsgLockedSend(ak types.AccountKeeper, bk keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		for _, acc := range accs {
			account := ak.GetAccount(ctx, acc.Address)
			mvacc, ok := account.(*vesting.ManualVestingAccount)
			if !ok || mvacc.OriginalVesting.IsEqual(mvacc.VestedCoins) {
				continue
			}

			from, _ := simulation.RandomAcc(r, accs)
			fromAcc := ak.GetAccount(ctx, from.Address)
			spendableCoins := fromAcc.SpendableCoins(ctx.BlockTime())
			sendCoins := simulation.RandSubsetCoins(r, spendableCoins)

			spendableCoins = spendableCoins.Sub(sendCoins)
			if sendCoins.Empty() {
				return simulation.NoOpMsg(cosmosBank.ModuleName), nil, nil
			}

			fees, err := simulation.RandomFees(r, ctx, spendableCoins)
			if err != nil {
				return simulation.NoOpMsg(cosmosBank.ModuleName), nil, err
			}

			msg := types.NewMsgLockedSend(fromAcc.GetAddress(), mvacc.Address, sdk.AccAddress{}, sendCoins)

			tx := helpers.GenTx(
				[]sdk.Msg{msg},
				fees,
				helpers.DefaultGenTxGas,
				chainID,
				[]uint64{fromAcc.GetAccountNumber()},
				[]uint64{fromAcc.GetSequence()},
				from.PrivKey,
			)

			_, _, err = app.Deliver(tx)
			if err != nil {
				return simulation.NoOpMsg(cosmosBank.ModuleName), nil, err
			}

			return simulation.NewOperationMsg(msg, true, ""), nil, nil
		}
		return simulation.NewOperationMsgBasic(cosmosBank.ModuleName,
			"NoOp: no available manual-vesting account found, skip this tx", "", false, nil), nil, nil
	}
}
