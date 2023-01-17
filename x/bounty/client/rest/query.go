package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// Get all delegations from a delegator
	r.HandleFunc(
		"/bounty/programs/{programID}",
		programsHandlerFn(clientCtx),
	).Methods("GET")
}

// HTTP request handler to query a delegator delegations
func programsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	// TODO: implement this
	return nil
}
