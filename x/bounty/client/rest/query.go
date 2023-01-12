package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/bounty/hosts", queryHostsHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/bounty/programs", queryProgramsWithParameterFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/bounty/programs/%s", RestProgramID), queryProgramsHandlerFn(clientCtx)).Methods("GET")
}

func queryHostsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	// TODO: implement this
	return nil
}

func queryProgramsWithParameterFn(clientCtx client.Context) http.HandlerFunc {
	// TODO: implement this
	return nil
}

func queryProgramsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	// TODO: implement this
	return nil
}
