package cert

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/x/cert/internal/keeper"
	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

// NewHandler returns a handler for "cert" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgCertifyValidator:
			return handleMsgCertifyValidator(ctx, k, msg)
		case types.MsgDecertifyValidator:
			return handleMsgDecertifyValidator(ctx, k, msg)
		case types.MsgCertifyPlatform:
			return handleMsgCertifyPlatform(ctx, k, msg)
		case types.MsgCertifyGeneral:
			return handleMsgCertifyGeneral(ctx, k, msg)
		case types.MsgCertifyCompilation:
			return handleMsgCertifyCompilation(ctx, k, msg)
		case types.MsgRevokeCertificate:
			return handleMsgRevokeCertificate(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Unrecognized cert Msg type: %v", msg.Type())
		}
	}
}

func handleMsgCertifyValidator(ctx sdk.Context, k Keeper, msg types.MsgCertifyValidator) (*sdk.Result, error) {
	if err := k.CertifyValidator(ctx, msg.Validator, msg.Certifier); err != nil {
		return nil, err
	}
	return &sdk.Result{}, nil
}

func handleMsgDecertifyValidator(ctx sdk.Context, k Keeper, msg types.MsgDecertifyValidator) (*sdk.Result, error) {
	if err := k.DecertifyValidator(ctx, msg.Validator, msg.Decertifier); err != nil {
		return nil, err
	}
	return &sdk.Result{}, nil
}

func handleMsgCertifyCompilation(ctx sdk.Context, k Keeper, msg types.MsgCertifyCompilation) (*sdk.Result, error) {
	certificate := types.NewCompilationCertificate(
		types.CertificateTypeCompilation,
		msg.SourceCodeHash,
		msg.Compiler,
		msg.BytecodeHash,
		msg.Description,
		msg.Certifier,
	)
	certificateID, err := k.IssueCertificate(ctx, certificate)
	if err != nil {
		return nil, err
	}
	certEvent := sdk.NewEvent(
		types.EventTypeCertifyCompilation,
		sdk.NewAttribute("certificate_id", certificateID.String()),
		sdk.NewAttribute("source_code_hash", msg.SourceCodeHash),
		sdk.NewAttribute("compiler", msg.Compiler),
		sdk.NewAttribute("bytecode_hash", msg.BytecodeHash),
		sdk.NewAttribute("certifier", msg.Certifier.String()),
	)
	ctx.EventManager().EmitEvent(certEvent)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgCertifyPlatform(ctx sdk.Context, k Keeper, msg types.MsgCertifyPlatform) (*sdk.Result, error) {
	if err := k.CertifyPlatform(ctx, msg.Certifier, msg.Validator, msg.Platform); err != nil {
		return nil, err
	}
	return &sdk.Result{}, nil
}

func handleMsgCertifyGeneral(ctx sdk.Context, k Keeper, msg types.MsgCertifyGeneral) (*sdk.Result, error) {
	certificate, err := types.NewGeneralCertificate(
		msg.CertificateType,
		msg.RequestContentType,
		msg.RequestContent,
		msg.Description,
		msg.Certifier,
	)
	if err != nil {
		return nil, err
	}
	certificateID, err := k.IssueCertificate(ctx, certificate)
	if err != nil {
		return nil, err
	}
	certEvent := sdk.NewEvent(
		types.EventTypeCertify,
		sdk.NewAttribute("certificate_id", certificateID.String()),
		sdk.NewAttribute("certificate_type", msg.CertificateType),
		sdk.NewAttribute("request_content_type", msg.RequestContentType),
		sdk.NewAttribute("request_content", msg.RequestContent),
		sdk.NewAttribute("description", msg.Description),
		sdk.NewAttribute("certifier", msg.Certifier.String()),
	)
	ctx.EventManager().EmitEvent(certEvent)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgRevokeCertificate(ctx sdk.Context, k Keeper, msg types.MsgRevokeCertificate) (*sdk.Result, error) {
	certificate, err := k.GetCertificateByID(ctx, msg.ID)
	if err != nil {
		return nil, err
	}
	if err := k.RevokeCertificate(ctx, certificate, msg.Revoker); err != nil {
		return nil, err
	}
	revokeEvent := sdk.NewEvent(
		types.EventTypeRevokeCertificate,
		sdk.NewAttribute("revoker", msg.Revoker.String()),
		sdk.NewAttribute("revoked_certificate", certificate.String()),
		sdk.NewAttribute("revoke_description", msg.Description),
	)
	ctx.EventManager().EmitEvent(revokeEvent)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func NewCertifierUpdateProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case types.CertifierUpdateProposal:
			return keeper.HandleCertifierUpdateProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized cert proposal content type: %T", c)
		}
	}
}
