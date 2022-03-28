package v231

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	v220 "github.com/certikfoundation/shentu/v2/x/gov/legacy/v220"
	"github.com/certikfoundation/shentu/v2/x/gov/types"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1alpha1"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
)

func MigrateShieldClaimProposal(store sdk.KVStore, cdc codec.BinaryCodec) error {
	iterator := sdk.KVStorePrefixIterator(store, govtypes.ProposalsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var proposal types.Proposal
		if err := cdc.Unmarshal(iterator.Value(), &proposal); err != nil {
			return err
		}

		v1, ok := proposal.GetContent().(*v1alpha1.ShieldClaimProposal)
		if ok {
			// proposal is a shield claim proposal v1
			v2 := &v1beta1.ShieldClaimProposal{
				ProposalId:  v1.ProposalId,
				PoolId:      v1.PoolId,
				Loss:        v1.Loss,
				Evidence:    v1.Evidence,
				Description: v1.Description,
				Proposer:    v1.Proposer,
			}
			if err := proposal.SetContent(v2); err != nil {
				return err
			}
			bz, err := cdc.Marshal(&proposal)
			if err != nil {
				return err
			}
			store.Set(iterator.Key(), bz)
		}
	}
	return nil
}

func MigrateLegacyDeposits(store sdk.KVStore, cdc codec.BinaryCodec) error {
	iterator := sdk.KVStorePrefixIterator(store, govtypes.DepositsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deposit govtypes.Deposit
		var legacyDeposit v220.Deposit
		if err := cdc.Unmarshal(iterator.Value(), &deposit); err != nil {
			if err := cdc.Unmarshal(iterator.Value(), &legacyDeposit); err != nil {
				return err
			}
			deposit = *legacyDeposit.Deposit
		}

		bz, err := cdc.Marshal(&deposit)
		if err != nil {
			return err
		}
		store.Set(iterator.Key(), bz)
	}
	return nil
}

func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	err := MigrateLegacyDeposits(store, cdc)
	if err != nil {
		return err
	}
	return MigrateShieldClaimProposal(store, cdc)
}
