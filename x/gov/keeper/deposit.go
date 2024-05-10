package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	v046 "github.com/cosmos/cosmos-sdk/x/gov/migrations/v046"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/shentufoundation/shentu/v2/x/gov/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"
)

// AddDeposit adds or updates a deposit of a specific depositor on a specific proposal.
// When the proposal type is ShieldClaim, it's not depositable.
// Activates voting period when appropriate.
func (k Keeper) AddDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) (bool, error) {
	// checks to see if proposal exists
	proposal, ok := k.GetProposal(ctx, proposalID)
	if !ok {
		return false, sdkerrors.Wrap(govtypes.ErrUnknownProposal, fmt.Sprint(proposalID))
	}

	// check if proposal is still depositable or if proposer is council member
	if proposal.Status != govtypesv1.StatusDepositPeriod {
		return false, sdkerrors.Wrap(govtypes.ErrAlreadyActiveProposal, fmt.Sprint(proposalID))
	}

	// update the governance module's account coins pool
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositorAddr, govtypes.ModuleName, depositAmount)
	if err != nil {
		return false, err
	}

	// update proposal
	proposal.TotalDeposit = sdk.NewCoins(proposal.TotalDeposit...).Add(depositAmount...)
	k.SetProposal(ctx, proposal)

	// check if deposit has provided sufficient total funds to transition the proposal into the voting period
	activatedVotingPeriod := false
	legacyProposal, err := v046.ConvertToLegacyProposal(proposal)
	if err != nil {
		return false, err
	}
	if proposal.Status == govtypesv1.StatusDepositPeriod && sdk.NewCoins(proposal.TotalDeposit...).IsAllGTE(k.GetDepositParams(ctx).MinDeposit) ||
		legacyProposal.ProposalType() == shieldtypes.ProposalTypeShieldClaim {
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
		deposit.Amount = sdk.NewCoins(deposit.Amount...).Add(depositAmount...)
	} else {
		deposit = govtypesv1.NewDeposit(proposalID, depositorAddr, depositAmount)
	}

	k.SetDeposit(ctx, deposit)
}
