package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/certikfoundation/shentu/x/shield/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
<<<<<<< HEAD
	r.HandleFunc(fmt.Sprintf("/%s/pool/id/{poolID}", types.QuerierRoute), queryPoolWithIDHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/pool/sponsor/{sponsor}", types.QuerierRoute), queryPoolWithSponsorHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/pools", types.QuerierRoute), queryPoolsHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/collaterals/{address}", types.QuerierRoute), queryCollateralsHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/purchase/{purchasetxhash}", types.QuerierRoute), queryPurchaseHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/purchases/{address}", types.QuerierRoute), queryPurchasesHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/pool/{poolID}/purchases", types.QuerierRoute), queryPoolPurchasesHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/pool/{poolID}/collaterals", types.QuerierRoute), queryPoolCollateralsHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/provider/{address}", types.QuerierRoute), queryProviderHandler(cliCtx)).Methods("GET")
=======
	r.HandleFunc(fmt.Sprintf("/%s/pools", types.QuerierRoute), queryPoolsHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/pool/id/{poolID}", types.QuerierRoute), queryPoolWithIDHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/pool/sponsor/{sponsor}", types.QuerierRoute), queryPoolWithSponsorHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/pool/{poolID}/collaterals", types.QuerierRoute), queryPoolCollateralsHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/pool/{poolID}/purchases", types.QuerierRoute), queryPoolPurchasesHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/pool/{poolID}/purchaser/{address}/purchases", types.QuerierRoute), queryPurchaseListHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/provider/{address}", types.QuerierRoute), queryProviderHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/purchaser/{address}/purchases", types.QuerierRoute), queryPurchaserPurchasesHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/collaterals/{address}", types.QuerierRoute), queryProviderCollateralsHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/pool_params", types.QuerierRoute), queryPoolParamsHandler(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/claim_params", types.QuerierRoute), queryClaimParamsHandler(cliCtx)).Methods("GET")
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
}

func queryPoolWithIDHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD

=======
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

<<<<<<< HEAD
		route := fmt.Sprintf("custom/%s/pool/id/%s", types.QuerierRoute, types.QueryPoolID)
=======
		vars := mux.Vars(r)
		poolID := vars["poolID"]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryPoolByID, poolID)
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

func queryPoolWithSponsorHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD

=======
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

<<<<<<< HEAD
		route := fmt.Sprintf("custom/%s/pool/sponsor/%s", types.QuerierRoute, types.QuerySponsor)
=======
		vars := mux.Vars(r)
		sponsor := vars["sponsor"]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryPoolBySponsor, sponsor)
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

func queryPoolsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

<<<<<<< HEAD
		route := fmt.Sprintf("custom/%s/pools", types.QuerierRoute)
=======
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPools)
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

<<<<<<< HEAD
func queryCollateralsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

=======
func queryProviderCollateralsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

<<<<<<< HEAD
		route := fmt.Sprintf("custom/%s/collaterals/%s", types.QuerierRoute, types.QueryAddress)
=======
		vars := mux.Vars(r)
		address := vars["address"]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryProviderCollaterals, address)
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

func queryPoolCollateralsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD

=======
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

<<<<<<< HEAD
		route := fmt.Sprintf("custom/%s/pool_collaterals/%s", types.QuerierRoute, types.QueryPoolID)
=======
		vars := mux.Vars(r)
		poolID := vars["poolID"]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryPoolCollaterals, poolID)
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

<<<<<<< HEAD
func queryPurchaseHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

=======
func queryPurchaseListHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

<<<<<<< HEAD
		route := fmt.Sprintf("custom/%s/purchase/%s", types.QuerierRoute, types.QueryPurchaseTxHash)
=======
		vars := mux.Vars(r)
		poolID := vars["poolID"]
		address := vars["address"]

		route := fmt.Sprintf("custom/%s/%s/%s/%s", types.QuerierRoute, types.QueryPurchaseList, poolID, address)
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

<<<<<<< HEAD
func queryPurchasesHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

=======
func queryPurchaserPurchasesHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

<<<<<<< HEAD
		route := fmt.Sprintf("custom/%s/purchases/%s", types.QuerierRoute, types.QueryAddress)
=======
		vars := mux.Vars(r)
		address := vars["address"]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryPurchaserPurchases, address)
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

func queryPoolPurchasesHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD
=======
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		poolID := vars["poolID"]
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryPoolPurchases, poolID)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryProviderHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

<<<<<<< HEAD
		route := fmt.Sprintf("custom/%s/pool_purchases/%s", types.QuerierRoute, types.QueryPoolID)
=======
		vars := mux.Vars(r)
		address := vars["address"]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryProvider, address)
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

<<<<<<< HEAD
func queryProviderHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
=======
func queryPoolParamsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPoolParams)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryClaimParamsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

<<<<<<< HEAD
		route := fmt.Sprintf("custom/%s/provider/%s", types.QuerierRoute, types.QueryAddress)
=======
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryClaimParams)
>>>>>>> aec1be1eb7334c836d6078a8eb77e82d81a46a30
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
