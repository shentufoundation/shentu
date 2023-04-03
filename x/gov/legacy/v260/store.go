package v260

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/shentufoundation/shentu/v2/x/gov/types"
)

// MigrateProposalStore performs migration of ProposalKey.Specifically, it performs:
// - Replace the old Proposal status
// - ProposalKey changed from SmallEndian to BigEndian
// - Delete old proposal data and add new proposal data
// nolint
func MigrateProposalStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, govtypes.ProposalsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var oldProposal Proposal
		if err := cdc.Unmarshal(iterator.Value(), &oldProposal); err != nil {
			return err
		}
		newProposal := govtypes.Proposal{
			ProposalId:       oldProposal.ProposalId,
			Content:          oldProposal.Content,
			FinalTallyResult: oldProposal.FinalTallyResult,
			SubmitTime:       oldProposal.SubmitTime,
			DepositEndTime:   oldProposal.DepositEndTime,
			TotalDeposit:     oldProposal.TotalDeposit,
			VotingStartTime:  oldProposal.VotingStartTime,
			VotingEndTime:    oldProposal.VotingEndTime,
		}

		switch oldProposal.Status {
		case StatusCertifierVotingPeriod:
		case StatusValidatorVotingPeriod:
			newProposal.Status = govtypes.StatusVotingPeriod
		case StatusPassed:
			newProposal.Status = govtypes.StatusPassed
		case StatusRejected:
			newProposal.Status = govtypes.StatusRejected
		case StatusFailed:
			newProposal.Status = govtypes.StatusFailed
		default:
			newProposal.Status = govtypes.ProposalStatus(oldProposal.Status)
		}

		store.Delete(iterator.Key())
		store.Set(govtypes.ProposalKey(oldProposal.ProposalId), iterator.Value())
	}

	return nil
}

func MigrateParams(ctx sdk.Context, paramSubspace types.ParamSubspace) error {
	var (
		depositParams  govtypes.DepositParams
		oldTallyParams TallyParams
	)

	paramSubspace.Get(ctx, govtypes.ParamStoreKeyDepositParams, &depositParams)
	tallyParamsBytes := paramSubspace.GetRaw(ctx, govtypes.ParamStoreKeyTallyParams)
	if err := json.Unmarshal(tallyParamsBytes, &oldTallyParams); err != nil {
		return err
	}

	// tallyParams
	defaultTally := oldTallyParams.DefaultTally

	securityVoteTally := oldTallyParams.CertifierUpdateSecurityVoteTally
	stakeVoteTally := oldTallyParams.CertifierUpdateStakeVoteTally
	tallyParams := govtypes.NewTallyParams(defaultTally.Quorum, defaultTally.Threshold, defaultTally.VetoThreshold)
	// customParams
	certifierUpdateSecurityVoteTally := govtypes.NewTallyParams(
		securityVoteTally.Quorum,
		securityVoteTally.Threshold,
		securityVoteTally.VetoThreshold,
	)
	certifierUpdateStakeVoteTally := govtypes.NewTallyParams(
		stakeVoteTally.Quorum,
		stakeVoteTally.Threshold,
		stakeVoteTally.VetoThreshold,
	)
	customParams := types.CustomParams{
		CertifierUpdateSecurityVoteTally: &certifierUpdateSecurityVoteTally,
		CertifierUpdateStakeVoteTally:    &certifierUpdateStakeVoteTally,
	}

	// set migrate params
	paramSubspace.Set(ctx, govtypes.ParamStoreKeyDepositParams, &depositParams)
	paramSubspace.Set(ctx, govtypes.ParamStoreKeyTallyParams, &tallyParams)
	paramSubspace.Set(ctx, types.ParamStoreKeyCustomParams, &customParams)

	return nil
}
