package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// RegisterRoutes registers the routes in main application.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
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
