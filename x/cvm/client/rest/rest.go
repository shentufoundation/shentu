package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterHandlers(clientCtx client.Context, rtr *mux.Router) {
	r := clientrest.WithHTTPDeprecationHeaders(rtr)
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)
}

type callReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Callee  string       `json:"callee"`
	Value   string       `json:"value"`
	Data    string       `json:"data"`
}

type deployReq struct {
	BaseReq      rest.BaseReq `json:"base_req"`
	Value        string       `json:"value"`
	Code         string       `json:"code"`
	DeployedCode string       `json:"deployed_code"`
	Data         string       `json:"data"`
	Abi          string       `json:"abi"`
	Meta         []string     `json:"meta"`
	IsEWASM      bool         `json:"is_ewasm"`
	IsRuntime    bool         `json:"is_runtime"`
}

type viewReq struct {
	Caller string `json:"caller"`
	Callee string `json:"callee"`
	Data   string `json:"data"`
}
