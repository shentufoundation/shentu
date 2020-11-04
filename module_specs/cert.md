# Cert

CertiK Chain has two governing bodies that make decisions for the chain. First, validators are chain users that stake tokens as part of CertiK Chain's Delegated Proof-of-Stake consensus protocol. A validator's voting is proportional to their amount staked. Second, certifiers are responsible for security-related votes. Each certifier gets one equally-weighted vote.

The `cert` module handles most of the certifier-related logic, include adding and removing certifiers, certifying validators, and issuing certificates.


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

## Messages

## Events

## Parameters
