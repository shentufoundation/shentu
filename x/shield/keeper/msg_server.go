package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the shield MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) v1beta1.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ v1beta1.MsgServer = msgServer{}

func (k msgServer) CreatePool(goCtx context.Context, msg *v1beta1.MsgCreatePool) (*v1beta1.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	poolID, err := k.Keeper.CreatePool(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			v1beta1.TypeMsgCreatePool,
			sdk.NewAttribute(types.AttributeKeySponsorAddress, msg.SponsorAddr),
			sdk.NewAttribute(types.AttributeKeyShieldRate, msg.ShieldRate.String()),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(poolID, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &v1beta1.MsgCreatePoolResponse{}, nil
}

func (k msgServer) UpdatePool(goCtx context.Context, msg *v1beta1.MsgUpdatePool) (*v1beta1.MsgUpdatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := k.Keeper.UpdatePool(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			v1beta1.TypeMsgUpdatePool,
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &v1beta1.MsgUpdatePoolResponse{}, nil
}

func (k msgServer) PausePool(goCtx context.Context, msg *v1beta1.MsgPausePool) (*v1beta1.MsgPausePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	_, err = k.Keeper.PausePool(ctx, fromAddr, msg.PoolId)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			v1beta1.TypeMsgPausePool,
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &v1beta1.MsgPausePoolResponse{}, nil
}

func (k msgServer) ResumePool(goCtx context.Context, msg *v1beta1.MsgResumePool) (*v1beta1.MsgResumePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	_, err = k.Keeper.ResumePool(ctx, fromAddr, msg.PoolId)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			v1beta1.TypeMsgResumePool,
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &v1beta1.MsgResumePoolResponse{}, nil
}

func (k msgServer) DepositCollateral(goCtx context.Context, msg *v1beta1.MsgDepositCollateral) (*v1beta1.MsgDepositCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	bondDenom := k.Keeper.BondDenom(ctx)
	for _, coin := range msg.Collateral {
		if coin.Denom != bondDenom {
			return nil, types.ErrCollateralBadDenom
		}
	}
	amount := msg.Collateral.AmountOf(bondDenom)

	if err := k.Keeper.DepositCollateral(ctx, fromAddr, amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			v1beta1.TypeMsgDepositCollateral,
			sdk.NewAttribute(types.AttributeKeyCollateral, amount.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &v1beta1.MsgDepositCollateralResponse{}, nil
}

func (k msgServer) WithdrawCollateral(goCtx context.Context, msg *v1beta1.MsgWithdrawCollateral) (*v1beta1.MsgWithdrawCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	bondDenom := k.Keeper.BondDenom(ctx)
	for _, coin := range msg.Collateral {
		if coin.Denom != bondDenom {
			return nil, types.ErrCollateralBadDenom
		}
	}
	amount := msg.Collateral.AmountOf(bondDenom)
	if err := k.Keeper.WithdrawCollateral(ctx, fromAddr, amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			v1beta1.TypeMsgWithdrawCollateral,
			sdk.NewAttribute(types.AttributeKeyCollateral, amount.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &v1beta1.MsgWithdrawCollateralResponse{}, nil
}

func (k msgServer) WithdrawRewards(goCtx context.Context, msg *v1beta1.MsgWithdrawRewards) (*v1beta1.MsgWithdrawRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	amount, err := k.Keeper.PayoutNativeRewards(ctx, fromAddr)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			v1beta1.TypeMsgWithdrawRewards,
			sdk.NewAttribute(types.AttributeKeyAccountAddress, msg.From),
			sdk.NewAttribute(types.AttributeKeyDenom, k.Keeper.BondDenom(ctx)),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &v1beta1.MsgWithdrawRewardsResponse{}, nil
}

func (k msgServer) UpdateSponsor(goCtx context.Context, msg *v1beta1.MsgUpdateSponsor) (*v1beta1.MsgUpdateSponsorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}
	sponsorAddr, err := sdk.AccAddressFromBech32(msg.SponsorAddr)
	if err != nil {
		return nil, err
	}

	pool, err := k.Keeper.UpdateSponsor(ctx, msg.PoolId, msg.Sponsor, sponsorAddr, fromAddr)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			v1beta1.TypeMsgUpdateSponsor,
			sdk.NewAttribute(types.AttributeKeySponsorAddress, pool.SponsorAddr),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(pool.Id, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &v1beta1.MsgUpdateSponsorResponse{}, nil
}

func (k msgServer) Purchase(goCtx context.Context, msg *v1beta1.MsgPurchase) (*v1beta1.MsgPurchaseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	purchase, err := k.Keeper.PurchaseShield(ctx, msg.PoolId, msg.Amount, msg.Description, fromAddr)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			v1beta1.TypeMsgStakeForShield,
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeKeyAccountAddress, msg.From),
			sdk.NewAttribute(types.AttributeKeyAmount, purchase.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &v1beta1.MsgPurchaseResponse{}, nil
}

func (k msgServer) Unstake(goCtx context.Context, msg *v1beta1.MsgUnstake) (*v1beta1.MsgUnstakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	err = k.Keeper.Unstake(ctx, msg.PoolId, fromAddr, msg.Amount)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			v1beta1.TypeMsgUnstakeFromShield,
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &v1beta1.MsgUnstakeResponse{}, nil
}

func (k msgServer) WithdrawForeignRewards(goCtx context.Context, msg *v1beta1.MsgWithdrawForeignRewards) (*v1beta1.MsgWithdrawForeignRewardsResponse, error) {
	return &v1beta1.MsgWithdrawForeignRewardsResponse{}, nil
}

func (k msgServer) Donate(goCtx context.Context, msg *v1beta1.MsgDonate) (*v1beta1.MsgDonateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	bondDenom := k.Keeper.BondDenom(ctx)
	for _, coin := range msg.Amount {
		if coin.Denom != bondDenom {
			return nil, types.ErrDonationBadDenom
		}
	}
	amount := msg.Amount.AmountOf(bondDenom)

	if err := k.Keeper.Donate(ctx, fromAddr, amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			v1beta1.TypeMsgDonate,
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &v1beta1.MsgDonateResponse{}, nil
}
