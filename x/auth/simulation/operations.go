package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/x/auth/internal/types"
	"github.com/certikfoundation/shentu/x/auth/vesting"
)

const OpWeightMsgUnlock = "op_weight_msg_create_operator"

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(appParams simulation.AppParams, cdc *codec.Codec, k auth.AccountKeeper) simulation.WeightedOperations {
	var weightMsgUnlock int
	appParams.GetOrGenerate(cdc, OpWeightMsgUnlock, &weightMsgUnlock, nil,
		func(_ *rand.Rand) {
			weightMsgUnlock = 15
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgUnlock,
			SimulateMsgUnlock(k),
		),
	}
}

func SimulateMsgUnlock(k auth.AccountKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		for _, acc := range accs {
			account := k.GetAccount(ctx, acc.Address)
			mvacc, ok := account.(*vesting.ManualVestingAccount)
			if !ok || mvacc.OriginalVesting.IsEqual(mvacc.VestedCoins) {
				continue
			}

			var unlockAmount sdk.Coins
			var err error
			if simulation.RandIntBetween(r, 0, 100) < 50 {
				unlockAmount = mvacc.OriginalVesting.Sub(mvacc.VestedCoins)
			} else {
				unlockAmount, err = simulation.RandomFees(r, ctx, mvacc.OriginalVesting.Sub(mvacc.VestedCoins))
				if err != nil {
					return simulation.NoOpMsg(types.ModuleName), nil, err
				}
			}

			fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
			if err != nil {
				return simulation.NoOpMsg(auth.ModuleName), nil, err
			}

			msg := types.NewMsgUnlock(acc.Address, acc.Address, unlockAmount)

			tx := helpers.GenTx(
				[]sdk.Msg{msg},
				fees,
				helpers.DefaultGenTxGas,
				chainID,
				[]uint64{account.GetAccountNumber()},
				[]uint64{account.GetSequence()},
				acc.PrivKey,
			)

			_, _, err = app.Deliver(tx)
			if err != nil {
				return simulation.NoOpMsg(types.ModuleName), nil, err
			}

			return simulation.NewOperationMsg(msg, true, ""), nil, nil
		}
		return simulation.NewOperationMsgBasic(types.ModuleName,
			"NoOp: no available manual-vesting account found, skip this tx", "", false, nil), nil, nil
	}
}
