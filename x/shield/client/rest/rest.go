package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	registerQueryRoutes(cliCtx, r, storeName)
	registerTxRoutes(cliCtx, r)
}

// ProposalRESTHandler returns a ProposalRESTHandler that exposes the shield claim REST handler with a given sub-route.
func ProposalRESTHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "shield_claim",
		Handler:  postProposalHandlerFn(cliCtx),
	}
}
