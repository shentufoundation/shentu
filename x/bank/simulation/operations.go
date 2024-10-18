package simulation

//
//import (
//	"math/rand"
//
//	"github.com/cosmos/cosmos-sdk/baseapp"
//	"github.com/cosmos/cosmos-sdk/codec"
//	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
//	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
//	banksim "github.com/cosmos/cosmos-sdk/x/bank/simulation"
//	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
//	"github.com/cosmos/cosmos-sdk/x/simulation"
//
//	simutil "github.com/shentufoundation/shentu/v2/x/auth/simulation"
//	vesting "github.com/shentufoundation/shentu/v2/x/auth/types"
//	"github.com/shentufoundation/shentu/v2/x/bank/keeper"
//	"github.com/shentufoundation/shentu/v2/x/bank/types"
//)
//
//const (
//	OpWeightMsgLockedSend      = "op_weight_msg_locked_send"
//	DefaultWeightMsgLockedSend = 10
//)
//
//func WeightedOperations(appParams simtypes.AppParams, cdc codec.JSONCodec, ak types.AccountKeeper, bk keeper.Keeper) simulation.WeightedOperations {
//	cosmosOps := banksim.WeightedOperations(appParams, cdc, ak, bk)
//
//	var weightMsgLockedSend int
//	appParams.GetOrGenerate(cdc, OpWeightMsgLockedSend, &weightMsgLockedSend, nil,
//		func(_ *rand.Rand) {
//			weightMsgLockedSend = DefaultWeightMsgLockedSend
//		},
//	)
//
//	op := simulation.NewWeightedOperation(weightMsgLockedSend, SimulateMsgLockedSend(ak, bk))
//	return append(cosmosOps, op)
//}
//
//// SimulateMsgLockedSend tests and runs a single msg send where both
//// accounts already exist.
//// nolint: funlen
//func SimulateMsgLockedSend(ak types.AccountKeeper, bk keeper.Keeper) simtypes.Operation {
//	return func(
//		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
//		accs []simtypes.Account, chainID string,
//	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
//		for _, acc := range accs {
//			account := ak.GetAccount(ctx, acc.Address)
//			mvacc, ok := account.(*vesting.ManualVestingAccount)
//			if !ok || mvacc.OriginalVesting.IsEqual(mvacc.VestedCoins) {
//				continue
//			}
//
//			from, _ := simtypes.RandomAcc(r, accs)
//			fromAcc := ak.GetAccount(ctx, from.Address)
//			spendableCoins := bk.SpendableCoins(ctx, fromAcc.GetAddress())
//			sendCoins := simtypes.RandSubsetCoins(r, spendableCoins)
//			if sendCoins.Empty() {
//				return simtypes.NoOpMsg(banktypes.ModuleName, types.TypeMsgLockedSend, "send coins empty"), nil, nil
//			}
//			spendableCoins = spendableCoins.Sub(sendCoins...)
//
//			fees, err := simutil.RandomReasonableFees(r, ctx, spendableCoins)
//			if err != nil {
//				return simtypes.NoOpMsg(banktypes.ModuleName, types.TypeMsgLockedSend, err.Error()), nil, err
//			}
//
//			toAddr, err := sdk.AccAddressFromBech32(mvacc.Address)
//			if err != nil {
//				return simtypes.NoOpMsg(banktypes.ModuleName, types.TypeMsgLockedSend, err.Error()), nil, err
//			}
//
//			msg := types.NewMsgLockedSend(fromAcc.GetAddress(), toAddr, "", sendCoins)
//
//			txGen := moduletestutil.MakeTestEncodingConfig().TxConfig
//			tx, err := simtestutil.GenSignedMockTx(
//				r,
//				txGen,
//				[]sdk.Msg{msg},
//				fees,
//				simtestutil.DefaultGenTxGas,
//				chainID,
//				[]uint64{fromAcc.GetAccountNumber()},
//				[]uint64{fromAcc.GetSequence()},
//				from.PrivKey,
//			)
//			if err != nil {
//				return simtypes.NoOpMsg(banktypes.ModuleName, msg.Type(), err.Error()), nil, err
//			}
//
//			_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)
//			if err != nil {
//				return simtypes.NoOpMsg(banktypes.ModuleName, msg.Type(), err.Error()), nil, err
//			}
//
//			return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
//		}
//		return simtypes.NewOperationMsgBasic(banktypes.ModuleName,
//			"NoOp: no available manual-vesting account found, skip this tx", "", false, nil), nil, nil
//	}
//}
