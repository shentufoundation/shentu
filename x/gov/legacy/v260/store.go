package v260

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// MigrateProposalStore performs migration of ProposalKey.Specifically, it performs:
// - Replace the old Proposal status
// - ProposalKey changed from SmallEndian to BigEndian
// - Delete old proposal data and add new proposal data
// nolint
func MigrateProposalStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, govtypes.ProposalsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var oldProposal Proposal
		if err := cdc.Unmarshal(iterator.Value(), &oldProposal); err != nil {
			return err
		}
		newProposal := govtypesv1beta1.Proposal{
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
			newProposal.Status = govtypesv1beta1.StatusVotingPeriod
		case StatusPassed:
			newProposal.Status = govtypesv1beta1.StatusPassed
		case StatusRejected:
			newProposal.Status = govtypesv1beta1.StatusRejected
		case StatusFailed:
			newProposal.Status = govtypesv1beta1.StatusFailed
		default:
			newProposal.Status = govtypesv1beta1.ProposalStatus(oldProposal.Status)
		}

		store.Delete(iterator.Key())
		bz, err := cdc.Marshal(&newProposal)
		if err != nil {
			return err
		}
		store.Set(govtypes.ProposalKey(newProposal.ProposalId), bz)
	}

	return nil
}

//func MigrateParams(ctx sdk.Context, paramSubspace exported.ParamSubspace) error {
//	var (
//		depositParams  govtypesv1beta1.DepositParams
//		oldTallyParams TallyParams
//	)
//
//	paramSubspace.Get(ctx, govtypesv1.ParamStoreKeyDepositParams, &depositParams)
//	tallyParamsBytes := paramSubspace.GetRaw(ctx, govtypesv1.ParamStoreKeyTallyParams)
//	if err := json.Unmarshal(tallyParamsBytes, &oldTallyParams); err != nil {
//		return err
//	}
//
//	// tallyParams
//	defaultTally := oldTallyParams.DefaultTally
//
//	securityVoteTally := oldTallyParams.CertifierUpdateSecurityVoteTally
//	stakeVoteTally := oldTallyParams.CertifierUpdateStakeVoteTally
//	tallyParams := govtypesv1beta1.NewTallyParams(defaultTally.Quorum, defaultTally.Threshold, defaultTally.VetoThreshold)
//	// customParams
//	certifierUpdateSecurityVoteTally := govtypesv1beta1.NewTallyParams(
//		securityVoteTally.Quorum,
//		securityVoteTally.Threshold,
//		securityVoteTally.VetoThreshold,
//	)
//	certifierUpdateStakeVoteTally := govtypesv1beta1.NewTallyParams(
//		stakeVoteTally.Quorum,
//		stakeVoteTally.Threshold,
//		stakeVoteTally.VetoThreshold,
//	)
//	customParams := v1alpha1.CustomParams{
//		CertifierUpdateSecurityVoteTally: &certifierUpdateSecurityVoteTally,
//		CertifierUpdateStakeVoteTally:    &certifierUpdateStakeVoteTally,
//	}
//
//	// set migrate params
//	paramSubspace.Set(ctx, govtypesv1.ParamStoreKeyDepositParams, &depositParams)
//	paramSubspace.Set(ctx, govtypesv1.ParamStoreKeyTallyParams, &tallyParams)
//	paramSubspace.Set(ctx, v1alpha1.ParamStoreKeyCustomParams, &customParams)
//
//	return nil
//}
