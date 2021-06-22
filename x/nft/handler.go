package nft

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	nfttypes "github.com/irisnet/irismod/modules/nft/types"

	"github.com/certikfoundation/shentu/x/nft/keeper"
	"github.com/certikfoundation/shentu/x/nft/types"
)

// NewHandler routes the messages to the handlers
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := k.NewMsgServerImpl()

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *nfttypes.MsgIssueDenom:
			if !k.CheckAdmin(ctx, msg.Sender) {
				return sdk.WrapServiceResult(ctx, nil, types.ErrAdminNotFound)
			}
			res, err := msgServer.IssueDenom(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *nfttypes.MsgMintNFT:
			if !k.CheckAdmin(ctx, msg.Sender) {
				return sdk.WrapServiceResult(ctx, nil, types.ErrAdminNotFound)
			}
			res, err := msgServer.MintNFT(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *nfttypes.MsgTransferNFT:
			res, err := msgServer.TransferNFT(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *nfttypes.MsgEditNFT:
			res, err := msgServer.EditNFT(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *nfttypes.MsgBurnNFT:
			res, err := msgServer.BurnNFT(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgCreateAdmin:
			res, err := msgServer.CreateAdmin(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgRevokeAdmin:
			res, err := msgServer.RevokeAdmin(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgIssueCertificate:
			res, err := msgServer.IssueCertificate(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgRevokeCertificate:
			res, err := msgServer.RevokeCertificate(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized nft message type: %T", msg)
		}
	}
}
