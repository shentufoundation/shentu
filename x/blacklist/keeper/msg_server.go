package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/shentufoundation/shentu/v2/x/blacklist/types"
)

type msgServer struct {
	k Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the blacklist MsgServer
// interface for the provided Keeper.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return msgServer{k: k}
}

// BanAddress implements the Msg/BanAddress RPC method. Two submitter lanes
// are accepted: the governance module account (proposal execution) and
// registered certifiers (direct emergency lane — TEMPORARY, see
// ValidateSubmitter).
func (ms msgServer) BanAddress(ctx context.Context, msg *types.MsgBanAddress) (*types.MsgBanAddressResponse, error) {
	if err := ms.k.ValidateSubmitter(ctx, msg.Authority); err != nil {
		return nil, err
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, msg.Address)
	}

	// Never allow blacklisting module accounts: it would freeze protocol
	// infrastructure (pools, fee collector, module operations) rather than a
	// malicious actor.
	if acc := ms.k.accountKeeper.GetAccount(ctx, addr); acc != nil {
		if _, ok := acc.(sdk.ModuleAccountI); ok {
			return nil, errorsmod.Wrapf(types.ErrModuleAccount, "%s is a module account", msg.Address)
		}
	}

	entry := types.Entry{
		Address:  msg.Address,
		Category: msg.Category,
	}
	if err := ms.k.SetEntry(ctx, entry); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBanAddress,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address),
			sdk.NewAttribute(types.AttributeKeyCategory, msg.Category.String()),
		),
	)

	return &types.MsgBanAddressResponse{}, nil
}

// UnbanAddress implements the Msg/UnbanAddress RPC method. Two submitter
// lanes are accepted: the governance module account (proposal execution) and
// registered certifiers (direct emergency lane — TEMPORARY, see
// ValidateSubmitter).
func (ms msgServer) UnbanAddress(ctx context.Context, msg *types.MsgUnbanAddress) (*types.MsgUnbanAddressResponse, error) {
	if err := ms.k.ValidateSubmitter(ctx, msg.Authority); err != nil {
		return nil, err
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	if _, err := ms.k.GetEntry(ctx, msg.Address); err != nil {
		return nil, errorsmod.Wrapf(types.ErrEntryNotFound, "%s is not blacklisted", msg.Address)
	}

	if err := ms.k.RemoveEntry(ctx, msg.Address); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUnbanAddress,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address),
		),
	)

	return &types.MsgUnbanAddressResponse{}, nil
}
