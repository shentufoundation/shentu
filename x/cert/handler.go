package cert

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/x/cert/keeper"
	"github.com/certikfoundation/shentu/x/cert/types"
)

// NewHandler returns a handler for "cert" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCertifyValidator:
			res, err := msgServer.CertifyValidator(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgDecertifyValidator:
			res, err := msgServer.DecertifyValidator(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Unrecognized cert Msg type: %v", msg.Type())
		}
	}
}

func NewCertifierUpdateProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.CertifierUpdateProposal:
			return keeper.HandleCertifierUpdateProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized cert proposal content type: %T", c)
		}
	}
}
