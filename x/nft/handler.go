package nft

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/nft/keeper"
	"github.com/irisnet/irismod/modules/nft/types"

	customkeeper "github.com/certikfoundation/shentu/x/nft/keeper"
	customtypes "github.com/certikfoundation/shentu/x/nft/types"
)

// NewHandler routes the messages to the handlers
func NewHandler(k customkeeper.Keeper) sdk.Handler {
	msgServer := k.NewMsgServerImpl()

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgIssueDenom:
			res, err := msgServer.IssueDenom(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgMintNFT:
			res, err := msgServer.MintNFT(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgTransferNFT:
			res, err := msgServer.TransferNFT(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgEditNFT:
			res, err := msgServer.EditNFT(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgBurnNFT:
			res, err := msgServer.BurnNFT(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *customtypes.MsgCreateNFTAdmin:
			res, err := msgServer.CreateNFTAdmin(sdk.WrapSDKContext(ctx), msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized nft message type: %T", msg)
		}
	}
}
