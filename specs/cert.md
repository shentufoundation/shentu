# Cert

The `cert` module handles most of the certifier-related logic, include adding and removing certifiers, certifying validators, and issuing certificates.

Certifiers are chain users that are responsible for overseeing the chain's security. CertiK Chain has two governing bodies that make decisions for the chain. First, validators are chain users that stake tokens as part of CertiK Chain's Delegated Proof-of-Stake consensus protocol. A validator's voting is proportional to their amount staked. Second, certifiers are responsible for security-related votes. Each certifier gets one equally-weighted vote. That being said, many of the actions that certifiers are privileged to take (see [Messages](#messages) below) do not require a voting process.

## State

### Certificates

A `Certificate` can be stored on-chain to signify that a smart contract has been audited or formally verified. All certificates are implementations of the following interface:


```go
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

	Bytes(*codec.Codec) []byte
	String() string

	SetCertificateID(CertificateID)
	SetTxHash(string)
}
```

There are currently two types of certificates, `CompilationCertificate`s and `GeneralCertificate`s:

```go
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

type GeneralCertificate struct {
	CertID          CertificateID   `json:"certificate_id"`
	CertType        CertificateType `json:"certificate_type"`
	ReqContent      RequestContent  `json:"request_content"`
	CertDescription string          `json:"description"`
	CertCertifier   sdk.AccAddress  `json:"certifier"`
	CertTxHash      string          `json:"txhash"`
}
```

### Certifiers

`Certifier` objects keep track of a certifier's information, including the certifier's alias and who proposed to add the certifier.

```go
type Certifier struct {
	Address     sdk.AccAddress `json:"certifier"`
	Alias       string         `json:"alias"`
	Proposer    sdk.AccAddress `json:"proposer"`
	Description string         `json:"description"`
}
```

### Validators

A validator's public key and the certifier that added the validator are stored in a `Validator` object.

```go
type Validator struct {
	PubKey    crypto.PubKey
	Certifier sdk.AccAddress
}
```

## Stores

`Certifier`s are stored both by their address (in `certifierStore`) and their alias (in `certifierAliasStore`).

```go
var (
	certifierStoreKeyPrefix      = []byte{0x0}
	validatorStoreKeyPrefix      = []byte{0x1}
	platformStoreKeyPrefix       = []byte{0x2}
	certificateStoreKeyPrefix    = []byte{0x5}
	libraryStoreKeyPrefix        = []byte{0x6}
	certifierAliasStoreKeyPrefix = []byte{0x7}
)
```

## Messages

`MsgProposeCertifier` must be proposed by a current certifier. It is first handled by the governance module for voting.

```go
type MsgProposeCertifier struct {
	Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
	Alias       string         `json:"alias" yaml:"alias"`
	Certifier   sdk.AccAddress `json:"certifier" yaml:"certifier"`
	Description string         `json:"description" yaml:"description"`
}
```

`MsgCertifyValidator` lets a certifier add a new validator.

```go
type MsgCertifyValidator struct {
	Certifier sdk.AccAddress `json:"certifier" yaml:"certifier"`
	Validator crypto.PubKey  `json:"validator" yaml:"validator"`
}
```

`MsgDecertifyValidator` lets a certifier remove an existing validator.

```go
type MsgDecertifyValidator struct {
	Decertifier sdk.AccAddress `json:"decertifier" yaml:"decertifier"`
	Validator   crypto.PubKey  `json:"validator" yaml:"validator"`
}
```

The following two messages create a new general or compilation certificate, respectively.

```go
type MsgCertifyGeneral struct {
	CertificateType    string         `json:"certificate_type" yaml:"certificate_type"`
	RequestContentType string         `json:"request_content_type" yaml:"request_content_type"`
	RequestContent     string         `json:"request_content" yaml:"request_content"`
	Description        string         `json:"description" yaml:"description"`
	Certifier          sdk.AccAddress `json:"certifier" yaml:"certiifer"`
}
type MsgCertifyCompilation struct {
	SourceCodeHash string         `json:"sourcecodehash" yaml:"sourcecodehash"`
	Compiler       string         `json:"compiler" yaml:"compiler"`
	BytecodeHash   string         `json:"bytecodehash" yaml:"bytecodehash"`
	Description    string         `json:"description" yaml:"description"`
	Certifier      sdk.AccAddress `json:"certifier" yaml:"certifier"`
}
```

`MsgRevokeCertificate` removes a certificate from the store.

```go
type MsgRevokeCertificate struct {
	Revoker     sdk.AccAddress `json:"revoker" yaml:"revoker"`
	ID          CertificateID  `json:"id" yaml:"id"`
	Description string         `json:"description" yaml:"description"`
}
```

`MsgCertifyPlatform` certifies a validator's host platform.

```go
type MsgCertifyPlatform struct {
	Certifier sdk.AccAddress `json:"certifier" yaml:"certifier"`
	Validator crypto.PubKey  `json:"validator" yaml:"validator"`
	Platform  string         `json:"platform" yaml:"platform"`
}
```
## Parameters

There are currently no parameters specific to the `cert` module.
