package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/certikfoundation/shentu/x/cvm/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	clientrest "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/gorilla/mux"
	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/txs/payload"
	"github.com/tendermint/crypto/sha3"
)

func registerTxRoutes(cliCtx client.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/call", types.QuerierRoute), callHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/deploy", types.QuerierRoute), deployHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/view", types.QuerierRoute), viewHandler(cliCtx)).Methods("POST")
}

func callHandler(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req callReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		caller, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		callee, err := sdk.AccAddressFromBech32(req.Callee)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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

		msg := types.NewMsgCall(caller, callee, value.Uint64(), data)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
func deployHandler(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req deployReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		caller, err := sdk.AccAddressFromBech32(req.BaseReq.From)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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

		msg := types.NewMsgDeploy(caller, value.Uint64(), code, string(abi), metas, req.IsEWASM, req.IsRuntime)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func viewHandler(cliCtx client.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req viewReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
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
