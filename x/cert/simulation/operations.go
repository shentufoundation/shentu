package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/x/cert/internal/keeper"
	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

const (
	OpWeightMsgProposeCertifier = "op_weight_msg_propose_certifier"
	OpWeightMsgCertifyValidator = "op_weight_msg_certify_validator"
	OpWeightMsgCertifyPlatform  = "op_weight_msg_certify_platform"
	OpWeightMsgCertifyAuditing  = "op_weight_msg_certify_auditing"
	OpWeightMsgCertifyProof     = "op_weight_msg_certify_proof"
)

// WeightedOperations creates an operation (with weight) for each type of message generators.
func WeightedOperations(appParams simulation.AppParams, cdc *codec.Codec, ak types.AccountKeeper,
	k keeper.Keeper) simulation.WeightedOperations {
	var weightMsgProposeCertifier int
	appParams.GetOrGenerate(cdc, OpWeightMsgProposeCertifier, &weightMsgProposeCertifier, nil,
		func(_ *rand.Rand) {
			weightMsgProposeCertifier = simappparams.DefaultWeightMsgSend
		})

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

	var weightMsgCertifyAuditing int
	appParams.GetOrGenerate(cdc, OpWeightMsgCertifyAuditing, &weightMsgCertifyAuditing, nil,
		func(_ *rand.Rand) {
			weightMsgCertifyAuditing = simappparams.DefaultWeightMsgSend
		})

	var weightMsgCertifyProof int
	appParams.GetOrGenerate(cdc, OpWeightMsgCertifyProof, &weightMsgCertifyProof, nil,
		func(_ *rand.Rand) {
			weightMsgCertifyProof = simappparams.DefaultWeightMsgSend
		})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgProposeCertifier, SimulateMsgProposeCertifier(k, ak)),
		simulation.NewWeightedOperation(weightMsgCertifyValidator, SimulateMsgCertifyValidator(ak, k)),
		simulation.NewWeightedOperation(weightMsgCertifyPlatform, SimulateMsgCertifyPlatform(ak, k)),
		simulation.NewWeightedOperation(weightMsgCertifyAuditing, SimulateMsgCertifyAuditing(ak, k)),
		simulation.NewWeightedOperation(weightMsgCertifyProof, SimulateMsgCertifyProof(ak, k)),
	}
}

// SimulateMsgProposeCertifier generates a MsgProposeCertifier object with all of its fields randomized.
func SimulateMsgProposeCertifier(k keeper.Keeper, ak types.AccountKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		proposer, _ := simulation.RandomAcc(r, accs)
		certifier := simulation.RandomAccounts(r, 1)[0]
		description := simulation.RandStringOfLength(r, 10)
		alias := simulation.RandStringOfLength(r, 5)

		msg := types.NewMsgProposeCertifier(proposer.Address, certifier.Address, alias, description)

		if len(k.GetAllCertifiers(ctx)) > 0 {
			return simulation.NewOperationMsgBasic(types.ModuleName, "NoOp: certifier already exists", "", false, nil), nil, nil
		}

		account := ak.GetAccount(ctx, proposer.Address)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			proposer.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgCertifyValidator generates a MsgCertifyValidator object which fields contain
// a randomly chosen existing certifier and randomized validator's PubKey.
func SimulateMsgCertifyValidator(ak types.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		certifier, _ := simulation.RandomAcc(r, accs)
		validator := simulation.RandomAccounts(r, 1)[0]

		msg := types.NewMsgCertifyValidator(certifier.Address, validator.PubKey)

		if !k.IsCertifier(ctx, certifier.Address) {
			return simulation.NewOperationMsgBasic(types.ModuleName, "NoOp: not a certifier", "", false, nil), nil, nil
		}

		account := ak.GetAccount(ctx, certifier.Address)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			certifier.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgCertifyPlatform generates a MsgCertifyPlatform object which fields contain
// a randomly chosen existing certifier, a randomized validator's PubKey and a random string description.
func SimulateMsgCertifyPlatform(ak types.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		certifier, _ := simulation.RandomAcc(r, accs)
		validator := simulation.RandomAccounts(r, 1)[0]
		platform := simulation.RandStringOfLength(r, 10)

		msg := types.NewMsgCertifyPlatform(certifier.Address, validator.PubKey, platform)

		if !k.IsCertifier(ctx, certifier.Address) {
			return simulation.NewOperationMsgBasic(types.ModuleName, "NoOp: not a certifier", "", false, nil), nil, nil
		}

		account := ak.GetAccount(ctx, certifier.Address)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			certifier.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgCertifyAuditing generates a MsgCertifyAuditing object which fields contain
// a randomly chosen existing certifer, a random contract and a random string description.
func SimulateMsgCertifyAuditing(ak types.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account,
		chainID string) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		certifier, _ := simulation.RandomAcc(r, accs)
		contract := simulation.RandomAccounts(r, 1)[0]
		description := simulation.RandStringOfLength(r, 10)

		msg := types.NewMsgCertifyGeneral("auditing", "address", contract.Address.String(), description, certifier.Address)

		if !k.IsCertifier(ctx, certifier.Address) {
			return simulation.NewOperationMsgBasic(types.ModuleName, "NoOp: not a certifier", "", false, nil), nil, nil
		}

		account := ak.GetAccount(ctx, certifier.Address)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			certifier.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgCertifyProof generates a MsgCertifyProof object which fields contain
// a randomly chosen existing certifer, a random contract and a random string description.
func SimulateMsgCertifyProof(ak types.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string) (
		simulation.OperationMsg, []simulation.FutureOperation, error) {
		certifier, _ := simulation.RandomAcc(r, accs)
		contract := simulation.RandomAccounts(r, 1)[0]
		description := simulation.RandStringOfLength(r, 10)

		msg := types.NewMsgCertifyGeneral("proof", "address", contract.Address.String(), description, certifier.Address)

		if !k.IsCertifier(ctx, certifier.Address) {
			return simulation.NewOperationMsgBasic(types.ModuleName, "NoOp: not a certifier", "", false, nil), nil, nil
		}

		account := ak.GetAccount(ctx, certifier.Address)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			certifier.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
