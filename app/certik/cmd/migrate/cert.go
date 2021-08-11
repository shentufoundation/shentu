package migrate

import (
	"encoding/hex"
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	certtypes "github.com/certikfoundation/shentu/v2/x/cert/types"
)

// CertificateType is the type for the type of a certificate.
type CertificateType byte

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

// Certifier is a type for certifier.
type Certifier struct {
	Address     sdk.AccAddress `json:"certifier"`
	Alias       string         `json:"alias"`
	Proposer    sdk.AccAddress `json:"proposer"`
	Description string         `json:"description"`
}

// Validator is a type for certified validator.
type Validator struct {
	PubKey    cryptotypes.PubKey
	Certifier sdk.AccAddress
}

// Platform is a genesis type for certified platform of a validator
type Platform struct {
	Validator   cryptotypes.PubKey
	Description string
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

// KVPair defines type for the key-value pair.
type KVPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RequestContent defines type for the request content.
type RequestContent struct {
	RequestContentType RequestContentType `json:"request_content_type"`
	RequestContent     string             `json:"request_content"`
}

// Certificate is the interface for all kinds of certificate
type Certificate interface {
	ID() CertificateID
	Type() CertificateType
	Certifier() sdk.AccAddress
	RequestContent() RequestContent
	CertificateContent() string
	FormattedCertificateContent() []KVPair
	Description() string
	TxHash() string

	Bytes(*codec.LegacyAmino) []byte
	String() string

	SetCertificateID(CertificateID)
	SetTxHash(string)
}

// Library is a type for certified libraries.
type Library struct {
	Address   sdk.AccAddress
	Publisher sdk.AccAddress
}

var _ Certificate = GeneralCertificate{}
var _ Certificate = CompilationCertificate{}

// GeneralCertificate defines the type for general certificate.
type GeneralCertificate struct {
	CertID          CertificateID   `json:"certificate_id"`
	CertType        CertificateType `json:"certificate_type"`
	ReqContent      RequestContent  `json:"request_content"`
	CertDescription string          `json:"description"`
	CertCertifier   sdk.AccAddress  `json:"certifier"`
	CertTxHash      string          `json:"txhash"`
}

// ID returns ID of the certificate.
func (c GeneralCertificate) ID() CertificateID {
	return c.CertID
}

// Type returns the certificate type.
func (c GeneralCertificate) Type() CertificateType {
	return c.CertType
}

// Certifier returns certifier account address of the certificate.
func (c GeneralCertificate) Certifier() sdk.AccAddress {
	return c.CertCertifier
}

// RequestContent returns request content of the certificate.
func (c GeneralCertificate) RequestContent() RequestContent {
	return c.ReqContent
}

// CertificateContent returns certificate content of the certificate.
func (c GeneralCertificate) CertificateContent() string {
	return "general certificate"
}

// FormattedCertificateContent returns formatted certificate content of the certificate.
func (c GeneralCertificate) FormattedCertificateContent() []KVPair {
	return nil
}

// Description returns description of the certificate.
func (c GeneralCertificate) Description() string {
	return c.CertDescription
}

// TxHash returns the hash of the tx when the certificate is issued.
func (c GeneralCertificate) TxHash() string {
	return c.CertTxHash
}

// Bytes returns a byte array for the certificate.
func (c GeneralCertificate) Bytes(cdc *codec.LegacyAmino) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(c)
}

// String returns a human readable string representation of the certificate.
func (c GeneralCertificate) String() string {
	return fmt.Sprintf("Compilation certificate\n"+
		"Certificate ID: %s\n"+
		"Certificate type: compilation\n"+
		"RequestContent:\n%s\n"+
		"CertificateContent:\n%s\n"+
		"Description: %s\n"+
		"Certifier: %s\n"+
		"TxHash: %s\n",
		c.CertID.String(), c.ReqContent.RequestContent, c.CertificateContent(),
		c.Description(), c.CertCertifier.String(), c.CertTxHash)
}

// SetCertificateID provides a method to set an ID for the certificate.
func (c GeneralCertificate) SetCertificateID(id CertificateID) {
	c.CertID = id
}

// SetTxHash provides a method to set txhash of the certificate.
func (c GeneralCertificate) SetTxHash(txhash string) {
	c.CertTxHash = txhash
}

// CompilationCertificateContent defines type for the compilation certificate content.
type CompilationCertificateContent struct {
	Compiler     string `json:"compiler"`
	BytecodeHash string `json:"bytecode_hash"`
}

// CompilationCertificate defines type for the compilation certificate.
type CompilationCertificate struct {
	IssueBlockHeight int64                         `json:"time_issued"`
	CertID           CertificateID                 `json:"certificate_id"`
	CertType         CertificateType               `json:"certificate_type"`
	ReqContent       RequestContent                `json:"request_content"`
	CertContent      CompilationCertificateContent `json:"certificate_content"`
	CertDescription  string                        `json:"description"`
	CertCertifier    sdk.AccAddress                `json:"certifier"`
	CertTxHash       string                        `json:"txhash"`
}

// ID returns ID of the certificate.
func (c CompilationCertificate) ID() CertificateID {
	return c.CertID
}

// Type returns the certificate type.
func (c CompilationCertificate) Type() CertificateType {
	return c.CertType
}

// Certifier returns certifier account address of the certificate.
func (c CompilationCertificate) Certifier() sdk.AccAddress {
	return c.CertCertifier
}

// RequestContent returns request content of the certificate.
func (c CompilationCertificate) RequestContent() RequestContent {
	return c.ReqContent
}

// CertificateContent returns certificate content of the certificate.
func (c CompilationCertificate) CertificateContent() string {
	return ""
}

// FormattedCertificateContent returns formatted certificate content of the certificate.
func (c CompilationCertificate) FormattedCertificateContent() []KVPair {
	return []KVPair{}
}

// Description returns description of the certificate.
func (c CompilationCertificate) Description() string {
	return c.CertDescription
}

// TxHash returns the hash of the tx when the certificate is issued.
func (c CompilationCertificate) TxHash() string {
	return c.CertTxHash
}

// Bytes returns a byte array for the certificate.
func (c CompilationCertificate) Bytes(cdc *codec.LegacyAmino) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(c)
}

// String returns a human readable string representation of the certificate.
func (c CompilationCertificate) String() string {
	return fmt.Sprintf("Compilation certificate\n"+
		"Certificate ID: %s\n"+
		"Certificate type: compilation\n"+
		"RequestContent:\n%s\n"+
		"CertificateContent:\n%s\n"+
		"Description: %s\n"+
		"Certifier: %s\n"+
		"TxHash: %s\n",
		c.CertID.String(), c.ReqContent.RequestContent, c.CertificateContent(),
		c.Description(), c.CertCertifier.String(), c.CertTxHash)
}

// SetCertificateID provides a method to set an ID for the certificate.
func (c CompilationCertificate) SetCertificateID(id CertificateID) {
	c.CertID = id
}

// SetTxHash provides a method to set txhash of the certificate.
func (c CompilationCertificate) SetTxHash(txhash string) {
	c.CertTxHash = txhash
}

// CertGenesisState - cert genesis state
type CertGenesisState struct {
	Certifiers   []Certifier   `json:"certifiers"`
	Validators   []Validator   `json:"validators"`
	Platforms    []Platform    `json:"platforms"`
	Certificates []Certificate `json:"certificates"`
	Libraries    []Library     `json:"libraries"`
}

func RegisterCertLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*Certificate)(nil), nil)
	cdc.RegisterConcrete(&GeneralCertificate{}, "cert/GeneralCertificate", nil)
	cdc.RegisterConcrete(&CompilationCertificate{}, "cert/CompilationCertificate", nil)
	cdc.RegisterConcrete(CertifierUpdateProposal{}, "cosmos-sdk/CertifierUpdateProposal", nil)
}

func migrateCert(oldGenState CertGenesisState) *certtypes.GenesisState {
	newCertifiers := make([]certtypes.Certifier, len(oldGenState.Certifiers))
	for i, c := range oldGenState.Certifiers {
		newCertifiers[i] = certtypes.Certifier{
			Alias:       c.Alias,
			Address:     c.Address.String(),
			Description: c.Description,
			Proposer:    c.Proposer.String(),
		}
	}

	newPlatforms := make([]certtypes.Platform, len(oldGenState.Platforms))
	for i, p := range oldGenState.Platforms {
		valPkAny := codectypes.UnsafePackAny(p.Validator)
		newPlatforms[i] = certtypes.Platform{
			ValidatorPubkey: valPkAny,
			Description:     p.Description,
		}
	}

	newCertificates := make([]certtypes.Certificate, len(oldGenState.Certificates))
	for i, c := range oldGenState.Certificates {
		content := AssembleContent(c.Type(), c.RequestContent().RequestContent)
		msg, ok := content.(proto.Message)
		if !ok {
			panic(fmt.Errorf("%T does not implement proto.Message", content))
		}
		any, err := codectypes.NewAnyWithValue(msg)
		if err != nil {
			panic(err)
		}
		newCertificates[i] = certtypes.Certificate{
			CertificateId:      uint64(i + 1),
			Content:            any,
			CompilationContent: &certtypes.CompilationContent{"", ""},
			Description:        c.Description(),
			Certifier:          c.Certifier().String(),
		}
	}

	newLibraries := make([]certtypes.Library, len(oldGenState.Libraries))
	for i, l := range oldGenState.Libraries {
		newLibraries[i] = certtypes.Library{
			Address:   l.Address.String(),
			Publisher: l.Publisher.String(),
		}
	}

	return &certtypes.GenesisState{
		Certifiers:        newCertifiers,
		Platforms:         newPlatforms,
		Certificates:      newCertificates,
		Libraries:         newLibraries,
		NextCertificateId: uint64(len(newCertificates) + 1),
	}
}

// AssembleContent constructs a struct instance that implements content interface.
func AssembleContent(certType CertificateType, content string) certtypes.Content {
	switch certType {
	case CertificateTypeCompilation:
		return &certtypes.Compilation{content}
	case CertificateTypeAuditing:
		return &certtypes.Auditing{content}
	case CertificateTypeProof:
		return &certtypes.Proof{content}
	case CertificateTypeOracleOperator:
		return &certtypes.OracleOperator{content}
	case CertificateTypeShieldPoolCreator:
		return &certtypes.ShieldPoolCreator{content}
	case CertificateTypeIdentity:
		return &certtypes.Identity{content}
	case CertificateTypeGeneral:
		return &certtypes.General{content}
	default:
		panic(certtypes.ErrInvalidCertificateType)
	}
}
