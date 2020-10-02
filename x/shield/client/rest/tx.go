package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/gov"

	"github.com/certikfoundation/shentu/x/shield/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/shield/deposit_collateral", depositCollateralHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/shield/withdraw_collateral", withdrawCollateralHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/shield/withdraw_rewards", withdrawRewardsHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/shield/withdraw_foreign_rewards", withdrawForeignRewardsHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/shield/purchase", purchaseHandlerFn(cliCtx)).Methods("POST")
}

type depositCollateralReq struct {
	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
	PoolID     uint64       `json:"pool_id" yaml:"pool_id"`
	Collateral sdk.Coins    `json:"collateral" yaml:"collateral"`
}

func depositCollateralHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req depositCollateralReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
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

		msg := types.NewMsgDepositCollateral(from, req.PoolID, req.Collateral)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

type withdrawCollateralReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	PoolID  uint64       `json:"pool_id" yaml:"pool_id"`
	Amount  sdk.Coins    `json:"amount" yaml:"amount"`
}

func withdrawCollateralHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req withdrawCollateralReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
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

		msg := types.NewMsgWithdrawCollateral(from, req.PoolID, req.Amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{})
	}
}

type withdrawRewardsReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
}

func withdrawRewardsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req withdrawRewardsReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
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

		msg := types.NewMsgWithdrawRewards(from)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

type withdrawForeignRewardsReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Denom   string       `json:"denom" yaml:"denom"`
	ToAddr  string       `json:"to_addr" yaml:"to_addr"`
}

func withdrawForeignRewardsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req withdrawForeignRewardsReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
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

		msg := types.NewMsgWithdrawForeignRewards(accAddr, req.Denom, req.ToAddr)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

type purchaseReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	PoolID      uint64       `json:"pool_id" yaml:"pool_id"`
	Shield      sdk.Coins    `json:"shield" yaml:"shield"`
	Description string       `json:"description" yaml:"description"`
}

func purchaseHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req purchaseReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
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
		msg := types.NewMsgPurchaseShield(req.PoolID, req.Shield, req.Description, from)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// ShieldClaimProposalReq defines a shield claim proposal request body.
type ShieldClaimProposalReq struct {
	BaseReq        rest.BaseReq `json:"base_req" yaml:"base_req"`
	PoolID         uint64       `json:"pool_id" yaml:"pool_id"`
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
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
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
