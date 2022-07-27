package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/certikfoundation/shentu/v2/x/cert/types"
)

func registerTxHandlers(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/propose/certifier", types.ModuleName),
		proposeCertifierHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/certify/platform", types.ModuleName),
		certifyPlatformHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/certify", types.ModuleName),
		issueCertificateHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/revoke/certificate", types.ModuleName),
		revokeCertificateHandler(cliCtx)).Methods("POST")
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

func issueCertificateHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req certifyGeneralReq
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

		var msg *types.MsgIssueCertificate
		content := types.AssembleContent(req.CertificateType, req.Content)
		certificateTypeString := strings.ToLower(req.CertificateType)
		if certificateTypeString == "compilation" {
			msg = types.NewMsgIssueCertificate(content, req.Compiler, req.BytecodeHash, req.Description, certifier)
		} else {
			msg = types.NewMsgIssueCertificate(content, "", "", req.Description, certifier)
		}
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func certifyPlatformHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req certifyPlatformReq
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

		var validator cryptotypes.PubKey
		err = cliCtx.Codec.UnmarshalJSON([]byte(req.Validator), validator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg, err := types.NewMsgCertifyPlatform(certifier, validator, req.Platform)
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

func revokeCertificateHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req revokeCertificateReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		revoker, err := sdk.AccAddressFromBech32(req.Revoker)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgRevokeCertificate(revoker, req.CertificateID, req.Description)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
