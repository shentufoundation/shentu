package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/legacy/types"
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
		if k.IsCertifier(ctx, certifierAddr) {
			return types.ErrCertifierAlreadyExists
		}
		if p.Alias != "" && k.HasCertifierAlias(ctx, p.Alias) {
			return types.ErrRepeatedAlias
		}

		certifier := types.NewCertifier(certifierAddr, p.Alias, proposerAddr, p.Description)
		k.SetCertifier(ctx, certifier)
		return nil
	case types.Remove:
		certifiers := k.GetAllCertifiers(ctx)
		if len(certifiers) == 1 {
			return types.ErrOnlyOneCertifier
		}
		return k.deleteCertifier(ctx, certifierAddr)
	default:
		return types.ErrAddOrRemove
	}
}
