package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/shield/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the shield MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}
	sponsorAddr, err := sdk.AccAddressFromBech32(msg.SponsorAddr)
	if err != nil {
		return nil, err
	}

	poolID, err := k.Keeper.CreatePool(ctx, fromAddr, msg.Shield, msg.Deposit, msg.Sponsor, sponsorAddr, msg.Description, msg.ShieldLimit)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeMsgCreatePool,
			sdk.NewAttribute(types.AttributeKeyShield, msg.Shield.String()),
			sdk.NewAttribute(types.AttributeKeyDeposit, msg.Deposit.String()),
			sdk.NewAttribute(types.AttributeKeySponsor, msg.Sponsor),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(poolID, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &types.MsgCreatePoolResponse{}, nil
}

func (k msgServer) UpdatePool(goCtx context.Context, msg *types.MsgUpdatePool) (*types.MsgUpdatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	_, err = k.Keeper.UpdatePool(ctx, msg.PoolId, msg.Description, fromAddr, msg.Shield, msg.ServiceFees, msg.ShieldLimit)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeMsgUpdatePool,
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &types.MsgUpdatePoolResponse{}, nil
}

func (k msgServer) PausePool(goCtx context.Context, msg *types.MsgPausePool) (*types.MsgPausePoolResponse, error) {
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
			types.TypeMsgPausePool,
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &types.MsgPausePoolResponse{}, nil
}

func (k msgServer) ResumePool(goCtx context.Context, msg *types.MsgResumePool) (*types.MsgResumePoolResponse, error) {
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
			types.TypeMsgResumePool,
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &types.MsgResumePoolResponse{}, nil
}

func (k msgServer) DepositCollateral(goCtx context.Context, msg *types.MsgDepositCollateral) (*types.MsgDepositCollateralResponse, error) {
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
			types.TypeMsgDepositCollateral,
			sdk.NewAttribute(types.AttributeKeyCollateral, amount.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &types.MsgDepositCollateralResponse{}, nil
}

func (k msgServer) WithdrawCollateral(goCtx context.Context, msg *types.MsgWithdrawCollateral) (*types.MsgWithdrawCollateralResponse, error) {
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
			types.TypeMsgWithdrawCollateral,
			sdk.NewAttribute(types.AttributeKeyCollateral, amount.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &types.MsgWithdrawCollateralResponse{}, nil
}

func (k msgServer) WithdrawRewards(goCtx context.Context, msg *types.MsgWithdrawRewards) (*types.MsgWithdrawRewardsResponse, error) {
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
			types.TypeMsgWithdrawRewards,
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

	return &types.MsgWithdrawRewardsResponse{}, nil
}

func (k msgServer) UpdateSponsor(goCtx context.Context, msg *types.MsgUpdateSponsor) (*types.MsgUpdateSponsorResponse, error) {
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
			types.TypeMsgUpdateSponsor,
			sdk.NewAttribute(types.AttributeKeySponsor, pool.Sponsor),
			sdk.NewAttribute(types.AttributeKeySponsorAddress, pool.SponsorAddr),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(pool.Id, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &types.MsgUpdateSponsorResponse{}, nil
}

func (k msgServer) StakeForShield(goCtx context.Context, msg *types.MsgStakeForShield) (*types.MsgStakeForShieldResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	purchase, err := k.Keeper.PurchaseShield(ctx, msg.PoolId, msg.Shield, msg.Description, fromAddr, true)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeMsgStakeForShield,
			sdk.NewAttribute(types.AttributeKeyPurchaseID, strconv.FormatUint(purchase.PurchaseId, 10)),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeKeyProtectionEndTime, purchase.ProtectionEndTime.String()),
			sdk.NewAttribute(types.AttributeKeyPurchaseDescription, purchase.Description),
			sdk.NewAttribute(types.AttributeKeyShield, purchase.Shield.String()),
			sdk.NewAttribute(types.AttributeKeyServiceFees, purchase.ServiceFees.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &types.MsgStakeForShieldResponse{}, nil
}

func (k msgServer) UnstakeFromShield(goCtx context.Context, msg *types.MsgUnstakeFromShield) (*types.MsgUnstakeFromShieldResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}
	amount := msg.Shield.AmountOf(k.Keeper.BondDenom(ctx))

	err = k.Keeper.UnstakeFromShield(ctx, msg.PoolId, fromAddr, amount)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeMsgUnstakeFromShield,
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &types.MsgUnstakeFromShieldResponse{}, nil
}

func (k msgServer) PurchaseShield(goCtx context.Context, msg *types.MsgPurchaseShield) (*types.MsgPurchaseShieldResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	purchase, err := k.Keeper.PurchaseShield(ctx, msg.PoolId, msg.Shield, msg.Description, fromAddr, false)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeMsgPurchaseShield,
			sdk.NewAttribute(types.AttributeKeyPurchaseID, strconv.FormatUint(purchase.PurchaseId, 10)),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeKeyProtectionEndTime, purchase.ProtectionEndTime.String()),
			sdk.NewAttribute(types.AttributeKeyPurchaseDescription, purchase.Description),
			sdk.NewAttribute(types.AttributeKeyShield, purchase.Shield.String()),
			sdk.NewAttribute(types.AttributeKeyServiceFees, purchase.ServiceFees.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &types.MsgPurchaseShieldResponse{}, nil
}

func (k msgServer) WithdrawReimbursement(goCtx context.Context, msg *types.MsgWithdrawReimbursement) (*types.MsgWithdrawReimbursementResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	amount, err := k.Keeper.WithdrawReimbursement(ctx, msg.ProposalId, fromAddr)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeMsgWithdrawReimbursement,
			sdk.NewAttribute(types.AttributeKeyPurchaseID, strconv.FormatUint(msg.ProposalId, 10)),
			sdk.NewAttribute(types.AttributeKeyCompensationAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyBeneficiary, msg.From),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &types.MsgWithdrawReimbursementResponse{}, nil
}

func (k msgServer) WithdrawForeignRewards(goCtx context.Context, msg *types.MsgWithdrawForeignRewards) (*types.MsgWithdrawForeignRewardsResponse, error) {
	return &types.MsgWithdrawForeignRewardsResponse{}, nil
}
