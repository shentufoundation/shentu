package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/auth/internal/types"
	"github.com/certikfoundation/shentu/x/auth/vesting"
)

// NewHandler returns a handler for "auth" type messages.
func NewHandler(ak AccountKeeper, ck types.CertKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgManualVesting:
			return handleMsgManualVesting(ctx, ak, ck, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Unrecognized cert Msg type: %v", msg.Type())
		}
	}
}

func handleMsgManualVesting(ctx sdk.Context, ak AccountKeeper, ck types.CertKeeper, msg types.MsgManualVesting) (*sdk.Result, error) {
	// preliminary checks
	if !ck.IsCertifier(ctx, msg.Certifier) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "not a valid certifier")
	}

	acc := ak.GetAccount(ctx, msg.Account)
	if acc == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", msg.Account)
	}

	mvacc, ok := acc.(*vesting.ManualVestingAccount)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "account does not appear to be a ManualVestingAccount")
	}

	if mvacc.VestedCoins.Add(msg.UnlockAmount).IsAnyGT(mvacc.OriginalVesting) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "cannot unlock more than the original vesting amount")
	}

	// update vested coins
	mvacc.VestedCoins = mvacc.VestedCoins.Add(msg.UnlockAmount)
	ak.SetAccount(ctx, mvacc)

	return &sdk.Result{}, nil
}
