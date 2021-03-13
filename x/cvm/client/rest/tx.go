package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/gorilla/mux"

	"github.com/tendermint/crypto/sha3"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/txs/payload"

	"github.com/certikfoundation/shentu/x/cvm/types"
)

func registerTxHandlers(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/call", types.QuerierRoute), callHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/deploy", types.QuerierRoute), deployHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/view", types.QuerierRoute), viewHandler(cliCtx)).Methods("POST")
}

func callHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req callReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		value, err := sdk.ParseUint(req.Value)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		data, err := hex.DecodeString(req.Data)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "cannot decode call data")
		}

		msg := types.NewMsgCall(req.BaseReq.From, req.Callee, value.Uint64(), data)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, &msg)
	}
}
func deployHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req deployReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		value, err := sdk.ParseUint(req.Value)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		code, err := hex.DecodeString(req.Code)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		data, err := hex.DecodeString(req.Data)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		code = append(code, data...)

		abi, err := hex.DecodeString(req.Abi)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var metas []*payload.ContractMeta

		for _, hexMetaData := range req.Meta {
			metasBytes, err := hex.DecodeString(hexMetaData)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			runtime, err := hex.DecodeString(req.DeployedCode)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			hash := sha3.NewLegacyKeccak256()
			hash.Write(runtime)

			var codeHash acmstate.CodeHash
			copy(codeHash[:], hash.Sum(nil))
			metas = append(metas, &payload.ContractMeta{
				CodeHash: codeHash.Bytes(),
				Meta:     string(metasBytes),
			})
		}

		msg := types.NewMsgDeploy(req.BaseReq.From, value.Uint64(), code, string(abi), metas, req.IsEWASM, req.IsRuntime)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, &msg)
	}
}

func viewHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req viewReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		// Caller is optional parameter, if "" it becomes the zero address
		if req.Caller != "" {
			_, err := sdk.AccAddressFromBech32(req.Caller)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		_, err := sdk.AccAddressFromBech32(req.Callee)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		data, err := hex.DecodeString(req.Data)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Cannot decode call data")
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/view/%s/%s", types.QuerierRoute, req.Caller, req.Callee), data)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
