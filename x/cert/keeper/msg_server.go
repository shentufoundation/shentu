package keeper

import (
	"context"
	"strconv"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/v2/x/cert/types"
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

func (k msgServer) ProposeCertifier(goCtx context.Context, msg *types.MsgProposeCertifier) (*types.MsgProposeCertifierResponse, error) {
	return &types.MsgProposeCertifierResponse{}, nil
}

func (k msgServer) IssueCertificate(goCtx context.Context, msg *types.MsgIssueCertificate) (*types.MsgIssueCertificateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	certificate := types.Certificate{
		Content:            msg.Content,
		CompilationContent: &types.CompilationContent{msg.Compiler, msg.BytecodeHash},
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

func (k msgServer) CertifyPlatform(goCtx context.Context, msg *types.MsgCertifyPlatform) (*types.MsgCertifyPlatformResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valPubKey, ok := msg.ValidatorPubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", valPubKey)
	}

	certifierAddr, err := sdk.AccAddressFromBech32(msg.Certifier)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.CertifyPlatform(ctx, certifierAddr, valPubKey, msg.Platform); err != nil {
		return nil, err
	}

	return &types.MsgCertifyPlatformResponse{}, nil
}
