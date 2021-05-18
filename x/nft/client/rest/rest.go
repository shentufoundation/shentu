package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"

	nftrest "github.com/irisnet/irismod/modules/nft/client/rest"
)

// RegisterHandlers registers the NFT REST routes.
func RegisterHandlers(cliCtx client.Context, r *mux.Router, queryRoute string) {
	nftrest.RegisterHandlers(cliCtx, r, queryRoute)
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r, queryRoute)
}
