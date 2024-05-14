package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/shentufoundation/shentu/v2/x/cert/client/cli"
)

// param change proposal handler
var (
	LegacyProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal)
)
