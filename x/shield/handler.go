package shield

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// NewHandler creates an sdk.Handler for all the shield type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCreatePool:
			return handleMsgCreatePool(ctx, msg, k)
		case types.MsgUpdatePool:
			return handleMsgUpdatePool(ctx, msg, k)
		case types.MsgPausePool:
			return handleMsgPausePool(ctx, msg, k)
		case types.MsgResumePool:
			return handleMsgResumePool(ctx, msg, k)
		case types.MsgWithdrawRewards:
			return handleMsgWithdrawRewards(ctx, msg, k)
		case types.MsgDepositCollateral:
			return handleMsgDepositCollateral(ctx, msg, k)
		case types.MsgWithdrawCollateral:
			return handleMsgWithdrawCollateral(ctx, msg, k)
		case types.MsgPurchaseShield:
			return handleMsgPurchaseShield(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func NewShieldClaimProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case types.ShieldClaimProposal:
			return handleShieldClaimProposal()
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized shield proposal content type: %T", c)
		}
	}
}

func handleShieldClaimProposal() error {
	return types.ErrOperationNotSupported
}

func handleMsgCreatePool(ctx sdk.Context, msg types.MsgCreatePool, k Keeper) (*sdk.Result, error) {
	poolID, err := k.CreatePool(ctx, msg.From, msg.Shield, msg.Deposit, msg.Sponsor, msg.SponsorAddr, msg.Description, msg.ShieldLimit)
	if err != nil {
		return &sdk.Result{Events: ctx.EventManager().Events()}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeKeyShield, msg.Shield.String()),
			sdk.NewAttribute(types.AttributeKeyDeposit, msg.Deposit.String()),
			sdk.NewAttribute(types.AttributeKeySponsor, msg.Sponsor),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(poolID, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgUpdatePool(ctx sdk.Context, msg types.MsgUpdatePool, k Keeper) (*sdk.Result, error) {
	_, err := k.UpdatePool(ctx, msg.PoolID, msg.Description, msg.From, msg.Shield, msg.ServiceFees, msg.ShieldLimit)
	if err != nil {
		return &sdk.Result{Events: ctx.EventManager().Events()}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdatePool,
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolID, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgPausePool(ctx sdk.Context, msg types.MsgPausePool, k Keeper) (*sdk.Result, error) {
	_, err := k.PausePool(ctx, msg.From, msg.PoolID)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePausePool,
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolID, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgResumePool(ctx sdk.Context, msg types.MsgResumePool, k Keeper) (*sdk.Result, error) {
	_, err := k.ResumePool(ctx, msg.From, msg.PoolID)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeResumePool,
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolID, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgWithdrawRewards(ctx sdk.Context, msg types.MsgWithdrawRewards, k Keeper) (*sdk.Result, error) {
	amount, err := k.PayoutNativeRewards(ctx, msg.From)
	if err != nil {
		return &sdk.Result{Events: ctx.EventManager().Events()}, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawRewards,
			sdk.NewAttribute(types.AttributeKeyAccountAddress, msg.From.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, k.BondDenom(ctx)),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDepositCollateral(ctx sdk.Context, msg types.MsgDepositCollateral, k Keeper) (*sdk.Result, error) {
	bondDenom := k.BondDenom(ctx)
	for _, coin := range msg.Collateral {
		if coin.Denom != bondDenom {
			return nil, types.ErrCollateralBadDenom
		}
	}
	amount := msg.Collateral.AmountOf(bondDenom)

	if err := k.DepositCollateral(ctx, msg.From, amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDepositCollateral,
			sdk.NewAttribute(types.AttributeKeyCollateral, amount.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgWithdrawCollateral(ctx sdk.Context, msg types.MsgWithdrawCollateral, k Keeper) (*sdk.Result, error) {
	bondDenom := k.BondDenom(ctx)
	for _, coin := range msg.Collateral {
		if coin.Denom != bondDenom {
			return nil, types.ErrCollateralBadDenom
		}
	}
	amount := msg.Collateral.AmountOf(bondDenom)
	if err := k.WithdrawCollateral(ctx, msg.From, amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawCollateral,
			sdk.NewAttribute(types.AttributeKeyCollateral, amount.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgPurchaseShield(ctx sdk.Context, msg types.MsgPurchaseShield, k Keeper) (*sdk.Result, error) {
	purchase, err := k.PurchaseShield(ctx, msg.PoolID, msg.Shield, msg.Description, msg.From)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePurchaseShield,
			sdk.NewAttribute(types.AttributeKeyPurchaseID, strconv.FormatUint(purchase.PurchaseID, 10)),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolID, 10)),
			sdk.NewAttribute(types.AttributeKeyProtectionEndTime, purchase.ProtectionEndTime.String()),
			sdk.NewAttribute(types.AttributeKeyPurchaseDescription, purchase.Description),
			sdk.NewAttribute(types.AttributeKeyShield, purchase.Shield.String()),
			sdk.NewAttribute(types.AttributeKeyServiceFees, purchase.ServiceFees.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
