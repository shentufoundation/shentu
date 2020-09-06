package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	clientrest "github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

func RegisterQueryRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/operator/{address}", storeName), operatorHandler(cliCtx, storeName)).Methods("Get")
	r.HandleFunc(fmt.Sprintf("/%s/operators", storeName), operatorsHandler(cliCtx, storeName)).Methods("Get")
	r.HandleFunc(fmt.Sprintf("/%s/withdraws", storeName), withdrawsHandler(cliCtx, storeName)).Methods("Get")

	r.HandleFunc(fmt.Sprintf("/%s/task", storeName), taskHandler(cliCtx, storeName)).Methods("Get")
	r.HandleFunc(fmt.Sprintf("/%s/response", storeName), responseHandler(cliCtx, storeName)).Methods("Get")
}

func operatorHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/operator/%s", storeName, vars["address"]), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func operatorsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/operators", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func withdrawsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/withdraws", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func createTaskHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createTaskReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		bounty, err := sdk.ParseCoins(req.Bounty)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bounty[0].Amount.IsPositive() {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "bounty amount is required to be positive")
			return
		}

		layout := "2006-01-02T03:04:05.000Z"
		var expiration time.Time
		if req.Expiration != "" {
			expiration, err = time.Parse(layout, req.Expiration)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		} else {
			expiration = time.Time{}
		}

		creator, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		closingHeight, err := strconv.ParseInt(req.ClosingHeight, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if closingHeight == int64(0) {
			closingHeight = int64(-1)
		}

		msg := types.NewMsgCreateTask(req.Contract, req.Function, bounty, req.Description, expiration, creator, closingHeight)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func taskHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
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
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/task", storeName)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func responseHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
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
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/response", storeName)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
