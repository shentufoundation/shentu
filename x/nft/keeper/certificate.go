package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/nft/types"
)

func (k Keeper) MarshalCertificate(ctx sdk.Context, certificate types.Certificate) string {
	return string(k.cdc.MustMarshalJSON(&certificate))
}

func (k Keeper) UnmarshalCertificate(ctx sdk.Context, tokenData string) types.Certificate {
	var certificate types.Certificate
	k.cdc.MustUnmarshalJSON([]byte(tokenData), &certificate)
	return certificate
}

// IssueCertificate issues a certificate.
func (k Keeper) IssueCertificate(ctx sdk.Context, denomID, tokenID, tokenNm, tokenURI string, certificate types.Certificate) error {
	if !k.certKeeper.IsCertifier(ctx, certificate.GetCertifier()) {
		return types.ErrUnqualifiedCertifier
	}
	certifier := certificate.GetCertifier()

	denomNm := types.GetCertDenomNm(denomID)
	if denomNm == "" {
		return types.ErrInvalidDenomID
	}
	if !k.HasDenomNm(ctx, denomNm) {
		if err := k.IssueDenom(ctx, denomID, denomNm, types.CertificateSchema, certifier); err != nil {
			return err
		}
	}

	tokenData := k.MarshalCertificate(ctx, certificate)
	return k.MintNFT(ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, certifier)
}

// RevokeCertificate revokes a certificate.
func (k Keeper) RevokeCertificate(ctx sdk.Context, denomID, tokenID string, revoker sdk.AccAddress) error {
	if !k.certKeeper.IsCertifier(ctx, revoker) {
		return types.ErrUnqualifiedRevoker
	}
	if types.GetCertDenomNm(denomID) == "" {
		return types.ErrInvalidDenomID
	}
	return k.BurnNFT(ctx, denomID, tokenID, revoker)
}
