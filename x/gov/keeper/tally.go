package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
)

// Tally iterates over the votes and updates the tally of a proposal based on the voting power of the
// voters
func (k Keeper) Tally(ctx context.Context, proposal govtypesv1.Proposal) (passes, burnDeposits bool, tallyResults govtypesv1.TallyResult, err error) {
	results := newResults()

	proposalMsgs, err := proposal.GetMsgs()
	if err != nil {
		return false, false, govtypesv1.TallyResult{}, err
	}

	totalVotingPower := math.LegacyZeroDec()
	currValidators := make(map[string]govtypesv1.ValidatorGovInfo)

	// fetch all the bonded validators, insert them into currValidators
	if err = fetchBondedValidators(ctx, k, currValidators); err != nil {
		return false, false, govtypesv1.TallyResult{}, err
	}

	rng := collections.NewPrefixedPairRange[uint64, sdk.AccAddress](proposal.Id)
	err = k.Votes.Walk(ctx, rng, func(key collections.Pair[uint64, sdk.AccAddress], vote govtypesv1.Vote) (bool, error) {
		// if validator, just record it in the map
		voter, err := k.authKeeper.AddressCodec().StringToBytes(vote.Voter)
		if err != nil {
			return false, err
		}

		valAddrStr, err := k.stakingKeeper.ValidatorAddressCodec().BytesToString(voter)
		if err != nil {
			return false, err
		}
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Options
			currValidators[valAddrStr] = val
		}

		// iterate over all delegations from voter, deduct from any delegated-to validators
		err = k.stakingKeeper.IterateDelegations(ctx, voter, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
			valAddrStr := delegation.GetValidatorAddr()

			if val, ok := currValidators[valAddrStr]; ok {
				// There is no need to handle the special case that validator address equal to voter address.
				// Because voter's voting power will tally again even if there will be deduction of voter's voting power from validator.
				val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
				currValidators[valAddrStr] = val

				// delegation shares * bonded / total shares
				votingPower := delegation.GetShares().MulInt(val.BondedTokens).Quo(val.DelegatorShares)

				for _, option := range vote.Options {
					weight, _ := math.LegacyNewDecFromStr(option.Weight)
					subPower := votingPower.Mul(weight)
					results[option.Option] = results[option.Option].Add(subPower)
				}
				totalVotingPower = totalVotingPower.Add(votingPower)
			}

			return false
		})
		if err != nil {
			return false, err
		}

		delegatorVoting(ctx, k, vote, currValidators, results, &totalVotingPower)

		return false, k.Votes.Remove(ctx, collections.Join(vote.ProposalId, sdk.AccAddress(voter)))
	})

	if err != nil {
		return false, false, tallyResults, err
	}

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if len(val.Vote) == 0 {
			continue
		}

		sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
		votingPower := sharesAfterDeductions.MulInt(val.BondedTokens).Quo(val.DelegatorShares)

		for _, option := range val.Vote {
			weight, _ := math.LegacyNewDecFromStr(option.Weight)
			subPower := votingPower.Mul(weight)
			results[option.Option] = results[option.Option].Add(subPower)
		}
		totalVotingPower = totalVotingPower.Add(votingPower)
	}

	params, _ := k.Params.Get(ctx)
	tallyParams := govtypesv1.TallyParams{
		Quorum:        params.Quorum,
		Threshold:     params.Threshold,
		VetoThreshold: params.VetoThreshold,
	}
	customParams, _ := k.GetCustomParams(ctx)
	tallyResults = govtypesv1.NewTallyResultFromMap(results)
	th := TallyHelper{
		totalVotingPower,
		tallyParams,
		results,
	}

	if len(proposalMsgs) == 1 {
		if legacyMsg, ok := proposalMsgs[0].(*govtypesv1.MsgExecLegacyContent); ok {
			// check that the content struct can be unmarshalled
			content, err := govtypesv1.LegacyContentFromMessage(legacyMsg)
			if err != nil {
				return false, false, govtypesv1.TallyResult{}, err
			}
			if content.ProposalType() == certtypes.ProposalTypeCertifierUpdate {
				th.tallyParams = *customParams.CertifierUpdateStakeVoteTally
			}
		}
	}

	pass, veto := passAndVetoStakeResult(k, ctx, th)

	return pass, veto, tallyResults, nil
}

// TallyHelper reduces number of arguments passed to passAndVetoStakeResult.
type TallyHelper struct {
	totalVotingPower math.LegacyDec
	tallyParams      govtypesv1.TallyParams
	results          map[govtypesv1.VoteOption]math.LegacyDec
}

func newResults() map[govtypesv1.VoteOption]math.LegacyDec {
	return map[govtypesv1.VoteOption]math.LegacyDec{
		govtypesv1.OptionYes:        math.LegacyZeroDec(),
		govtypesv1.OptionAbstain:    math.LegacyZeroDec(),
		govtypesv1.OptionNo:         math.LegacyZeroDec(),
		govtypesv1.OptionNoWithVeto: math.LegacyZeroDec(),
	}
}

// fetchBondedValidators fetches all the bonded validators, insert them into currValidators.
func fetchBondedValidators(ctx context.Context, k Keeper, validators map[string]govtypesv1.ValidatorGovInfo) error {
	err := k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		valBz, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(validator.GetOperator())
		if err != nil {
			return false
		}
		validators[validator.GetOperator()] = govtypesv1.NewValidatorGovInfo(
			valBz,
			validator.GetBondedTokens(),
			validator.GetDelegatorShares(),
			math.LegacyZeroDec(),
			govtypesv1.WeightedVoteOptions{},
		)

		return false
	})
	if err != nil {
		return err
	}
	return nil
}

func delegatorVoting(ctx context.Context, k Keeper, vote govtypesv1.Vote, validators map[string]govtypesv1.ValidatorGovInfo, results map[govtypesv1.VoteOption]math.LegacyDec, totalVotingPower *math.LegacyDec) {
	voter, err := sdk.AccAddressFromBech32(vote.Voter)
	if err != nil {
		panic(err)
	}

	// iterate over all delegations from voter, deduct from any delegated-to validators
	k.stakingKeeper.IterateDelegations(ctx, voter, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
		valAddrStr := delegation.GetValidatorAddr()

		if val, ok := validators[valAddrStr]; ok {
			// There is no need to handle the special case that validator address equal to voter address.
			// Because voter's voting power will tally again even if there will deduct voter's voting power from validator.
			val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
			validators[valAddrStr] = val

			// delegation shares * bonded / total shares
			votingPower := delegation.GetShares().MulInt(val.BondedTokens).Quo(val.DelegatorShares)

			for _, option := range vote.Options {
				weight, _ := math.LegacyNewDecFromStr(option.Weight)
				subPower := votingPower.Mul(weight)
				results[option.Option] = results[option.Option].Add(subPower)
			}
			*totalVotingPower = totalVotingPower.Add(votingPower)
		}

		return false
	})
}

func passAndVetoStakeResult(k Keeper, ctx context.Context, th TallyHelper) (pass bool, veto bool) {
	// If there is no staked coins, the proposal fails.
	totalBonded, err := k.stakingKeeper.TotalBondedTokens(ctx)
	if err != nil {
		return false, false
	}
	if totalBonded.IsZero() {
		return false, false
	}

	// If there is not enough quorum of votes, the proposal fails.
	percentVoting := th.totalVotingPower.Quo(math.LegacyNewDecFromInt(totalBonded))
	quorum, _ := math.LegacyNewDecFromStr(th.tallyParams.Quorum)
	if percentVoting.LT(quorum) {
		return false, false
	}

	// If no one votes (everyone abstains), proposal fails.
	if th.totalVotingPower.Sub(th.results[govtypesv1.OptionAbstain]).Equal(math.LegacyZeroDec()) {
		return false, false
	}

	// If more than 1/3 of voters veto, proposal fails.
	vetoThreshold, _ := math.LegacyNewDecFromStr(th.tallyParams.VetoThreshold)
	if th.results[govtypesv1.OptionNoWithVeto].Quo(th.totalVotingPower).GT(vetoThreshold) {
		return false, true
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes.
	threshold, _ := math.LegacyNewDecFromStr(th.tallyParams.Threshold)
	if th.results[govtypesv1.OptionYes].Quo(th.totalVotingPower.Sub(th.results[govtypesv1.OptionAbstain])).
		GT(threshold) {
		return true, false
	}

	// If more than 1/2 of non-abstaining voters vote No, proposal fails.
	return false, false
}

// passAndVetoSecurityResult has two storeKey differences from passAndVetoStakeResult:
//  1. Every certifier has equal voting power (1 head =  1 vote)
func passAndVetoSecurityResult(k Keeper, ctx context.Context, th TallyHelper) (pass bool) {
	// If no one votes (everyone abstains), proposal fails.
	if th.totalVotingPower.IsZero() {
		return false
	}

	nCertifiers := math.LegacyNewDec(int64(len(k.certKeeper.GetAllCertifiers(ctx))))
	// If there are no certifiers, the proposal fails.
	if nCertifiers.IsZero() {
		return false
	}

	// If there is not enough quorum of votes, the proposal fails.
	quorum := math.LegacyMustNewDecFromStr(th.tallyParams.Quorum)
	if th.totalVotingPower.Quo(nCertifiers).LT(quorum) {
		return false
	}

	// If percentage of "yes" votes is above threshold, proposal passes.
	threshold := math.LegacyMustNewDecFromStr(th.tallyParams.Threshold)
	return th.results[govtypesv1.OptionYes].Quo(th.totalVotingPower).GT(threshold)
}

// SecurityTally only gets called if the proposal is a software upgrade or
// certifier update and if it is the certifier round. If the proposal passes,
// we set up the validator voting round and the calling function EndBlocker
// continues to the next iteration. If it fails, the proposal is removed by the
// logic in EndBlocker.
func SecurityTally(ctx context.Context, k Keeper, proposal govtypesv1.Proposal) (bool, bool, govtypesv1.TallyResult) {
	proposalMsgs, err := proposal.GetMsgs()
	if err != nil {
		return false, false, govtypesv1.TallyResult{}
	}

	results := newResults()
	totalHeadCounts := math.LegacyZeroDec()

	var currVotes []govtypesv1.Vote
	rng := collections.NewPrefixedPairRange[uint64, sdk.AccAddress](proposal.Id)
	err = k.Votes.Walk(ctx, rng, func(key collections.Pair[uint64, sdk.AccAddress], vote govtypesv1.Vote) (bool, error) {
		currVotes = append(currVotes, vote)
		return true, nil
	})

	for _, vote := range currVotes {
		if len(vote.Options) != 1 {
			continue
		}
		results[vote.Options[0].Option] = results[vote.Options[0].Option].Add(math.LegacyNewDec(1))
		totalHeadCounts = totalHeadCounts.Add(math.LegacyNewDec(1))
	}
	customParams, err := k.GetCustomParams(ctx)
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

	var endVoting, isCertifierUpdateProposal bool
	// For CertifierUpdateProposal: If security round didn't pass, continue to
	// stake voting (must pass one of the two rounds).
	//
	// For other proposal types (SoftwareUpgrade, etc.): Only continue to stake
	// round if security round passed (must pass both rounds).
	if len(proposalMsgs) == 1 {
		if legacyMsg, ok := proposalMsgs[0].(*govtypesv1.MsgExecLegacyContent); ok {
			// check that the content struct can be unmarshalled
			content, err := govtypesv1.LegacyContentFromMessage(legacyMsg)
			if err != nil {
				return false, false, govtypesv1.TallyResult{}
			}
			if content.ProposalType() == certtypes.ProposalTypeCertifierUpdate {
				isCertifierUpdateProposal = true
			}
		}
	}

	endVoting = (pass && isCertifierUpdateProposal) || (!pass && !isCertifierUpdateProposal)
	return pass, endVoting, tallyResults
}
