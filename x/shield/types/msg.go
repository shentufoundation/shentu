package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// MsgBeginRedelegate defines the attributes of a bonding transaction.
type MsgCreatePool struct {
	CreatorAddress sdk.AccAddress `json:"creator_address" yaml:"creator_address"`
	Coverage       sdk.Coins      `json:"amount" yaml:"amount"`
	Deposit        sdk.Coins      `json:"deposit" yaml:"deposit"`
}

// NewMsgBeginRedelegate creates a new MsgBeginRedelegate instance.
func NewMsgCreatePool(accAddr sdk.AccAddress, coverage sdk.Coins, deposit sdk.Coins) (MsgCreatePool, error) {
	return MsgCreatePool{
		CreatorAddress: accAddr,
		Coverage:       coverage,
		Deposit:        deposit,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgCreatePool) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg MsgCreatePool) Type() string { return "create_pool" }

// GetSigners implements the sdk.Msg interface
func (msg MsgCreatePool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.CreatorAddress}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgCreatePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreatePool) ValidateBasic() error {
	if msg.CreatorAddress.Empty() {
		return staking.ErrEmptyValidatorAddr
	}
	return nil
}
