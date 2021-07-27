package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/certikfoundation/shentu/x/cert/types"
)

func registerTxHandlers(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/propose/certifier", types.ModuleName),
		proposeCertifierHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/certify/validator", types.ModuleName),
		certifyValidatorHandler(cliCtx)).Methods("POST")
}

func proposeCertifierHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req proposeCertifierReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		proposer, err := sdk.AccAddressFromBech32(req.Proposer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		certifier, err := sdk.AccAddressFromBech32(req.Certifier)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgProposeCertifier(proposer, certifier, req.Alias, req.Description)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func certifyValidatorHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req certifyValidatorReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		certifier, err := sdk.AccAddressFromBech32(req.Certifier)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		validator, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, req.Validator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg, err := types.NewMsgCertifyValidator(certifier, validator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
