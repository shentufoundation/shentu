package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/certikfoundation/shentu/v2/x/shield/client/cli"
	"github.com/certikfoundation/shentu/v2/x/shield/client/rest"
)

var (
	// shield claim proposal handler
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
