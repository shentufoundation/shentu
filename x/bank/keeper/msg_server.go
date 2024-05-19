package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	vesting "github.com/shentufoundation/shentu/v2/x/auth/types"
	"github.com/shentufoundation/shentu/v2/x/bank/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) LockedSend(goCtx context.Context, msg *types.MsgLockedSend) (*types.MsgLockedSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}
	toAddr, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, err
	}
	var unlocker sdk.AccAddress
	if msg.UnlockerAddress != "" {
		unlocker, err = sdk.AccAddressFromBech32(msg.UnlockerAddress)
		if err != nil {
			return nil, err
		}
	}

	// preliminary checks
	from := k.ak.GetAccount(ctx, fromAddr)
	if from == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "sender account %s does not exist", msg.FromAddress)
	}
	if toAddr.Equals(unlocker) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "recipient cannot be the unlocker")
	}

	acc := k.ak.GetAccount(ctx, toAddr)

	var toAcc *vesting.ManualVestingAccount
	if acc == nil {
		acc = k.ak.NewAccountWithAddress(ctx, toAddr)
		toAddr, err := sdk.AccAddressFromBech32(msg.ToAddress)
		if err != nil {
			panic(err)
		}

		baseAcc := authtypes.NewBaseAccount(toAddr, acc.GetPubKey(), acc.GetAccountNumber(), acc.GetSequence())
		if unlocker.Empty() {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid unlocker address provided")
		}
		toAcc = vesting.NewManualVestingAccount(baseAcc, sdk.NewCoins(), sdk.NewCoins(), unlocker)
	} else {
		var ok bool
		toAcc, ok = acc.(*vesting.ManualVestingAccount)
		if !ok {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "receiver account is not a ManualVestingAccount")
		}
		if !unlocker.Empty() {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "cannot change the unlocker for existing ManualVestingAccount")
		}
	}

	// send from sender account to receiver account
	// but make the added amount vesting (OV := Vesting + Vested)
	err = k.SendCoins(ctx, fromAddr, toAddr, msg.Amount)
	if err != nil {
		return nil, err
	}

	toAcc.OriginalVesting = toAcc.OriginalVesting.Add(msg.Amount...)
	k.ak.SetAccount(ctx, toAcc)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLockedSendToVestingAccount,
			sdk.NewAttribute(bankTypes.AttributeKeySender, msg.FromAddress),
			sdk.NewAttribute(bankTypes.AttributeKeyRecipient, msg.ToAddress),
			sdk.NewAttribute(types.AttributeKeyUnlocker, msg.UnlockerAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	)
	return &types.MsgLockedSendResponse{}, nil
}
