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
		isCertifier, err := k.IsCertifier(ctx, certifierAddr)
		if err != nil {
			return err
		}
		if !isCertifier {
			return types.ErrCertifierAlreadyExists
		}
		found, err := k.HasCertifierAlias(ctx, p.Alias)
		if err != nil {
			return err
		}
		if found {
			return types.ErrRepeatedAlias
		}

		certifier := types.NewCertifier(certifierAddr, p.Alias, proposerAddr, p.Description)
		if err := k.SetCertifier(ctx, certifier); err != nil {
			return err
		}
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
