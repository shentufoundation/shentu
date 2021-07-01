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

func (k msgServer) IssueCertificate(goCtx context.Context, msg *types.MsgIssueCertificate) (*types.MsgIssueCertificateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	certificate := types.Certificate{
		Content:     msg.Content,
		Description: msg.Description,
		Certifier:   msg.Certifier,
	}

	if err := k.Keeper.IssueCertificate(ctx, msg.DenomId, msg.TokenId, msg.Name, msg.Uri, certificate); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueCertificate,
			sdk.NewAttribute("denom_id", msg.DenomId),
			sdk.NewAttribute("token_id", msg.TokenId),
			sdk.NewAttribute("name", msg.Name),
			sdk.NewAttribute("uri", msg.Uri),
			sdk.NewAttribute("content", msg.Content),
			sdk.NewAttribute("description", msg.Description),
			sdk.NewAttribute("certifier", msg.Certifier),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, nfttypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Certifier),
		),
	})

	return &types.MsgIssueCertificateResponse{}, nil
}

func (k msgServer) EditCertificate(goCtx context.Context, msg *types.MsgEditCertificate) (*types.MsgEditCertificateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	certificate := types.Certificate{
		Content:     msg.Content,
		Description: msg.Description,
		Certifier:   msg.Owner,
	}

	if err := k.Keeper.EditCertificate(ctx, msg.DenomId, msg.TokenId, msg.Name, msg.Uri, certificate); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditCertificate,
			sdk.NewAttribute("denom_id", msg.DenomId),
			sdk.NewAttribute("token_id", msg.TokenId),
			sdk.NewAttribute("name", msg.Name),
			sdk.NewAttribute("uri", msg.Uri),
			sdk.NewAttribute("content", msg.Content),
			sdk.NewAttribute("description", msg.Description),
			sdk.NewAttribute("owner", msg.Owner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, nfttypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})

	return &types.MsgEditCertificateResponse{}, nil
}

func (k msgServer) RevokeCertificate(goCtx context.Context, msg *types.MsgRevokeCertificate) (*types.MsgRevokeCertificateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	revokerAddr, err := sdk.AccAddressFromBech32(msg.Revoker)
	if err != nil {
		panic(err)
	}

	if err := k.Keeper.RevokeCertificate(ctx, msg.DenomId, msg.TokenId, revokerAddr); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRevokeCertificate,
			sdk.NewAttribute("denom_id", msg.DenomId),
			sdk.NewAttribute("token_id", msg.TokenId),
			sdk.NewAttribute("revoker", msg.Revoker),
			sdk.NewAttribute("description", msg.Description),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, nfttypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Revoker),
		),
	})

	return &types.MsgRevokeCertificateResponse{}, nil
}

var _ types.MsgServer = msgServer{}

func (k Keeper) NewMsgServerImpl() types.MsgServer {
	basicServer := keeper.NewMsgServerImpl(k.Keeper)
	return msgServer{
		MsgServer: basicServer,
		Keeper:    k,
	}
}
