package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// QueryCertificate is the query endpoint for a certificate.
	QueryCertificate = "certificate"

	// QueryCertificate is the query endpoint for a certificate type.
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

// NewQueryCertificatesParams creates a new instance of QueryCertificatesParams.
func NewQueryCertificatesParams(page, limit int, certifier sdk.AccAddress, CertType CertificateType) QueryCertificatesParams {
	return QueryCertificatesParams{
		Page:            page,
		Limit:           limit,
		Certifier:       certifier,
		CertificateType: CertType,
	}
}
