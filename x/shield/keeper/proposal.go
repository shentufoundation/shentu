package keeper

import (
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/certikfoundation/shentu/x/shield/types"
)

func (k Keeper) ClaimLock(ctx sdk.Context, proposalID, poolID uint64, purchaser sdk.AccAddress, purchaseID uint64, loss sdk.Coins, lockPeriod time.Duration) error {
	return nil
}

func (k Keeper) ClaimUnlock(ctx sdk.Context, proposalID, poolID uint64, loss sdk.Coins) error {
	return nil
}

func (k Keeper) RestoreShield(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, id uint64, loss sdk.Coins) error {
	return nil
}

// SetReimbursement sets a reimbursement in store.
func (k Keeper) SetReimbursement(ctx sdk.Context, proposalID uint64, payout types.Reimbursement) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(payout)
	store.Set(types.GetReimbursementKey(proposalID), bz)
}

// GetReimbursement get a reimbursement in store.
func (k Keeper) GetReimbursement(ctx sdk.Context, proposalID uint64) (types.Reimbursement, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetReimbursementKey(proposalID))
	if bz != nil {
		var reimbursement types.Reimbursement
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &reimbursement)
		return reimbursement, nil
	}
	return types.Reimbursement{}, types.ErrCompensationNotFound
}

// DeleteReimbursement deletes a reimbursement.
func (k Keeper) DeleteReimbursement(ctx sdk.Context, proposalID uint64) error {
	store := ctx.KVStore(k.storeKey)
	if _, err := k.GetReimbursement(ctx, proposalID); err != nil {
		return err
	}
	store.Delete(types.GetReimbursementKey(proposalID))
	return nil
}

// CreateReimbursement creates a reimbursement.
func (k Keeper) CreateReimbursement(ctx sdk.Context, proposalID uint64, poolID uint64, rmb sdk.Coins, beneficiary sdk.AccAddress) error {
	totalCollateral := k.GetTotalCollateral(ctx)
	rmbAmt := rmb.AmountOf(k.BondDenom(ctx))
	payoutRatio := rmbAmt.ToDec().Quo(totalCollateral.ToDec())
	// FIXME: Consider leftovers.
	// FIXME: Accumulative errors could be large when the number of providers is large.
	for _, provider := range k.GetAllProviders(ctx) {
		// FIXME: Use truncate int or dec?
		payoutAmt := provider.Collateral.ToDec().Mul(payoutRatio).TruncateInt()

		if err := k.PayoutByProvider(ctx, provider.Address, payoutAmt); err != nil {
			panic(err)
		}
	}
	reimbursement := types.NewReimbursement(rmb, beneficiary, ctx.BlockTime().Add(k.GetClaimProposalParams(ctx).PayoutPeriod))
	k.SetReimbursement(ctx, proposalID, reimbursement)
	return nil
}

// PayoutByProvider undelegates delegations and send coins the the shield module.
func (k Keeper) PayoutByProvider(ctx sdk.Context, delAddr sdk.AccAddress, loss sdk.Int) error {
	delegations := k.sk.GetAllDelegatorDelegations(ctx, delAddr)
	totalDelAmountDec := sdk.ZeroDec()
	for _, del := range delegations {
		val, found := k.sk.GetValidator(ctx, del.GetValidatorAddr())
		if !found {
			panic("validator is not found")
		}
		totalDelAmountDec = totalDelAmountDec.Add(val.TokensFromShares(del.GetShares()))
	}

	// Start with bonded delegations.
	lossAmountDec := loss.ToDec()
	remainingDec := lossAmountDec
	for i := range delegations {
		val, found := k.sk.GetValidator(ctx, delegations[i].GetValidatorAddr())
		if !found {
			panic("validator is not found")
		}
		delAmountDec := val.TokensFromShares(delegations[i].GetShares())
		var ubdAmountDec sdk.Dec
		if totalDelAmountDec.GT(lossAmountDec) {
			if i == len(delegations)-1 {
				ubdAmountDec = remainingDec
			} else {
				ubdAmountDec = lossAmountDec.Mul(delAmountDec).Quo(totalDelAmountDec)
				remainingDec = remainingDec.Sub(ubdAmountDec)
			}
		} else {
			ubdAmountDec = delAmountDec
		}
		ubdShares, err := val.SharesFromTokens(ubdAmountDec.TruncateInt())
		if err != nil {
			panic(err)
		}
		k.UndelegateShares(ctx, delegations[i].DelegatorAddress, delegations[i].ValidatorAddress, ubdShares)
	}
	if totalDelAmountDec.GTE(lossAmountDec) {
		return nil
	}

	// If bonded delegations are not enough, track unbonding delegations.
	unbondingDelegations := k.GetSortedUnbondingDelegations(ctx, delAddr)
	for _, ubd := range unbondingDelegations {
		entry := 0
		if !remainingDec.IsPositive() {
			return nil
		}
		for i := range ubd.Entries {
			entry = i
			ubdAmountDec := ubd.Entries[i].InitialBalance.ToDec()
			if ubdAmountDec.GT(remainingDec) {
				overflowCoins := sdk.NewDecCoins(sdk.NewDecCoin(k.sk.BondDenom(ctx), ubdAmountDec.Sub(remainingDec).TruncateInt()))
				overflowMixedCoins := types.MixedDecCoins{Native: overflowCoins}
				k.AddRewards(ctx, delAddr, overflowMixedCoins)
				break
			}
			remainingDec = remainingDec.Sub(ubdAmountDec)
		}
		k.RedirectUnbondingEntryToShieldModule(ctx, ubd, entry)
	}

	if remainingDec.IsPositive() {
		panic("not enough bonded stake")
	}
	return nil
}

// UndelegateShares undelegates delegations of a delegator to a validator by shares.
func (k Keeper) UndelegateShares(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec) {
	delegation, found := k.sk.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		panic("delegation is not found")
	}
	k.sk.BeforeDelegationSharesModified(ctx, delAddr, valAddr)

	// Undelegate coins from the staking module account to the shield module account.
	validator, found := k.sk.GetValidator(ctx, valAddr)
	if !found {
		panic("validator was not found")
	}
	ubdCoins := sdk.NewCoins(sdk.NewCoin(k.sk.BondDenom(ctx), validator.TokensFromShares(shares).TruncateInt()))
	if err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, staking.BondedPoolName, types.ModuleName, ubdCoins); err != nil {
		panic(err)
	}

	// Update delegation records.
	delegation.Shares = delegation.Shares.Sub(shares)
	k.sk.SetDelegation(ctx, delegation)
	k.sk.AfterDelegationModified(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)

	// Update the validator.
	k.sk.RemoveValidatorTokensAndShares(ctx, validator, shares)
}

// GetSortedUnbondingDelegations gets unbonding delegations sorted by completion time.
func (k Keeper) GetSortedUnbondingDelegations(ctx sdk.Context, delAddr sdk.AccAddress) []staking.UnbondingDelegation {
	ubds := k.sk.GetAllUnbondingDelegations(ctx, delAddr)
	var unbondingDelegations []staking.UnbondingDelegation
	for _, ubd := range ubds {
		for _, entry := range ubd.Entries {
			unbondingDelegations = append(
				unbondingDelegations,
				types.NewUnbondingDelegation(ubd.DelegatorAddress, ubd.ValidatorAddress, entry),
			)
		}
	}
	sort.SliceStable(unbondingDelegations, func(i, j int) bool {
		return unbondingDelegations[i].Entries[0].CompletionTime.After(unbondingDelegations[j].Entries[0].CompletionTime)
	})
	return unbondingDelegations
}

// RedirectUnbondingEntryToShieldModule changes the recipient of an unbonding delegation to the shield module.
func (k Keeper) RedirectUnbondingEntryToShieldModule(ctx sdk.Context, ubd staking.UnbondingDelegation, endIndex int) {
	delAddr := ubd.DelegatorAddress
	valAddr := ubd.ValidatorAddress
	shieldAddr := k.supplyKeeper.GetModuleAddress(types.ModuleName)

	// Iterate through entries and add it to shield unbonding entries.
	shieldUbd, _ := k.sk.GetUnbondingDelegation(ctx, shieldAddr, valAddr)
	for _, entry := range ubd.Entries[:endIndex+1] {
		shieldUbd.AddEntry(entry.CreationHeight, entry.CompletionTime, entry.Balance)
		timeSlice := k.sk.GetUBDQueueTimeSlice(ctx, entry.CompletionTime)
		for i := 0; i < len(timeSlice); i++ {
			if timeSlice[i].DelegatorAddress.Equals(delAddr) && timeSlice[i].ValidatorAddress.Equals(valAddr) {
				timeSlice = append(timeSlice[:i], timeSlice[i+1:]...)
				k.sk.SetUBDQueueTimeSlice(ctx, entry.CompletionTime, timeSlice)
				break
			}
		}
		k.sk.InsertUBDQueue(ctx, shieldUbd, entry.CompletionTime)
	}
	ubd.Entries = ubd.Entries[endIndex+1:]
	k.sk.SetUnbondingDelegation(ctx, ubd)
	k.sk.SetUnbondingDelegation(ctx, shieldUbd)
}

// WithdrawReimbursement withdraws a reimbursement made for a beneficiary.
func (k Keeper) WithdrawReimbursement(ctx sdk.Context, proposalID uint64, beneficiary sdk.AccAddress) (sdk.Coins, error) {
	reimbursement, err := k.GetReimbursement(ctx, proposalID)
	if err != nil {
		return sdk.Coins{}, err
	}

	// check beneficiary and time
	if !reimbursement.Beneficiary.Equals(beneficiary) {
		return sdk.Coins{}, types.ErrInvalidBeneficiary
	}
	if reimbursement.PayoutTime.After(ctx.BlockTime()) {
		return sdk.Coins{}, types.ErrNotPayoutTime
	}

	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, beneficiary, reimbursement.Amount); err != nil {
		return sdk.Coins{}, types.ErrNotPayoutTime
	}
	if err := k.DeleteReimbursement(ctx, proposalID); err != nil {
		return sdk.Coins{}, err
	}
	return reimbursement.Amount, nil
}
