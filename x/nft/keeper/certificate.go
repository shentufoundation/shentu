package keeper

import (
	"github.com/cosmos/cosmos-sdk/client"
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
func (k Keeper) IssueCertificate(
	ctx sdk.Context, denomID, tokenID, tokenNm, tokenURI string,
	certificate types.Certificate,
) error {
	certifier := certificate.GetCertifier()
	if !k.certKeeper.IsCertifier(ctx, certifier) {
		return types.ErrUnqualifiedCertifier
	}
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

// GetCertificatesFiltered gets certificates filtered.
func (k Keeper) GetCertificatesFiltered(ctx sdk.Context, params types.QueryCertificatesParams) (uint64, []types.Certificate, error) {
	certNFTs := k.GetNFTs(ctx, params.DenomID)
	filteredCertificates := []types.Certificate{}
	for i := 0; i < len(certNFTs); i++ {
		certificate := k.UnmarshalCertificate(ctx, certNFTs[i].GetData())
		if len(params.Certifier) == 0 || certificate.GetCertifier().Equals(params.Certifier) {
			filteredCertificates = append(filteredCertificates, certificate)
		}
	}

	// Post-processing
	start, end := client.Paginate(len(filteredCertificates), params.Page, params.Limit, 100)
	if start < 0 || end < 0 {
		filteredCertificates = []types.Certificate{}
	} else {
		filteredCertificates = filteredCertificates[start:end]
	}

	return uint64(len(filteredCertificates)), filteredCertificates, nil
}

// EditCertificate edits the certificate nft.
func (k Keeper) EditCertificate(
	ctx sdk.Context, denomID, tokenID, tokenNm,
	tokenURI string, certificate types.Certificate,
) error {
	denomNm := types.GetCertDenomNm(denomID)
	if denomNm == "" {
		return types.ErrInvalidDenomID
	}
	owner := certificate.GetCertifier()
	if !k.certKeeper.IsCertifier(ctx, owner) {
		return types.ErrUnqualifiedCertifier
	}
	tokenData := k.MarshalCertificate(ctx, certificate)
	return k.EditNFT(ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, owner)
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
