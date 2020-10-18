package keeper

import (
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// ClaimLock locks collaterals after a claim proposal is submitted.
func (k Keeper) ClaimLock(ctx sdk.Context, proposalID uint64, poolID uint64, purchaser sdk.AccAddress, purchaseID uint64, loss sdk.Coins, lockPeriod time.Duration) error {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.ErrNoPoolFound
	}
	if !pool.Shield.IsAllGTE(loss) {
		panic(types.ErrNotEnoughShield)
	}

	lossAmt := loss.AmountOf(k.sk.BondDenom(ctx))
	poolValidCollateral := sdk.ZeroInt()
	collaterals := k.GetAllPoolCollaterals(ctx, pool)
	for _, collateral := range collaterals {
		if collateral.Amount.IsPositive() {
			poolValidCollateral = poolValidCollateral.Add(collateral.Amount)
		}
	}
	if poolValidCollateral.LT(lossAmt) {
		panic(types.ErrNotEnoughCollateral)
	}

	// Update shield amount and delete time of the purchase.
	purchaseList, found := k.GetPurchaseList(ctx, poolID, purchaser)
	if !found {
		return types.ErrPurchaseNotFound
	}
	var index int
	for i, entry := range purchaseList.Entries {
		if entry.PurchaseID == purchaseID {
			index = i
			break
		}
	}
	k.DequeuePurchase(ctx, purchaseList, purchaseList.Entries[index].DeleteTime)

	if !purchaseList.Entries[index].Shield.IsAllGTE(loss) {
		return types.ErrNotEnoughShield
	}
	purchaseList.Entries[index].Shield = purchaseList.Entries[index].Shield.Sub(loss)
	votingEndTime := ctx.BlockTime().Add(lockPeriod)
	if purchaseList.Entries[index].DeleteTime.Before(votingEndTime) {
		purchaseList.Entries[index].DeleteTime = votingEndTime
	}
	k.SetPurchaseList(ctx, purchaseList)
	k.InsertPurchaseQueue(ctx, purchaseList, purchaseList.Entries[index].DeleteTime)

	// Update locked collaterals for community members and shield admin.
	proportionDec := lossAmt.ToDec().Quo(poolValidCollateral.ToDec())
	remaining := lossAmt
	for i := range collaterals {
		if !collaterals[i].Amount.IsPositive() {
			continue
		}
		var lockAmt sdk.Int
		if i < len(collaterals)-1 {
			lockAmt = collaterals[i].Amount.ToDec().Mul(proportionDec).TruncateInt()
			if lockAmt.LT(collaterals[i].Amount) && lockAmt.LT(remaining) {
				lockAmt = lockAmt.Add(sdk.OneInt())
			}
			remaining = remaining.Sub(lockAmt)
		} else {
			lockAmt = remaining
		}
		lockedCollateral := types.NewLockedCollateral(proposalID, lockAmt)
		collaterals[i].LockedCollaterals = append(collaterals[i].LockedCollaterals, lockedCollateral)
		collaterals[i].Amount = collaterals[i].Amount.Sub(lockAmt)
		collaterals[i].TotalLocked = collaterals[i].TotalLocked.Add(lockAmt)
		k.SetCollateral(ctx, pool, collaterals[i].Provider, collaterals[i])
		k.LockProvider(ctx, collaterals[i].Provider, lockAmt, lockPeriod)
	}

	// Update the pool.
	pool.Shield = pool.Shield.Sub(loss)
	pool.TotalCollateral = pool.TotalCollateral.Sub(lossAmt)
	pool.TotalLocked = pool.TotalLocked.Add(lossAmt)
	k.SetPool(ctx, pool)

	return nil
}

// LockProvider checks if delegations of an account can cover the loss.
// It modifies unbonding time if the totals delegations cannot cover the loss.
func (k Keeper) LockProvider(ctx sdk.Context, delAddr sdk.AccAddress, amount sdk.Int, lockPeriod time.Duration) {
	provider, found := k.GetProvider(ctx, delAddr)
	if !found {
		panic(types.ErrProviderNotFound)
	}

	// Update the provider.
	provider.TotalLocked = provider.TotalLocked.Add(amount)
	provider.Collateral = provider.Collateral.Sub(amount)
	k.SetProvider(ctx, delAddr, provider)

	// If there are enough delegations backing collaterals, do nothing.
	if provider.DelegationBonded.GTE(provider.TotalLocked) {
		return
	}

	// If there are not enough delegations, extend unbondings.
	unbondingDelegations := k.GetSortedUnbondingDelegations(ctx, delAddr)
	remaining := provider.TotalLocked.Sub(provider.DelegationBonded)
	endTime := ctx.BlockTime().Add(lockPeriod)
	for _, ubd := range unbondingDelegations {
		if !remaining.IsPositive() {
			return
		}
		entry := ubd.Entries[0]
		if entry.CompletionTime.Before(endTime) {
			// Change unbonding completion time.
			timeSlice := k.sk.GetUBDQueueTimeSlice(ctx, entry.CompletionTime)
			if len(timeSlice) > 1 {
				for i := 0; i < len(timeSlice); i++ {
					if timeSlice[i].DelegatorAddress.Equals(delAddr) && timeSlice[i].ValidatorAddress.Equals(ubd.ValidatorAddress) {
						timeSlice = append(timeSlice[:i], timeSlice[i+1:]...)
						k.sk.SetUBDQueueTimeSlice(ctx, entry.CompletionTime, timeSlice)
						break
					}
				}
			} else {
				k.sk.RemoveUBDQueue(ctx, entry.CompletionTime)
			}

			unbonding, found := k.sk.GetUnbondingDelegation(ctx, ubd.DelegatorAddress, ubd.ValidatorAddress)
			if !found {
				panic("unbonding delegation was not found")
			}
			found = false
			for i := 0; i < len(unbonding.Entries); i++ {
				if !found && unbonding.Entries[i].CreationHeight == entry.CreationHeight && unbonding.Entries[i].InitialBalance.Equal(entry.InitialBalance) {
					unbonding.Entries[i].CompletionTime = endTime
					found = true
				} else if found && unbonding.Entries[i].CompletionTime.Before(unbonding.Entries[i-1].CompletionTime) {
					unbonding.Entries[i-1], unbonding.Entries[i] = unbonding.Entries[i], unbonding.Entries[i-1]
				} else if found {
					break
				}
			}
			k.sk.SetUnbondingDelegation(ctx, unbonding)
			k.sk.InsertUBDQueue(ctx, unbonding, endTime)
		}
		remaining = remaining.Sub(entry.Balance)
	}
	if remaining.IsPositive() {
		panic("not enough bonded and unbonding delegations")
	}
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

// ClaimUnlock unlocks locked collaterals.
func (k Keeper) ClaimUnlock(ctx sdk.Context, proposalID uint64, poolID uint64, loss sdk.Coins) error {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.ErrNoPoolFound
	}
	lossAmt := loss.AmountOf(k.sk.BondDenom(ctx))
	pool.TotalCollateral = pool.TotalCollateral.Add(lossAmt)
	pool.TotalLocked = pool.TotalLocked.Sub(lossAmt)
	k.SetPool(ctx, pool)

	// Update collaterals and providers.
	collaterals := k.GetAllPoolCollaterals(ctx, pool)
	for _, collateral := range collaterals {
		for j := range collateral.LockedCollaterals {
			if collateral.LockedCollaterals[j].ProposalID == proposalID {
				lockedAmount := collateral.LockedCollaterals[j].Amount
				provider, found := k.GetProvider(ctx, collateral.Provider)
				if !found {
					panic(types.ErrProviderNotFound)
				}
				provider.Collateral = provider.Collateral.Add(lockedAmount)
				provider.TotalLocked = provider.TotalLocked.Sub(lockedAmount)
				k.SetProvider(ctx, collateral.Provider, provider)

				collateral.Amount = collateral.Amount.Add(lockedAmount)
				collateral.TotalLocked = collateral.TotalLocked.Sub(lockedAmount)
				collateral.LockedCollaterals = append(collateral.LockedCollaterals[:j], collateral.LockedCollaterals[j+1:]...)
				k.SetCollateral(ctx, pool, collateral.Provider, collateral)
				break
			}
		}
	}

	return nil
}

// RestoreShield restores shield for proposer.
func (k Keeper) RestoreShield(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, id uint64, loss sdk.Coins) error {
	// Update shield of the pool.
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.ErrNoPoolFound
	}
	pool.Shield = pool.Shield.Add(loss...)
	k.SetPool(ctx, pool)

	// Update shield of the purchase.
	purchaseList, found := k.GetPurchaseList(ctx, poolID, purchaser)
	if !found {
		return types.ErrPurchaseNotFound
	}
	for i := range purchaseList.Entries {
		if purchaseList.Entries[i].PurchaseID == id {
			purchaseList.Entries[i].Shield = purchaseList.Entries[i].Shield.Add(loss...)
			break
		}
	}
	k.SetPurchaseList(ctx, purchaseList)

	return nil
}

// UndelegateCoinsToShieldModule undelegates delegations and send coins the the shield module.
func (k Keeper) UndelegateCoinsToShieldModule(ctx sdk.Context, delAddr sdk.AccAddress, loss sdk.Int) error {
	delegations := k.sk.GetAllDelegatorDelegations(ctx, delAddr)
	totalDelAmountDec := sdk.ZeroDec()
	for _, del := range delegations {
		val, found := k.sk.GetValidator(ctx, del.GetValidatorAddr())
		if !found {
			panic("validator was not found")
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
				// FIXME not a good way to go maybe?
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

// SetReimbursement sets a reimbursement in store.
func (k Keeper) SetReimbursement(ctx sdk.Context, proposalID uint64, payout types.Reimbursement) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(payout)
	store.Set(types.GetReimbursementKey(proposalID), bz)
}

// GetReimbursement get a reimbursement in store.
func (k Keeper) GetReimbursement(ctx sdk.Context, proposalID uint64) (types.Reimbursement, error) {
	store := ctx.KVStore(k.storeKey)
	if bz := store.Get(types.GetReimbursementKey(proposalID)); bz != nil {
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
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.ErrNoPoolFound
	}
	pool.TotalLocked = pool.TotalLocked.Sub(rmb.AmountOf(k.BondDenom(ctx)))

	// For each community member, get coins from delegations.
	poolTotal := sdk.ZeroInt()
	collaterals := k.GetAllPoolCollaterals(ctx, pool)
	for _, collateral := range collaterals {
		for j := range collateral.LockedCollaterals {
			if collateral.LockedCollaterals[j].ProposalID == proposalID {
				lockedAmount := collateral.LockedCollaterals[j].Amount
				provider, found := k.GetProvider(ctx, collateral.Provider)
				if !found {
					panic(types.ErrProviderNotFound)
				}
				provider.TotalLocked = provider.TotalLocked.Sub(lockedAmount)
				k.SetProvider(ctx, collateral.Provider, provider)

				if err := k.UndelegateCoinsToShieldModule(ctx, collateral.Provider, lockedAmount); err != nil {
					panic(err)
				}
				collateral.TotalLocked = collateral.TotalLocked.Sub(lockedAmount)
				if collateral.Amount.Add(collateral.TotalLocked).IsPositive() {
					poolTotal = poolTotal.Add(collateral.Amount.Add(collateral.TotalLocked))
				}
				collateral.LockedCollaterals = append(collateral.LockedCollaterals[:j], collateral.LockedCollaterals[j+1:]...)
				k.SetCollateral(ctx, pool, collateral.Provider, collateral)
				break
			}
		}
	}
	pool.TotalCollateral = poolTotal.Sub(pool.TotalLocked)
	k.SetPool(ctx, pool)

	proposalParams := k.GetClaimProposalParams(ctx)
	reimbursement := types.NewReimbursement(rmb, beneficiary, ctx.BlockTime().Add(proposalParams.PayoutPeriod))
	k.SetReimbursement(ctx, proposalID, reimbursement)
	return nil
}

// WithdrawReimbursement checks a reimbursement and pays the beneficiary.
func (k Keeper) WithdrawReimbursement(ctx sdk.Context, proposalID uint64, beneficiary sdk.AccAddress) (sdk.Coins, error) {
	reimbursement, err := k.GetReimbursement(ctx, proposalID)
	if err != nil {
		return sdk.Coins{}, err
	}

	// Check beneficiary and time.
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
