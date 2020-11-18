package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/auth/internal/types"
	"github.com/certikfoundation/shentu/x/auth/vesting"
)

// NewHandler returns a handler for "auth" type messages.
func NewHandler(ak AccountKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgUnlock:
			return handleMsgUnlock(ctx, ak, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Unrecognized cert Msg type: %v", msg.Type())
		}
	}
}

func handleMsgUnlock(ctx sdk.Context, ak AccountKeeper, msg types.MsgUnlock) (*sdk.Result, error) {
	// preliminary checks
	acc := ak.GetAccount(ctx, msg.Account)
	if acc == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", msg.Account)
	}

	mvacc, ok := acc.(*vesting.ManualVestingAccount)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "receiver account is not a manual vesting account")
	}

	if !msg.Issuer.Equals(mvacc.Unlocker) {
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
	ak.SetAccount(ctx, mvacc)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUnlock,
			sdk.NewAttribute("unlocker", msg.Issuer.String()),
			sdk.NewAttribute("account", msg.Account.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.UnlockAmount.String()),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
