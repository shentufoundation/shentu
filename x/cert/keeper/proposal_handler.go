package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/cert/types"
)

// HandleCertifierUpdateProposal is a handler for executing a passed certifier update proposal
func HandleCertifierUpdateProposal(ctx sdk.Context, k Keeper, p types.CertifierUpdateProposal) error {
	switch p.AddOrRemove {
	case types.Add:
		if k.IsCertifier(ctx, p.Certifier) {
			return types.ErrCertifierAlreadyExists
		}
		if p.Alias != "" && k.HasCertifierAlias(ctx, p.Alias) {
			return types.ErrRepeatedAlias
		}

		certifier := types.NewCertifier(p.Certifier, p.Alias, p.Proposer, p.Description)
		k.SetCertifier(ctx, certifier)
		return nil
	case types.Remove:
		certifiers := k.GetAllCertifiers(ctx)
		if len(certifiers) == 1 {
			return types.ErrOnlyOneCertifier
		}
		return k.deleteCertifier(ctx, p.Certifier)
	default:
		return types.ErrAddOrRemove
	}
}
