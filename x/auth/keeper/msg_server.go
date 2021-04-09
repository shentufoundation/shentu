package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/auth/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the auth MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Unlock(goCtx context.Context, msg *types.MsgUnlock) (*types.MsgUnlockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	issuerAddr, err := sdk.AccAddressFromBech32(msg.Issuer)
	if err != nil {
		return nil, err
	}
	accountAddr, err := sdk.AccAddressFromBech32(msg.Account)
	if err != nil {
		return nil, err
	}

	// preliminary checks
	acc := k.ak.GetAccount(ctx, accountAddr)
	if acc == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", accountAddr)
	}

	mvacc, ok := acc.(*types.ManualVestingAccount)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "receiver account is not a manual vesting account")
	}

	unlocker, err := sdk.AccAddressFromBech32(mvacc.Unlocker)
	if err != nil {
		return nil, err
	}
	if !issuerAddr.Equals(unlocker) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "sender of this transaction is not the designated unlocker")
	}

	if mvacc.VestedCoins.Add(msg.UnlockAmount...).IsAnyGT(mvacc.OriginalVesting) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "cannot unlock more than the original vesting amount")
	}

	// update vested coins
	mvacc.VestedCoins = mvacc.VestedCoins.Add(msg.UnlockAmount...)
	if mvacc.DelegatedVesting.IsAllGT(mvacc.OriginalVesting.Sub(mvacc.VestedCoins)) {
		unlockedDelegated := mvacc.DelegatedVesting.Sub(mvacc.OriginalVesting.Sub(mvacc.VestedCoins))
		mvacc.DelegatedVesting = mvacc.DelegatedVesting.Sub(unlockedDelegated)
		mvacc.DelegatedFree = mvacc.DelegatedFree.Add(unlockedDelegated...)
	}
	k.ak.SetAccount(ctx, mvacc)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUnlock,
			sdk.NewAttribute("unlocker", issuerAddr.String()),
			sdk.NewAttribute("account", accountAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.UnlockAmount.String()),
		),
	)

	return &types.MsgUnlockResponse{}, nil
}
