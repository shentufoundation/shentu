package types

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CertificateType is the type for the type of a certificate.
type CertificateType byte

// Certificate types
const (
	CertificateTypeNil CertificateType = iota
	CertificateTypeCompilation
	CertificateTypeAuditing
	CertificateTypeProof
	CertificateTypeOracleOperator
	CertificateTypeShieldPoolCreator
	CertificateTypeIdentity
	CertificateTypeGeneral
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

// String returns the string for a certificate type.
func (c CertificateType) String() string {
	switch c {
	case CertificateTypeCompilation:
		return "Compilation"
	case CertificateTypeAuditing:
		return "Auditing"
	case CertificateTypeProof:
		return "Proof"
	case CertificateTypeOracleOperator:
		return "OracleOperator"
	case CertificateTypeShieldPoolCreator:
		return "ShieldPoolCreator"
	case CertificateTypeIdentity:
		return "Identity"
	case CertificateTypeGeneral:
		return "General"
	default:
		return "UnknownCertificateType"
	}
}

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

// Certificate is the interface for all kinds of certificate
type Certificate interface {
	ID() uint64
	Type() CertificateType
	Certifier() sdk.AccAddress
	RequestContent() RequestContent
	CertificateContent() string
	FormattedCertificateContent() []KVPair
	Description() string
	TxHash() string

	Bytes(*codec.Codec) []byte
	String() string

	SetCertificateID(uint64)
	SetTxHash(string)
}

// RequestContentType is the type for requestContent
type RequestContentType byte

// RequestContent types
const (
	RequestContentTypeNil RequestContentType = iota
	RequestContentTypeSourceCodeHash
	RequestContentTypeAddress
	RequestContentTypeBytecodeHash
	RequestContentTypeGeneral
)

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

// RequestContent defines type for the request content.
type RequestContent struct {
	RequestContentType RequestContentType `json:"request_content_type"`
	RequestContent     string             `json:"request_content"`
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

// GeneralCertificate defines the type for general certificate.
type GeneralCertificate struct {
	CertID          uint64          `json:"certificate_id"`
	CertType        CertificateType `json:"certificate_type"`
	ReqContent      RequestContent  `json:"request_content"`
	CertDescription string          `json:"description"`
	CertCertifier   sdk.AccAddress  `json:"certifier"`
	CertTxHash      string          `json:"txhash"`
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
		ReqContent:      reqContent,
		CertDescription: description,
		CertCertifier:   certifier,
	}, nil
}

// ID returns ID of the certificate.
func (c *GeneralCertificate) ID() uint64 {
	return c.CertID
}

// Type returns the certificate type.
func (c *GeneralCertificate) Type() CertificateType {
	return c.CertType
}

// Certifier returns certifier account address of the certificate.
func (c *GeneralCertificate) Certifier() sdk.AccAddress {
	return c.CertCertifier
}

// RequestContent returns request content of the certificate.
func (c *GeneralCertificate) RequestContent() RequestContent {
	return c.ReqContent
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
func (c *GeneralCertificate) Bytes(cdc *codec.Codec) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(c)
}

// String returns a human readable string representation of the certificate.
func (c *GeneralCertificate) String() string {
	return fmt.Sprintf("General certificate\n"+
		"Certificate ID: %s\n"+
		"Certificate type: %s\n"+
		"RequestContent:\n%s\n"+
		"Description: %s\n"+
		"Certifier: %s\n"+
		"TxHash: %s\n",
		strconv.FormatUint(c.CertID, 10), c.CertType.String(), c.ReqContent.RequestContent, c.CertDescription, c.CertCertifier.String(), c.CertTxHash)
}

// SetCertificateID provides a method to set an ID for the certificate.
func (c *GeneralCertificate) SetCertificateID(id uint64) {
	c.CertID = id
}

// SetTxHash provides a method to set txhash of the certificate.
func (c *GeneralCertificate) SetTxHash(txhash string) {
	c.CertTxHash = txhash
}

// CompilationCertificateContent defines type for the compilation certificate content.
type CompilationCertificateContent struct {
	Compiler     string `json:"compiler"`
	BytecodeHash string `json:"bytecode_hash"`
}

// NewCompilationCertificateContent returns a new compilation certificate content.
func NewCompilationCertificateContent(compiler, bytecodeHash string) CompilationCertificateContent {
	return CompilationCertificateContent{Compiler: compiler, BytecodeHash: bytecodeHash}
}

// String returns string of the compilation certificate content.
func (c CompilationCertificateContent) String() string {
	return fmt.Sprintf("Compilation certificate content:\n"+
		"Compiler: %s\n"+
		"Bytecode Hash: %s",
		c.Compiler, c.BytecodeHash)
}

// KVPair defines type for the key-value pair.
type KVPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewKVPair returns a new key-value pair.
func NewKVPair(key string, value string) KVPair {
	return KVPair{Key: key, Value: value}
}

// CompilationCertificate defines type for the compilation certificate.
type CompilationCertificate struct {
	IssueBlockHeight int64                         `json:"time_issued"`
	CertID           uint64                        `json:"certificate_id"`
	CertType         CertificateType               `json:"certificate_type"`
	ReqContent       RequestContent                `json:"request_content"`
	CertContent      CompilationCertificateContent `json:"certificate_content"`
	CertDescription  string                        `json:"description"`
	CertCertifier    sdk.AccAddress                `json:"certifier"`
	CertTxHash       string                        `json:"txhash"`
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
		ReqContent:      requestContent,
		CertContent:     certificateContent,
		CertDescription: description,
		CertCertifier:   certifier,
	}
}

// ID returns ID of the certificate.
func (c *CompilationCertificate) ID() uint64 {
	return c.CertID
}

// Type returns the certificate type.
func (c *CompilationCertificate) Type() CertificateType {
	return c.CertType
}

// Certifier returns certifier account address of the certificate.
func (c *CompilationCertificate) Certifier() sdk.AccAddress {
	return c.CertCertifier
}

// RequestContent returns request content of the certificate.
func (c *CompilationCertificate) RequestContent() RequestContent {
	return c.ReqContent
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
func (c *CompilationCertificate) Bytes(cdc *codec.Codec) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(c)
}

// String returns a human readable string representation of the certificate.
func (c *CompilationCertificate) String() string {
	return fmt.Sprintf("Compilation certificate\n"+
		"Certificate ID: %s\n"+
		"Certificate type: compilation\n"+
		"RequestContent:\n%s\n"+
		"CertificateContent:\n%s\n"+
		"Description: %s\n"+
		"Certifier: %s\n"+
		"TxHash: %s\n",
		strconv.FormatUint(c.CertID, 10), c.ReqContent.RequestContent, c.CertificateContent(),
		c.Description(), c.CertCertifier.String(), c.CertTxHash)
}

// SetCertificateID provides a method to set an ID for the certificate.
func (c *CompilationCertificate) SetCertificateID(id uint64) {
	c.CertID = id
}

// SetTxHash provides a method to set txhash of the certificate.
func (c *CompilationCertificate) SetTxHash(txhash string) {
	c.CertTxHash = txhash
}
