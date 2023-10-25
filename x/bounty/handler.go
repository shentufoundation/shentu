package bounty

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateProgram:
			res, err := msgServer.CreateProgram(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgEditProgram:
			res, err := msgServer.EditProgram(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgModifyProgramStatus:
			if msg.Status == types.ProgramStatusActive {
				res, err := msgServer.OpenProgram(sdk.WrapSDKContext(ctx), msg)
				return sdk.WrapServiceResult(ctx, res, err)
			}
			res, err := msgServer.CloseProgram(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSubmitFinding:
			res, err := msgServer.SubmitFinding(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgModifyFindingStatus:
			if msg.Status == types.FindingStatusClosed {
				res, err := msgServer.RejectFinding(sdk.WrapSDKContext(ctx), msg)
				return sdk.WrapServiceResult(ctx, res, err)
			}
			if msg.Status == types.FindingStatusConfirmed {
				res, err := msgServer.AcceptFinding(sdk.WrapSDKContext(ctx), msg)
				return sdk.WrapServiceResult(ctx, res, err)
			}
			res, err := msgServer.CancelFinding(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgReleaseFinding:
			res, err := msgServer.ReleaseFinding(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
