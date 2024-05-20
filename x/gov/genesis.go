package gov

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

// InitGenesis stores genesis parameters.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, ak govtypes.AccountKeeper, bk govtypes.BankKeeper, data *typesv1.GenesisState) {
	k.SetProposalID(ctx, data.StartingProposalId)
	k.SetDepositParams(ctx, *data.DepositParams)
	k.SetVotingParams(ctx, *data.VotingParams)
	k.SetTallyParams(ctx, *data.TallyParams)
	k.SetCustomParams(ctx, *data.CustomParams)

	// check if the deposits pool account exists
	moduleAcc := k.GetGovernanceAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", govtypes.ModuleName))
	}

	var totalDeposits sdk.Coins
	for _, deposit := range data.Deposits {
		k.SetDeposit(ctx, *deposit)
		totalDeposits = totalDeposits.Add(deposit.Amount...)
	}

	for _, vote := range data.Votes {
		k.SetVote(ctx, *vote)
	}

	for _, proposalID := range data.CertVotedProposalIds {
		k.SetCertVote(ctx, proposalID)
	}

	for _, proposal := range data.Proposals {
		switch proposal.Status {
		case govtypesv1.StatusDepositPeriod:
			k.InsertInactiveProposalQueue(ctx, proposal.Id, *proposal.DepositEndTime)
		case govtypesv1.StatusVotingPeriod:
			k.InsertActiveProposalQueue(ctx, proposal.Id, *proposal.VotingEndTime)
		}
		k.SetProposal(ctx, *proposal)
	}

	// if account has zero balance it probably means it's not set, so we set it
	balance := bk.GetAllBalances(ctx, moduleAcc.GetAddress())
	if balance.IsZero() {
		ak.SetModuleAccount(ctx, moduleAcc)
	}

	// check if total deposits equals balance, if it doesn't panic because there were export/import errors
	if !balance.IsEqual(totalDeposits) {
		panic(fmt.Sprintf("expected module account was %s but we got %s", balance.String(), totalDeposits.String()))
	}
}

// ExportGenesis writes the current store values to a genesis file, which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *typesv1.GenesisState {
	startingProposalID, _ := k.GetProposalID(ctx)
	depositParams := k.GetDepositParams(ctx)
	votingParams := k.GetVotingParams(ctx)
	tallyParams := k.GetTallyParams(ctx)
	proposals := k.GetProposals(ctx)
	customParams := k.GetCustomParams(ctx)

	var genState typesv1.GenesisState

	for _, proposal := range proposals {
		genState.Deposits = append(genState.Deposits, k.GetDeposits(ctx, proposal.Id)...)
		genState.Votes = append(genState.Votes, k.GetVotes(ctx, proposal.Id)...)
		if k.GetCertifierVoted(ctx, proposal.Id) {
			genState.CertVotedProposalIds = append(genState.CertVotedProposalIds, proposal.Id)
		}
	}
	genState.StartingProposalId = startingProposalID
	genState.Proposals = proposals
	genState.DepositParams = &depositParams
	genState.VotingParams = &votingParams
	genState.TallyParams = &tallyParams
	genState.CustomParams = &customParams

	return &genState
}
