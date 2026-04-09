# Cert

The `cert` module is responsible for two domain concepts:

- `certifier`: a governance-controlled roster of privileged security reviewers
- `certificate`: an on-chain record issued by a certifier

This module no longer treats validator certification, platform certification, or library registration as first-class scope. The current refactor direction is to remove `platform` and `library`, keep `certifier` and `certificate`, and improve certificate storage/query performance before any deeper certificate model redesign.

## Domain Boundary

### Certifiers

`Certifier` is a permission roster, not a parameter set.

It exists as module state because:

- other modules depend on `IsCertifier`
- certifiers are governance-controlled actors, not runtime config
- the roster may carry governance metadata and lifecycle checks

The target certifier model is intentionally small:

```go
type Certifier struct {
    Address     string `json:"address"`
    Proposer    string `json:"proposer"`
    Description string `json:"description"`
}
```

Notes:

- `alias` is being removed from the certifier model
- certifier identity is address-based only
- any human-friendly labeling should be handled off-chain or as optional non-indexed metadata later

### Certificates

`Certificate` remains the core business object of the module.

At the current refactor stage, the chain keeps the existing certificate data model and focuses first on:

- storage/index refactor
- query refactor
- runtime/API cleanup

This means the current `Any`-based content model is still preserved during the architecture refactor, even though it may be simplified in a later follow-up.

## State

### Certificates

Certificates remain stored by ID as the canonical primary record.

- Certificate: `0x5 | LittleEndian(CertificateId) -> certificate`
- NextCertificateID: `0x8 -> NextCertificateID`

The architecture refactor adds secondary indexes so common reads no longer require full-store scans.

Required index directions:

- `certifier -> certificate ids`
- `certificate type -> certificate ids`
- `certifier + certificate type -> certificate ids`
- `content hash -> certificate ids`
- `certificate type + content hash -> certificate ids`

These indexes support:

- certificate list queries
- `IsCertified`
- `IsContentCertified`
- `IsBountyAdmin`

### Certifiers

Certifiers remain stored as dedicated module state.

- Certifier: `0x0 | Address -> certifier`

The previous alias index is being removed as part of the certifier refactor.

## Governance And Mutation Model

### Certifier updates

Adding or removing certifiers should be governance-authorized and use a single execution path.

The refactor direction is:

- remove the ineffective `MsgProposeCertifier` flow
- stop exposing a user-facing message that appears successful but performs no state transition
- replace certifier mutation with an authority-controlled path executed by governance proposal messages

The preferred end state is an authority-gated message such as `MsgUpdateCertifier`, called only by the governance authority after proposal passage.

Target execution flow:

1. a user submits a governance proposal
2. the proposal contains `MsgUpdateCertifier`
3. governance passes the proposal
4. the governance module executes `MsgUpdateCertifier` using module authority
5. `cert` keeper applies the certifier roster change

That path should own:

- add certifier
- remove certifier
- validation of uniqueness and minimum-roster safety

### Certificate updates

Certificate mutation remains message-driven:

- `MsgIssueCertificate`
- `MsgRevokeCertificate`

Revocation behavior is part of the storage/query refactor scope because current revocation semantics interact with indexed lookups.

## Queries

The module should expose:

- query certifier by address
- query all certifiers
- query certificate by ID
- query certificates with indexed filtering

The previous certifier query-by-alias path is being removed together with `alias`.

The previous platform query is being removed together with `platform`.

## Messages

Current target public message surface after refactor:

- governance-authorized certifier update message
- `MsgIssueCertificate`
- `MsgRevokeCertificate`

Messages planned for removal:

- `MsgProposeCertifier`
- `MsgCertifyPlatform`

## Refactor Direction

The current execution order is:

1. shrink module scope to `certifier` and `certificate`
2. remove `platform` and `library`
3. refactor keeper structure by responsibility
4. add certificate secondary indexes
5. replace scan-based queries with indexed lookups
6. migrate existing state and update tests/docs
7. evaluate a later follow-up to simplify the certificate content model

## Parameters

There are currently no parameters specific to the `cert` module.

`certifier` must not be moved into params. It is governance-controlled state, not module configuration.
