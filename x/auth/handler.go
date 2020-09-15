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
		case types.MsgTriggerVesting:
			return handleMsgTriggerVesting(ctx, ak, ck, msg)
		case types.MsgManualVesting:
			return handleMsgManualVesting(ctx, ak, ck, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Unrecognized cert Msg type: %v", msg.Type())
		}
	}
}

func handleMsgTriggerVesting(ctx sdk.Context, ak AccountKeeper, ck types.CertKeeper, msg types.MsgTriggerVesting) (*sdk.Result, error) {
	if !ck.IsCertifier(ctx, msg.Certifier) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "not a valid certifier")
	}

	acc := ak.GetAccount(ctx, msg.Account)
	if acc == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", msg.Account)
	}
	tvacc, ok := acc.(*vesting.TriggeredVestingAccount)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "account does not appear to be a TriggeredVestingAccount")
	}

	startTime := ctx.BlockTime().Unix()
	endTime := startTime
	for _, p := range tvacc.VestingPeriods {
		endTime += p.Length
	}
	tvacc.BaseVestingAccount.EndTime = endTime

	newacc := vesting.NewTriggeredVestingAccountRaw(
		tvacc.BaseVestingAccount,
		startTime,
		tvacc.VestingPeriods,
		true,
	)
	ak.SetAccount(ctx, newacc)
	return &sdk.Result{}, nil
}

func handleMsgManualVesting(ctx sdk.Context, ak AccountKeeper, ck types.CertKeeper, msg types.MsgManualVesting) (*sdk.Result, error) {
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

	newacc := vesting.NewManualVestingAccountRaw(
		mvacc.BaseVestingAccount,
		mvacc.VestedCoins.Add(msg.UnlockAmount),
	)
	ak.SetAccount(ctx, newacc)
	return &sdk.Result{}, nil
}
