package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v046 "github.com/cosmos/cosmos-sdk/x/gov/migrations/v046"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

// validatorGovInfo used for tallying
type validatorGovInfo struct {
	Address             sdk.ValAddress                 // address of the validator operator
	BondedTokens        math.Int                       // power of a Validator
	DelegatorShares     sdk.Dec                        // total outstanding delegator shares
	DelegatorDeductions sdk.Dec                        // delegator deductions from validator's delegators voting independently
	Vote                govtypesv1.WeightedVoteOptions // vote of the validator
}

// Tally counts the votes and returns whether the proposal passes and/or if tokens should be burned.
func Tally(ctx sdk.Context, k Keeper, proposal govtypesv1.Proposal) (pass bool, veto bool, tallyResults govtypesv1.TallyResult) {
	results := newResults()

	totalVotingPower := sdk.ZeroDec()
	currValidators := make(map[string]*validatorGovInfo)

	fetchBondedValidators(ctx, k, currValidators)

	k.IterateVotes(ctx, proposal.Id, func(vote govtypesv1.Vote) bool {
		voter, err := sdk.AccAddressFromBech32(vote.Voter)
		if err != nil {
			panic(err)
		}

		// if validator, just record it in the map
		// if delegator tally voting power
		valAddrStr := sdk.ValAddress(voter.Bytes()).String()
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Options
			currValidators[valAddrStr] = val
		}

		delegatorVoting(ctx, k, vote, currValidators, results, &totalVotingPower)

		k.deleteVote(ctx, vote.ProposalId, voter)
		return false
	})

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if len(val.Vote) == 0 {
			continue
		}

		sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
		votingPower := sharesAfterDeductions.MulInt(val.BondedTokens).Quo(val.DelegatorShares)

		for _, option := range val.Vote {
			weight, _ := sdk.NewDecFromStr(option.Weight)
			subPower := votingPower.Mul(weight)
			results[option.Option] = results[option.Option].Add(subPower)
		}
		totalVotingPower = totalVotingPower.Add(votingPower)
	}

	tallyParams := k.GetTallyParams(ctx)
	customParams := k.GetCustomParams(ctx)
	tallyResults = govtypesv1.NewTallyResultFromMap(results)

	var tp govtypesv1.TallyParams
	legacyProposal, err := v046.ConvertToLegacyProposal(proposal)
	if err != nil {
		return false, false, govtypesv1.TallyResult{}
	}
	switch legacyProposal.GetContent().(type) {
	case *certtypes.CertifierUpdateProposal:
		tp = *customParams.CertifierUpdateStakeVoteTally
	default:
		tp = tallyParams
	}

	th := TallyHelper{
		totalVotingPower,
		tp,
		results,
	}
	if legacyProposal.GetContent().ProposalType() == shieldtypes.ProposalTypeShieldClaim {
		pass, veto = passAndVetoStakeResultForShieldClaim(k, ctx, th)
	} else {
		pass, veto = passAndVetoStakeResult(k, ctx, th)
	}

	return pass, veto, tallyResults
}

// TallyHelper reduces number of arguments passed to passAndVetoStakeResult.
type TallyHelper struct {
	totalVotingPower sdk.Dec
	tallyParams      govtypesv1.TallyParams
	results          map[govtypesv1.VoteOption]sdk.Dec
}

func newResults() map[govtypesv1.VoteOption]sdk.Dec {
	return map[govtypesv1.VoteOption]sdk.Dec{
		govtypesv1.OptionYes:        sdk.ZeroDec(),
		govtypesv1.OptionAbstain:    sdk.ZeroDec(),
		govtypesv1.OptionNo:         sdk.ZeroDec(),
		govtypesv1.OptionNoWithVeto: sdk.ZeroDec(),
	}
}

func newValidatorGovInfo(address sdk.ValAddress, bondedTokens math.Int, delegatorShares,
	delegatorDeductions sdk.Dec, vote govtypesv1.WeightedVoteOptions) *validatorGovInfo {
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
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		validators[validator.GetOperator().String()] = newValidatorGovInfo(
			validator.GetOperator(),
			validator.GetBondedTokens(),
			validator.GetDelegatorShares(),
			sdk.ZeroDec(),
			govtypesv1.WeightedVoteOptions{},
		)
		return false
	})
}

func delegatorVoting(ctx sdk.Context, k Keeper, vote govtypesv1.Vote, validators map[string]*validatorGovInfo, results map[govtypesv1.VoteOption]sdk.Dec, totalVotingPower *sdk.Dec) {
	voter, err := sdk.AccAddressFromBech32(vote.Voter)
	if err != nil {
		panic(err)
	}

	// iterate over all delegations from voter, deduct from any delegated-to validators
	k.stakingKeeper.IterateDelegations(ctx, voter, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
		valAddrStr := delegation.GetValidatorAddr().String()

		if val, ok := validators[valAddrStr]; ok {
			// There is no need to handle the special case that validator address equal to voter address.
			// Because voter's voting power will tally again even if there will deduct voter's voting power from validator.
			val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
			validators[valAddrStr] = val

			// delegation shares * bonded / total shares
			votingPower := delegation.GetShares().MulInt(val.BondedTokens).Quo(val.DelegatorShares)

			for _, option := range vote.Options {
				weight, _ := sdk.NewDecFromStr(option.Weight)
				subPower := votingPower.Mul(weight)
				results[option.Option] = results[option.Option].Add(subPower)
			}
			*totalVotingPower = totalVotingPower.Add(votingPower)
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
	percentVoting := th.totalVotingPower.Quo(sdk.NewDecFromInt(k.stakingKeeper.TotalBondedTokens(ctx)))
	quorum, _ := sdk.NewDecFromStr(th.tallyParams.Quorum)
	if percentVoting.LT(quorum) {
		return false, false
	}

	// If no one votes (everyone abstains), proposal fails.
	if th.totalVotingPower.Sub(th.results[govtypesv1.OptionAbstain]).Equal(sdk.ZeroDec()) {
		return false, false
	}

	// If more than 1/3 of voters veto, proposal fails.
	vetoThreshold, _ := sdk.NewDecFromStr(th.tallyParams.VetoThreshold)
	if th.results[govtypesv1.OptionNoWithVeto].Quo(th.totalVotingPower).GT(vetoThreshold) {
		return false, true
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes.
	threshold, _ := sdk.NewDecFromStr(th.tallyParams.Threshold)
	if th.results[govtypesv1.OptionYes].Quo(th.totalVotingPower.Sub(th.results[govtypesv1.OptionAbstain])).
		GT(threshold) {
		return true, false
	}

	// If more than 1/2 of non-abstaining voters vote No, proposal fails.
	return false, false
}

func passAndVetoStakeResultForShieldClaim(k Keeper, ctx sdk.Context, th TallyHelper) (pass bool, veto bool) {
	totalBondedByCertifiedIdentities := k.TotalBondedByCertifiedIdentities(ctx)
	// If there is no staked coins, the proposal fails.
	if totalBondedByCertifiedIdentities.IsZero() {
		return false, false
	}

	// If there is not enough quorum of votes, the proposal fails.
	percentVoting := th.totalVotingPower.Quo(sdk.NewDecFromInt(totalBondedByCertifiedIdentities))
	quorum, _ := sdk.NewDecFromStr(th.tallyParams.Quorum)
	if percentVoting.LT(quorum) {
		return false, false
	}

	// If no one votes (everyone abstains), proposal fails.
	if th.totalVotingPower.Sub(th.results[govtypesv1.OptionAbstain]).Equal(sdk.ZeroDec()) {
		return false, false
	}

	// If more than 1/3 of voters veto, proposal fails.
	vetoThreshold, _ := sdk.NewDecFromStr(th.tallyParams.VetoThreshold)
	if th.results[govtypesv1.OptionNoWithVeto].Quo(th.totalVotingPower).GT(vetoThreshold) {
		return false, true
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes.
	threshold, _ := sdk.NewDecFromStr(th.tallyParams.Threshold)
	if th.results[govtypesv1.OptionYes].Quo(th.totalVotingPower.Sub(th.results[govtypesv1.OptionAbstain])).
		GT(threshold) {
		return true, false
	}

	// If more than 1/2 of non-abstaining voters vote No, proposal fails.
	return false, false
}

// passAndVetoSecurityResult has two storeKey differences from passAndVetoStakeResult:
//  1. Every certifier has equal voting power (1 head =  1 vote)
func passAndVetoSecurityResult(k Keeper, ctx sdk.Context, th TallyHelper) (pass bool) {
	nCertifiers := sdk.NewDec(int64(len(k.CertKeeper.GetAllCertifiers(ctx))))

	// If there are no certifiers, the proposal fails.
	if nCertifiers.IsZero() {
		return false
	}

	// If there is not enough quorum of votes, the proposal fails.
	quorum, _ := sdk.NewDecFromStr(th.tallyParams.Quorum)
	percentVoting := th.totalVotingPower.Quo(nCertifiers)
	if percentVoting.LT(quorum) {
		return false
	}

	// If no one votes (everyone abstains), proposal fails.
	if th.totalVotingPower.IsZero() {
		return false
	}

	// If percentage of "yes" votes is above threshold, proposal passes.
	threshold, _ := sdk.NewDecFromStr(th.tallyParams.Threshold)
	if th.results[govtypesv1.OptionYes].Quo(th.totalVotingPower).
		GT(threshold) {
		return true
	}

	return th.results[govtypesv1.OptionYes].Quo(th.totalVotingPower).GT(threshold)
}

// SecurityTally only gets called if the proposal is a software upgrade or
// certifier update and if it is the certifier round. If the proposal passes,
// we setup the validator voting round and the calling function EndBlocker
// continues to the next iteration. If it fails, the proposal is removed by the
// logic in EndBlocker.
func SecurityTally(ctx sdk.Context, k Keeper, proposal govtypesv1.Proposal) (bool, bool, govtypesv1.TallyResult) {
	results := newResults()
	totalHeadCounts := sdk.ZeroDec()

	currVotes := k.GetVotes(ctx, proposal.Id)
	for _, vote := range currVotes {
		if len(vote.Options) != 1 {
			continue
		}

		results[vote.Options[0].Option] = results[vote.Options[0].Option].Add(sdk.NewDec(1))
		totalHeadCounts = totalHeadCounts.Add(sdk.NewDec(1))
	}
	customParams := k.GetCustomParams(ctx)
	tallyResults := govtypesv1.NewTallyResultFromMap(results)

	ctally := customParams.CertifierUpdateSecurityVoteTally

	th := TallyHelper{
		totalHeadCounts,
		govtypesv1.TallyParams{
			Quorum:        ctally.Quorum,
			Threshold:     ctally.Threshold,
			VetoThreshold: ctally.VetoThreshold,
		},
		results,
	}
	pass := passAndVetoSecurityResult(k, ctx, th)

	var endVoting bool

	// For CertifierUpdateProposal: If security round didn't pass, continue to
	// stake voting (must pass one of the two rounds).
	//
	// For other proposal types (SoftwareUpgrade, etc.): Only continue to stake
	// round if security round passed (must pass both rounds).
	legacyProposal, err := v046.ConvertToLegacyProposal(proposal)
	if err != nil {
		return false, false, govtypesv1.TallyResult{}
	}
	_, isCert := legacyProposal.GetContent().(*certtypes.CertifierUpdateProposal)
	endVoting = (pass && isCert) || (!pass && !isCert)

	return pass, endVoting, tallyResults
}
