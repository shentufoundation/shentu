package types

import (
	"fmt"
	"strings"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/cosmos/gogoproto/proto"
)

// Certificate types implement UnpackInterfaceMessages to unpack Content field.
var _ codecTypes.UnpackInterfacesMessage = Certificate{}

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
	case "BOUNTYADMIN", "CERT_TYPE_BOUNTYADMIN":
		return CertificateTypeBountyAdmin
	default:
		return CertificateTypeNil
	}
}

// TranslateCertificateType determines certificate type based on content interface type switch.
func TranslateCertificateType(certificate Certificate) CertificateType {
	switch certificate.GetContent().(type) {
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
	case *BountyAdmin:
		return CertificateTypeBountyAdmin
	default:
		return CertificateTypeNil
	}
}

// Content is the interface for all kinds of certificate content.
type Content interface {
	proto.Message

	GetContent() string
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
	case CertificateTypeBountyAdmin:
		return &BountyAdmin{content}
	default:
		return nil
	}
}

// NewCompilationCertificateContent returns a new compilation content.
func NewCompilationCertificateContent(compiler, bytecodeHash string) CompilationContent {
	return CompilationContent{Compiler: compiler, BytecodeHash: bytecodeHash}
}

// NewCertificate returns a new certificate.
func NewCertificate(
	certTypeStr, contStr, compiler, bytecodeHash, description string, certifier sdk.AccAddress,
) (Certificate, error) {
	content := AssembleContent(certTypeStr, contStr)
	msg, ok := content.(proto.Message)
	if !ok {
		return Certificate{}, fmt.Errorf("%T does not implement proto.Message", content)
	}
	any, err := codecTypes.NewAnyWithValue(msg)
	if err != nil {
		return Certificate{}, err
	}
	compilationContent := NewCompilationCertificateContent(compiler, bytecodeHash)
	return Certificate{
		Content:            any,
		CompilationContent: &compilationContent,
		Description:        description,
		Certifier:          certifier.String(),
	}, nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces.
func (c Certificate) UnpackInterfaces(unpacker codecTypes.AnyUnpacker) error {
	var content Content
	return unpacker.UnpackAny(c.Content, &content)
}

// GetContent returns content of the certificate.
func (c Certificate) GetContent() Content {
	content, ok := c.Content.GetCachedValue().(Content)
	if !ok {
		return nil
	}
	return content
}

// GetContentString returns string content of the certificate.
func (c Certificate) GetContentString() string {
	return c.GetContent().GetContent()
}

// FormattedCompilationContent returns formatted certificate content of the certificate.
func (c Certificate) FormattedCompilationContent() []KVPair {
	return []KVPair{
		NewKVPair("compiler", c.CompilationContent.Compiler),
		NewKVPair("bytecodeHash", c.CompilationContent.BytecodeHash),
	}
}

// GetCertifier returns certificer of the certificate.
func (c Certificate) GetCertifier() sdk.AccAddress {
	certifierAddr, err := sdk.AccAddressFromBech32(c.Certifier)
	if err != nil {
		panic(err)
	}
	return certifierAddr
}

func (c Certificate) ToString() string {
	certStr := fmt.Sprintf("certificate_id:%d content:%s description:%s certifier:%s compilation_content:",
		c.CertificateId, c.GetContentString(), c.Description, c.Certifier)
	if c.CompilationContent != nil {
		return certStr + "<" + c.CompilationContent.String() + ">"
	}
	return certStr + "<>"
}

// NewKVPair returns a new key-value pair.
func NewKVPair(key string, value string) KVPair {
	return KVPair{Key: key, Value: value}
}
