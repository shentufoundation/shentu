package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/x/cert/keeper"
	"github.com/certikfoundation/shentu/x/cert/types"
)

const (
	OpWeightMsgCertifyValidator = "op_weight_msg_certify_validator"
	OpWeightMsgCertifyPlatform  = "op_weight_msg_certify_platform"
)

// Default simulation operation weights for messages.
const (
	DefaultWeightMsgCertify int = 20
)

// WeightedOperations creates an operation (with weight) for each type of message generators.
func WeightedOperations(appParams simtypes.AppParams, cdc codec.JSONMarshaler, ak types.AccountKeeper,
	bk types.BankKeeper, k keeper.Keeper) simulation.WeightedOperations {
	var weightMsgCertifyValidator int
	appParams.GetOrGenerate(cdc, OpWeightMsgCertifyValidator, &weightMsgCertifyValidator, nil,
		func(_ *rand.Rand) {
			weightMsgCertifyValidator = simappparams.DefaultWeightMsgSend
		})
	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgCertifyValidator, SimulateMsgCertifyValidator(ak, bk, k)),
	}
}

// SimulateMsgCertifyValidator generates a MsgCertifyValidator object which fields contain
// a randomly chosen existing certifier and randomized validator's PubKey.
func SimulateMsgCertifyValidator(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		certifiers := k.GetAllCertifiers(ctx)
		certifier := certifiers[r.Intn(len(certifiers))]
		certifierAddr, err := sdk.AccAddressFromBech32(certifier.Address)
		if err != nil {
			panic(err)
		}
		var certifierAcc simtypes.Account
		for _, acc := range accs {
			if acc.Address.Equals(certifierAddr) {
				certifierAcc = acc
				break
			}
		}
		validator := simtypes.RandomAccounts(r, 1)[0]

		msg, err := types.NewMsgCertifyValidator(certifierAddr, validator.PubKey)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCertifyValidator, err.Error()), nil, err
		}

		account := ak.GetAccount(ctx, certifierAddr)
		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, account.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCertifyValidator, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			certifierAcc.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCertifyValidator, err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
