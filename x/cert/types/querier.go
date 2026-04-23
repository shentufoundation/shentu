package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	qtypes "github.com/cosmos/cosmos-sdk/types/query"
)

const (
	// QueryCertifier is the query endpoint for certifier information.
	QueryCertifier = "certifier"

	// QueryCertifiers is the query endpoint for all certifiers information.
	QueryCertifiers = "certifiers"

	// QueryCertificate is the query endpoint for a certificate.
	QueryCertificate = "certificate"

	// QueryCertificates is the query endpoint for certificates.
	QueryCertificates = "certificates"
)

// QueryCertificatesParams is the type for parameters of querying certificates.
type QueryCertificatesParams struct {
	Pagination      *qtypes.PageRequest
	Certifier       sdk.AccAddress
	CertificateType CertificateType
	Content         string
}

// NewQueryCertificatesParams creates a new instance of QueryCertificatesParams.
func NewQueryCertificatesParams(page, limit int, certifier sdk.AccAddress, certType CertificateType, content ...string) QueryCertificatesParams {
	if page <= 0 {
		page = 1
	}
	if limit < 0 {
		limit = 0
	}
	var contentFilter string
	if len(content) > 0 {
		contentFilter = content[0]
	}
	return QueryCertificatesParams{
		Pagination: &qtypes.PageRequest{
			Offset:     uint64((page - 1) * limit), //nolint:gosec // page>=1, limit>=0
			Limit:      uint64(limit),              //nolint:gosec // limit>=0
			CountTotal: true,
		},
		Certifier:       certifier,
		CertificateType: certType,
		Content:         contentFilter,
	}
}
