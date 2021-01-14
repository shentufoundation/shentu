package types

import (
	"encoding/hex"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/gogo/protobuf/proto"
)

// CertificateTypes is an array of all certificate types.
var CertificateTypes = [...]CertificateType{
	CertificateTypeNil,
	CertificateTypeCompilation,
	CertificateTypeAuditing,
	CertificateTypeProof,
	CertificateTypeOracleOperator,
	CertificateTypeShieldPoolCreator,
	CertificateTypeIdentity,
	CertificateTypeGeneral,
}

// Bytes returns the byte array for a certificate type.
func (c CertificateType) Bytes() []byte {
	return []byte{byte(c)}
}

// TODO
// // String returns the string for a certificate type.
// func (c CertificateType) String() string {
// 	switch c {
// 	case CertificateTypeCompilation:
// 		return "Compilation"
// 	case CertificateTypeAuditing:
// 		return "Auditing"
// 	case CertificateTypeProof:
// 		return "Proof"
// 	case CertificateTypeOracleOperator:
// 		return "OracleOperator"
// 	case CertificateTypeShieldPoolCreator:
// 		return "ShieldPoolCreator"
// 	case CertificateTypeIdentity:
// 		return "Identity"
// 	case CertificateTypeGeneral:
// 		return "General"
// 	default:
// 		return "UnknownCertificateType"
// 	}
// }

// CertificateTypeFromString returns a certificate type by parsing a string.
func CertificateTypeFromString(s string) CertificateType {
	switch strings.ToUpper(s) {
	case "COMPILATION":
		return CertificateTypeCompilation
	case "AUDITING":
		return CertificateTypeAuditing
	case "PROOF":
		return CertificateTypeProof
	case "ORACLEOPERATOR":
		return CertificateTypeOracleOperator
	case "SHIELDPOOLCREATOR":
		return CertificateTypeShieldPoolCreator
	case "IDENTITY":
		return CertificateTypeIdentity
	case "GENERAL":
		return CertificateTypeGeneral
	default:
		return CertificateTypeNil
	}
}

// CertificateID is the type for the ID of a certificate.
type CertificateID string

// Bytes returns the byte array for a certificate ID.
func (id CertificateID) Bytes() []byte {
	decoded, err := hex.DecodeString(id.String())
	if err != nil {
		panic(err)
	}
	return decoded
}

func (id CertificateID) String() string {
	return string(id)
}

// Certificate is the interface for all kinds of certificate
type Certificate interface {
	proto.Message

	ID() CertificateID
	Type() CertificateType
	Certifier() sdk.AccAddress
	RequestContent() RequestContent
	CertificateContent() string
	FormattedCertificateContent() []KVPair
	Description() string
	TxHash() string

	String() string

	SetCertificateID(CertificateID)
	SetTxHash(string)
}

// RequestContentTypes is an array of all request content types.
var RequestContentTypes = [...]RequestContentType{
	RequestContentTypeNil,
	RequestContentTypeSourceCodeHash,
	RequestContentTypeAddress,
	RequestContentTypeBytecodeHash,
	RequestContentTypeGeneral,
}

// Bytes returns the byte array for a request content type.
func (c RequestContentType) Bytes() []byte {
	return []byte{byte(c)}
}

// String returns string of the request content type.
func (c RequestContentType) String() string {
	switch c {
	case RequestContentTypeSourceCodeHash:
		return "SourceCodeHash"
	case RequestContentTypeAddress:
		return "Address"
	case RequestContentTypeBytecodeHash:
		return "BytecodeHash"
	case RequestContentTypeGeneral:
		return "General"
	default:
		return "UnknownRequestContentType"
	}
}

// RequestContentTypeFromString returns the request content type by parsing a string.
func RequestContentTypeFromString(s string) RequestContentType {
	switch strings.ToUpper(s) {
	case "SOURCECODEHASH":
		return RequestContentTypeSourceCodeHash
	case "ADDRESS":
		return RequestContentTypeAddress
	case "BYTECODEHASH":
		return RequestContentTypeBytecodeHash
	case "GENERAL":
		return RequestContentTypeGeneral
	default:
		return RequestContentTypeNil
	}
}

// NewRequestContent returns a new request content.
func NewRequestContent(
	requestContentTypeString string,
	requestContent string,
) (RequestContent, error) {
	requestContentType := RequestContentTypeFromString(requestContentTypeString)
	if requestContentType == RequestContentTypeNil {
		return RequestContent{}, ErrInvalidRequestContentType
	}
	return RequestContent{RequestContentType: requestContentType, RequestContent: requestContent}, nil
}

// NewGeneralCertificate returns a new general certificate.
func NewGeneralCertificate(
	certTypeStr, reqContTypeStr, reqContStr, description string, certifier sdk.AccAddress,
) (*GeneralCertificate, error) {
	certType := CertificateTypeFromString(certTypeStr)
	if certType == CertificateTypeNil {
		return nil, ErrInvalidCertificateType
	}
	reqContent, err := NewRequestContent(reqContTypeStr, reqContStr)
	if err != nil {
		return nil, err
	}
	return &GeneralCertificate{
		CertType:        certType,
		ReqContent:      &reqContent,
		CertDescription: description,
		CertCertifier:   certifier.String(),
	}, nil
}

// ID returns ID of the certificate.
func (c *GeneralCertificate) ID() CertificateID {
	return c.CertId
}

// Type returns the certificate type.
func (c *GeneralCertificate) Type() CertificateType {
	return c.CertType
}

// Certifier returns certifier account address of the certificate.
func (c *GeneralCertificate) Certifier() sdk.AccAddress {
	certifierAddr, err := sdk.AccAddressFromBech32(c.CertCertifier)
	if err != nil {
		panic(err)
	}
	return certifierAddr
}

// RequestContent returns request content of the certificate.
func (c *GeneralCertificate) RequestContent() RequestContent {
	return *c.ReqContent
}

// CertificateContent returns certificate content of the certificate.
func (c *GeneralCertificate) CertificateContent() string {
	return "general certificate"
}

// FormattedCertificateContent returns formatted certificate content of the certificate.
func (c *GeneralCertificate) FormattedCertificateContent() []KVPair {
	return nil
}

// Description returns description of the certificate.
func (c *GeneralCertificate) Description() string {
	return c.CertDescription
}

// TxHash returns the hash of the tx when the certificate is issued.
func (c *GeneralCertificate) TxHash() string {
	return c.CertTxHash
}

// Bytes returns a byte array for the certificate.
func (c *GeneralCertificate) Bytes(cdc *codec.LegacyAmino) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(c)
}

// SetCertificateID provides a method to set an ID for the certificate.
func (c *GeneralCertificate) SetCertificateID(id CertificateID) {
	c.CertId = id
}

// SetTxHash provides a method to set txhash of the certificate.
func (c *GeneralCertificate) SetTxHash(txhash string) {
	c.CertTxHash = txhash
}

// NewCompilationCertificateContent returns a new compilation certificate content.
func NewCompilationCertificateContent(compiler, bytecodeHash string) CompilationCertificateContent {
	return CompilationCertificateContent{Compiler: compiler, BytecodeHash: bytecodeHash}
}

// NewKVPair returns a new key-value pair.
func NewKVPair(key string, value string) KVPair {
	return KVPair{Key: key, Value: value}
}

// NewCompilationCertificate returns a new compilation certificate
func NewCompilationCertificate(
	certificateType CertificateType,
	sourceCodeHash string,
	compiler string,
	bytecodeHash string,
	description string,
	certifier sdk.AccAddress,
) *CompilationCertificate {
	requestContent, _ := NewRequestContent("sourcecodehash", sourceCodeHash)
	certificateContent := NewCompilationCertificateContent(compiler, bytecodeHash)
	return &CompilationCertificate{
		CertType:        certificateType,
		ReqContent:      &requestContent,
		CertContent:     &certificateContent,
		CertDescription: description,
		CertCertifier:   certifier.String(),
	}
}

// ID returns ID of the certificate.
func (c *CompilationCertificate) ID() CertificateID {
	return c.CertId
}

// Type returns the certificate type.
func (c *CompilationCertificate) Type() CertificateType {
	return c.CertType
}

// Certifier returns certifier account address of the certificate.
func (c *CompilationCertificate) Certifier() sdk.AccAddress {
	certifierAddr, err := sdk.AccAddressFromBech32(c.CertCertifier)
	if err != nil {
		panic(err)
	}
	return certifierAddr
}

// RequestContent returns request content of the certificate.
func (c *CompilationCertificate) RequestContent() RequestContent {
	return *c.ReqContent
}

// CertificateContent returns certificate content of the certificate.
func (c *CompilationCertificate) CertificateContent() string {
	return c.CertContent.String()
}

// FormattedCertificateContent returns formatted certificate content of the certificate.
func (c *CompilationCertificate) FormattedCertificateContent() []KVPair {
	return []KVPair{
		NewKVPair("compiler", c.CertContent.Compiler),
		NewKVPair("bytecodeHash", c.CertContent.BytecodeHash),
	}
}

// Description returns description of the certificate.
func (c *CompilationCertificate) Description() string {
	return c.CertDescription
}

// TxHash returns the hash of the tx when the certificate is issued.
func (c *CompilationCertificate) TxHash() string {
	return c.CertTxHash
}

// Bytes returns a byte array for the certificate.
func (c *CompilationCertificate) Bytes(cdc *codec.LegacyAmino) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(c)
}

// SetCertificateID provides a method to set an ID for the certificate.
func (c *CompilationCertificate) SetCertificateID(id CertificateID) {
	c.CertId = id
}

// SetTxHash provides a method to set txhash of the certificate.
func (c *CompilationCertificate) SetTxHash(txhash string) {
	c.CertTxHash = txhash
}
