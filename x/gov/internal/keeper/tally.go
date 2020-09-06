package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"

	"github.com/certikfoundation/shentu/x/cert"
	"github.com/certikfoundation/shentu/x/gov/internal/types"
)

// validatorGovInfo used for tallying
type validatorGovInfo struct {
	Address             sdk.ValAddress      // address of the validator operator
	BondedTokens        sdk.Int             // power of a Validator
	DelegatorShares     sdk.Dec             // total outstanding delegator shares
	DelegatorDeductions sdk.Dec             // delegator deductions from validator's delegators voting independently
	Vote                govTypes.VoteOption // vote of the validator
}

// Tally counts the votes and returns whether the proposal passes and/or if tokens should be burned.
func Tally(ctx sdk.Context, k Keeper, proposal types.Proposal) (pass bool, veto bool, tallyResults govTypes.TallyResult) {
	results := newResults()

	totalVotingPower := sdk.ZeroDec()
	currValidators := make(map[string]*validatorGovInfo)

	fetchBondedValidators(ctx, k, currValidators)

	k.IterateVotes(ctx, proposal.ProposalID, func(vote types.Vote) bool {
		// if validator, just record it in the map
		// if delegator tally voting power
		valAddrStr := sdk.ValAddress(vote.Voter).String()
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Option
			currValidators[valAddrStr] = val
		} else {
			delegatorVoting(ctx, k, vote, currValidators, results, &totalVotingPower)
		}
		k.deleteVote(ctx, vote.ProposalID, vote.Voter)
		return false
	})

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if val.Vote == govTypes.OptionEmpty {
			continue
		}

		sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
		fractionAfterDeductions := sharesAfterDeductions.Quo(val.DelegatorShares)
		votingPower := fractionAfterDeductions.MulInt(val.BondedTokens)

		results[val.Vote] = results[val.Vote].Add(votingPower)
		totalVotingPower = totalVotingPower.Add(votingPower)
	}

	tallyParams := k.GetTallyParams(ctx)
	tallyResults = govTypes.NewTallyResultFromMap(results)

	var tp govTypes.TallyParams
	switch proposal.Content.(type) {
	case cert.CertifierUpdateProposal:
		tp = tallyParams.CertifierUpdateStakeVoteTally
	default:
		tp = tallyParams.DefaultTally
	}

	th := TallyHelper{
		totalVotingPower,
		tp,
		results,
	}
	pass, veto = passAndVetoStakeResult(k, ctx, th)

	return pass, veto, tallyResults
}

// TallyHelper reduces number of arguments passed to passAndVetoStakeResult.
type TallyHelper struct {
	totalVotingPower sdk.Dec
	tallyParams      govTypes.TallyParams
	results          map[govTypes.VoteOption]sdk.Dec
}

func newResults() map[govTypes.VoteOption]sdk.Dec {
	return map[govTypes.VoteOption]sdk.Dec{
		govTypes.OptionYes:        sdk.ZeroDec(),
		govTypes.OptionAbstain:    sdk.ZeroDec(),
		govTypes.OptionNo:         sdk.ZeroDec(),
		govTypes.OptionNoWithVeto: sdk.ZeroDec(),
	}
}

func newValidatorGovInfo(address sdk.ValAddress, bondedTokens sdk.Int, delegatorShares,
	delegatorDeductions sdk.Dec, vote govTypes.VoteOption) *validatorGovInfo {
	return &validatorGovInfo{
		Address:             address,
		BondedTokens:        bondedTokens,
		DelegatorShares:     delegatorShares,
		DelegatorDeductions: delegatorDeductions,
		Vote:                vote,
	}
}

// fetchBondedValidators fetches all the bonded validators, insert them into currValidators.
func fetchBondedValidators(ctx sdk.Context, k Keeper, validators map[string]*validatorGovInfo) {
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator exported.ValidatorI) (stop bool) {
		validators[validator.GetOperator().String()] = newValidatorGovInfo(
			validator.GetOperator(),
			validator.GetBondedTokens(),
			validator.GetDelegatorShares(),
			sdk.ZeroDec(),
			govTypes.OptionEmpty,
		)
		return false
	})
}

func delegatorVoting(ctx sdk.Context, k Keeper, vote types.Vote, validators map[string]*validatorGovInfo,
	results map[govTypes.VoteOption]sdk.Dec, totalVotingPower *sdk.Dec) {
	// iterate over all delegations from voter, deduct from any delegated-to validators
	k.stakingKeeper.IterateDelegations(ctx, vote.Voter, func(index int64, delegation exported.DelegationI) (stop bool) {
		valAddrStr := delegation.GetValidatorAddr().String()

		if val, ok := validators[valAddrStr]; ok {
			val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
			validators[valAddrStr] = val

			delegatorShare := delegation.GetShares().Quo(val.DelegatorShares)
			votingPower := delegatorShare.MulInt(val.BondedTokens)

			results[vote.Option] = results[vote.Option].Add(votingPower)
			// false positive
			*totalVotingPower = (*totalVotingPower).Add(votingPower)
		}
		return false
	})
}

func passAndVetoStakeResult(k Keeper, ctx sdk.Context, th TallyHelper) (pass bool, veto bool) {
	// If there is no staked coins, the proposal fails.
	if k.stakingKeeper.TotalBondedTokens(ctx).IsZero() {
		return false, false
	}

	// If there is not enough quorum of votes, the proposal fails.
	percentVoting := th.totalVotingPower.Quo(k.stakingKeeper.TotalBondedTokens(ctx).ToDec())
	if percentVoting.LT(th.tallyParams.Quorum) {
		return false, false
	}

	// If no one votes (everyone abstains), proposal fails.
	if th.totalVotingPower.Sub(th.results[govTypes.OptionAbstain]).Equal(sdk.ZeroDec()) {
		return false, false
	}

	// If more than 1/3 of voters veto, proposal fails.
	if th.results[govTypes.OptionNoWithVeto].Quo(th.totalVotingPower).GT(th.tallyParams.Veto) {
		return false, true
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes.
	if th.results[govTypes.OptionYes].Quo(th.totalVotingPower.Sub(th.results[govTypes.OptionAbstain])).
		GT(th.tallyParams.Threshold) {
		return true, false
	}

	// If more than 1/2 of non-abstaining voters vote No, proposal fails.
	return false, false
}

// passAndVetoSecurityResult has two storeKey differences from passAndVetoStakeResult:
//		1. Every certifier has equal voting power (1 head =  1 vote)
//		2. The only voting options are "yes" and "no".
func passAndVetoSecurityResult(k Keeper, ctx sdk.Context, th TallyHelper) (pass bool) {
	nCertifiers := sdk.NewDec(int64(len(k.CertKeeper.GetAllCertifiers(ctx))))

	// If there are no certifiers, the proposal fails.
	if nCertifiers.IsZero() {
		return false
	}

	// If there is not enough quorum of votes, the proposal fails.
	percentVoting := th.totalVotingPower.Quo(nCertifiers)
	if percentVoting.LT(th.tallyParams.Quorum) {
		return false
	}

	// If no one votes (everyone abstains), proposal fails.
	if th.totalVotingPower.IsZero() {
		return false
	}

	// If percentage of "yes" votes is above threshold, proposal passes.
	if th.results[govTypes.OptionYes].Quo(th.totalVotingPower).
		GT(th.tallyParams.Threshold) {
		return true
	}

	// If there are not enough "yes" votes, proposal fails.
	return false
}

// SecurityTally only gets called if the proposal is a software upgrade or
// certifier update and if it is the certifier round. If the proposal passes,
// we setup the validator voting round and the calling function EndBlocker
// continues to the next iteration. If it fails, the proposal is removed by the
// logic in EndBlocker.
func SecurityTally(ctx sdk.Context, k Keeper, proposal types.Proposal) (bool, bool, govTypes.TallyResult) {
	results := newResults()
	totalHeadCounts := sdk.ZeroDec()

	currVotes := k.GetAllVotes(ctx)
	for _, val := range currVotes {
		if val.Option == govTypes.OptionEmpty {
			continue
		}
		results[val.Option] = results[val.Option].Add(sdk.NewDec(1))
		totalHeadCounts = totalHeadCounts.Add(sdk.NewDec(1))
	}
	tallyParams := k.GetTallyParams(ctx)
	tallyResults := govTypes.NewTallyResultFromMap(results)

	th := TallyHelper{
		totalHeadCounts,
		tallyParams.CertifierUpdateSecurityVoteTally,
		results,
	}
	pass := passAndVetoSecurityResult(k, ctx, th)

	var endVoting bool

	// For CertifierUpdateProposal: If security round didn't pass, continue to
	// stake voting (must pass one of the two rounds).
	//
	// For other proposal types (SoftwareUpgrade, etc.): Only continue to stake
	// round if security round passed (must pass both rounds).
	_, isCert := proposal.Content.(cert.CertifierUpdateProposal)
	endVoting = (pass && isCert) || (!pass && !isCert)

	return pass, endVoting, tallyResults
}

// TODO:
//		Query tally in certifier round should show headcount, not amount staked
