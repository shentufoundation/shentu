package types

import (
	"github.com/tendermint/tendermint/crypto"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Validators is a collection of Validator objects.
type Validators []Validator

// TmConsPubKey casts Validator.ConsensusPubkey to crypto.PubKey
func (v Validator) TmConsPubKey() (crypto.PubKey, error) {
	pk, ok := v.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting crypto.PubKey, got %T", pk)
	}

	// The way things are refactored now, v.ConsensusPubkey is sometimes a TM
	// ed25519 pubkey, sometimes our own ed25519 pubkey. This is very ugly and
	// inconsistent.
	// Luckily, here we coerce it into a TM ed25519 pubkey always, as this
	// pubkey will be passed into TM (eg calling encoding.PubKeyToProto).
	if intoTmPk, ok := pk.(cryptotypes.IntoTmPubKey); ok {
		return intoTmPk.AsTmPubKey(), nil
	}
	return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "Logic error: ConsensusPubkey must be an SDK key and SDK PubKey types must be convertible to tendermint PubKey; got: %T", pk)
}
