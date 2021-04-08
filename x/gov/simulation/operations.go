package simulation

import (
	"math"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	certtypes "github.com/certikfoundation/shentu/x/cert/types"
	"github.com/certikfoundation/shentu/x/gov/keeper"
	"github.com/certikfoundation/shentu/x/gov/types"
	shieldtypes "github.com/certikfoundation/shentu/x/shield/types"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(appParams simtypes.AppParams, cdc codec.JSONMarshaler, ak govtypes.AccountKeeper, bk govtypes.BankKeeper, ck types.CertKeeper,
	k keeper.Keeper, wContents []simtypes.WeightedProposalContent) simulation.WeightedOperations {
	// generate the weighted operations for the proposal contents
	var wProposalOps simulation.WeightedOperations

	for _, wContent := range wContents {
		wContent := wContent // pin variable
		var weight int
		appParams.GetOrGenerate(cdc, wContent.AppParamsKey(), &weight, nil,
			func(_ *rand.Rand) { weight = wContent.DefaultWeight() })

		wProposalOps = append(
			wProposalOps,
			simulation.NewWeightedOperation(
				weight,
				SimulateSubmitProposal(ak, bk, ck, k, wContent.ContentSimulatorFn()),
			),
		)
	}

	return wProposalOps
}

// SimulateSubmitProposal simulates creating a msg Submit Proposal
// voting on the proposal, and subsequently slashing the proposal. It is implemented using
// future operations.
func SimulateSubmitProposal(
	ak govtypes.AccountKeeper, bk govtypes.BankKeeper, ck types.CertKeeper, k keeper.Keeper, contentSim simtypes.ContentSimulatorFn,
) simtypes.Operation {
	// The states are:
	// column 1: All validators vote
	// column 2: 90% vote
	// column 3: 75% vote
	// column 4: 40% vote
	// column 5: 15% vote
	// column 6: noone votes
	// All columns sum to 100 for simplicity, values chosen by @valardragon semi-arbitrarily,
	// feel free to change.
	numVotesTransitionMatrix, _ := simulation.CreateTransitionMatrix([][]int{
		{20, 10, 0, 0, 0, 0},
		{55, 50, 20, 10, 0, 0},
		{25, 25, 30, 25, 30, 15},
		{0, 15, 30, 25, 30, 30},
		{0, 0, 20, 30, 30, 30},
		{0, 0, 0, 10, 10, 25},
	})

	statePercentageArray := []float64{1, .9, .75, .4, .15, 0}
	curNumVotesState := 1

	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// 1) submit proposal now
		content := contentSim(r, ctx, accs)
		if content == nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgSubmitProposal, ""), nil, nil
		}

		var (
			deposit sdk.Coins
			err     error
		)
		var simAccount simtypes.Account
		if content.ProposalType() == shieldtypes.ProposalTypeShieldClaim {
			c := content.(*shieldtypes.ShieldClaimProposal)
			for _, simAcc := range accs {
				proposerAddr, _ := sdk.AccAddressFromBech32(c.Proposer)
				if simAcc.Address.Equals(proposerAddr) {
					simAccount = simAcc
					break
				}
			}
			spendable := bk.SpendableCoins(ctx, simAccount.Address)
			if spendable == nil {
				return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgSubmitProposal, ""), nil, nil
			}
			denom := sdk.DefaultBondDenom
			lossAmountDec := c.Loss.AmountOf(denom).ToDec()
			claimProposalParams := k.ShieldKeeper.GetClaimProposalParams(ctx)
			depositRate := claimProposalParams.DepositRate
			minDepositAmountDec := sdk.MaxDec(claimProposalParams.MinDeposit.AmountOf(denom).ToDec(), lossAmountDec.Mul(depositRate))
			minDepositAmount := minDepositAmountDec.Ceil().RoundInt()
			if minDepositAmount.GT(spendable.AmountOf(denom)) {
				return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgSubmitProposal, ""), nil, nil
			}
			deposit = sdk.NewCoins(sdk.NewCoin(denom, minDepositAmount))
		} else {
			simAccount, _ = simtypes.RandomAcc(r, accs)
			spendable := bk.SpendableCoins(ctx, simAccount.Address)
			minDeposit := k.GetDepositParams(ctx).MinDeposit
			if spendable.AmountOf(sdk.DefaultBondDenom).LT(minDeposit.AmountOf(sdk.DefaultBondDenom)) {
				deposit = simtypes.RandSubsetCoins(r, spendable)
			} else {
				deposit = simtypes.RandSubsetCoins(r, minDeposit)
			}
		}

		minInitialDeposit := k.GetDepositParams(ctx).MinInitialDeposit
		if deposit.AmountOf(sdk.DefaultBondDenom).LT(minInitialDeposit.AmountOf(sdk.DefaultBondDenom)) &&
			!k.IsCouncilMember(ctx, simAccount.Address) {
			return simtypes.NewOperationMsgBasic(govtypes.ModuleName,
				"NoOp: insufficient initial deposit amount, skip this tx", "", false, nil), nil, nil
		}

		msg, _ := govtypes.NewMsgSubmitProposal(content, deposit, simAccount.Address)

		account := ak.GetAccount(ctx, simAccount.Address)
		coins := bk.SpendableCoins(ctx, simAccount.Address)

		var fees sdk.Coins
		coins, hasNeg := coins.SafeSub(deposit)
		if !hasNeg {
			fees, err = simtypes.RandomFees(r, ctx, coins)
			if err != nil {
				return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgSubmitProposal, ""), nil, err
			}
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas*100,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		// get the submitted proposal ID
		proposalID, err := k.GetProposalID(ctx)
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgSubmitProposal, ""), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgSubmitProposal, ""), nil, err
		}

		opMsg := simtypes.NewOperationMsg(msg, true, "")

		var fops []simtypes.FutureOperation

		// 2) Schedule deposit operations
		if content.ProposalType() != shieldtypes.ProposalTypeShieldClaim {
			for i := 0; i < 10; i++ {
				fops = append(fops, simtypes.FutureOperation{
					BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 1, 5),
					Op:          SimulateMsgDeposit(ak, bk, k, proposalID),
				})
			}
		}

		// 3) Schedule operations for certifier voting
		if content.ProposalType() == shieldtypes.ProposalTypeShieldClaim ||
			content.ProposalType() == certtypes.ProposalTypeCertifierUpdate ||
			content.ProposalType() == upgradetypes.ProposalTypeSoftwareUpgrade {
			for _, acc := range accs {
				if ck.IsCertifier(ctx, acc.Address) {
					fops = append(fops, simtypes.FutureOperation{
						BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 5, 10),
						Op:          SimulateCertifierMsgVote(ak, bk, ck, k, acc, proposalID),
					})
				}
			}
		}

		// 4) Schedule operations for validator/delegator voting
		if content.ProposalType() == shieldtypes.ProposalTypeShieldClaim {
			for _, acc := range accs {
				if k.IsCertifiedIdentity(ctx, acc.Address) {
					fops = append(fops, simtypes.FutureOperation{
						BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 10, 15),
						Op:          SimulateMsgVote(ak, bk, k, acc, proposalID),
					})
				}
			}
		} else {
			// 4.1) first pick a number of people to vote.
			curNumVotesState = numVotesTransitionMatrix.NextState(r, curNumVotesState)
			numVotes := int(math.Ceil(float64(len(accs)) * statePercentageArray[curNumVotesState]))
			// 4.2) select who votes and when
			whoVotes := r.Perm(len(accs))
			whoVotes = whoVotes[:numVotes]

			for i := 0; i < numVotes; i++ {
				fops = append(fops, simtypes.FutureOperation{
					BlockHeight: int(ctx.BlockHeight()) + simtypes.RandIntBetween(r, 10, 15),
					Op:          SimulateMsgVote(ak, bk, k, accs[whoVotes[i]], proposalID),
				})
			}
		}

		return opMsg, fops, nil
	}
}

func SimulateMsgVote(ak govtypes.AccountKeeper, bk govtypes.BankKeeper, k keeper.Keeper,
	simAccount simtypes.Account, proposalID uint64) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		proposal, ok := k.GetProposal(ctx, proposalID)
		if !ok {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgVote, ""), nil, nil
		}

		if proposal.Status != types.StatusValidatorVotingPeriod {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgVote, ""), nil, nil
		}

		option := randomVotingOption(r)

		msg := govtypes.NewMsgVote(simAccount.Address, proposalID, option)

		account := ak.GetAccount(ctx, simAccount.Address)
		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, simAccount.Address))
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgVote, ""), nil, err
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
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgVote, ""), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func SimulateCertifierMsgVote(ak govtypes.AccountKeeper, bk govtypes.BankKeeper, ck types.CertKeeper, k keeper.Keeper,
	simAccount simtypes.Account, proposalID uint64) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if !ck.IsCertifier(ctx, simAccount.Address) {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgVote, ""), nil, nil
		}

		proposal, ok := k.GetProposal(ctx, proposalID)
		if !ok {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgVote, ""), nil, nil
		}

		if proposal.Status != types.StatusCertifierVotingPeriod {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgVote, ""), nil, nil
		}

		var option govtypes.VoteOption
		if simtypes.RandIntBetween(r, 0, 100) < 70 {
			option = govtypes.OptionYes
		} else {
			option = govtypes.OptionNo
		}

		msg := govtypes.NewMsgVote(simAccount.Address, proposalID, option)

		account := ak.GetAccount(ctx, simAccount.Address)
		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, simAccount.Address))
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgVote, ""), nil, err
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
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgVote, ""), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func SimulateMsgDeposit(ak govtypes.AccountKeeper, bk govtypes.BankKeeper, k keeper.Keeper, proposalID uint64) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		proposal, ok := k.GetProposal(ctx, proposalID)
		if !ok {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgDeposit, ""), nil, nil
		}

		if proposal.Status != types.StatusDepositPeriod {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgDeposit, ""), nil, nil
		}

		simAcc, _ := simtypes.RandomAcc(r, accs)
		acc := ak.GetAccount(ctx, simAcc.Address)
		spendable := bk.SpendableCoins(ctx, simAcc.Address)
		minDeposit := k.GetDepositParams(ctx).MinDeposit
		var deposit sdk.Coins
		if spendable.AmountOf(sdk.DefaultBondDenom).LT(minDeposit.AmountOf(sdk.DefaultBondDenom)) {
			deposit = simtypes.RandSubsetCoins(r, spendable)
		} else {
			deposit = simtypes.RandSubsetCoins(r, minDeposit)
		}

		msg := govtypes.NewMsgDeposit(simAcc.Address, proposalID, deposit)

		fees, err := simtypes.RandomFees(r, ctx, spendable.Sub(deposit))
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgDeposit, ""), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{acc.GetAccountNumber()},
			[]uint64{acc.GetSequence()},
			simAcc.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgDeposit, ""), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// Pick a random voting option
func randomVotingOption(r *rand.Rand) govtypes.VoteOption {
	switch r.Intn(12) {
	case 0:
		return govtypes.OptionAbstain
	case 1:
		return govtypes.OptionNo
	case 2:
		return govtypes.OptionNoWithVeto
	default:
		return govtypes.OptionYes
	}
}
