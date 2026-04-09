package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// QueryCertifier is the query endpoint for certifier information.
	QueryCertifier = "certifier"

	// QueryCertifiers is the query endpoint for all certifiers information.
	QueryCertifiers = "certifiers"

	// QueryCertificate is the query endpoint for a certificate.
	QueryCertificate = "certificate"

	// QueryCertificateType is the query endpoint for a certificate type.
	QueryCertificateType = "certificateType"

	// QueryCertificates is the query endpoint for certificates.
	QueryCertificates = "certificates"
)

// QueryCertificatesParams is the type for parameters of querying certificates.
type QueryCertificatesParams struct {
	Page            int
	Limit           int
	Certifier       sdk.AccAddress
	CertificateType CertificateType
}

// QueryResCertifiers is the query result payload for all certifiers.
type QueryResCertifiers struct {
	Certifiers Certifiers `json:"certifiers"`
}

// String implements fmt.Stringer.
func (q QueryResCertifiers) String() string {
	return q.Certifiers.String()
}

// NewQueryCertificatesParams creates a new instance of QueryCertificatesParams.
func NewQueryCertificatesParams(page, limit int, certifier sdk.AccAddress, certType CertificateType) QueryCertificatesParams {
	return QueryCertificatesParams{
		Page:            page,
		Limit:           limit,
		Certifier:       certifier,
		CertificateType: certType,
	}
}
