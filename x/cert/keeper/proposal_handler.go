package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// HandleCertifierUpdateProposal is a handler for executing a passed certifier update proposal
func HandleCertifierUpdateProposal(ctx sdk.Context, k Keeper, p *types.CertifierUpdateProposal) error {
	certifierAddr, err := sdk.AccAddressFromBech32(p.Certifier)
	if err != nil {
		panic(err)
	}
	proposerAddr, err := sdk.AccAddressFromBech32(p.Proposer)
	if err != nil {
		panic(err)
	}

	switch p.AddOrRemove {
	case types.Add:
		certifier := types.NewCertifier(certifierAddr, proposerAddr, p.Description)
		return k.UpdateCertifier(ctx, types.Add, certifier)
	case types.Remove:
		certifier := types.NewCertifier(certifierAddr, proposerAddr, p.Description)
		return k.UpdateCertifier(ctx, types.Remove, certifier)
	default:
		return types.ErrAddOrRemove
	}
}
