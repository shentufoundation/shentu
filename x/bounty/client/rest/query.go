package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/bounty/hosts", queryHostsHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/bounty/programs", queryProgramsWithParameterFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/bounty/programs/{%s}", RestProgramID), queryProgramHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/bounty/findings", queryFindingsWithParameterFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/bounty/findings/{%s}", RestFindingID), queryFindingsHandlerFn(clientCtx)).Methods("GET")
}

func queryHostsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	// TODO: implement this
	return func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("queryHostsHandlerFn Not Implemented")
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
	}
}

func queryProgramsWithParameterFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		// TODO add filter
		params := types.NewQueryProgramsParams(page, limit)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPrograms)
		res, height, err := clientCtx.QueryWithData(route, bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryProgramHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProgramID := vars[RestProgramID]

		if len(strProgramID) == 0 {
			err := errors.New("programID required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		programID, ok := rest.ParseUint64OrReturnBadRequest(w, strProgramID)
		if !ok {
			return
		}

		clientCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryProgramIDParams(programID)

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryProgram)
		res, height, err := clientCtx.QueryWithData(route, bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryFindingsWithParameterFn(clientCtx client.Context) http.HandlerFunc {
	// TODO: implement this
	return func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("queryFindingsWithParameterFn Not Implemented")
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
	}
}

func queryFindingsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	// TODO: implement this
	return func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("queryFindingsHandlerFn Not Implemented")
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
	}
}
