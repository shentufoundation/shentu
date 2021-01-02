package gov

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/x/gov/keeper"
	"github.com/certikfoundation/shentu/x/gov/types"
)

// InitGenesis stores genesis parameters.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, ak govTypes.AccountKeeper, bk govTypes.BankKeeper, data types.GenesisState) {
	k.SetProposalID(ctx, data.StartingProposalId)
	k.SetDepositParams(ctx, data.DepositParams)
	k.SetVotingParams(ctx, data.VotingParams)
	k.SetTallyParams(ctx, data.TallyParams)

	// check if the deposits pool account exists
	moduleAcc := k.GetGovernanceAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", govTypes.ModuleName))
	}

	var totalDeposits sdk.Coins
	for _, deposit := range data.Deposits {
		k.SetDeposit(ctx, deposit)
		totalDeposits = totalDeposits.Add(deposit.Amount...)
	}

	for _, vote := range data.Votes {
		k.SetVote(ctx, vote)
	}

	for _, proposal := range data.Proposals {
		switch proposal.Status {
		case types.StatusDepositPeriod:
			k.InsertInactiveProposalQueue(ctx, proposal.ProposalId, proposal.DepositEndTime)
		case types.StatusCertifierVotingPeriod:
			k.InsertActiveProposalQueue(ctx, proposal.ProposalId, proposal.VotingEndTime)
		case types.StatusValidatorVotingPeriod:
			k.InsertActiveProposalQueue(ctx, proposal.ProposalId, proposal.VotingEndTime)
		}
		k.SetProposal(ctx, proposal)
	}

	// add coins if not provided on genesis
	if bk.GetAllBalances(ctx, moduleAcc.GetAddress()).IsZero() {
		if err := bk.SetBalances(ctx, moduleAcc.GetAddress(), totalDeposits); err != nil {
			panic(err)
		}

		ak.SetModuleAccount(ctx, moduleAcc)
	}
}

// ExportGenesis writes the current store values to a genesis file, which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	startingProposalID, _ := k.GetProposalID(ctx)
	depositParams := k.GetDepositParams(ctx)
	votingParams := k.GetVotingParams(ctx)
	tallyParams := k.GetTallyParams(ctx)
	proposals := k.GetProposals(ctx)

	var proposalsDeposits types.Deposits
	var proposalsVotes types.Votes

	for _, proposal := range proposals {
		deposits := k.GetDepositsByProposalID(ctx, proposal.ProposalId)
		proposalsDeposits = append(proposalsDeposits, deposits...)

		votes := k.GetVotes(ctx, proposal.ProposalId)
		proposalsVotes = append(proposalsVotes, votes...)
	}

	return &types.GenesisState{
		StartingProposalId: startingProposalID,
		Deposits:           proposalsDeposits,
		Votes:              proposalsVotes,
		Proposals:          proposals,
		DepositParams:      depositParams,
		VotingParams:       votingParams,
		TallyParams:        tallyParams,
	}
}
