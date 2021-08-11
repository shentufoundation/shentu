package rest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/certikfoundation/shentu/v2/x/cert/types"
)

type (
	// CertifierUpdateProposalReq defines a certifier update proposal request body.
	CertifierUpdateProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string            `json:"title" yaml:"title"`
		Description string            `json:"description" yaml:"description"`
		Certifier   sdk.AccAddress    `json:"certifier" yaml:"certifier"`
		Alias       string            `json:"alias" yaml:"alias"`
		AddOrRemove types.AddOrRemove `json:"add_or_remove" yaml:"add_or_remove"`
		Deposit     sdk.Coins         `json:"deposit" yaml:"deposit"`
	}
)
