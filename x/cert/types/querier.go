package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// QueryCertifier is the query endpoint for certifier information.
	QueryCertifier = "certifier"

	// QueryCertifiers is the query endpoint for all certifiers information.
	QueryCertifiers = "certifiers"

	// QueryCertifierByAlias
	QueryCertifierByAlias = "certifieralias"

	// QueryValidator is the query endpoint for validator node certification.
	QueryCertifiedValidator = "validator"

	// QueryValidators is the query endpoint for all certified validator nodes.
	QueryCertifiedValidators = "validators"
)

// QueryResCertifiers is the query result payload for all certifiers.
type QueryResCertifiers struct {
	Certifiers Certifiers `json:"certifiers"`
}

// String implements fmt.Stringer.
func (q QueryResCertifiers) String() string {
	return q.Certifiers.String()
}

// QueryResValidator is the query result payload for a certified validator query.
type QueryResValidator struct {
	Certifier sdk.AccAddress `json:"certifier"`
}

// String implements fmt.Stringer.
func (q QueryResValidator) String() string {
	return q.Certifier.String()
}

// QueryResValidators is the query result payload for all certified validators.
type QueryResValidators struct {
	Validators []string `json:"validators"`
}

// String implements fmt.Stringer.
func (q QueryResValidators) String() string {
	validatorBech32s := make([]string, len(q.Validators))
	validatorBech32s = append(validatorBech32s, q.Validators...)
	return strings.Join(validatorBech32s, ", ")
}
