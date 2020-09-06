// Package client specifies the client for the gov module.
package client

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/certikfoundation/shentu/x/gov/client/rest"
)

// RESTHandlerFn defines the rest handler.
type RESTHandlerFn func(context.CLIContext) rest.ProposalRESTHandler

// CLIHandlerFn defines the cli handler.
type CLIHandlerFn func(*codec.Codec) *cobra.Command

// ProposalHandler is the combined type for a proposal handler for both cli and rest.
type ProposalHandler struct {
	CLIHandler  CLIHandlerFn
	RESTHandler RESTHandlerFn
}

// NewProposalHandler creates a new ProposalHandler object.
func NewProposalHandler(cliHandler CLIHandlerFn, restHandler RESTHandlerFn) ProposalHandler {
	return ProposalHandler{
		CLIHandler:  cliHandler,
		RESTHandler: restHandler,
	}
}
