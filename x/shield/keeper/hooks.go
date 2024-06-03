package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Wrapper struct
type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

// Create new distribution hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

// - when a delegator delegates
// - when a delegator sends an unbond tx. before the delegation is moved to unbonding queue
// - when a validator self-delegates to become a validator
// - when a delegator sends a redelegate tx
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, _ sdk.ValAddress) error {
	h.k.UpdateDelegationAmount(ctx, delAddr)
	return nil
}

// - when a delegator sends an unbond tx that removes a delegation
// - when a delegator sends a redelegate tx
func (h Hooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	h.k.RemoveDelegation(ctx, delAddr, valAddr)
	return nil
}

// unused hooks
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error { return nil }
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error {
	return nil
}
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}
