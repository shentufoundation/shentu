package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/certikfoundation/shentu/x/shield/types"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type createPoolReq struct {
	BaseReq          rest.BaseReq     `json:"base_req" yaml:"base_req"`
	Shield           sdk.Coins        `json:"shield" yaml:"shield"`
	Deposit          types.MixedCoins `json:"deposit" yaml:"deposit"`
	Sponsor          string           `json:"sponsor" yaml:"sponsor"`
	TimeOfCoverage   int64            `json:"time_of_coverage" yaml:"time_of_coverage"`
	BlocksOfCoverage int64            `json:"blocks_of_coverage" yaml:"blocks_of_coverage"`
}

// ShieldClaimProposalReq defines a shield claim proposal request body.
type ShieldClaimProposalReq struct {
	BaseReq        rest.BaseReq `json:"base_req" yaml:"base_req"`
	PoolID         uint64       `json:"pool_id" yaml:"pool_id"`
	Loss           sdk.Coins    `json:"loss" yaml:"loss"`
	Evidence       string       `json:"evidence" yaml:"evidence"`
	PurchaseTxHash string       `json:"purchase_txash" yaml:"purchase_txash"`
	Description    string       `json:"description" yaml:"description"`
	Deposit        sdk.Coins    `json:"deposit" yaml:"deposit"`
}

type withdrawForeignRewardsReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Denom   string       `json:"denom" yaml:"denom"`
	ToAddr  string       `json:"to_addr" yaml:"to_addr"`
}
