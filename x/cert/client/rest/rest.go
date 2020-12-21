// Package rest defines the RESTful service for the cert module.
package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

// RegisterRoutes registers the routes in main application.
func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	RegisterTxRoutes(cliCtx, r)
	RegisterQueryRoutes(cliCtx, r)
}

type proposeCertifierReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Proposer    string       `json:"proposer"`
	Certifier   string       `json:"certifier"`
	Alias       string       `json:"alias"`
	Description string       `json:"description"`
}

type certifyValidatorReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Certifier string       `json:"certifier"`
	Validator string       `json:"validator"`
}

type certifyGeneralReq struct {
	BaseReq         rest.BaseReq `json:"base_req"`
	CertificateType string       `json:"certificate_type"`
	ContentType     string       `json:"content_type"`
	Content         string       `json:"content"`
	Description     string       `json:"description"`
	Certifier       string       `json:"certifier"`
}

type certifyCompilationReq struct {
	BaseReq        rest.BaseReq `json:"base_req"`
	SourceCodeHash string       `json:"source_code_hash"`
	Compiler       string       `json:"compiler"`
	BytecodeHash   string       `json:"bytecode_hash"`
	Description    string       `json:"description"`
}

type certifyPlatformReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Certifier string       `json:"certifier"`
	Validator string       `json:"validator"`
	Platform  string       `json:"platform"`
}

type revokeCertificateReq struct {
	BaseReq       rest.BaseReq `json:"base_req"`
	Revoker       string       `json:"revoker"`
	CertificateID string       `json:"certificate_id"`
	Description   string       `json:"description"`
}

// ProposalRESTHandler returns a ProposalRESTHandler that exposes the community pool spend REST handler with a given sub-route.
func ProposalRESTHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "certifier_update",
		Handler:  postProposalHandlerFn(cliCtx),
	}
}

func postProposalHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CertifierUpdateProposalReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		from, err := sdk.AccAddressFromHex(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		content := types.NewCertifierUpdateProposal(
			req.Title,
			req.Description,
			req.Certifier,
			req.Alias,
			from,
			req.AddOrRemove,
		)

		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, from)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
