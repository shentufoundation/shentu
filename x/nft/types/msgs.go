package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/irisnet/irismod/modules/nft/types"
)

// NewMsgCreateNFTAdmin returns a new certifier proposal message.
func NewMsgCreateNFTAdmin(proposer, certifier sdk.AccAddress, alias string, description string) *MsgCreateNFTAdmin {
	return &MsgCreateNFTAdmin{
		Proposer:    proposer.String(),
		Certifier:   certifier.String(),
		Alias:       alias,
		Description: description,
	}
}

// Route returns the module name.
func (m MsgCreateNFTAdmin) Route() string { return types.ModuleName }

// Type returns the action name.
func (m MsgCreateNFTAdmin) Type() string { return "propose_certifier" }

// ValidateBasic runs stateless checks on the message.
func (m MsgCreateNFTAdmin) ValidateBasic() error {
	certifierAddr, err := sdk.AccAddressFromBech32(m.Issuer)
	if err != nil {
		panic(err)
	}
	if certifierAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, certifierAddr.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing.
func (m MsgCreateNFTAdmin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners defines whose signature is required.
func (m MsgCreateNFTAdmin) GetSigners() []sdk.AccAddress {
	proposerAddr, err := sdk.AccAddressFromBech32(m.Proposer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{proposerAddr}
}
