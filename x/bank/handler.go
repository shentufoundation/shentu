package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/certikfoundation/shentu/x/bank/internal/keeper"
	"github.com/certikfoundation/shentu/x/bank/internal/types"
)

// NewHandler returns a handler for "auth" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case *types.MsgLockedSend:
			res, err := msgServer.LockedSend(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			cosmosHandler := bank.NewHandler(k)
			return cosmosHandler(ctx, msg)
		}
	}
}
