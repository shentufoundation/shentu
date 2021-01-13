package keeper

import (
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/types"
)

// Max timestamp supported by Amino.
// Dec 31, 9999 - 23:59:59 GMT
const MaxTimestamp = 253402300799

// SetValidator sets a certified validator.
func (k Keeper) SetValidator(ctx sdk.Context, validator cryptotypes.PubKey, certifier sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	pkAny, err := codectypes.NewAnyWithValue(validator)
	if err != nil {
		panic(err)
	}
	validatorData := types.Validator{Pubkey: pkAny, Certifier: certifier.String()}
	store.Set(types.ValidatorStoreKey(validator), k.cdc.MustMarshalBinaryLengthPrefixed(&validatorData))
}

// deleteValidator removes a validator from being certified.
func (k Keeper) deleteValidator(ctx sdk.Context, validator cryptotypes.PubKey) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ValidatorStoreKey(validator))
}

// IsValidatorCertified checks if the validator is current certified.
func (k Keeper) IsValidatorCertified(ctx sdk.Context, validator cryptotypes.PubKey) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ValidatorStoreKey(validator))
}

// GetValidatorCertifier gets the original certifier of the validator.
func (k Keeper) GetValidatorCertifier(ctx sdk.Context, validator cryptotypes.PubKey) (sdk.AccAddress, error) {
	store := ctx.KVStore(k.storeKey)
	if validatorData := store.Get(types.ValidatorStoreKey(validator)); validatorData != nil {
		var validator types.Validator
		k.cdc.MustUnmarshalBinaryLengthPrefixed(validatorData, &validator)

		certifierAddr, err := sdk.AccAddressFromBech32(validator.Certifier)
		if err != nil {
			panic(err)
		}
		return certifierAddr, nil
	}
	return nil, types.ErrValidatorUncertified
}

// GetValidator returns the validator.
func (k Keeper) GetValidator(ctx sdk.Context, validator cryptotypes.PubKey) ([]byte, bool) {
	if validatorData := ctx.KVStore(k.storeKey).Get(types.ValidatorStoreKey(validator)); validatorData != nil {
		return validatorData, true
	}
	return nil, false
}

// CertifyValidator certifies a validator.
func (k Keeper) CertifyValidator(ctx sdk.Context, validator cryptotypes.PubKey, certifier sdk.AccAddress) error {
	if !k.IsCertifier(ctx, certifier) {
		return types.ErrUnqualifiedCertifier
	}
	if k.IsValidatorCertified(ctx, validator) {
		return types.ErrValidatorCertified
	}
	k.SetValidator(ctx, validator, certifier)
	return nil
}

// DecertifyValidator de-certifies a certified validator.
func (k Keeper) DecertifyValidator(ctx sdk.Context, valPubKey cryptotypes.PubKey, decertifier sdk.AccAddress) error {
	if !k.IsCertifier(ctx, decertifier) {
		return types.ErrUnqualifiedCertifier
	}
	certifier, err := k.GetValidatorCertifier(ctx, valPubKey)
	if err != nil {
		return types.ErrValidatorUncertified
	}
	// Can only be de-certified if it's the original certifier, or that the original certifier is no longer certifier.
	if !decertifier.Equals(certifier) && k.IsCertifier(ctx, certifier) {
		return types.ErrUnqualifiedCertifier
	}
	k.deleteValidator(ctx, valPubKey)

	// Now jail the validator until year 9999. This will begin unbonding the
	// validator if not already unbonding (tombstoned).
	consAddr := sdk.GetConsAddress(valPubKey)
	validator := k.stakingKeeper.ValidatorByConsAddr(ctx, consAddr)
	if validator == nil {
		return types.ErrMissingValidator
	}
	if k.slashingKeeper.IsTombstoned(ctx, consAddr) {
		return types.ErrTombstonedValidator
	}
	if !validator.IsJailed() {
		k.slashingKeeper.Jail(ctx, consAddr)
	}
	k.slashingKeeper.JailUntil(ctx, consAddr, time.Unix(MaxTimestamp, 0))
	k.slashingKeeper.Tombstone(ctx, consAddr)
	return nil
}

// IterateAllValidators iterates over the all the stored validators and performs a callback function.
func (k Keeper) IterateAllValidators(ctx sdk.Context, callback func(validator types.Validator) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorsStoreKey())

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var validator types.Validator
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &validator)

		if callback(validator) {
			break
		}
	}
}

// GetAllValidators gets all validators.
func (k Keeper) GetAllValidators(ctx sdk.Context) (validators types.Validators) {
	k.IterateAllValidators(ctx, func(validator types.Validator) bool {
		validators = append(validators, validator)
		return false
	})
	return
}

// GetAllValidatorPubkeys gets all validator pubkeys.
func (k Keeper) GetAllValidatorPubkeys(ctx sdk.Context) (validatorAddresses []string) {
	k.IterateAllValidators(ctx, func(validator types.Validator) bool {
		pk, err := validator.ConsPubKey()
		if err != nil {
			panic(err)
		}
		validatorAddresses = append(validatorAddresses, sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, pk))
		return false
	})
	return
}
