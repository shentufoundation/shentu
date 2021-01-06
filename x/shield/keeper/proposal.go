package keeper

import (
	"encoding/binary"
	"fmt"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	vestexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	vesting "github.com/certikfoundation/shentu/x/auth/types"
	"github.com/certikfoundation/shentu/x/shield/types"
)

// SecureCollaterals is called after a claim is submitted to secure
// the given amount of collaterals for the duration and adjust shield
// module states accordingly.
func (k Keeper) SecureCollaterals(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, purchaseID uint64, loss sdk.Coins, duration time.Duration) error {
	lossAmt := loss.AmountOf(k.sk.BondDenom(ctx))

	// Verify shield.
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.ErrNoPoolFound
	}
	if lossAmt.GT(pool.Shield) {
		return types.ErrNotEnoughShield
	}

	// Verify collateral availability.
	totalCollateral := k.GetTotalCollateral(ctx)
	totalSecureAmt := k.GetTotalClaimed(ctx).Add(lossAmt)
	if totalSecureAmt.GT(totalCollateral) {
		panic("total secure amount surpassed total collateral")
	}

	// Verify purchase.
	purchaseList, found := k.GetPurchaseList(ctx, poolID, purchaser)
	if !found {
		return types.ErrPurchaseNotFound
	}
	var index int
	for i, entry := range purchaseList.Entries {
		if entry.PurchaseId == purchaseID {
			index = i
			break
		}
	}
	purchase := &purchaseList.Entries[index]
	if lossAmt.GT(purchase.Shield) {
		return types.ErrNotEnoughShield
	}

	// Secure the updated loss ratio from each provider to cover total claimed.
	providers := k.GetAllProviders(ctx)
	claimedRatio := totalSecureAmt.ToDec().Quo(totalCollateral.ToDec())
	remaining := totalSecureAmt
	for i := range providers {
		secureAmt := sdk.MinInt(providers[i].Collateral.ToDec().Mul(claimedRatio).TruncateInt(), remaining)

		// Require each provider to secure one more unit, if possible,
		// so that the last provider does not have to cover combined
		// truncated amounts.
		if secureAmt.LT(remaining) && secureAmt.LT(providers[i].Collateral) {
			secureAmt = secureAmt.Add(sdk.OneInt())
		}
		k.SecureFromProvider(ctx, providers[i], secureAmt, duration)
		remaining = remaining.Sub(secureAmt)
	}

	// Update purchase states.
	purchase.Shield = purchase.Shield.Sub(lossAmt)
	votingEndTime := ctx.BlockTime().Add(duration)
	if purchase.DeletionTime.Before(votingEndTime) {
		purchase.DeletionTime = votingEndTime
	}
	k.SetPurchaseList(ctx, purchaseList)

	// Update pool and global pool states.
	pool.Shield = pool.Shield.Sub(lossAmt)
	k.SetPool(ctx, pool)

	totalShield := k.GetTotalShield(ctx)
	totalShield = totalShield.Sub(lossAmt)
	k.SetTotalShield(ctx, totalShield)
	k.SetTotalClaimed(ctx, totalSecureAmt)

	return nil
}

// SecureFromProvider secures the specified amount of collaterals from
// the provider for the duration. If necessary, it extends withdrawing
// collaterals and, if exist, their linked unbondings as well.
func (k Keeper) SecureFromProvider(ctx sdk.Context, provider types.Provider, amount sdk.Int, duration time.Duration) {
	providerAddr, err := sdk.AccAddressFromBech32(provider.Address)
	if err != nil {
		panic(err)
	}

	// Lenient check:
	// We are done if non-withdrawing, bonded delegation-backed
	// collaterals can cover the amount.
	if provider.Collateral.Sub(provider.Withdrawing).GTE(amount) && provider.DelegationBonded.GTE(amount) {
		return
	}

	// Secure the given amount of collaterals until the end of the
	// lock period by delaying withdrawals, if necessary.
	// availableCollateralByEndTime = ProviderTotalCollateral - WithdrawnByEndTime
	endTime := ctx.BlockTime().Add(duration)
	availableCollateralByEndTime := provider.Collateral.Sub(k.ComputeWithdrawAmountByTime(ctx, provider.Address, endTime))

	// Secure the given amount of staking (bonded or unbonding) until
	// the end of the lock period by delaying unbondings, if necessary.
	// availableDelegationByEndTime = Bonded + Unbonding - UnbondedByEndTime
	availableDelegationByEndTime := provider.DelegationBonded.Add(k.ComputeTotalUnbondingAmount(ctx, providerAddr).Sub(k.ComputeUnbondingAmountByTime(ctx, providerAddr, endTime)))

	// Collaterals that won't be withdrawn until the end time must be
	// backed by staking that won't be unbonded until the end time.
	if !availableCollateralByEndTime.LTE(availableDelegationByEndTime) {
		panic("notWithdrawnSoon must be less than or equal to notUnbondedSoon")
	}

	if amount.GT(availableCollateralByEndTime) {
		withdrawDelayAmt := amount.Sub(availableCollateralByEndTime)
		k.DelayWithdraws(ctx, provider.Address, withdrawDelayAmt, endTime)
		if amount.GT(availableDelegationByEndTime) {
			unbondingDelayAmt := amount.Sub(availableDelegationByEndTime)
			k.DelayUnbonding(ctx, providerAddr, unbondingDelayAmt, endTime)
		}
	}
}

// ClaimEnd ends a claim process by updating the total claimed amount.
func (k Keeper) ClaimEnd(ctx sdk.Context, id, poolID uint64, loss sdk.Coins) {
	lossAmt := loss.AmountOf(k.sk.BondDenom(ctx))
	totalClaimed := k.GetTotalClaimed(ctx).Sub(lossAmt)
	k.SetTotalClaimed(ctx, totalClaimed)
}

// RestoreShield restores shield-related states as they were prior to
// the claim proposal submission.
func (k Keeper) RestoreShield(ctx sdk.Context, poolID uint64, purchaser sdk.AccAddress, id uint64, loss sdk.Coins) error {
	lossAmt := loss.AmountOf(k.sk.BondDenom(ctx))

	// Update the total shield.
	totalShield := k.GetTotalShield(ctx).Add(lossAmt)
	k.SetTotalShield(ctx, totalShield)

	// Update shield of the pool.
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.ErrNoPoolFound
	}
	pool.Shield = pool.Shield.Add(lossAmt)
	k.SetPool(ctx, pool)

	// Update shield of the purchase.
	purchaseList, found := k.GetPurchaseList(ctx, poolID, purchaser)
	if !found {
		return types.ErrPurchaseNotFound
	}
	for i := range purchaseList.Entries {
		if purchaseList.Entries[i].PurchaseId == id {
			purchaseList.Entries[i].Shield = purchaseList.Entries[i].Shield.Add(lossAmt)
			break
		}
	}
	k.SetPurchaseList(ctx, purchaseList)

	return nil
}

// SetReimbursement sets a reimbursement in store.
func (k Keeper) SetReimbursement(ctx sdk.Context, proposalID uint64, payout types.Reimbursement) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(&payout)
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
	return types.Reimbursement{}, types.ErrReimbursementNotFound
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

// IterateReimbursements iterates through all reimbursements.
func (k Keeper) IterateReimbursements(ctx sdk.Context, callback func(rmb types.Reimbursement) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ReimbursementKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rmb types.Reimbursement
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &rmb)

		if callback(rmb) {
			break
		}
	}
}

// GetAllReimbursements retrieves all reimbursements.
func (k Keeper) GetAllReimbursements(ctx sdk.Context) (rmbs []types.Reimbursement) {
	k.IterateReimbursements(ctx, func(rmb types.Reimbursement) bool {
		rmbs = append(rmbs, rmb)
		return false
	})
	return
}

// GetAllProposalIDReimbursementPairs retrieves all proposal ID and reimbursement pairs.
func (k Keeper) GetAllProposalIDReimbursementPairs(ctx sdk.Context) []types.ProposalIDReimbursementPair {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ReimbursementKey)

	var pRPairs []types.ProposalIDReimbursementPair
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID := binary.LittleEndian.Uint64(iterator.Key()[len(types.ReimbursementKey):])
		var reimbursement types.Reimbursement
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &reimbursement)

		pRPairs = append(pRPairs, types.NewProposalIDReimbursementPair(proposalID, reimbursement))
	}
	return pRPairs
}

// CreateReimbursement creates a reimbursement.
func (k Keeper) CreateReimbursement(ctx sdk.Context, proposalID uint64, amount sdk.Coins, beneficiary sdk.AccAddress) error {
	bondDenom := k.BondDenom(ctx)
	totalCollateral := k.GetTotalCollateral(ctx)
	totalPurchased := k.GetTotalShield(ctx)
	totalPayout := amount.AmountOf(bondDenom)
	purchaseRatio := totalPurchased.ToDec().Quo(totalCollateral.ToDec())
	payoutRatio := totalPayout.ToDec().Quo(totalCollateral.ToDec())
	for _, provider := range k.GetAllProviders(ctx) {
		if !totalPayout.IsPositive() {
			break
		}

		providerAddr, err := sdk.AccAddressFromBech32(provider.Address)
		if err != nil {
			panic(err)
		}

		purchased := provider.Collateral.ToDec().Mul(purchaseRatio).TruncateInt()
		if purchased.GT(totalPurchased) {
			purchased = totalPurchased
		}
		payout := provider.Collateral.ToDec().Mul(payoutRatio).TruncateInt()
		if payout.GT(totalPayout) {
			payout = totalPayout
		}

		// Require providers to cover (purchased + 1) and (payout + 1) if it's possible,
		// so that the last provider will not be asked to cover all truncated amount.
		if purchased.LT(totalPurchased) && provider.Collateral.GT(payout.Add(purchased)) {
			purchased = purchased.Add(sdk.OneInt())
		}
		if payout.LT(totalPayout) && provider.Collateral.GT(payout.Add(purchased)) {
			payout = payout.Add(sdk.OneInt())
		}

		if err := k.UpdateProviderCollateralForPayout(ctx, providerAddr, purchased, payout); err != nil {
			panic(err)
		}

		if err := k.MakePayoutByProviderDelegations(ctx, providerAddr, purchased, payout); err != nil {
			panic(err)
		}

		totalPurchased = totalPurchased.Sub(purchased)
		totalPayout = totalPayout.Sub(payout)
	}
	if totalPayout.IsPositive() {
		panic("not enough payout made")
	}
	reimbursement := types.NewReimbursement(amount, beneficiary, ctx.BlockTime().Add(k.GetClaimProposalParams(ctx).PayoutPeriod))
	k.SetReimbursement(ctx, proposalID, reimbursement)

	totalCollateral = totalCollateral.Sub(amount.AmountOf(bondDenom))
	totalClaimed := k.GetTotalClaimed(ctx)
	totalClaimed = totalClaimed.Sub(amount.AmountOf(bondDenom))
	k.SetTotalCollateral(ctx, totalCollateral)
	k.SetTotalClaimed(ctx, totalClaimed)

	return nil
}

// UpdateProviderCollateralForPayout updates a provider's collateral and withdraws according to the payout.
func (k Keeper) UpdateProviderCollateralForPayout(ctx sdk.Context, providerAddr sdk.AccAddress, purchased, payout sdk.Int) error {
	provider, found := k.GetProvider(ctx, providerAddr)
	if !found {
		return types.ErrProviderNotFound
	}
	totalWithdrawing := k.GetTotalWithdrawing(ctx)

	uncoveredPurchase := sdk.ZeroInt()
	payoutFromCollateral := sdk.ZeroInt()
	if provider.Collateral.Sub(provider.Withdrawing).GTE(purchased.Add(payout)) {
		// If collateral - withdraw >= purchased + payout:
		//     purchased       payout
		//   ----------------|--------|
		//       collateral - withdraw                withdraw
		// -------------------------------|---------------------------------
		payoutFromCollateral = payout
	} else if provider.Collateral.Sub(provider.Withdrawing).GTE(purchased) {
		// If purchased <= collateral - withdraw < purchased + payout:
		//               purchased       payout
		//             ----------------|--------|
		//       collateral - withdraw                withdraw
		// -------------------------------|---------------------------------
		payoutFromCollateral = provider.Collateral.Sub(provider.Withdrawing).Sub(purchased)
	} else {
		// If collateral - withdraw < purchased:
		//                      purchased       payout
		//                    ----------------|--------|
		//       collateral - withdraw                withdraw
		// -------------------------------|---------------------------------
		uncoveredPurchase = purchased.Sub(provider.Collateral.Sub(provider.Withdrawing))
	}
	payoutFromWithdraw := payout.Sub(payoutFromCollateral)

	// Update provider's collateral and total withdraw.
	provider.Collateral = provider.Collateral.Sub(payout)
	provider.Withdrawing = provider.Withdrawing.Sub(payoutFromWithdraw)
	totalWithdrawing = totalWithdrawing.Sub(payoutFromWithdraw)

	// Update provider's withdraws from latest to oldest.
	withdraws := k.GetWithdrawsByProvider(ctx, provider.Address)
	for i := len(withdraws) - 1; i >= 0 && payoutFromWithdraw.IsPositive(); i-- {
		// If purchased is not fully covered, cover purchased first.
		remainingWithdraw := sdk.MaxInt(withdraws[i].Amount.Sub(uncoveredPurchase), sdk.ZeroInt())
		uncoveredPurchase = sdk.MaxInt(uncoveredPurchase.Sub(withdraws[i].Amount), sdk.ZeroInt())
		if remainingWithdraw.IsZero() {
			continue
		}

		// Update the withdraw based on payout after purchased is fully covered.
		payoutFromThisWithdraw := sdk.MinInt(payoutFromWithdraw, remainingWithdraw)
		payoutFromWithdraw = payoutFromWithdraw.Sub(payoutFromThisWithdraw)
		timeSlice := k.GetWithdrawQueueTimeSlice(ctx, withdraws[i].CompletionTime)
		for j := range timeSlice {
			if timeSlice[j].Address != withdraws[i].Address || !timeSlice[j].Amount.Equal(withdraws[i].Amount) {
				continue
			}

			if withdraws[i].Amount.Equal(payoutFromThisWithdraw) {
				if len(timeSlice) == 1 {
					k.RemoveTimeSliceFromWithdrawQueue(ctx, withdraws[i].CompletionTime)
				} else {
					timeSlice = append(timeSlice[:j], timeSlice[j+1:]...)
					k.SetWithdrawQueueTimeSlice(ctx, withdraws[i].CompletionTime, timeSlice)
				}
			} else {
				timeSlice[j].Amount = withdraws[i].Amount.Sub(payoutFromThisWithdraw)
				k.SetWithdrawQueueTimeSlice(ctx, withdraws[i].CompletionTime, timeSlice)
			}
			break
		}
	}
	if !payoutFromWithdraw.IsZero() {
		panic("exact payout was not made from withdrawals")
	}

	k.SetTotalWithdrawing(ctx, totalWithdrawing)
	k.SetProvider(ctx, providerAddr, provider)

	return nil
}

// MakePayoutByProviderDelegations undelegates the provider's delegations and transfers tokens from the staking module account to the shield module account.
func (k Keeper) MakePayoutByProviderDelegations(ctx sdk.Context, providerAddr sdk.AccAddress, purchased, payout sdk.Int) error {
	provider, found := k.GetProvider(ctx, providerAddr)
	if !found {
		return types.ErrProviderNotFound
	}

	uncoveredPurchase := sdk.ZeroInt()
	payoutFromDelegation := sdk.ZeroInt()
	if provider.DelegationBonded.GTE(purchased.Add(payout)) {
		// If delegation >= purchased + payout:
		//     purchased       payout
		//   ----------------|--------|
		//            delegations                     unbondings
		// -------------------------------|---------------------------------
		payoutFromDelegation = payout
	} else if provider.DelegationBonded.GTE(purchased) {
		// If purchased <= delegation < purchased + payout:
		//               purchased       payout
		//             ----------------|--------|
		//            delegations                     unbondings
		// -------------------------------|---------------------------------
		payoutFromDelegation = provider.DelegationBonded.Sub(purchased)
	} else {
		// If delegation < purchased:
		//                      purchased       payout
		//                    ----------------|--------|
		//            delegations                     unbondings
		// -------------------------------|---------------------------------
		uncoveredPurchase = purchased.Sub(provider.DelegationBonded)
	}
	payoutFromUnbonding := payout.Sub(payoutFromDelegation)

	if payoutFromDelegation.IsPositive() {
		k.PayFromDelegation(ctx, providerAddr, payoutFromDelegation)
	}

	if payoutFromUnbonding.IsZero() {
		return nil
	}

	unbondingDelegations := k.GetSortedUnbondingDelegations(ctx, providerAddr)
	for _, ubd := range unbondingDelegations {
		if !payoutFromUnbonding.IsPositive() {
			break
		}
		entry := ubd.Entries[0]

		// If purchased is not fully covered, cover purchased first.
		remainingUbd := sdk.MaxInt(entry.Balance.Sub(uncoveredPurchase), sdk.ZeroInt())
		uncoveredPurchase = sdk.MaxInt(uncoveredPurchase.Sub(entry.Balance), sdk.ZeroInt())
		if remainingUbd.IsZero() {
			continue
		}

		// Make payout after purchased is fully covered.
		payoutFromThisUbd := sdk.MinInt(payoutFromUnbonding, remainingUbd)
		k.PayFromUnbondings(ctx, ubd, payoutFromThisUbd)

		payoutFromUnbonding = payoutFromUnbonding.Sub(payoutFromThisUbd)
	}
	if !payoutFromUnbonding.IsZero() {
		panic("exact pay out was not made from unbondings")
	}

	return nil
}

// PayFromDelegation reduce provider's delegations and transfer tokens to the shield module account.
func (k Keeper) PayFromDelegation(ctx sdk.Context, delAddr sdk.AccAddress, payout sdk.Int) {
	provider, found := k.GetProvider(ctx, delAddr)
	if !found {
		panic(types.ErrProviderNotFound)
	}
	totalDelAmount := provider.DelegationBonded

	delegations := k.sk.GetAllDelegatorDelegations(ctx, delAddr)
	payoutRatio := payout.ToDec().Quo(totalDelAmount.ToDec())
	remaining := payout
	for i := range delegations {
		if !remaining.IsPositive() {
			return
		}

		val, found := k.sk.GetValidator(ctx, delegations[i].GetValidatorAddr())
		if !found {
			panic("validator is not found")
		}
		delAmount := val.TokensFromShares(delegations[i].GetShares()).TruncateInt()
		var ubdAmount sdk.Int
		if i == len(delegations)-1 {
			ubdAmount = remaining
		} else {
			ubdAmount = sdk.MinInt(payoutRatio.MulInt(delAmount).TruncateInt(), remaining)
			if ubdAmount.LT(remaining) && ubdAmount.LT(delAmount) {
				ubdAmount = ubdAmount.Add(sdk.OneInt())
			}
			remaining = remaining.Sub(ubdAmount)
		}
		ubdShares, err := val.SharesFromTokens(ubdAmount)
		if err != nil {
			panic(err)
		}

		delAddr, err := sdk.AccAddressFromBech32(delegations[i].DelegatorAddress)
		if err != nil {
			panic(err)
		}
		valAddr, err := sdk.ValAddressFromBech32(delegations[i].ValidatorAddress)
		if err != nil {
			panic(err)
		}
		k.UndelegateShares(ctx, delAddr, valAddr, ubdShares)
	}
}

// PayFromUnbondings reduce provider's unbonding delegations and transfer tokens to the shield module account.
func (k Keeper) PayFromUnbondings(ctx sdk.Context, ubd stakingtypes.UnbondingDelegation, payout sdk.Int) {
	delAddr, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	valAddr, err := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	unbonding, found := k.sk.GetUnbondingDelegation(ctx, delAddr, valAddr)
	if !found {
		panic("unbonding delegation is not found")
	}

	// Update unbonding delegations between the delegator and the validator.
	for i := range unbonding.Entries {
		if !unbonding.Entries[i].Balance.Equal(ubd.Entries[0].Balance) || !unbonding.Entries[i].CompletionTime.Equal(ubd.Entries[0].CompletionTime) {
			continue
		}

		if unbonding.Entries[i].Balance.Equal(payout) {
			// Update the unbonding queue and remove the entry.
			timeSlice := k.sk.GetUBDQueueTimeSlice(ctx, unbonding.Entries[i].CompletionTime)
			if len(timeSlice) > 1 {
				for j, slice := range timeSlice {
					if slice.DelegatorAddress == ubd.DelegatorAddress && slice.ValidatorAddress == ubd.ValidatorAddress {
						timeSlice = append(timeSlice[:j], timeSlice[j+1:]...)
						k.sk.SetUBDQueueTimeSlice(ctx, unbonding.Entries[i].CompletionTime, timeSlice)
						break
					}
				}
			} else {
				k.sk.RemoveUBDQueue(ctx, unbonding.Entries[i].CompletionTime)
			}
			unbonding.RemoveEntry(int64(i))
		} else {
			unbonding.Entries[i].Balance = unbonding.Entries[i].Balance.Sub(payout)
		}
		if len(unbonding.Entries) == 0 {
			k.sk.RemoveUnbondingDelegation(ctx, unbonding)
		} else {
			k.sk.SetUnbondingDelegation(ctx, unbonding)
		}
		break
	}

	// Transfer tokens from staking module's not bonded pool.
	payoutCoins := sdk.NewCoins(sdk.NewCoin(k.sk.BondDenom(ctx), payout))
	if err := k.UndelegateFromAccountToShieldModule(ctx, stakingtypes.NotBondedPoolName, delAddr, payoutCoins); err != nil {
		panic(err)
	}
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
		panic("validator is not found")
	}

	// Update delegation records.
	delegation.Shares = delegation.Shares.Sub(shares)

	isValidatorOperator := delegation.GetDelegatorAddr().Equals(validator.GetOperator())
	if isValidatorOperator && !validator.Jailed && validator.TokensFromShares(delegation.Shares).TruncateInt().LT(validator.MinSelfDelegation) {
		validator.Jailed = true
		k.sk.SetValidator(ctx, validator)
		k.sk.DeleteValidatorByPowerIndex(ctx, validator)
		validator, found = k.sk.GetValidator(ctx, valAddr)
		if !found {
			panic(fmt.Sprintf("validator record not found for address: %X\n", valAddr))
		}
	}

	if delegation.Shares.IsZero() {
		k.sk.RemoveDelegation(ctx, delegation)
	} else {
		k.sk.SetDelegation(ctx, delegation)
		k.sk.AfterDelegationModified(ctx, delegation.GetDelegatorAddr(), delegation.GetValidatorAddr())
	}

	validator, amount := k.sk.RemoveValidatorTokensAndShares(ctx, validator, shares)
	if validator.DelegatorShares.IsZero() && validator.IsUnbonded() {
		k.sk.RemoveValidator(ctx, validator.GetOperator())
	}

	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), amount))
	srcPool := stakingtypes.NotBondedPoolName
	if validator.IsBonded() {
		srcPool = stakingtypes.BondedPoolName
	}
	err := k.UndelegateFromAccountToShieldModule(ctx, srcPool, delAddr, coins)
	if err != nil {
		panic(err)
	}
}

// UndelegateFromAccountToShieldModule performs undelegations from a delegator's staking to the shield module.
func (k Keeper) UndelegateFromAccountToShieldModule(ctx sdk.Context, senderModule string, delAddr sdk.AccAddress, amt sdk.Coins) error {
	delAcc := k.ak.GetAccount(ctx, delAddr)
	if delAcc == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", delAddr)
	}

	if vacc, ok := delAcc.(vestexported.VestingAccount); ok {
		originalDelegatedVesting := vacc.GetDelegatedVesting()
		vacc.TrackUndelegation(amt)
		updatedDelegatedVesting := vacc.GetDelegatedVesting()
		updateAmt := originalDelegatedVesting.Sub(updatedDelegatedVesting)
		if mvacc, ok := delAcc.(*vesting.ManualVestingAccount); ok {
			var unlockAmt sdk.Coins
			if mvacc.OriginalVesting.Sub(mvacc.VestedCoins).IsAllGT(updateAmt) {
				unlockAmt = updateAmt
			} else {
				unlockAmt = mvacc.OriginalVesting.Sub(mvacc.VestedCoins)
			}
			mvacc.VestedCoins = mvacc.VestedCoins.Add(unlockAmt...)
		}
		k.ak.SetAccount(ctx, delAcc)
	}

	if err := k.bk.SendCoinsFromModuleToModule(ctx, senderModule, types.ModuleName, amt); err != nil {
		panic(err)
	}

	return nil
}

// GetSortedUnbondingDelegations gets unbonding delegations sorted by completion time from latest to earliest.
func (k Keeper) GetSortedUnbondingDelegations(ctx sdk.Context, delAddr sdk.AccAddress) []stakingtypes.UnbondingDelegation {
	ubds := k.sk.GetAllUnbondingDelegations(ctx, delAddr)
	var unbondingDelegations []stakingtypes.UnbondingDelegation
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

// WithdrawReimbursement withdraws a reimbursement made for a beneficiary.
func (k Keeper) WithdrawReimbursement(ctx sdk.Context, proposalID uint64, beneficiary sdk.AccAddress) (sdk.Coins, error) {
	reimbursement, err := k.GetReimbursement(ctx, proposalID)
	if err != nil {
		return sdk.Coins{}, err
	}

	// check beneficiary and time
	if reimbursement.Beneficiary != beneficiary.String() {
		return sdk.Coins{}, types.ErrInvalidBeneficiary
	}
	if reimbursement.PayoutTime.After(ctx.BlockTime()) {
		return sdk.Coins{}, types.ErrNotPayoutTime
	}

	if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, beneficiary, reimbursement.Amount); err != nil {
		return sdk.Coins{}, types.ErrNotPayoutTime
	}
	if err := k.DeleteReimbursement(ctx, proposalID); err != nil {
		return sdk.Coins{}, err
	}
	return reimbursement.Amount, nil
}
