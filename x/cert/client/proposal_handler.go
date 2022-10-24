package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/shentufoundation/shentu/v2/x/cert/client/cli"
	"github.com/shentufoundation/shentu/v2/x/cert/client/rest"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
