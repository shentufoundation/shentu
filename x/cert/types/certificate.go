package types

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
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

type Content interface {
	proto.Message

	GetType() RequestContentType
	GetContent() string
}

// Certificate is the interface for all kinds of certificate
type Certificate interface {
	proto.Message
	codecTypes.UnpackInterfacesMessage

	ID() uint64
	Type() CertificateType
	Certifier() sdk.AccAddress
	Content() Content
	FormattedContent() []KVPair
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

func AssembleContent(certTypeStr, reqContTypeStr, reqContStr string) Content {
	certType := CertificateTypeFromString(certTypeStr)
	reqContType := RequestContentTypeFromString(reqContTypeStr)
	switch certType {
	case CertificateTypeCompilation:
		return &Compilation{reqContType, reqContStr}
	case CertificateTypeAuditing:
		return &Auditing{reqContType, reqContStr}
	case CertificateTypeIdentity:
		return &Identity{reqContType, reqContStr}

	// TODO: more types to come

	default:
		return nil
	}
}

// NewGeneralCertificate returns a new general certificate.
func NewGeneralCertificate(
	certTypeStr, reqContTypeStr, reqContStr, description string, certifier sdk.AccAddress,
) (*GeneralCertificate, error) {
	certType := CertificateTypeFromString(certTypeStr)
	if certType == CertificateTypeNil {
		return nil, ErrInvalidCertificateType
	}
	content := AssembleContent(certTypeStr, reqContTypeStr, reqContStr)
	msg, ok := content.(proto.Message)
	if !ok {
		return &GeneralCertificate{}, fmt.Errorf("%T does not implement proto.Message", content)
	}
	any, err := codecTypes.NewAnyWithValue(msg)
	if err != nil {
		return &GeneralCertificate{}, err
	}
	return &GeneralCertificate{
		CertType:        certType,
		ReqContent:      any,
		CertDescription: description,
		CertCertifier:   certifier.String(),
	}, nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (c *GeneralCertificate) UnpackInterfaces(unpacker codecTypes.AnyUnpacker) error {
	var content Content
	return unpacker.UnpackAny(c.ReqContent, &content)
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
func (c *GeneralCertificate) Content() Content {
	content, ok := c.ReqContent.GetCachedValue().(Content)
	if !ok {
		return nil
	}
	return content
}

// FormattedCertificateContent returns formatted certificate content of the certificate.
func (c *GeneralCertificate) FormattedContent() []KVPair {
	return []KVPair{
		NewKVPair("content_type", c.Content().GetType().String()),
		NewKVPair("content", c.Content().GetContent()),
	}
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
	content := AssembleContent("COMPILATION", "SOURCECODEHASH", sourceCodeHash)
	msg, ok := content.(proto.Message)
	if !ok {
		return &CompilationCertificate{}
	}
	any, err := codecTypes.NewAnyWithValue(msg)
	if err != nil {
		return &CompilationCertificate{}
	}
	certificateContent := NewCompilationCertificateContent(compiler, bytecodeHash)
	return &CompilationCertificate{
		CertType:        certificateType,
		ReqContent:      any,
		CertContent:     &certificateContent,
		CertDescription: description,
		CertCertifier:   certifier.String(),
	}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (c *CompilationCertificate) UnpackInterfaces(unpacker codecTypes.AnyUnpacker) error {
	var content Content
	return unpacker.UnpackAny(c.ReqContent, &content)
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
func (c *CompilationCertificate) Content() Content { // TODO: problematic
	content, ok := c.ReqContent.GetCachedValue().(Content)
	if !ok {
		return nil
	}
	return content
}

// FormattedCertificateContent returns formatted certificate content of the certificate.
func (c *CompilationCertificate) FormattedContent() []KVPair {
	return []KVPair{
		NewKVPair("content_type", c.Content().GetType().String()),
		NewKVPair("content", c.Content().GetContent()),
	}
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
