package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

func RegisterQueryRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/operator/{address}", types.QuerierRoute), operatorHandler(cliCtx)).Methods("Get")
	r.HandleFunc(fmt.Sprintf("/%s/operators", types.QuerierRoute), operatorsHandler(cliCtx)).Methods("Get")
	r.HandleFunc(fmt.Sprintf("/%s/withdraws", types.QuerierRoute), withdrawsHandler(cliCtx)).Methods("Get")

	r.HandleFunc(fmt.Sprintf("/%s/task", types.QuerierRoute), taskHandler(cliCtx)).Methods("Get")
	r.HandleFunc(fmt.Sprintf("/%s/response", types.QuerierRoute), responseHandler(cliCtx)).Methods("Get")
}

func operatorHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		address := vars["address"]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryOperator, address)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func operatorsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryOperators)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func withdrawsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryWithdrawals)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func taskHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		contract := r.URL.Query().Get("contract")
		if contract == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "contract is require to submit a task")
		}
		function := r.URL.Query().Get("function")
		if function == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "function is require to submit a task")
		}

		params := types.NewQueryTaskParams(contract, function)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTask)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func responseHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		contract := r.URL.Query().Get("contract")
		if contract == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "contract is require to respond to a task")
		}
		function := r.URL.Query().Get("function")
		if function == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "function is require to respond to a task")
		}
		var err error
		var operatorAddress sdk.AccAddress
		if operator := r.URL.Query().Get("operator"); operator != "" {
			operatorAddress, err = sdk.AccAddressFromBech32(operator)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		params := types.NewQueryResponseParams(contract, function, operatorAddress)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryResponse)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
