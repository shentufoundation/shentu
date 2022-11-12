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

	"github.com/shentufoundation/shentu/v2/x/cert/keeper"
	"github.com/shentufoundation/shentu/v2/x/cert/types"
	simutil "github.com/shentufoundation/shentu/v2/x/cvm/simulation"
)

const (
	OpWeightMsgCertifyValidator = "op_weight_msg_certify_validator"
	OpWeightMsgCertifyPlatform  = "op_weight_msg_certify_platform"
	OpWeightMsgIssueCertificate = "op_weight_msg_issue_certificate"
)

// Default simulation operation weights for messages.
const (
	DefaultWeightMsgCertify int = 20
)

// WeightedOperations creates an operation (with weight) for each type of message generators.
func WeightedOperations(appParams simtypes.AppParams, cdc codec.JSONCodec, ak types.AccountKeeper,
	bk types.BankKeeper, k keeper.Keeper) simulation.WeightedOperations {
	var weightMsgCertifyValidator int
	appParams.GetOrGenerate(cdc, OpWeightMsgCertifyValidator, &weightMsgCertifyValidator, nil,
		func(_ *rand.Rand) {
			weightMsgCertifyValidator = simappparams.DefaultWeightMsgSend
		})

	var weightMsgCertifyPlatform int
	appParams.GetOrGenerate(cdc, OpWeightMsgCertifyPlatform, &weightMsgCertifyPlatform, nil,
		func(_ *rand.Rand) {
			weightMsgCertifyPlatform = simappparams.DefaultWeightMsgSend
		})

	var weightMsgIssueCertificate int
	appParams.GetOrGenerate(cdc, OpWeightMsgIssueCertificate, &weightMsgIssueCertificate, nil,
		func(_ *rand.Rand) {
			weightMsgIssueCertificate = simappparams.DefaultWeightMsgSend
		})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgCertifyPlatform, SimulateMsgCertifyPlatform(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgIssueCertificate, SimulateMsgIssueCertificates(ak, bk, k)),
	}
}

// SimulateMsgCertifyPlatform generates a MsgCertifyPlatform object which fields contain
// a randomly chosen existing certifier, a randomized validator's PubKey and a random string description.
func SimulateMsgCertifyPlatform(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
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
		platform := simtypes.RandStringOfLength(r, 10)

		msg, err := types.NewMsgCertifyPlatform(certifierAddr, validator.PubKey, platform)
		if err != nil {
			panic(err)
		}

		account := ak.GetAccount(ctx, certifierAddr)
		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, account.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCertifyPlatform, err.Error()), nil, err
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgIssueCertificates generates a MsgCertifyGeneral object which field values.
func SimulateMsgIssueCertificates(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account,
		chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
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

		certType := types.CertificateType_name[r.Int31n(7)+1]
		contentStr := simtypes.RandStringOfLength(r, 20)
		compiler := simtypes.RandStringOfLength(r, 5)
		bytecodeHash := simtypes.RandStringOfLength(r, 20)
		description := simtypes.RandStringOfLength(r, 10)

		content := types.AssembleContent(certType, contentStr)
		msg := types.NewMsgIssueCertificate(content, compiler, bytecodeHash, description, certifierAddr)

		account := ak.GetAccount(ctx, certifierAddr)
		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, account.GetAddress()))
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
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			certifierAcc.PrivKey,
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
