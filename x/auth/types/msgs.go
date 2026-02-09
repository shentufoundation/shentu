package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = ModuleName

	TypeMsgUnlock = "unlock"
)

// NewMsgUnlock returns a MsgUnlock object.
func NewMsgUnlock(issuer, account sdk.AccAddress, unlockAmount sdk.Coins) *MsgUnlock {
	return &MsgUnlock{
		Issuer:       issuer.String(),
		Account:      account.String(),
		UnlockAmount: unlockAmount,
	}
}

// Route returns the name of the module.
func (m MsgUnlock) Route() string { return ModuleName }

// Type returns a human-readable string for the message.
func (m MsgUnlock) Type() string { return "unlock" }

// GetSigners defines whose signature is required.
func (m MsgUnlock) GetSigners() []sdk.AccAddress {
	issuer, err := sdk.AccAddressFromBech32(m.Issuer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{issuer}
}
