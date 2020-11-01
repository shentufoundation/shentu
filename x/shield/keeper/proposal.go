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
	totalPurchased := k.GetTotalShield(ctx)
	totalPayout := rmb.AmountOf(k.BondDenom(ctx))
	purchaseRatio := totalPurchased.ToDec().Quo(totalCollateral.ToDec())
	payoutRatio := totalPayout.ToDec().Quo(totalCollateral.ToDec())
	for _, provider := range k.GetAllProviders(ctx) {
		if !totalPayout.IsPositive() {
			break
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
		if err := k.MakePayoutByProvider(ctx, provider.Address, purchased, payout); err != nil {
			panic(err)
		}
		totalPurchased = totalPurchased.Sub(purchased)
		totalPayout = totalPayout.Sub(payout)
	}
	reimbursement := types.NewReimbursement(rmb, beneficiary, ctx.BlockTime().Add(k.GetClaimProposalParams(ctx).PayoutPeriod))
	k.SetReimbursement(ctx, proposalID, reimbursement)
	return nil
}

// PayFromDelegation reduce provider's delegations and transfer tokens to the shield module account.
func (k Keeper) PayFromDelegation(ctx sdk.Context, providerAddr sdk.AccAddress, payout sdk.Int) {
	delegations := k.sk.GetAllDelegatorDelegations(ctx, providerAddr)
	totalDelAmountDec := sdk.ZeroDec()
	for _, del := range delegations {
		val, found := k.sk.GetValidator(ctx, del.GetValidatorAddr())
		if !found {
			panic("validator is not found")
		}
		totalDelAmountDec = totalDelAmountDec.Add(val.TokensFromShares(del.GetShares()))
	}

	payoutDec := payout.ToDec()
	remainingDec := payoutDec
	for i := range delegations {
		val, found := k.sk.GetValidator(ctx, delegations[i].GetValidatorAddr())
		if !found {
			panic("validator is not found")
		}
		delAmountDec := val.TokensFromShares(delegations[i].GetShares())
		var ubdAmountDec sdk.Dec
		if totalDelAmountDec.GT(payoutDec) {
			// FIXME: Corner case: not enough amount in the last delegation.
			if i == len(delegations)-1 {
				ubdAmountDec = remainingDec
			} else {
				ubdAmountDec = payoutDec.Mul(delAmountDec).Quo(totalDelAmountDec)
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
}

// PayFromUnbonding reduce provider's unbonding delegation and transfer tokens to the shield module account.
func (k Keeper) PayFromUnbonding(ctx sdk.Context, ubd staking.UnbondingDelegation, payout sdk.Int) {
	delAddr := ubd.DelegatorAddress
	valAddr := ubd.ValidatorAddress
	unbonding, found := k.sk.GetUnbondingDelegation(ctx, delAddr, valAddr)
	if !found {
		panic("unbonding delegation is not found")
	}

	// Update unbonding delegations between the delegator and the validator.
	for i := range unbonding.Entries {
		if unbonding.Entries[i].Balance.Equal(ubd.Entries[0].Balance) && unbonding.Entries[i].CompletionTime.Equal(ubd.Entries[0].CompletionTime) {
			unbonding.Entries[i].Balance = unbonding.Entries[i].Balance.Sub(payout)
			k.sk.SetUnbondingDelegation(ctx, unbonding)
			break
		}
	}

	// FIXME: Update the unbonding queue only if entry is removed.

	// Transfer tokens from the staking module to the shield module.
	payoutCoins := sdk.NewCoins(sdk.NewCoin(k.sk.BondDenom(ctx), payout))
	if err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, staking.NotBondedPoolName, types.ModuleName, payoutCoins); err != nil {
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

func (k Keeper) PayFromWithdraw(ctx sdk.Context, withdraw types.Withdraw, payout, purchased sdk.Int) {
	provider, found := k.GetProvider(ctx, withdraw.Address)
	if !found {
		panic(types.ErrProviderNotFound)
	}

	payoutFromDelegation := sdk.ZeroInt()
	payoutFromUnbonding := sdk.ZeroInt()
	if provider.DelegationBonded.GTE(purchased.Add(payout)) {
		// If delegation >= purchased + payout:
		//        |        withdraw      |
		//     purchased     | payout |
		//   -----|----------'--------'--|-----
		//            delegations                     unbondings
		// -------------------------------|---------------------------------
		payoutFromDelegation = payout
	} else if provider.DelegationBonded.GTE(purchased) {
		// If purchased <= delegation < purchased + payout:
		//                  |        withdraw      |
		//               purchased     | payout |
		//             -----|----------'--------'--|-----
		//            delegations                     unbondings
		// -------------------------------|---------------------------------
		payoutFromDelegation = provider.DelegationBonded.Sub(purchased)
		payoutFromUnbonding = payout.Sub(payoutFromDelegation)
	} else {
		// If delegation < purchased:
		//                         |        withdraw      |
		//                      purchased     | payout |
		//                    -----|----------'--------'--|-----
		//            delegations                     unbondings
		// -------------------------------|---------------------------------
		// or
		//                                   |        withdraw      |
		//                                purchased     | payout |
		//                              -----|----------'--------'--|-----
		//            delegations                     unbondings
		// -------------------------------|---------------------------------
		payoutFromUnbonding = payout
	}

	if payoutFromDelegation.IsPositive() {
		k.PayFromDelegation(ctx, provider.Address, payoutFromDelegation)
	}

	if payoutFromUnbonding.IsPositive() {
		uncoveredPurchase := sdk.MaxInt(sdk.ZeroInt(), purchased.Sub(provider.DelegationBonded))
		unbondingDelegations := k.GetSortedUnbondingDelegations(ctx, provider.Address)
		for _, ubd := range unbondingDelegations {
			entry := ubd.Entries[0]
			// If purchased is not fully covered, cover purchased first.
			remainingUbd := entry.Balance
			if uncoveredPurchase.IsPositive() {
				if uncoveredPurchase.GTE(entry.Balance) {
					uncoveredPurchase = uncoveredPurchase.Sub(entry.Balance)
					remainingUbd = sdk.ZeroInt()
				} else {
					remainingUbd = entry.Balance.Sub(uncoveredPurchase)
					uncoveredPurchase = sdk.ZeroInt()
				}
			}

			// Make payout after purchased is fully covered.
			if remainingUbd.IsPositive() {
				if remainingUbd.GTE(payout) {
					k.PayFromUnbonding(ctx, ubd, payout)
					return
				}
				k.PayFromUnbonding(ctx, ubd, remainingUbd)
				payoutFromUnbonding = payoutFromUnbonding.Sub(remainingUbd)
			}
		}
	}

	// Update collateral and withdraw.
	provider.Collateral = provider.Collateral.Sub(payout)
	provider.Withdrawing = provider.Withdrawing.Sub(payout)
	k.SetProvider(ctx, provider.Address, provider)

	// Update withdraw queue.
	withdraws := k.GetWithdrawQueueTimeSlice(ctx, withdraw.CompletionTime)
	for i := range withdraws {
		if withdraws[i].Address.Equals(withdraw.Address) && withdraws[i].Amount.Equal(withdraw.Amount) {
			withdraws[i].Amount = withdraws[i].Amount.Sub(payout)
			break
		}
	}
	k.SetWithdrawQueueTimeSlice(ctx, withdraw.CompletionTime, withdraws)
}

// GetSortedUnbondingDelegations gets unbonding delegations sorted by completion time from latest to earliest.
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

// MakePayoutByProvider undelegates delegations and send coins the the shield module.
func (k Keeper) MakePayoutByProvider(ctx sdk.Context, providerAddr sdk.AccAddress, purchased, payout sdk.Int) error {
	provider, found := k.GetProvider(ctx, providerAddr)
	if !found {
		return types.ErrProviderNotFound
	}

	// If collateral - withdraw >= purchased + payout, make payouts from collateral - withdraw.
	// If purchased <= collateral - withdraw < purchased + payout, make payouts from collateral - withdraw and withdraw.
	// Otherwise, make payouts from withdraw.
	uncoveredPurchase := sdk.ZeroInt()
	payoutFromCollateral := sdk.ZeroInt()
	if provider.Collateral.Sub(provider.Withdrawing).GTE(purchased.Add(payout)) {
		payoutFromCollateral = payout
	} else if provider.Collateral.Sub(provider.Withdrawing).GTE(purchased) {
		payoutFromCollateral = provider.Collateral.Sub(provider.Withdrawing).Sub(purchased)
	} else {
		uncoveredPurchase = purchased.Sub(provider.Collateral.Sub(provider.Withdrawing))
	}
	payoutFromWithdraw := payout.Sub(payoutFromCollateral)

	// Make payout from collateral - withdraw.
	if payoutFromCollateral.IsPositive() {
		k.PayFromDelegation(ctx, providerAddr, payoutFromCollateral)
	}

	// If no payout needs to be made from withdraw, finish payout.
	if payoutFromWithdraw.IsZero() {
		return nil
	}

	// Make payout from withdraw.
	withdraws := k.GetWithdrawsByProvider(ctx, providerAddr)
	// From latest to oldest.
	for i := len(withdraws) - 1; i >= 0; i-- {
		// If purchased is not fully covered, cover purchased first.
		remainingWithdraw := withdraws[i].Amount
		if uncoveredPurchase.IsPositive() {
			if uncoveredPurchase.GTE(withdraws[i].Amount) {
				uncoveredPurchase = uncoveredPurchase.Sub(withdraws[i].Amount)
				remainingWithdraw = sdk.ZeroInt()
			} else {
				remainingWithdraw = withdraws[i].Amount.Sub(uncoveredPurchase)
				uncoveredPurchase = sdk.ZeroInt()
			}
		}
		// Make payout after purchased is fully covered.
		if remainingWithdraw.IsPositive() {
			if remainingWithdraw.GTE(payoutFromWithdraw) {
				k.PayFromWithdraw(ctx, withdraws[i], payoutFromWithdraw, purchased)
				return nil
			}
			k.PayFromWithdraw(ctx, withdraws[i], remainingWithdraw, purchased)
			payoutFromWithdraw = payoutFromCollateral.Sub(remainingWithdraw)
		}
	}
	if payoutFromWithdraw.IsPositive() {
		panic("payout is not covered")
	}
	return nil
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
