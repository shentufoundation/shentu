package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// QueryCertifier is the query endpoint for certifier information.
	QueryCertifier = "certifier"

	// QueryCertifiers is the query endpoint for all certifiers information.
	QueryCertifiers = "certifiers"

	// QueryCertifierByAlias
	QueryCertifierByAlias = "certifieralias"
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
