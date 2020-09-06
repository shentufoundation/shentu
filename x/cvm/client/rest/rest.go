package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tendermint/crypto/sha3"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	clientrest "github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/txs/payload"

	"github.com/certikfoundation/shentu/x/cvm/internal/types"
)

// RegisterRoutes registers the routes in main application.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/call", storeName), callHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/deploy", storeName), deployHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/view", storeName), viewHandler(cliCtx, storeName)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/code/{address}", storeName), codeHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/storage/{address}/{key}", storeName), storageHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/abi/{address}", storeName), abiHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/address-meta/{address}", storeName), addressMetaHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/meta/{hash}", storeName), metaHandler(cliCtx, storeName)).Methods("GET")
}

type callReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Callee  string       `json:"callee"`
	Value   string       `json:"value"`
	Data    string       `json:"data"`
}

func callHandler(cliCtx context.CLIContext) http.HandlerFunc {
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

type deployReq struct {
	BaseReq      rest.BaseReq `json:"base_req"`
	Value        string       `json:"value"`
	Code         string       `json:"code"`
	DeployedCode string       `json:"deployed_code"`
	Data         string       `json:"data"`
	Abi          string       `json:"abi"`
	Meta         []string     `json:"meta"`
}

func deployHandler(cliCtx context.CLIContext) http.HandlerFunc {
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

		msg := types.NewMsgDeploy(caller, value.Uint64(), code, string(abi), metas)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func codeHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/code/%s", storeName, vars["address"]), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func storageHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/storage/%s/%s", storeName, vars["address"], vars["key"]), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func abiHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/abi/%s", storeName, vars["address"]), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

type viewReq struct {
	Caller string `json:"caller"`
	Callee string `json:"callee"`
	Data   string `json:"data"`
}

func viewHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
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

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/view/%s/%s", storeName, req.Caller, req.Callee), data)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func addressMetaHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/address-meta/%s", storeName, vars["address"]), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func metaHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/meta/%s", storeName, vars["hash"]), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
