package keeper

import (
	"context"
	"strconv"

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
		Content:            msg.Content,
		CompilationContent: &types.CompilationContent{Compiler: msg.Compiler, BytecodeHash: msg.BytecodeHash},
		Description:        msg.Description,
		Certifier:          msg.Certifier,
	}

	certificateID, err := k.Keeper.IssueCertificate(ctx, certificate)
	if err != nil {
		return nil, err
	}
	certEvent := sdk.NewEvent(
		types.EventTypeCertify,
		sdk.NewAttribute("certificate_id", strconv.FormatUint(certificateID, 10)),
		sdk.NewAttribute("certificate_type", types.TranslateCertificateType(certificate).String()),
		sdk.NewAttribute("content", certificate.GetContentString()),
		sdk.NewAttribute("compiler", msg.Compiler),
		sdk.NewAttribute("bytecode_hash", msg.BytecodeHash),
		sdk.NewAttribute("description", msg.Description),
		sdk.NewAttribute("certifier", msg.Certifier),
	)
	ctx.EventManager().EmitEvent(certEvent)

	return &types.MsgIssueCertificateResponse{}, nil
}

func (k msgServer) RevokeCertificate(goCtx context.Context, msg *types.MsgRevokeCertificate) (*types.MsgRevokeCertificateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	certificate, err := k.Keeper.GetCertificateByID(ctx, msg.Id)
	if err != nil {
		return nil, err
	}

	revokerAddr, err := sdk.AccAddressFromBech32(msg.Revoker)
	if err != nil {
		panic(err)
	}

	if err := k.Keeper.RevokeCertificate(ctx, certificate, revokerAddr); err != nil {
		return nil, err
	}
	revokeEvent := sdk.NewEvent(
		types.EventTypeRevokeCertificate,
		sdk.NewAttribute("revoker", msg.Revoker),
		sdk.NewAttribute("revoked_certificate", certificate.String()),
		sdk.NewAttribute("revoke_description", msg.Description),
	)
	ctx.EventManager().EmitEvent(revokeEvent)

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
