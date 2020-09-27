package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"

	"github.com/certikfoundation/shentu/x/shield/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/shield/create_pool", createPoolHandlerFn(cliCtx)).Methods("POST")
}

type createPoolReq struct {
	BaseReq          rest.BaseReq     `json:"base_req" yaml:"base_req"`
	Shield           sdk.Coins        `json:"shield" yaml:"shield"`
	Deposit          types.MixedCoins `json:"deposit" yaml:"deposit"`
	Sponsor          string           `json:"sponsor" yaml:"sponsor"`
	TimeOfCoverage   int64            `json:"time_of_coverage" yaml:"time_of_coverage"`
	BlocksOfCoverage int64            `json:"blocks_of_coverage" yaml:"blocks_of_coverage"`
}

func createPoolHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createPoolReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		accAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		msg, err := types.NewMsgCreatePool(accAddr, req.Shield, req.Deposit, req.Sponsor, req.TimeOfCoverage, req.BlocksOfCoverage)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// ProposalRESTHandler returns a ProposalRESTHandler that exposes the shield claim REST handler with a given sub-route.
func ProposalRESTHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "shield_claim",
		Handler:  postProposalHandlerFn(cliCtx),
	}
}

// ShieldClaimProposalReq defines a shield claim proposal request body.
type ShieldClaimProposalReq struct {
	BaseReq        rest.BaseReq `json:"base_req" yaml:"base_req"`
	PoolID         int64        `json:"pool_id" yaml:"pool_id"`
	Loss           sdk.Coins    `json:"loss" yaml:"loss"`
	Evidence       string       `json:"evidence" yaml:"evidence"`
	PurchaseTxHash string       `json:"purchase_txash" yaml:"purchase_txash"`
	Description    string       `json:"description" yaml:"description"`
	Deposit        sdk.Coins    `json:"deposit" yaml:"deposit"`
}

func postProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ShieldClaimProposalReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
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

		content := types.NewShieldClaimProposal(req.PoolID, req.Loss, req.Evidence,
			req.PurchaseTxHash, req.Description, from, req.Deposit)

		msg := gov.NewMsgSubmitProposal(content, req.Deposit, from)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
