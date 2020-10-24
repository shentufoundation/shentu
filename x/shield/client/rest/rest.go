package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

// ProposalRESTHandler returns a ProposalRESTHandler that exposes the shield claim REST handler with a given sub-route.
func ProposalRESTHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "shield_claim",
		Handler:  postProposalHandlerFn(cliCtx),
	}
}

type depositCollateralReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Coins    `json:"amount" yaml:"amount"`
}

type withdrawCollateralReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Coins    `json:"amount" yaml:"amount"`
}

type withdrawRewardsReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
}

type withdrawForeignRewardsReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Denom   string       `json:"denom" yaml:"denom"`
	ToAddr  string       `json:"to_addr" yaml:"to_addr"`
}

type withdrawReimbursementReq struct {
	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
	ProposalID uint64       `json:"proposal_id" yaml:"proposal_id"`
}

type purchaseReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	PoolID      uint64       `json:"pool_id" yaml:"pool_id"`
	Shield      sdk.Coins    `json:"shield" yaml:"shield"`
	Description string       `json:"description" yaml:"description"`
}

// ShieldClaimProposalReq defines a shield claim proposal request body.
type ShieldClaimProposalReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	PoolID      uint64       `json:"pool_id" yaml:"pool_id"`
	Loss        sdk.Coins    `json:"loss" yaml:"loss"`
	Evidence    string       `json:"evidence" yaml:"evidence"`
	PurchaseID  uint64       `json:"purchase_id" yaml:"purchase_id"`
	Description string       `json:"description" yaml:"description"`
	Deposit     sdk.Coins    `json:"deposit" yaml:"deposit"`
}
