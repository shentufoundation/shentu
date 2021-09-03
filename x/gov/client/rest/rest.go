// Package rest defines the RESTful service for the gov module.
package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govRest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"

	"github.com/certikfoundation/shentu/v2/x/gov/types"
)

// REST Variable names
const (
	RestParamsType     = "type"
	RestProposalID     = "proposal-id"
	RestDepositor      = "depositor"
	RestVoter          = "voter"
	RestProposalStatus = "status"
)

// ProposalRESTHandler defines a REST handler implemented in another module. The
// sub-route is mounted on the governance REST handler.
type ProposalRESTHandler struct {
	SubRoute string
	Handler  func(http.ResponseWriter, *http.Request)
}

// RegisterRoutes is the central function to define routes that get registered by the main application.
func RegisterRoutes(cliCtx client.Context, r *mux.Router, phs []govRest.ProposalRESTHandler) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r, phs)
}

type VoteWithPower struct {
	types.Vote
	VotingPower sdk.Dec `json:"voting_power" yaml:"voting_power"`
}

type VotesWithPower []VoteWithPower
