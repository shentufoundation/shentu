package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client"

	"github.com/certikfoundation/shentu/x/bank/internal/types"
)

// RegisterRoutes registers custom REST routes.
func RegisterRoutes(cliCtx client.CLIContext, r *mux.Router) {
	r.HandleFunc("/bank/accounts/{address}/locked_transfers", LockedSendRequestHandlerFn(cliCtx)).Methods("POST")
}

// LockedSendReq defines the properties of a send request's body.
type LockedSendReq struct {
	BaseReq  rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount   sdk.Coins    `json:"amount" yaml:"amount"`
	Unlocker string       `json:"unlocker" yaml:"unlocker"`
}

// LockedSendRequestHandlerFn is an http request handler to send coins
// to a manual vesting account and have them locked (vesting).
func LockedSendRequestHandlerFn(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		toAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req LockedSendReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		unlocker, err := sdk.AccAddressFromBech32(req.Unlocker)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgLockedSend(fromAddr, toAddr, unlocker, req.Amount)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
