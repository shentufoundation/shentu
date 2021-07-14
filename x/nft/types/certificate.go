package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetCertDenomNm returns the DenomNm of the certificate NFT, if valid.
func GetCertDenomNm(denomID string) string {
	switch denomID {
	case "certificateauditing":
		return "Auditing"
	case "certificateidentity":
		return "Identity"
	case "certificategeneral":
		return "General"
	default:
		return ""
	}
}

// GetCertifier returns certifier of the certificate.
func (c Certificate) GetCertifier() sdk.AccAddress {
	certifierAddr, err := sdk.AccAddressFromBech32(c.Certifier)
	if err != nil {
		panic(err)
	}
	return certifierAddr
}

const (
	CertificateSchema = `
	{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"title": "certificate-schema",
		"description": "Certificate NFT Schema",
		"type": "object",
		"properties": {
			"content": {
				"description": "content of certificate",
				"type": "string",
			},
			"description": {
				"description": "description of certificate",
				"type": "string",
			},
			"certifier": {
				"description": "certifier address",
				"type": "string",
			}
		},
		"additionalProperties": false,
		"required": [
			"content",
			"description",
			"certifier"
		]
	}`
)
