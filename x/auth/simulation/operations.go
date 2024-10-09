package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/shentufoundation/shentu/v2/x/auth/types"
)

const OpWeightMsgUnlock = "op_weight_msg_create_operator"

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(appParams simtypes.AppParams, cdc codec.JSONCodec, k types.AccountKeeper, bk types.BankKeeper) simulation.WeightedOperations {
	var weightMsgUnlock int
	appParams.GetOrGenerate(OpWeightMsgUnlock, &weightMsgUnlock, nil,
		func(_ *rand.Rand) {
			weightMsgUnlock = 15
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgUnlock,
			SimulateMsgUnlock(k, bk),
		),
	}
}

func SimulateMsgUnlock(k types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		for _, acc := range accs {
			account := k.GetAccount(ctx, acc.Address)
			mvacc, ok := account.(*types.ManualVestingAccount)
			if !ok || mvacc.OriginalVesting.Equal(mvacc.VestedCoins) {
				continue
			}

			var unlockAmount sdk.Coins
			var err error
			if simtypes.RandIntBetween(r, 0, 100) < 50 {
				unlockAmount = mvacc.OriginalVesting.Sub(mvacc.VestedCoins...)
			} else {
				unlockAmount, err = RandomReasonableFees(r, ctx, mvacc.OriginalVesting.Sub(mvacc.VestedCoins...))
				if err != nil {
					return simtypes.NoOpMsg(authtypes.ModuleName, types.TypeMsgUnlock, err.Error()), nil, err
				}
			}

			fees, err := RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, acc.Address))
			if err != nil {
				return simtypes.NoOpMsg(authtypes.ModuleName, types.TypeMsgUnlock, err.Error()), nil, err
			}

			msg := types.NewMsgUnlock(acc.Address, acc.Address, unlockAmount)
			txGen := moduletestutil.MakeTestEncodingConfig().TxConfig
			tx, err := simtestutil.GenSignedMockTx(
				r,
				txGen,
				[]sdk.Msg{msg},
				fees,
				simtestutil.DefaultGenTxGas,
				chainID,
				[]uint64{account.GetAccountNumber()},
				[]uint64{account.GetSequence()},
				acc.PrivKey,
			)
			if err != nil {
				return simtypes.NoOpMsg(authtypes.ModuleName, msg.Type(), err.Error()), nil, err
			}

			_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)
			if err != nil {
				return simtypes.NoOpMsg(authtypes.ModuleName, msg.Type(), err.Error()), nil, err
			}

			return simtypes.NewOperationMsg(msg, true, ""), nil, nil
		}
		return simtypes.NewOperationMsgBasic(authtypes.ModuleName,
			"NoOp: no available manual-vesting account found, skip this tx", "", false, nil), nil, nil
	}
}
