package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	simutil "github.com/shentufoundation/shentu/v2/x/auth/simulation"
	"github.com/shentufoundation/shentu/v2/x/cert/keeper"
	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

const (
	OpWeightMsgIssueCertificate  = "op_weight_msg_issue_certificate"
	OpWeightMsgRevokeCertificate = "op_weight_msg_revoke_certificate"

	DefaultWeightMsgIssueCertificate  = 30
	DefaultWeightMsgRevokeCertificate = 15
)

// WeightedOperations returns all the operations from the module with their respective weights.
func WeightedOperations(
	appParams simtypes.AppParams,
	_ codec.JSONCodec,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simulation.WeightedOperations {
	var weightMsgIssue, weightMsgRevoke int

	appParams.GetOrGenerate(OpWeightMsgIssueCertificate, &weightMsgIssue, nil,
		func(_ *rand.Rand) { weightMsgIssue = DefaultWeightMsgIssueCertificate },
	)
	appParams.GetOrGenerate(OpWeightMsgRevokeCertificate, &weightMsgRevoke, nil,
		func(_ *rand.Rand) { weightMsgRevoke = DefaultWeightMsgRevokeCertificate },
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgIssue, SimulateMsgIssueCertificate(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgRevoke, SimulateMsgRevokeCertificate(ak, bk, k)),
	}
}

// SimulateMsgIssueCertificate generates a MsgIssueCertificate with random fields.
func SimulateMsgIssueCertificate(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msgType := sdk.MsgTypeURL(&types.MsgIssueCertificate{})

		certifiers := k.GetAllCertifiers(ctx)
		if len(certifiers) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, msgType, "no certifiers"), nil, nil
		}

		// Pick a random certifier and find the matching sim account.
		certifier := certifiers[r.Intn(len(certifiers))]
		certifierAddr, err := sdk.AccAddressFromBech32(certifier.Address)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msgType, "invalid certifier address"), nil, nil
		}

		var simAccount simtypes.Account
		var found bool
		for _, acc := range accs {
			if acc.Address.Equals(certifierAddr) {
				simAccount = acc
				found = true
				break
			}
		}
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, msgType, "certifier not in sim accounts"), nil, nil
		}

		// Build the certificate message.
		certType := randCertificateType(r)
		contentStr := simtypes.RandStringOfLength(r, 20)
		content := types.AssembleContent(certType, contentStr)
		description := simtypes.RandStringOfLength(r, 10)

		msg := types.NewMsgIssueCertificate(content, "", "", description, certifierAddr)

		account := ak.GetAccount(ctx, certifierAddr)
		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, account.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msgType, err.Error()), nil, err
		}

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig
		tx, err := simtestutil.GenSignedMockTx(
			r, txGen, []sdk.Msg{msg}, fees, simtestutil.DefaultGenTxGas, chainID,
			[]uint64{account.GetAccountNumber()}, []uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msgType, err.Error()), nil, err
		}

		_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msgType, err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgRevokeCertificate generates a MsgRevokeCertificate for a random existing certificate.
func SimulateMsgRevokeCertificate(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msgType := sdk.MsgTypeURL(&types.MsgRevokeCertificate{})

		certs := k.GetAllCertificates(ctx)
		if len(certs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, msgType, "no certificates to revoke"), nil, nil
		}

		// Pick a random certificate.
		cert := certs[r.Intn(len(certs))]

		// Use a current certifier as the revoker, not the original issuer.
		// The issuer may have been removed by a governance proposal, but any
		// active certifier is allowed to revoke.
		certifiers := k.GetAllCertifiers(ctx)
		if len(certifiers) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, msgType, "no certifiers"), nil, nil
		}
		revoker := certifiers[r.Intn(len(certifiers))]
		revokerAddr, err := sdk.AccAddressFromBech32(revoker.Address)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msgType, "invalid certifier address"), nil, nil
		}

		var simAccount simtypes.Account
		var found bool
		for _, acc := range accs {
			if acc.Address.Equals(revokerAddr) {
				simAccount = acc
				found = true
				break
			}
		}
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, msgType, "certifier not in sim accounts"), nil, nil
		}

		description := simtypes.RandStringOfLength(r, 10)
		msg := types.NewMsgRevokeCertificate(revokerAddr, cert.CertificateId, description)

		account := ak.GetAccount(ctx, revokerAddr)
		fees, err := simutil.RandomReasonableFees(r, ctx, bk.SpendableCoins(ctx, account.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msgType, err.Error()), nil, err
		}

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig
		tx, err := simtestutil.GenSignedMockTx(
			r, txGen, []sdk.Msg{msg}, fees, simtestutil.DefaultGenTxGas, chainID,
			[]uint64{account.GetAccountNumber()}, []uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msgType, err.Error()), nil, err
		}

		_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msgType, err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
