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

// Content is the interface for all kins of certificate content.
type Content interface {
	proto.Message

	GetType() ContentType
	GetContent() string
}

// Certificate is the interface for all kinds of certificate
type Certificate interface {
	proto.Message
	codecTypes.UnpackInterfacesMessage

	ID() uint64
	Certifier() sdk.AccAddress
	Content() Content
	FormattedContent() []KVPair
	CertificateContent() string
	FormattedCertificateContent() []KVPair
	Description() string

	String() string

	SetCertificateID(uint64)
}

// TranslateCertificateType determines certificate type based on content interface type switch.
func TranslateCertificateType(certificate Certificate) CertificateType {
	switch certificate.Content().(type) {
	case *Compilation:
		return CertificateTypeCompilation
	case *Auditing:
		return CertificateTypeAuditing
	case *Proof:
		return CertificateTypeProof
	case *OracleOperator:
		return CertificateTypeOracleOperator
	case *ShieldPoolCreator:
		return CertificateTypeShieldPoolCreator
	case *Identity:
		return CertificateTypeIdentity
	case *General:
		return CertificateTypeGeneral
	default:
		return CertificateTypeNil
	}
}

// ContentTypes is an array of all content types.
var ContentTypes = [...]ContentType{
	ContentTypeNil,
	ContentTypeSourceCodeHash,
	ContentTypeAddress,
	ContentTypeBytecodeHash,
	ContentTypeGeneral,
}

// Bytes returns the byte array for a content type.
func (c ContentType) Bytes() []byte {
	return []byte{byte(c)}
}

// ContentTypeFromString returns the content type by parsing a string.
func ContentTypeFromString(s string) ContentType {
	switch strings.ToUpper(s) {
	case "SOURCECODEHASH", "CONTENT_TYPE_SOURCE_CODE_HASH":
		return ContentTypeSourceCodeHash
	case "ADDRESS", "CONTENT_TYPE_ADDRESS":
		return ContentTypeAddress
	case "BYTECODEHASH", "CONTENT_TYPE_BYTECODE_HASH":
		return ContentTypeBytecodeHash
	case "GENERAL", "CONTENT_TYPE_GENERAL":
		return ContentTypeGeneral
	default:
		return ContentTypeNil
	}
}

// AssembleContent constructs a struct instance that implements content interface.
func AssembleContent(certTypeStr, contTypeStr, content string) Content {
	certType := CertificateTypeFromString(certTypeStr)
	contentType := ContentTypeFromString(contTypeStr)
	switch certType {
	case CertificateTypeCompilation:
		return &Compilation{contentType, content}
	case CertificateTypeAuditing:
		return &Auditing{contentType, content}
	case CertificateTypeProof:
		return &Proof{contentType, content}
	case CertificateTypeOracleOperator:
		return &OracleOperator{contentType, content}
	case CertificateTypeShieldPoolCreator:
		return &ShieldPoolCreator{contentType, content}
	case CertificateTypeIdentity:
		return &Identity{contentType, content}
	case CertificateTypeGeneral:
		return &General{contentType, content}
	default:
		return nil
	}
}

// NewGeneralCertificate returns a new general certificate.
func NewGeneralCertificate(
	certTypeStr, contTypeStr, contStr, description string, certifier sdk.AccAddress,
) (*GeneralCertificate, error) {
	content := AssembleContent(certTypeStr, contTypeStr, contStr)
	msg, ok := content.(proto.Message)
	if !ok {
		return &GeneralCertificate{}, fmt.Errorf("%T does not implement proto.Message", content)
	}
	any, err := codecTypes.NewAnyWithValue(msg)
	if err != nil {
		return &GeneralCertificate{}, err
	}
	return &GeneralCertificate{
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

// Certifier returns certifier account address of the certificate.
func (c *GeneralCertificate) Certifier() sdk.AccAddress {
	certifierAddr, err := sdk.AccAddressFromBech32(c.CertCertifier)
	if err != nil {
		panic(err)
	}
	return certifierAddr
}

// Content returns content of the certificate.
func (c *GeneralCertificate) Content() Content {
	content, ok := c.ReqContent.GetCachedValue().(Content)
	if !ok {
		return nil
	}
	return content
}

// FormattedContent returns formatted content of the certificate.
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

// Certifier returns certifier account address of the certificate.
func (c *CompilationCertificate) Certifier() sdk.AccAddress {
	certifierAddr, err := sdk.AccAddressFromBech32(c.CertCertifier)
	if err != nil {
		panic(err)
	}
	return certifierAddr
}

// Content returns content of the certificate.
func (c *CompilationCertificate) Content() Content {
	content, ok := c.ReqContent.GetCachedValue().(Content)
	if !ok {
		return nil
	}
	return content
}

// FormattedContent returns formatted content of the certificate.
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
