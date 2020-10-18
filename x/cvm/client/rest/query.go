package rest

import (
	"fmt"
	"net/http"

	"github.com/certikfoundation/shentu/x/cvm/internal/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/code/{address}", types.QuerierRoute), codeHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/storage/{address}/{key}", types.QuerierRoute), storageHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/abi/{address}", types.QuerierRoute), abiHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/address-meta/{address}", types.QuerierRoute), addressMetaHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/meta/{hash}", types.QuerierRoute), metaHandler(cliCtx)).Methods("GET")
}

func codeHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD

=======
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

<<<<<<< HEAD
		route := fmt.Sprintf("custom/%s/code/%s", types.QuerierRoute, types.QueryAddress)
=======
		vars := mux.Vars(r)
		address := vars["address"]

		route := fmt.Sprintf("custom/%s/code/%s", types.QuerierRoute, address)
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func storageHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD

=======
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

<<<<<<< HEAD
		route := fmt.Sprintf("custom/%s/storage/%s/%s", types.QuerierRoute, types.QueryAddress, types.QueryKey)
=======
		vars := mux.Vars(r)
		address := vars["address"]
		key := vars["key"]

		route := fmt.Sprintf("custom/%s/storage/%s/%s", types.QuerierRoute, address, key)
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func abiHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD

=======
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

<<<<<<< HEAD
		route := fmt.Sprintf("custom/%s/abi/%s", types.QuerierRoute, types.QueryAddress)
=======
		vars := mux.Vars(r)
		address := vars["address"]

		route := fmt.Sprintf("custom/%s/abi/%s", types.QuerierRoute, address)
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func addressMetaHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD

=======
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

<<<<<<< HEAD
		route := fmt.Sprintf("custom/%s/address-meta/%s", types.QuerierRoute, types.QueryAddress)
=======
		vars := mux.Vars(r)
		address := vars["address"]

		route := fmt.Sprintf("custom/%s/address-meta/%s", types.QuerierRoute, address)
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func metaHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD

=======
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

<<<<<<< HEAD
		route := fmt.Sprintf("custom/%s/meta/%s", types.QuerierRoute, types.QueryHash)
=======
		vars := mux.Vars(r)
		hash := vars["hash"]

		route := fmt.Sprintf("custom/%s/meta/%s", types.QuerierRoute, hash)
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
