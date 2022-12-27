package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	v220 "github.com/shentufoundation/shentu/v2/x/gov/legacy/v220"
	"github.com/shentufoundation/shentu/v2/x/gov/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

// AddDeposit adds or updates a deposit of a specific depositor on a specific proposal.
// When proposer is a council member, it's not depositable.
// Activates voting period when appropriate.
func (k Keeper) AddDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) (bool, error) {
	// checks to see if proposal exists
	proposal, ok := k.GetProposal(ctx, proposalID)
	if !ok {
		return false, sdkerrors.Wrap(govtypes.ErrUnknownProposal, fmt.Sprint(proposalID))
	}
	// check if proposal is still depositable or if proposer is council member
	if proposal.Status != govtypes.StatusDepositPeriod {
		return false, sdkerrors.Wrap(govtypes.ErrAlreadyActiveProposal, fmt.Sprint(proposalID))
	}

	// update the governance module's account coins pool
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositorAddr, govtypes.ModuleName, depositAmount)
	if err != nil {
		return false, err
	}
	// update proposal
	proposal.TotalDeposit = proposal.TotalDeposit.Add(depositAmount...)
	k.SetProposal(ctx, proposal)

	// check if deposit has provided sufficient total funds to transition the proposal into the voting period
	activatedVotingPeriod := false
	if proposal.Status == govtypes.StatusDepositPeriod && proposal.TotalDeposit.IsAllGTE(k.GetDepositParams(ctx).MinDeposit) ||
		proposal.ProposalType() == shieldtypes.ProposalTypeShieldClaim {
		k.ActivateVotingPeriod(ctx, proposal)
		activatedVotingPeriod = true
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			govtypes.EventTypeProposalDeposit,
			sdk.NewAttribute(sdk.AttributeKeyAmount, depositAmount.String()),
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyDepositor, depositorAddr.String()),
		),
	)

	k.upsertDeposit(ctx, proposalID, depositorAddr, depositAmount)

	return activatedVotingPeriod, nil
}

// upsertDeposit updates or inserts a deposit to a proposal.
func (k Keeper) upsertDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) {
	// add or update deposit object
	deposit, found := k.GetDeposit(ctx, proposalID, depositorAddr)
	if found {
		deposit.Amount = deposit.Amount.Add(depositAmount...)
	} else {
		deposit = govtypes.NewDeposit(proposalID, depositorAddr, depositAmount)
	}

	k.SetDeposit(ctx, deposit)
}

// GetAllDeposits returns all the deposits from the store.
func (k Keeper) GetAllDeposits(ctx sdk.Context) (deposits govtypes.Deposits) {
	k.IterateAllDeposits(ctx, func(deposit govtypes.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})
	return
}

// GetDepositsByProposalID returns all the deposits from a proposal.
func (k Keeper) GetDepositsByProposalID(ctx sdk.Context, proposalID uint64) (deposits govtypes.Deposits) {
	k.IterateDeposits(ctx, proposalID, func(deposit govtypes.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})
	return
}

// GetDepositsIteratorByProposalID gets all the deposits on a specific proposal as an sdk.Iterator.
func (k Keeper) GetDepositsIteratorByProposalID(ctx sdk.Context, proposalID uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, govtypes.DepositsKey(proposalID))
}

// RefundDepositsByProposalID refunds and deletes all the deposits on a specific proposal.
func (k Keeper) RefundDepositsByProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	k.IterateDeposits(ctx, proposalID, func(deposit govtypes.Deposit) bool {
		depositor, err := sdk.AccAddressFromBech32(deposit.Depositor)
		if err != nil {
			panic(err)
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, govtypes.ModuleName, depositor, deposit.Amount)
		if err != nil {
			panic(err)
		}

		store.Delete(govtypes.DepositKey(proposalID, depositor))
		return false
	})
}

// DeleteDepositsByProposalID deletes all the deposits on a specific proposal without refunding them.
func (k Keeper) DeleteDepositsByProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	k.IterateDeposits(ctx, proposalID, func(deposit govtypes.Deposit) bool {
		err := k.bankKeeper.BurnCoins(ctx, govtypes.ModuleName, deposit.Amount)
		if err != nil {
			panic(err)
		}

		depositor, err := sdk.AccAddressFromBech32(deposit.Depositor)
		if err != nil {
			panic(err)
		}

		store.Delete(govtypes.DepositKey(proposalID, depositor))
		return false
	})
}

// GetDeposits returns all the deposits from a proposal.
func (k Keeper) GetDeposits(ctx sdk.Context, proposalID uint64) (deposits govtypes.Deposits) {
	k.IterateDeposits(ctx, proposalID, func(deposit govtypes.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})
	return
}

// IterateDeposits iterates over the all the proposals deposits and performs a callback function.
func (k Keeper) IterateDeposits(ctx sdk.Context, proposalID uint64, cb func(deposit govtypes.Deposit) (stop bool)) {
	iterator := k.GetDepositsIteratorByProposalID(ctx, proposalID)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deposit govtypes.Deposit
		err := k.cdc.Unmarshal(iterator.Value(), &deposit)
		if err != nil {
			var legacydeposit v220.Deposit
			k.cdc.MustUnmarshal(iterator.Value(), &legacydeposit)
			deposit = *legacydeposit.Deposit
		}
		if cb(deposit) {
			break
		}
	}
}
