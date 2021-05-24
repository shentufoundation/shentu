package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/nft/keeper"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"

	"github.com/certikfoundation/shentu/x/nft/types"
)

type msgServer struct {
	nfttypes.MsgServer
	Keeper
}

func (m msgServer) CreateAdmin(ctx context.Context, msg *types.MsgCreateAdmin) (*types.MsgIssueAdminResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}
	_, err = m.certKeeper.GetCertifier(sdkContext, creator)
	if err != nil {
		return nil, err
	}
	admin, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}
	m.SetAdmin(sdkContext, admin)

	sdkContext.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateAdmin,
			sdk.NewAttribute(types.AttributeKeyAdminCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyCreated, msg.Address),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, nfttypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	})
	return &types.MsgIssueAdminResponse{}, nil
}

func (m msgServer) RevokeAdmin(ctx context.Context, msg *types.MsgRevokeAdmin) (*types.MsgRevokeAdminResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	creator, err := sdk.AccAddressFromBech32(msg.Revoker)
	if err != nil {
		return nil, err
	}
	_, err = m.certKeeper.GetCertifier(sdkContext, creator)
	if err != nil {
		return nil, err
	}
	admin, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}
	m.DeleteAdmin(sdkContext, admin)

	sdkContext.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRevokeAdmin,
			sdk.NewAttribute(types.AttributeKeyAdminRevoker, msg.Revoker),
			sdk.NewAttribute(types.AttributeKeyRevoked, msg.Address),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, nfttypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Revoker),
		),
	})
	return &types.MsgRevokeAdminResponse{}, nil
}

var _ types.MsgServer = msgServer{}

func (k Keeper) NewMsgServerImpl() types.MsgServer {
	basicServer := keeper.NewMsgServerImpl(k.Keeper)
	return msgServer{
		MsgServer: basicServer,
		Keeper:    k,
	}
}
