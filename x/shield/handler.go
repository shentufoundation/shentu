package shield

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/x/shield/types"
)

var (
	AEnabled = false
	PEnabled = false
)

// NewHandler creates an sdk.Handler for all the shield type messages
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
		case types.MsgWithdrawForeignRewards:
			return handleMsgWithdrawForeignRewards(ctx, msg, k)
		case types.MsgDepositCollateral:
			return handleMsgDepositCollateral(ctx, msg, k)
		case types.MsgWithdrawCollateral:
			return handleMsgWithdrawCollateral(ctx, msg, k)
		case types.MsgPurchaseShield:
			return handleMsgPurchaseShield(ctx, msg, k)
		case types.MsgWithdrawReimbursement:
			return handleMsgWithdrawReimbursement(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func NewShieldClaimProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case types.ShieldClaimProposal:
			return handleShieldClaimProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized shield proposal content type: %T", c)
		}
	}
}

func handleShieldClaimProposal(ctx sdk.Context, k Keeper, p types.ShieldClaimProposal) error {
	if err := k.CreateReimbursement(ctx, p.ProposalID, p.PoolID, p.Loss, p.Proposer); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateReimbursement,
			sdk.NewAttribute(types.AttributeKeyPurchaseID, strconv.FormatUint(p.PurchaseID, 10)),
			sdk.NewAttribute(types.AttributeKeyCompensationAmount, p.Loss.String()),
			sdk.NewAttribute(types.AttributeKeyBeneficiary, p.Proposer.String()),
		),
	})
	return nil
}

func handleMsgCreatePool(ctx sdk.Context, msg types.MsgCreatePool, k Keeper) (*sdk.Result, error) {
	pool, err := k.CreatePool(ctx, msg.From, msg.Shield, msg.Deposit, msg.Sponsor, msg.SponsorAddr, msg.TimeOfCoverage)
	if err != nil {
		return &sdk.Result{Events: ctx.EventManager().Events()}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeKeyShield, msg.Shield.String()),
			sdk.NewAttribute(types.AttributeKeyDeposit, msg.Deposit.String()),
			sdk.NewAttribute(types.AttributeKeySponsor, msg.Sponsor),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(pool.PoolID, 10)),
			sdk.NewAttribute(types.AttributeKeyTimeOfCoverage, msg.TimeOfCoverage.String()),
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
	_, err := k.UpdatePool(ctx, msg.From, msg.Shield, msg.Deposit, msg.PoolID, msg.AdditionalTime, msg.Description)
	if err != nil {
		return &sdk.Result{Events: ctx.EventManager().Events()}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdatePool,
			sdk.NewAttribute(types.AttributeKeyShield, msg.Shield.String()),
			sdk.NewAttribute(types.AttributeKeyDeposit, msg.Deposit.String()),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolID, 10)),
			sdk.NewAttribute(types.AttributeKeyAdditionalTime, msg.AdditionalTime.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgWithdrawCollateral(ctx sdk.Context, msg types.MsgWithdrawCollateral, k Keeper) (*sdk.Result, error) {
	if msg.Collateral.Denom != k.BondDenom(ctx) {
		return nil, types.ErrCollateralBadDenom
	}

	if err := k.WithdrawCollateral(ctx, msg.From, msg.PoolID, msg.Collateral.Amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawCollateral,
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolID, 10)),
			sdk.NewAttribute(types.AttributeKeyCollateral, msg.Collateral.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDepositCollateral(ctx sdk.Context, msg types.MsgDepositCollateral, k Keeper) (*sdk.Result, error) {
	if !AEnabled {
		return &sdk.Result{}, types.ErrOperationNotSupported
	}

	if msg.Collateral.Denom != k.BondDenom(ctx) {
		return nil, types.ErrCollateralBadDenom
	}
	if err := k.DepositCollateral(ctx, msg.From, msg.PoolID, msg.Collateral.Amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDepositCollateral,
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolID, 10)),
			sdk.NewAttribute(types.AttributeKeyCollateral, msg.Collateral.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgPausePool(ctx sdk.Context, msg types.MsgPausePool, k Keeper) (*sdk.Result, error) {
	if _, err := k.PausePool(ctx, msg.From, msg.PoolID); err != nil {
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
	if _, err := k.ResumePool(ctx, msg.From, msg.PoolID); err != nil {
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

func handleMsgPurchaseShield(ctx sdk.Context, msg types.MsgPurchaseShield, k Keeper) (*sdk.Result, error) {
	if !PEnabled {
		return &sdk.Result{}, types.ErrOperationNotSupported
	}

	_, err := k.PurchaseShield(ctx, msg.PoolID, msg.Shield, msg.Description, msg.From)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePurchaseShield,
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolID, 10)),
			sdk.NewAttribute(types.AttributeKeyShield, msg.Shield.String()),
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

func handleMsgWithdrawForeignRewards(ctx sdk.Context, msg types.MsgWithdrawForeignRewards, k Keeper) (*sdk.Result, error) {
	rewards := k.GetRewards(ctx, msg.From)
	amount := rewards.Foreign.AmountOf(msg.Denom)
	if amount.Equal(sdk.ZeroDec()) {
		return &sdk.Result{Events: ctx.EventManager().Events()}, types.ErrNoRewards
	}
	rewards.Foreign = rewards.Foreign.Sub(
		sdk.DecCoins{sdk.NewDecCoinFromDec(msg.Denom, amount)})
	k.SetRewards(ctx, msg.From, rewards)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawForeignRewards,
			sdk.NewAttribute(types.AttributeKeyToAddr, msg.ToAddr),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
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

func handleMsgWithdrawReimbursement(ctx sdk.Context, msg types.MsgWithdrawReimbursement, k Keeper) (*sdk.Result, error) {
	if !PEnabled {
		return &sdk.Result{}, types.ErrOperationNotSupported
	}

	amount, err := k.WithdrawReimbursement(ctx, msg.ProposalID, msg.From)
	if err != nil {
		return &sdk.Result{Events: ctx.EventManager().Events()}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawReimbursement,
			sdk.NewAttribute(types.AttributeKeyPurchaseID, strconv.FormatUint(msg.ProposalID, 10)),
			sdk.NewAttribute(types.AttributeKeyCompensationAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyBeneficiary, msg.From.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
