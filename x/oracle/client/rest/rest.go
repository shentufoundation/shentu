package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rest"
	sdk "github.com/cosmos/cosmos-sdk/types"
	resttypes "github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterHandlers(cliCtx client.Context, rtr *mux.Router) {
	r := rest.WithHTTPDeprecationHeaders(rtr)
	registerQueryRoutes(cliCtx, r)
	registerTxHandlers(cliCtx, r)
}

type createOperatorReq struct {
	BaseReq    resttypes.BaseReq `json:"base_req"`
	Address    string            `json:"address"`
	Collateral sdk.Coins         `json:"collateral"`
	Proposer   string            `json:"proposer"`
	Name       string            `json:"name"`
}

type removeOperatorReq struct {
	BaseReq  resttypes.BaseReq `json:"base_req"`
	Address  string            `json:"address"`
	Proposer string            `json:"proposer"`
}

type depositCollateralReq struct {
	BaseReq             resttypes.BaseReq `json:"base_req"`
	Address             string            `json:"address"`
	CollateralIncrement sdk.Coins         `json:"collateral_increment"`
}

type withdrawCollateralReq struct {
	BaseReq             resttypes.BaseReq `json:"base_req"`
	Address             string            `json:"address"`
	CollateralDecrement sdk.Coins         `json:"collateral_decrement"`
}

type claimRewardReq struct {
	BaseReq resttypes.BaseReq `json:"base_req"`
	Address string            `json:"address"`
}

type createTaskReq struct {
	BaseReq       resttypes.BaseReq `json:"base_req"`
	Contract      string            `json:"contract"`
	Function      string            `json:"function"`
	Bounty        string            `json:"bounty"`
	Description   string            `json:"description"`
	Wait          string            `json:"wait"`
	ValidDuration string            `json:"valid_duration"`
}

type respondToTaskReq struct {
	BaseReq  resttypes.BaseReq `json:"base_req"`
	Contract string            `json:"contract"`
	Function string            `json:"function"`
	Score    string            `json:"score"`
	Operator string            `json:"operator"`
}

type deleteTaskReq struct {
	BaseReq  resttypes.BaseReq `json:"base_req"`
	Contract string            `json:"contract"`
	Function string            `json:"function"`
	Force    string            `json:"force"`
}
