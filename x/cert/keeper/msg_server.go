package keeper

import (
	"context"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/cert/types"
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

func (k msgServer) CertifyValidator(goCtx context.Context, msg *types.MsgCertifyValidator) (*types.MsgCertifyValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valPubKey, ok := msg.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", valPubKey)
	}

	certifierAddr, err := sdk.AccAddressFromBech32(msg.Certifier)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.CertifyValidator(ctx, valPubKey, certifierAddr); err != nil {
		return nil, err
	}

	return &types.MsgCertifyValidatorResponse{}, nil
}

func (k msgServer) DecertifyValidator(goCtx context.Context, msg *types.MsgDecertifyValidator) (*types.MsgDecertifyValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valPubKey, ok := msg.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", valPubKey)
	}

	decertifierAddr, err := sdk.AccAddressFromBech32(msg.Decertifier)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.DecertifyValidator(ctx, valPubKey, decertifierAddr); err != nil {
		return nil, err
	}

	return &types.MsgDecertifyValidatorResponse{}, nil
}
