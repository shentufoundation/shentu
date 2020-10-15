package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stk "github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/certikfoundation/shentu/x/staking/internal/types"
)

// ErrUnauthorizedValidator is the SDK error for creating validators from non-certified full nodes.
var ErrUnauthorizedValidator = sdkerrors.Register(stakingTypes.ModuleName, 201, "only certified full nodes can become validators")

// NewHandler returns a customized staking message handler.
func NewHandler(k keeper.Keeper, ck types.CertKeeper) sdk.Handler {
	handler := stk.NewHandler(k)
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case stakingTypes.MsgCreateValidator:
			return handleMsgCreateValidator(ctx, msg, ck, handler)
		case stakingTypes.MsgEditValidator:
			return handleMsgEditValidator(ctx, msg, handler)
		default:
			return handler(ctx, msg)
		}
	}
}

func handleMsgCreateValidator(ctx sdk.Context, msg stakingTypes.MsgCreateValidator, ck types.CertKeeper, handler sdk.Handler) (*sdk.Result, error) {
	if _, ok := ck.GetValidator(ctx, msg.PubKey); !ok && ctx.BlockHeight() > 0 {
		return &sdk.Result{Events: ctx.EventManager().Events()}, ErrUnauthorizedValidator
	}
	if msg.Commission.Rate.LT(sdk.OneDec()) {
		return &sdk.Result{Events: ctx.EventManager().Events()}, sdkerrors.Register(stakingTypes.ModuleName, 99, "commission cannot be less than 100%")
	}
	return handler(ctx, msg)
}

func handleMsgEditValidator(ctx sdk.Context, msg stakingTypes.MsgEditValidator, handler sdk.Handler) (*sdk.Result, error) {
	if msg.CommissionRate.LT(sdk.OneDec()) {
		return &sdk.Result{Events: ctx.EventManager().Events()}, sdkerrors.Register(stakingTypes.ModuleName, 99, "commission cannot be less than 100%")
	}
	return handler(ctx, msg)
}
