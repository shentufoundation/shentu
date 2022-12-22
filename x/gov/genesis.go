package gov

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
)

// InitGenesis stores genesis parameters.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, ak govTypes.AccountKeeper, bk govTypes.BankKeeper, data govTypes.GenesisState) {
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
		case govTypes.StatusDepositPeriod:
			k.InsertInactiveProposalQueue(ctx, proposal.ProposalId, proposal.DepositEndTime)
		case govTypes.StatusVotingPeriod:
			k.InsertActiveProposalQueue(ctx, proposal.ProposalId, proposal.VotingEndTime)
		}
		k.SetProposal(ctx, proposal)
	}

	// if account has zero balance it probably means it's not set, so we set it
	balance := bk.GetAllBalances(ctx, moduleAcc.GetAddress())

	if balance.IsZero() {
		ak.SetModuleAccount(ctx, moduleAcc)
	}

	if !balance.IsEqual(totalDeposits) {
		panic(fmt.Sprintf("expected module account was %s but we got %s", balance.String(), totalDeposits.String()))
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *govTypes.GenesisState {
	startingProposalID, _ := k.GetProposalID(ctx)
	depositParams := k.GetDepositParams(ctx)
	votingParams := k.GetVotingParams(ctx)
	tallyParams := k.GetTallyParams(ctx)
	proposals := k.GetProposals(ctx)

	var proposalsDeposits govTypes.Deposits
	var proposalsVotes govTypes.Votes
	for _, proposal := range proposals {
		deposits := k.GetDeposits(ctx, proposal.ProposalId)
		proposalsDeposits = append(proposalsDeposits, deposits...)

		votes := k.GetVotes(ctx, proposal.ProposalId)
		proposalsVotes = append(proposalsVotes, votes...)
	}

	return &govTypes.GenesisState{
		StartingProposalId: startingProposalID,
		Deposits:           proposalsDeposits,
		Votes:              proposalsVotes,
		Proposals:          proposals,
		DepositParams:      depositParams,
		VotingParams:       votingParams,
		TallyParams:        tallyParams,
	}
}
