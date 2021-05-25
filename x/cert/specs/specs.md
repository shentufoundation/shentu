# Cert

The `cert` module handles most of the certifier-related logic, include adding and removing certifiers, certifying validators, and issuing certificates.

Certifiers are chain users that are responsible for overseeing the chain's security. CertiK Chain has two governing bodies that make decisions for the chain. First, validators are chain users that stake tokens as part of CertiK Chain's Delegated Proof-of-Stake consensus protocol. A validator's voting is proportional to their amount staked. Second, certifiers are responsible for security-related votes. Each certifier gets one equally-weighted vote. That being said, many of the actions that certifiers are privileged to take (see [Messages](#messages) below) do not require a voting process.

## State

### Certificates

A `Certificate` can be stored on-chain to signify that a smart contract has been audited or formally verified.

- Certificate: `0x5 | LittleEndian(CertificateId) -> amino(certificate)`
- NextCertificateID: `0x8 -> NextCertificateID`

```go
type Certificate struct {
    CertificateId      uint64				`json:"certificate_id"`
    Content            *types.Any			`json:"content"`
    CompilationContent *CompilationContent	`json:"compilation_content"`
    Description        string				`json:"description"`
    Certifier          string				`json:"certifier"`
}
```

There are currently seven types of supported certificates:

- `Compilation`
- `Auditing`
- `Proof`
- `OracleOperator`
- `ShieldPoolCreator`
- `Identity`
- `General`

In addition to `CertificateId`, all certificates can be looked up by their `Certifier` and `Content`. Certificate types are determined by parsing the content string, in which contents are assembled into `Content` struct instances of their respective types.

```go
// Content is the interface for all kinds of certificate content.
type Content interface {
    proto.Message

    GetContent() string
}
```

In addition to the source code, additional information, such as compiler version and optimization settings, needs to be provided by the certifier. `CompilationContent` contains the compilation info for a smart contract in need of certification.

```go
type CompilationContent struct {
    Compiler		string
    BytecodeHash	string
}
```

### Certifiers

`Certifier` objects keep track of a certifier's information, including the certifier's alias and who proposed to add the certifier.

- Certifier: `0x0 | Address -> amino(certifier)`
- CertifierByAlias: `0x7 | Alias -> amino(certifier)`

```go
type Certifier struct {
    Address		string	`json:"certifier"`
    Alias		string	`json:"alias"`
    Proposer	string	`json:"proposer"`
    Description	string	`json:"description"`
}
```

### Validators

A `Validator` is a validator node, which can be certified by an existing certifier. When de-certified, the validator will then be unbonded, which prevents it from signing blocks or earning rewards.

- Validator: `0x1 | Pubkey -> amino(validator)`

```go
type Validator struct {
    Pubkey    *types.Any `json:"pubkey"`
    Certifier string     `json:"certifier"`
}
```

### Platforms

`Platform` objects are validators' host platforms that can be certified.

- Platform: `0x2 | ValidatorPubkey -> amino(platform)`

```go
type Platform struct {
    ValidatorPubkey	*types.Any	`json:"validator_pubkey"`
    Description		string		`json:"description"`
}
```

### Libraries

A `Library` object can be certified as well. It stores the library address and its publisher.

- Library: `0x6 | Address -> amino(library)`

```go
type Library struct {
    Address		string	`json:"address"`
    Publisher	string	`json:"publisher"`
}
```

## Messages

`MsgProposeCertifier` must be proposed by a current certifier. It is first handled by the governance module for voting.

```go
type MsgProposeCertifier struct {
    Proposer	string	`json:"proposer" yaml:"proposer"`
    Alias		string	`json:"alias" yaml:"alias"`
    Certifier	string	`json:"certifier" yaml:"certifier"`
    Description	string	`json:"description" yaml:"description"`
}
```

`MsgCertifyValidator` and `MsgDecertifyValidator` certifies and de-certifies a validator, respectively.

```go
type MsgCertifyValidator struct {
    Certifier	string		`json:"certifier" yaml:"certifier"`
    Pubkey		*types.Any	`json:"pubkey"`
}
```

```go
type MsgDecertifyValidator struct {
    Decertifier	string		`json:"decertifier" yaml:"decertifier"`
    Pubkey		*types.Any	`json:"pubkey"`
}
```

`MsgIssueCertificate` issues a certificate. It fails when the given certifier does not exist.

```go
type MsgIssueCertificate struct {
    Content			*types.Any	`json:"content"`
    Compiler		string		`json:"compiler" yaml:"compiler"`
    BytecodeHash	string		`json:"bytecode_hash" yaml:"bytecodehash"`
    Description		string		`json:"description" yaml:"description"`
    Certifier		string		`json:"certifier" yaml:"certifier"`
}
```

`MsgRevokeCertificate` removes a certificate from the store.

```go
type MsgRevokeCertificate struct {
    Revoker		string	`json:"revoker" yaml:"revoker"`
    Id			uint64	`json:"id" yaml:"id"`
    Description	string	`json:"description" yaml:"description"`
}
```

`MsgCertifyPlatform` certifies a validator's host platform.

```go
type MsgCertifyPlatform struct {
    Certifier		string     `json:"certifier" yaml:"certifier"`
    ValidatorPubkey *types.Any `json:"validator_pubkey"`
    Platform        string     `json:"platform" yaml:"platform"`
}
```

## Parameters

There are currently no parameters specific to the `cert` module.
