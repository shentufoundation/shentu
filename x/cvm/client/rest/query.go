package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	auth_types "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/certikfoundation/shentu/x/cvm/client/utils"
	"github.com/certikfoundation/shentu/x/cvm/internal/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/code/{address}", types.QuerierRoute), codeHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/storage/{address}/{key}", types.QuerierRoute), storageHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/abi/{address}", types.QuerierRoute), abiHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/address-meta/{address}", types.QuerierRoute), addressMetaHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/meta/{hash}", types.QuerierRoute), metaHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/contract/{address}", types.QuerierRoute), contractHandler(cliCtx)).Methods("GET")
}

func codeHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		address := vars["address"]

		route := fmt.Sprintf("custom/%s/code/%s", types.QuerierRoute, address)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func storageHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		address := vars["address"]
		key := vars["key"]

		route := fmt.Sprintf("custom/%s/storage/%s/%s", types.QuerierRoute, address, key)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func abiHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		address := vars["address"]

		route := fmt.Sprintf("custom/%s/abi/%s", types.QuerierRoute, address)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func addressMetaHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		address := vars["address"]

		route := fmt.Sprintf("custom/%s/address-meta/%s", types.QuerierRoute, address)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func metaHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		hash := vars["hash"]

		route := fmt.Sprintf("custom/%s/meta/%s", types.QuerierRoute, hash)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func contractHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		var baseAcc *auth_types.BaseAccount
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/cvm/account/%s", bech32addr), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx.LegacyAmino.MustUnmarshalJSON(res, &baseAcc)

		cvmAcc, err := utils.QueryCVMAccount(cliCtx, bech32addr, baseAcc)
		if err == nil {
			rest.PostProcessResponse(w, cliCtx, cvmAcc)
			return
		}

		rest.PostProcessResponse(w, cliCtx, baseAcc)
	}
}
