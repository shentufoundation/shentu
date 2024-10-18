package gov

import (
	"fmt"

	"cosmossdk.io/collections"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
	typesv1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

// InitGenesis stores genesis parameters.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, ak govtypes.AccountKeeper, bk govtypes.BankKeeper, data *typesv1.GenesisState) {
	err := k.ProposalID.Set(ctx, data.StartingProposalId)
	if err != nil {
		panic(err)
	}

	err = k.Params.Set(ctx, *data.Params)
	if err != nil {
		panic(err)
	}

	err = k.Constitution.Set(ctx, data.Constitution)
	if err != nil {
		panic(err)
	}

	err = k.SetCustomParams(ctx, *data.CustomParams)
	if err != nil {
		return
	}

	// check if the deposits pool account exists
	moduleAcc := k.GetGovernanceAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", govtypes.ModuleName))
	}

	var totalDeposits sdk.Coins
	for _, deposit := range data.Deposits {
		err := k.SetDeposit(ctx, *deposit)
		if err != nil {
			return
		}
		totalDeposits = totalDeposits.Add(deposit.Amount...)
	}

	for _, vote := range data.Votes {
		addr, err := ak.AddressCodec().StringToBytes(vote.Voter)
		if err != nil {
			panic(err)
		}
		err = k.Votes.Set(ctx, collections.Join(vote.ProposalId, sdk.AccAddress(addr)), *vote)
		if err != nil {
			panic(err)
		}
	}

	for _, proposalID := range data.CertVotedProposalIds {
		err := k.SetCertVote(ctx, proposalID)
		if err != nil {
			return
		}
	}

	for _, proposal := range data.Proposals {
		switch proposal.Status {
		case govtypesv1.StatusDepositPeriod:
			err := k.InactiveProposalsQueue.Set(ctx, collections.Join(*proposal.DepositEndTime, proposal.Id), proposal.Id)
			if err != nil {
				panic(err)
			}
		case govtypesv1.StatusVotingPeriod:
			err := k.ActiveProposalsQueue.Set(ctx, collections.Join(*proposal.VotingEndTime, proposal.Id), proposal.Id)
			if err != nil {
				panic(err)
			}
		}
		err := k.SetProposal(ctx, *proposal)
		if err != nil {
			panic(err)
		}
	}

	// if account has zero balance it probably means it's not set, so we set it
	balance := bk.GetAllBalances(ctx, moduleAcc.GetAddress())
	if balance.IsZero() {
		ak.SetModuleAccount(ctx, moduleAcc)
	}

	// check if total deposits equals balance, if it doesn't panic because there were export/import errors
	if !balance.Equal(totalDeposits) {
		panic(fmt.Sprintf("expected module account was %s but we got %s", balance.String(), totalDeposits.String()))
	}
}

// ExportGenesis writes the current store values to a genesis file, which can be imported again with InitGenesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) (*typesv1.GenesisState, error) {
	startingProposalID, err := k.ProposalID.Peek(ctx)
	if err != nil {
		return nil, err
	}

	var proposals govtypesv1.Proposals
	err = k.Proposals.Walk(ctx, nil, func(_ uint64, value govtypesv1.Proposal) (stop bool, err error) {
		proposals = append(proposals, &value)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	constitution, err := k.Constitution.Get(ctx)
	if err != nil {
		return nil, err
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	customParams, err := k.GetCustomParams(ctx)
	if err != nil {
		return nil, err
	}

	var proposalsDeposits govtypesv1.Deposits
	err = k.Deposits.Walk(ctx, nil, func(_ collections.Pair[uint64, sdk.AccAddress], value govtypesv1.Deposit) (stop bool, err error) {
		proposalsDeposits = append(proposalsDeposits, &value)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	// export proposals votes
	var proposalsVotes govtypesv1.Votes
	err = k.Votes.Walk(ctx, nil, func(_ collections.Pair[uint64, sdk.AccAddress], value govtypesv1.Vote) (stop bool, err error) {
		proposalsVotes = append(proposalsVotes, &value)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	var certVotedProposalIds []uint64
	for _, proposal := range proposals {
		voted, err := k.GetCertifierVoted(ctx, proposal.Id)
		if err != nil {
			return nil, err
		}
		if voted {
			certVotedProposalIds = append(certVotedProposalIds, proposal.Id)
		}
	}

	return &typesv1.GenesisState{
		StartingProposalId:   startingProposalID,
		Deposits:             proposalsDeposits,
		Votes:                proposalsVotes,
		Proposals:            proposals,
		Params:               &params,
		Constitution:         constitution,
		CustomParams:         &customParams,
		CertVotedProposalIds: certVotedProposalIds,
	}, nil
}
