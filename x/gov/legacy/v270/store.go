package v270

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	v220 "github.com/shentufoundation/shentu/v2/x/gov/legacy/v220"
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
		var oldProposal v220.Proposal
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
		case v220.StatusCertifierVotingPeriod:
		case v220.StatusValidatorVotingPeriod:
			newProposal.Status = govtypes.StatusVotingPeriod
		case v220.StatusPassed:
			newProposal.Status = govtypes.StatusPassed
		case v220.StatusRejected:
			newProposal.Status = govtypes.StatusRejected
		case v220.StatusFailed:
			newProposal.Status = govtypes.StatusFailed
		default:
			newProposal.Status = govtypes.ProposalStatus(oldProposal.Status)
		}

		store.Delete(iterator.Key())
		store.Set(govtypes.ProposalKey(oldProposal.ProposalId), iterator.Value())
	}

	return nil
}
