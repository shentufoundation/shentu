# Cert

The `cert` module handles most of the certifier-related logic, include adding and removing certifiers, and certifying validators.

Certifiers are chain users that are responsible for overseeing the chain's security. CertiK Chain has two governing bodies that make decisions for the chain. First, validators are chain users that stake tokens as part of CertiK Chain's Delegated Proof-of-Stake consensus protocol. A validator's voting is proportional to their amount staked. Second, certifiers are responsible for security-related votes. Each certifier gets one equally-weighted vote. That being said, many of the actions that certifiers are privileged to take (see [Messages](#messages) below) do not require a voting process.

## State

### Certifiers

`Certifier` objects keep track of a certifier's information, including the certifier's alias and who proposed to add the certifier.

- Certifier: `0x0 | Address -> amino(certifier)`
- CertifierByAlias: `0x7 | Alias -> amino(certifier)`

```go
type Certifier struct {
    Address     string  `json:"certifier"`
    Alias       string  `json:"alias"`
    Proposer    string  `json:"proposer"`
    Description string  `json:"description"`
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

## Messages

`MsgProposeCertifier` must be proposed by a current certifier. It is first handled by the governance module for voting.

```go
type MsgProposeCertifier struct {
    Proposer    string  `json:"proposer" yaml:"proposer"`
    Alias       string  `json:"alias" yaml:"alias"`
    Certifier   string  `json:"certifier" yaml:"certifier"`
    Description string  `json:"description" yaml:"description"`
}
```

`MsgCertifyValidator` and `MsgDecertifyValidator` certifies and de-certifies a validator, respectively.

```go
type MsgCertifyValidator struct {
    Certifier   string      `json:"certifier" yaml:"certifier"`
    Pubkey      *types.Any  `json:"pubkey"`
}
```

```go
type MsgDecertifyValidator struct {
    Decertifier string      `json:"decertifier" yaml:"decertifier"`
    Pubkey      *types.Any  `json:"pubkey"`
}
```

## Parameters

There are currently no parameters specific to the `cert` module.
