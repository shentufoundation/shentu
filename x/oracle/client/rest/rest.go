package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterRoutes(cliCtx client.CLIContext, r *mux.Router) {
	RegisterTxRoutes(cliCtx, r)
	RegisterQueryRoutes(cliCtx, r)
}

type inquiryTaskReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Contract string       `json:"contract"`
	Function string       `json:"function"`
	TxHash   string       `json:"txhash"`
}

type createOperatorReq struct {
	BaseReq    rest.BaseReq `json:"base_req"`
	Address    string       `json:"address"`
	Collateral sdk.Coins    `json:"collateral"`
	Proposer   string       `json:"proposer"`
	Name       string       `json:"name"`
}

type removeOperatorReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Address  string       `json:"address"`
	Proposer string       `json:"proposer"`
}

type depositCollateralReq struct {
	BaseReq             rest.BaseReq `json:"base_req"`
	Address             string       `json:"address"`
	CollateralIncrement sdk.Coins    `json:"collateral_increment"`
}

type withdrawCollateralReq struct {
	BaseReq             rest.BaseReq `json:"base_req"`
	Address             string       `json:"address"`
	CollateralDecrement sdk.Coins    `json:"collateral_decrement"`
}

type claimRewardReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Address string       `json:"address"`
}

type createTaskReq struct {
	BaseReq       rest.BaseReq `json:"base_req"`
	Contract      string       `json:"contract"`
	Function      string       `json:"function"`
	Bounty        string       `json:"bounty"`
	Description   string       `json:"description"`
	Wait          string       `json:"wait"`
	ValidDuration string       `json:"valid_duration"`
}

type respondToTaskReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Contract string       `json:"contract"`
	Function string       `json:"function"`
	Score    string       `json:"score"`
	Operator string       `json:"operator"`
}

type deleteTaskReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Contract string       `json:"contract"`
	Function string       `json:"function"`
	Force    string       `json:"force"`
}
