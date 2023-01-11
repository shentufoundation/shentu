package types

import (
	"fmt"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgCreateProgram = "create_program"
)

// NewMsgCreateProgram creates a new NewMsgCreateProgram instance.
// Delegator address and validator address are the same.
func NewMsgCreateProgram(
	creatorAddress string, description string, encKey []byte, commissionRate sdk.Dec, deposit sdk.Coins,
	submissionEndTime, judgingEndTime, claimEndTime time.Time,
) (*MsgCreateProgram, error) {
	var encAny *codectypes.Any
	if encKey != nil {
		encKeyMsg := EciesPubKey{
			PubKey: encKey,
		}

		var err error
		if encAny, err = codectypes.NewAnyWithValue(&encKeyMsg); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("encKey is empty")
	}

	return &MsgCreateProgram{
		Description:       description,
		CommissionRate:    commissionRate,
		SubmissionEndTime: submissionEndTime,
		CreatorAddress:    creatorAddress,
		EncryptionKey:     encAny,
		Deposit:           deposit,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgCreateProgram) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgCreateProgram) Type() string { return TypeMsgCreateProgram }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as delegator's, then the validator must
// sign the msg as well.
func (msg MsgCreateProgram) GetSigners() []sdk.AccAddress {
	// creator should sign the message
	cAddr, err := sdk.AccAddressFromBech32(msg.CreatorAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cAddr}
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCreateProgram) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreateProgram) ValidateBasic() error {
	// TODO: implement ValidateBasic
	return nil
}
