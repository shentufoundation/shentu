package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simTypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/x/auth"
	"github.com/certikfoundation/shentu/x/auth/types"
	"github.com/certikfoundation/shentu/x/auth/vesting"
)

const OpWeightMsgUnlock = "op_weight_msg_create_operator"

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(appParams simTypes.AppParams, cdc codec.JSONMarshaler, k auth.AccountKeeper) simulation.WeightedOperations {
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

func SimulateMsgUnlock(k auth.AccountKeeper) simTypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simTypes.Account, chainID string) (
		simTypes.OperationMsg, []simTypes.FutureOperation, error) {
		for _, acc := range accs {
			account := k.GetAccount(ctx, acc.Address)
			mvacc, ok := account.(*vesting.ManualVestingAccount)
			if !ok || mvacc.OriginalVesting.IsEqual(mvacc.VestedCoins) {
				continue
			}

			var unlockAmount sdk.Coins
			var err error
			if simTypes.RandIntBetween(r, 0, 100) < 50 {
				unlockAmount = mvacc.OriginalVesting.Sub(mvacc.VestedCoins)
			} else {
				unlockAmount, err = simTypes.RandomFees(r, ctx, mvacc.OriginalVesting.Sub(mvacc.VestedCoins))
				if err != nil {
					return simTypes.NoOpMsg(authTypes.ModuleName, types.TypeMsgUnlock, err.Error()), nil, err
				}
			}

			fees, err := simTypes.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
			if err != nil {
				return simTypes.NoOpMsg(authTypes.ModuleName, types.TypeMsgUnlock, err.Error()), nil, err
			}

			msg := types.NewMsgUnlock(acc.Address, acc.Address, unlockAmount)
			txGen := simappparams.MakeTestEncodingConfig().TxConfig
			tx, err := helpers.GenTx(
				txGen,
				[]sdk.Msg{msg},
				fees,
				helpers.DefaultGenTxGas,
				chainID,
				[]uint64{account.GetAccountNumber()},
				[]uint64{account.GetSequence()},
				acc.PrivKey,
			)
			if err != nil {
				return simTypes.NoOpMsg(authTypes.ModuleName, msg.Type(), err.Error()), nil, err
			}

			_, _, err = app.Deliver(txGen.TxEncoder(), tx)
			if err != nil {
				return simTypes.NoOpMsg(authTypes.ModuleName, msg.Type(), err.Error()), nil, err
			}

			return simTypes.NewOperationMsg(msg, true, ""), nil, nil
		}
		return simTypes.NewOperationMsgBasic(authTypes.ModuleName,
			"NoOp: no available manual-vesting account found, skip this tx", "", false, nil), nil, nil
	}
}
