package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/certikfoundation/shentu/x/oracle/types"
)

func registerTxHandlers(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/create-operator", types.ModuleName), createOperatorHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/remove-operator", types.ModuleName), removeOperatorHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/deposit-collateral", types.ModuleName), depositCollateralHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/withdraw-collateral", types.ModuleName), withdrawCollateralHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/claim-reward", types.ModuleName), claimRewardHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/create-task", types.ModuleName), createTaskHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/respond-to-task", types.ModuleName), respondToTaskHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/inquiry-task", types.ModuleName), inquireTaskHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/delete-task", types.ModuleName), deleteTaskHandler(cliCtx)).Methods("POST")
}

func createTaskHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createTaskReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		bounty, err := sdk.ParseCoinsNormalized(req.Bounty)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bounty[0].Amount.IsPositive() {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "bounty amount is required to be positive")
			return
		}

		creator, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		wait, err := strconv.ParseInt(req.Wait, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		hours, err := strconv.ParseInt(req.ValidDuration, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		validDuration := time.Duration(hours) * time.Hour

		msg := types.NewMsgCreateTask(req.Contract, req.Function, bounty, req.Description, creator, wait, validDuration)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
	}
}

func inquireTaskHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req inquiryTaskReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		inquirer, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgInquiryTask(req.Contract, req.Function, req.TxHash, inquirer)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
	}
}

func createOperatorHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createOperatorReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		address, err := sdk.AccAddressFromBech32(req.Address)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposer, err := sdk.AccAddressFromBech32(req.Proposer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCreateOperator(address, req.Collateral, proposer, req.Name)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
	}
}

func removeOperatorHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req removeOperatorReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		address, err := sdk.AccAddressFromBech32(req.Address)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposer, err := sdk.AccAddressFromBech32(req.Proposer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgRemoveOperator(address, proposer)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
	}
}

func depositCollateralHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req depositCollateralReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		address, err := sdk.AccAddressFromBech32(req.Address)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgAddCollateral(address, req.CollateralIncrement)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
	}
}

func withdrawCollateralHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req withdrawCollateralReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		address, err := sdk.AccAddressFromBech32(req.Address)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgReduceCollateral(address, req.CollateralDecrement)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
	}
}

func claimRewardHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req claimRewardReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		address, err := sdk.AccAddressFromBech32(req.Address)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgWithdrawReward(address)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
	}
}

func respondToTaskHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req respondToTaskReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		score, err := strconv.ParseInt(req.Score, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		operator, err := sdk.AccAddressFromBech32(req.Operator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgTaskResponse(req.Contract, req.Function, score, operator)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
	}
}

func deleteTaskHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req deleteTaskReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		deleter, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		force, err := strconv.ParseBool(req.Force)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgDeleteTask(req.Contract, req.Function, force, deleter)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
	}
}
