package migrate

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	certtypes "github.com/certikfoundation/shentu/x/cert/types"
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
	PubKey    crypto.PubKey
	Certifier sdk.AccAddress
}

// Platform is a genesis type for certified platform of a validator
type Platform struct {
	Validator   crypto.PubKey
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
	ID() uint64
	Type() CertificateType
	Certifier() sdk.AccAddress
	RequestContent() RequestContent
	CertificateContent() string
	FormattedCertificateContent() []KVPair
	Description() string

	Bytes(*codec.LegacyAmino) []byte
	String() string

	SetCertificateID(uint64)
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
	CertID          uint64          `json:"certificate_id"`
	CertType        CertificateType `json:"certificate_type"`
	ReqContent      RequestContent  `json:"request_content"`
	CertDescription string          `json:"description"`
	CertCertifier   sdk.AccAddress  `json:"certifier"`
}

// ID returns ID of the certificate.
func (c GeneralCertificate) ID() uint64 {
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

// Bytes returns a byte array for the certificate.
func (c GeneralCertificate) Bytes(cdc *codec.LegacyAmino) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(c)
}

// String returns a human readable string representation of the certificate.
func (c GeneralCertificate) String() string {
	return fmt.Sprintf("Compilation certificate\n"+
		"Certificate ID: %d\n"+
		"Certificate type: compilation\n"+
		"RequestContent:\n%s\n"+
		"CertificateContent:\n%s\n"+
		"Description: %s\n"+
		"Certifier: %s\n",
		c.CertID, c.ReqContent.RequestContent, c.CertificateContent(),
		c.Description(), c.CertCertifier.String())
}

// SetCertificateID provides a method to set an ID for the certificate.
func (c GeneralCertificate) SetCertificateID(id uint64) {
	c.CertID = id
}

// CompilationCertificateContent defines type for the compilation certificate content.
type CompilationCertificateContent struct {
	Compiler     string `json:"compiler"`
	BytecodeHash string `json:"bytecode_hash"`
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
}

// ID returns ID of the certificate.
func (c CompilationCertificate) ID() uint64 {
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

// Bytes returns a byte array for the certificate.
func (c CompilationCertificate) Bytes(cdc *codec.LegacyAmino) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(c)
}

// String returns a human readable string representation of the certificate.
func (c CompilationCertificate) String() string {
	return fmt.Sprintf("Compilation certificate\n"+
		"Certificate ID: %d\n"+
		"Certificate type: compilation\n"+
		"RequestContent:\n%s\n"+
		"CertificateContent:\n%s\n"+
		"Description: %s\n"+
		"Certifier: %s\n",
		c.CertID, c.ReqContent.RequestContent, c.CertificateContent(),
		c.Description(), c.CertCertifier.String())
}

// SetCertificateID provides a method to set an ID for the certificate.
func (c CompilationCertificate) SetCertificateID(id uint64) {
	c.CertID = id
}

// CertGenesisState - cert genesis state
type CertGenesisState struct {
	Certifiers        []Certifier   `json:"certifiers"`
	Validators        []Validator   `json:"validators"`
	Platforms         []Platform    `json:"platforms"`
	Certificates      []Certificate `json:"certificates"`
	Libraries         []Library     `json:"libraries"`
	NextCertificateId uint64        `json:"next_certificate_id"`
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

	newValidators := make([]certtypes.Validator, len(oldGenState.Validators))
	for i, v := range oldGenState.Validators {
		pkAny := codectypes.UnsafePackAny(v.PubKey)
		newValidators[i] = certtypes.Validator{
			Pubkey:    pkAny,
			Certifier: v.Certifier.String(),
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

	newCertificates := make([]*codectypes.Any, len(oldGenState.Certificates))
	for i, c := range oldGenState.Certificates {
		var newCert certtypes.Certificate
		reqContent := certtypes.RequestContent{
			RequestContentType: certtypes.RequestContentType(c.RequestContent().RequestContentType),
			RequestContent:     c.RequestContent().RequestContent,
		}
		newCert = &certtypes.GeneralCertificate{
			CertId:          c.ID(),
			CertType:        certtypes.CertificateType(c.Type()),
			ReqContent:      &reqContent,
			CertDescription: c.Description(),
			CertCertifier:   c.Certifier().String(),
		}
		msg, ok := newCert.(proto.Message)
		if !ok {
			panic(sdkerrors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", newCert))
		}
		certAny, err := codectypes.NewAnyWithValue(msg)
		if err != nil {
			panic(err)
		}
		newCertificates[i] = certAny
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
		Validators:        newValidators,
		Platforms:         newPlatforms,
		Certificates:      newCertificates,
		Libraries:         newLibraries,
		NextCertificateId: oldGenState.NextCertificateId,
	}
}
