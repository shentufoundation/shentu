package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rest"
	sdk "github.com/cosmos/cosmos-sdk/types"
	resttypes "github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
)

func RegisterHandlers(cliCtx client.Context, rtr *mux.Router) {
	r := rest.WithHTTPDeprecationHeaders(rtr)
	registerQueryRoutes(cliCtx, r)
	registerTxHandlers(cliCtx, r)
}

// ProposalRESTHandler returns a ProposalRESTHandler that exposes the shield claim REST handler with a given sub-route.
func ProposalRESTHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "shield_claim",
		Handler:  postProposalHandlerFn(cliCtx),
	}
}

type depositCollateralReq struct {
	BaseReq resttypes.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Coins         `json:"amount" yaml:"amount"`
}

type withdrawCollateralReq struct {
	BaseReq resttypes.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Coins         `json:"amount" yaml:"amount"`
}

type withdrawRewardsReq struct {
	BaseReq resttypes.BaseReq `json:"base_req" yaml:"base_req"`
}

type withdrawForeignRewardsReq struct {
	BaseReq resttypes.BaseReq `json:"base_req" yaml:"base_req"`
	Denom   string            `json:"denom" yaml:"denom"`
	ToAddr  string            `json:"to_addr" yaml:"to_addr"`
}

type withdrawReimbursementReq struct {
	BaseReq    resttypes.BaseReq `json:"base_req" yaml:"base_req"`
	ProposalID uint64            `json:"proposal_id" yaml:"proposal_id"`
}

type purchaseReq struct {
	BaseReq     resttypes.BaseReq `json:"base_req" yaml:"base_req"`
	PoolID      uint64            `json:"pool_id" yaml:"pool_id"`
	Shield      sdk.Coins         `json:"shield" yaml:"shield"`
	Description string            `json:"description" yaml:"description"`
}

type withdrawFromShieldReq struct {
	BaseReq resttypes.BaseReq `json:"base_req" yaml:"base_req"`
	PoolID  uint64            `json:"pool_id" yaml:"pool_id"`
	Amount  sdk.Coins         `json:"shield" yaml:"shield"`
}

// ShieldClaimProposalReq defines a shield claim proposal request body.
type ShieldClaimProposalReq struct {
	BaseReq     resttypes.BaseReq `json:"base_req" yaml:"base_req"`
	PoolID      uint64            `json:"pool_id" yaml:"pool_id"`
	Loss        sdk.Coins         `json:"loss" yaml:"loss"`
	Evidence    string            `json:"evidence" yaml:"evidence"`
	Description string            `json:"description" yaml:"description"`
	Deposit     sdk.Coins         `json:"deposit" yaml:"deposit"`
}
