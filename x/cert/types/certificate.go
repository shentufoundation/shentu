package types

import (
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

// CertificateTypeFromString returns a certificate type by parsing a string.
func CertificateTypeFromString(s string) CertificateType {
	switch strings.ToUpper(s) {
	case "COMPILATION", "CERT_TYPE_COMPILATION":
		return CertificateTypeCompilation
	case "AUDITING", "CERT_TYPE_AUDITING":
		return CertificateTypeAuditing
	case "PROOF", "CERT_TYPE_PROOF":
		return CertificateTypeProof
	case "ORACLEOPERATOR", "CERT_TYPE_ORACLE_OPERATOR":
		return CertificateTypeOracleOperator
	case "SHIELDPOOLCREATOR", "CERT_TYPE_SHIELD_POOL_CREATOR":
		return CertificateTypeShieldPoolCreator
	case "IDENTITY", "CERT_TYPE_IDENTITY":
		return CertificateTypeIdentity
	case "GENERAL", "CERT_TYPE_GENERAL":
		return CertificateTypeGeneral
	default:
		return CertificateTypeNil
	}
}

// Certificate is the interface for all kinds of certificate
type Certificate interface {
	proto.Message

	ID() uint64
	Type() CertificateType
	Certifier() sdk.AccAddress
	RequestContent() RequestContent
	CertificateContent() string
	FormattedCertificateContent() []KVPair
	Description() string

	String() string

	SetCertificateID(uint64)
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

// RequestContentTypeFromString returns the request content type by parsing a string.
func RequestContentTypeFromString(s string) RequestContentType {
	switch strings.ToUpper(s) {
	case "SOURCECODEHASH", "REQ_CONTENT_TYPE_SOURCE_CODE_HASH":
		return RequestContentTypeSourceCodeHash
	case "ADDRESS", "REQ_CONTENT_TYPE_ADDRESS":
		return RequestContentTypeAddress
	case "BYTECODEHASH", "REQ_CONTENT_TYPE_BYTECODE_HASH":
		return RequestContentTypeBytecodeHash
	case "GENERAL", "REQ_CONTENT_TYPE_GENERAL":
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
func (c *GeneralCertificate) ID() uint64 {
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

// Bytes returns a byte array for the certificate.
func (c *GeneralCertificate) Bytes(cdc *codec.LegacyAmino) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(c)
}

// SetCertificateID provides a method to set an ID for the certificate.
func (c *GeneralCertificate) SetCertificateID(id uint64) {
	c.CertId = id
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
func (c *CompilationCertificate) ID() uint64 {
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

// Bytes returns a byte array for the certificate.
func (c *CompilationCertificate) Bytes(cdc *codec.LegacyAmino) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(c)
}

// SetCertificateID provides a method to set an ID for the certificate.
func (c *CompilationCertificate) SetCertificateID(id uint64) {
	c.CertId = id
}
