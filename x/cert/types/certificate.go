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

	GetContent() string
}

// Certificate is the interface for all kinds of certificate
type Certificate interface {
	proto.Message
	codecTypes.UnpackInterfacesMessage

	ID() uint64
	Certifier() sdk.AccAddress
	Content() Content
	CompilationContent() string
	FormattedCompilationContent() []KVPair
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

// AssembleContent constructs a struct instance that implements content interface.
func AssembleContent(certTypeStr, content string) Content {
	certType := CertificateTypeFromString(certTypeStr)
	switch certType {
	case CertificateTypeCompilation:
		return &Compilation{content}
	case CertificateTypeAuditing:
		return &Auditing{content}
	case CertificateTypeProof:
		return &Proof{content}
	case CertificateTypeOracleOperator:
		return &OracleOperator{content}
	case CertificateTypeShieldPoolCreator:
		return &ShieldPoolCreator{content}
	case CertificateTypeIdentity:
		return &Identity{content}
	case CertificateTypeGeneral:
		return &General{content}
	default:
		return nil
	}
}

// NewGeneralCertificate returns a new general certificate.
func NewGeneralCertificate(
	certTypeStr, contStr, description string, certifier sdk.AccAddress,
) (*GeneralCertificate, error) {
	content := AssembleContent(certTypeStr, contStr)
	msg, ok := content.(proto.Message)
	if !ok {
		return &GeneralCertificate{}, fmt.Errorf("%T does not implement proto.Message", content)
	}
	any, err := codecTypes.NewAnyWithValue(msg)
	if err != nil {
		return &GeneralCertificate{}, err
	}
	return &GeneralCertificate{
		CertContent:     any,
		CertDescription: description,
		CertCertifier:   certifier.String(),
	}, nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (c *GeneralCertificate) UnpackInterfaces(unpacker codecTypes.AnyUnpacker) error {
	var content Content
	return unpacker.UnpackAny(c.CertContent, &content)
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
	content, ok := c.CertContent.GetCachedValue().(Content)
	if !ok {
		return nil
	}
	return content
}

// CertificateContent returns certificate content of the certificate.
func (c *GeneralCertificate) CompilationContent() string {
	return "general certificate"
}

// FormattedCertificateContent returns formatted certificate content of the certificate.
func (c *GeneralCertificate) FormattedCompilationContent() []KVPair {
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
func NewCompilationCertificateContent(compiler, bytecodeHash string) CompilationContent {
	return CompilationContent{Compiler: compiler, BytecodeHash: bytecodeHash}
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
	content := AssembleContent("COMPILATION", sourceCodeHash)
	msg, ok := content.(proto.Message)
	if !ok {
		return &CompilationCertificate{}
	}
	any, err := codecTypes.NewAnyWithValue(msg)
	if err != nil {
		return &CompilationCertificate{}
	}
	compilationContent := NewCompilationCertificateContent(compiler, bytecodeHash)
	return &CompilationCertificate{
		CertContent:     any,
		CompContent:     &compilationContent,
		CertDescription: description,
		CertCertifier:   certifier.String(),
	}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (c *CompilationCertificate) UnpackInterfaces(unpacker codecTypes.AnyUnpacker) error {
	var content Content
	return unpacker.UnpackAny(c.CertContent, &content)
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
	content, ok := c.CertContent.GetCachedValue().(Content)
	if !ok {
		return nil
	}
	return content
}

// CertificateContent returns certificate content of the certificate.
func (c *CompilationCertificate) CompilationContent() string {
	return c.CompContent.String()
}

// FormattedCertificateContent returns formatted certificate content of the certificate.
func (c *CompilationCertificate) FormattedCompilationContent() []KVPair {
	return []KVPair{
		NewKVPair("compiler", c.CompContent.Compiler),
		NewKVPair("bytecodeHash", c.CompContent.BytecodeHash),
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
