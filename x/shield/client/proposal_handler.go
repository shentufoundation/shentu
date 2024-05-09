package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/shentufoundation/shentu/v2/x/shield/client/cli"
)

var (
	// shield claim proposal handler
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal)
)
