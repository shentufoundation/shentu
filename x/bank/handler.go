package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/certikfoundation/shentu/x/bank/internal/types"
)

// NewHandler returns a handler for "auth" type messages.
func NewHandler(k Keeper, ak types.AccountKeeper) sdk.Handler {
	cosmosHandler := bank.NewHandler(k)
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case *bankTypes.MsgSend:
			res, err := msgServer.Send(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *bankTypes.MsgMultiSend:
			res, err := msgServer.MultiSend(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgLockedSend:
			res, err := msgServer.LockedSend(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return cosmosHandler(ctx, msg)
		}
	}
}

func handleMsgLockedSend(ctx sdk.Context, k Keeper, ak types.AccountKeeper, msg types.MsgLockedSend) (*sdk.Result, error) {
}
