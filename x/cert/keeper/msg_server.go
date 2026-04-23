package keeper

import (
	"context"
	"strconv"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the cert MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) UpdateCertifier(goCtx context.Context, msg *types.MsgUpdateCertifier) (*types.MsgUpdateCertifierResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.Keeper.authority != msg.Authority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "invalid authority; expected %s, got %s", k.Keeper.authority, msg.Authority)
	}

	operation, err := types.AddOrRemoveFromProto(msg.Operation)
	if err != nil {
		return nil, err
	}

	certifierAddr, err := sdk.AccAddressFromBech32(msg.Certifier)
	if err != nil {
		return nil, err
	}

	certifier := types.NewCertifier(certifierAddr, msg.Description)
	if err := k.Keeper.UpdateCertifier(ctx, operation, certifier); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"update_certifier",
			sdk.NewAttribute("operation", operation.String()),
			sdk.NewAttribute("certifier", msg.Certifier),
			sdk.NewAttribute("description", msg.Description),
		),
	)

	return &types.MsgUpdateCertifierResponse{}, nil
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
		sdk.NewAttribute("revoked_certificate", certificate.ToString()),
		sdk.NewAttribute("revoke_description", msg.Description),
	)
	ctx.EventManager().EmitEvent(revokeEvent)

	return &types.MsgRevokeCertificateResponse{}, nil
}
