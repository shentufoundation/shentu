package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreatePool             = "create_pool"
	TypeMsgUpdatePool             = "update_pool"
	TypeMsgPausePool              = "pause_pool"
	TypeMsgResumePool             = "resume_pool"
	TypeMsgDepositCollateral      = "deposit_collateral"
	TypeMsgWithdrawCollateral     = "withdraw_collateral"
	TypeMsgWithdrawRewards        = "withdraw_rewards"
	TypeMsgWithdrawForeignRewards = "withdraw_foreign_rewards"
	TypeMsgPurchaseShield         = "purchase_shield"
	TypeMsgWithdrawReimbursement  = "withdraw_reimbursement"
	TypeMsgStakeForShield         = "stake_for_shield"
	TypeMsgUnstakeFromShield      = "unstake_from_shield"
	TypeMsgUpdateSponsor          = "update_sponsor"
	TypeMsgDonate                 = "donate"
)

// NewMsgCreatePool creates a new NewMsgCreatePool instance.
func NewMsgCreatePool(accAddr, sponsorAddr sdk.AccAddress, description string, shieldRate sdk.Dec, shieldLimit sdk.Coins) *MsgCreatePool {
	return &MsgCreatePool{
		From:        accAddr.String(),
		SponsorAddr: sponsorAddr.String(),
		Description: description,
		ShieldRate:  shieldRate,
		ShieldLimit: shieldLimit,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgCreatePool) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgCreatePool) Type() string { return TypeMsgCreatePool }

// GetSigners implements the sdk.Msg interface.
func (msg MsgCreatePool) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgCreatePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreatePool) ValidateBasic() error {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return err
	}
	if from.Empty() {
		return ErrEmptySender
	}

	if !msg.ShieldRate.GTE(sdk.NewDec(1)) {
		return ErrInvalidShieldRate
	}
	return nil
}

// NewMsgUpdatePool creates a new MsgUpdatePool instance.
func NewMsgUpdatePool(accAddr sdk.AccAddress, id uint64, description string, active bool, shieldRate sdk.Dec, shieldLimit sdk.Coins) *MsgUpdatePool {
	return &MsgUpdatePool{
		From:        accAddr.String(),
		PoolId:      id,
		Description: description,
		Active:      active,
		ShieldRate:  shieldRate,
		ShieldLimit: shieldLimit,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdatePool) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdatePool) Type() string { return TypeMsgUpdatePool }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdatePool) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdatePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdatePool) ValidateBasic() error {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	if from.Empty() {
		return ErrEmptySender
	}

	if msg.PoolId == 0 {
		return ErrInvalidPoolID
	}
	if !msg.ShieldRate.IsPositive() {
		return ErrInvalidShieldRate
	}
	return nil
}

// NewMsgPausePool creates a new NewMsgPausePool instance.
func NewMsgPausePool(accAddr sdk.AccAddress, id uint64) *MsgPausePool {
	return &MsgPausePool{
		From:   accAddr.String(),
		PoolId: id,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgPausePool) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgPausePool) Type() string { return TypeMsgPausePool }

// GetSigners implements the sdk.Msg interface.
func (msg MsgPausePool) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgPausePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgPausePool) ValidateBasic() error {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	if from.Empty() {
		return ErrEmptySender
	}

	if msg.PoolId == 0 {
		return ErrInvalidPoolID
	}
	return nil
}

func NewMsgResumePool(accAddr sdk.AccAddress, id uint64) *MsgResumePool {
	return &MsgResumePool{
		From:   accAddr.String(),
		PoolId: id,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgResumePool) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgResumePool) Type() string { return TypeMsgResumePool }

// GetSigners implements the sdk.Msg interface.
func (msg MsgResumePool) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgResumePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgResumePool) ValidateBasic() error {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	if from.Empty() {
		return ErrEmptySender
	}

	if msg.PoolId == 0 {
		return ErrInvalidPoolID
	}
	return nil
}

// NewMsgDonate creates a new MsgDonate instance.
func NewMsgDonate(sender sdk.AccAddress, amount sdk.Coins) *MsgDonate {
	return &MsgDonate{
		From:   sender.String(),
		Amount: amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgDonate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgDonate) Type() string { return "donate" }

// GetSigners implements the sdk.Msg interface.
func (msg MsgDonate) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgDonate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgDonate) ValidateBasic() error {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	if from.Empty() {
		return ErrEmptySender
	}

	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "Donation amount: %s", msg.Amount)
	}
	return nil
}

// NewMsgDepositCollateral creates a new MsgDepositCollateral instance.
func NewMsgDepositCollateral(sender sdk.AccAddress, collateral sdk.Coins) *MsgDepositCollateral {
	return &MsgDepositCollateral{
		From:       sender.String(),
		Collateral: collateral,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgDepositCollateral) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgDepositCollateral) Type() string { return "deposit_collateral" }

// GetSigners implements the sdk.Msg interface.
func (msg MsgDepositCollateral) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgDepositCollateral) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgDepositCollateral) ValidateBasic() error {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	if from.Empty() {
		return ErrEmptySender
	}

	if !msg.Collateral.IsValid() || msg.Collateral.IsZero() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "Collateral amount: %s", msg.Collateral)
	}
	return nil
}

// NewMsgWithdrawCollateral creates a new NewMsgWithdrawCollateral instance.
func NewMsgWithdrawCollateral(sender sdk.AccAddress, collateral sdk.Coins) *MsgWithdrawCollateral {
	return &MsgWithdrawCollateral{
		From:       sender.String(),
		Collateral: collateral,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgWithdrawCollateral) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgWithdrawCollateral) Type() string { return "withdraw_collateral" }

// GetSigners implements the sdk.Msg interface.
func (msg MsgWithdrawCollateral) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgWithdrawCollateral) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgWithdrawCollateral) ValidateBasic() error {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	if from.Empty() {
		return ErrEmptySender
	}

	if !msg.Collateral.IsValid() || msg.Collateral.IsZero() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "Collateral amount: %s", msg.Collateral)
	}
	return nil
}

// NewMsgWithdrawRewards creates a new MsgWithdrawRewards instance.
func NewMsgWithdrawRewards(sender sdk.AccAddress) *MsgWithdrawRewards {
	return &MsgWithdrawRewards{
		From: sender.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgWithdrawRewards) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgWithdrawRewards) Type() string { return TypeMsgWithdrawRewards }

// GetSigners implements the sdk.Msg interface.
func (msg MsgWithdrawRewards) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgWithdrawRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgWithdrawRewards) ValidateBasic() error {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	if from.Empty() {
		return ErrEmptySender
	}

	return nil
}

// NewMsgWithdrawForeignRewards creates a new MsgWithdrawForeignRewards instance.
func NewMsgWithdrawForeignRewards(sender sdk.AccAddress, denom, toAddr string) *MsgWithdrawForeignRewards {
	return &MsgWithdrawForeignRewards{
		From:   sender.String(),
		Denom:  denom,
		ToAddr: toAddr,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgWithdrawForeignRewards) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg MsgWithdrawForeignRewards) Type() string { return TypeMsgWithdrawForeignRewards }

// GetSigners implements the sdk.Msg interface
func (msg MsgWithdrawForeignRewards) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgWithdrawForeignRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgWithdrawForeignRewards) ValidateBasic() error {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	if from.Empty() {
		return ErrEmptySender
	}
	if strings.TrimSpace(msg.ToAddr) == "" {
		return ErrInvalidToAddr
	}
	return nil
}

// NewMsgPurchase creates a new MsgPurchaseShield instance.
func NewMsgPurchase(poolID uint64, shield sdk.Coins, description string, from sdk.AccAddress) *MsgPurchase {
	return &MsgPurchase{
		PoolId:      poolID,
		Amount:      shield,
		Description: description,
		From:        from.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgPurchase) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgPurchase) Type() string { return TypeMsgStakeForShield }

// GetSigners implements the sdk.Msg interface.
func (msg MsgPurchase) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgPurchase) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgPurchase) ValidateBasic() error {
	return nil
}

// NewMsgUnstake creates a new MsgPurchaseShield instance.
func NewMsgUnstake(poolID uint64, shield sdk.Coins, from sdk.AccAddress) *MsgUnstake {
	return &MsgUnstake{
		PoolId: poolID,
		Amount: shield,
		From:   from.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUnstake) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUnstake) Type() string { return TypeMsgUnstakeFromShield }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUnstake) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUnstake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUnstake) ValidateBasic() error {
	return nil
}

// NewMsgUpdateSponsor creates a new NewMsgUpdateSponsor instance.
func NewMsgUpdateSponsor(poolID uint64, sponsor string, sponsorAddr, fromAddr sdk.AccAddress) *MsgUpdateSponsor {
	return &MsgUpdateSponsor{
		PoolId:      poolID,
		Sponsor:     sponsor,
		SponsorAddr: sponsorAddr.String(),
		From:        fromAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateSponsor) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateSponsor) Type() string { return TypeMsgUpdateSponsor }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateSponsor) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdateSponsor) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateSponsor) ValidateBasic() error {
	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	if fromAddr.Empty() {
		return ErrEmptySender
	}

	sponsorAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	if sponsorAddr.Empty() || strings.TrimSpace(msg.Sponsor) == "" {
		return ErrEmptySponsor
	}
	return nil
}
