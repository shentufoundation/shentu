# Cert

The `cert` module is responsible for two domain concepts:

- `certifier`: a governance-controlled roster of privileged security reviewers
- `certificate`: an on-chain record issued by a certifier

Validator certification, platform certification, and library registration have been removed. The module scope is now limited to `certifier` and `certificate`.

## Domain Boundary

### Certifiers

`Certifier` is a permission roster, not a parameter set.

It exists as module state because:

- other modules depend on `IsCertifier`
- certifiers are governance-controlled actors, not runtime config
- the roster may carry governance metadata and lifecycle checks

The certifier model:

```go
type Certifier struct {
    Address     string `json:"address"`
    Proposer    string `json:"proposer"`
    Description string `json:"description"`
}
```

Notes:

- `alias` has been removed from the certifier model
- certifier identity is address-based only
- any human-friendly labeling should be handled off-chain or as optional non-indexed metadata later

### Certificates

`Certificate` is the core business object of the module.

The current `Any`-based content model is preserved. A follow-up patch may simplify it later.

## State

### Certificates

Certificates are stored by ID as the canonical primary record.

- Certificate: `0x5 | LittleEndian(CertificateId) -> ProtoBuf(Certificate)`
- NextCertificateID: `0x8 -> uint64`

### Certificate Secondary Indexes

Five secondary indexes eliminate full-store scans for common reads. Each index key stores a sentinel value (`[]byte{0x1}`); the certificate ID is embedded in the key for extraction during iteration.

| Prefix | Key Layout | Purpose |
|--------|-----------|---------|
| `0x10` | `[0x10][20B certifier_addr][8B cert_id_BE]` | Lookup by certifier |
| `0x11` | `[0x11][1B cert_type][8B cert_id_BE]` | Lookup by certificate type |
| `0x12` | `[0x12][20B certifier_addr][1B cert_type][8B cert_id_BE]` | Lookup by certifier + type |
| `0x13` | `[0x13][32B content_sha256][8B cert_id_BE]` | Lookup by content hash |
| `0x14` | `[0x14][1B cert_type][32B content_sha256][8B cert_id_BE]` | Lookup by type + content hash |

These indexes support:

- `GetCertificatesFiltered` (paginated certificate list queries)
- `GetCertificatesByCertifier`
- `IsCertified` (type + content hash prefix scan)
- `IsContentCertified` (content hash prefix scan)
- `IsBountyAdmin` (type + content hash prefix scan)

Index consistency is maintained by `writeCertificateIndexes` and `deleteCertificateIndexes`, called from `SetCertificate`, `DeleteCertificate`, and the v2→v3 migration.

### Certifiers

Certifiers are stored as dedicated module state.

- Certifier: `0x0 | Address -> LengthPrefixed(Certifier)`

The alias index (prefix `0x7`) has been removed.

### Removed Prefixes

The following store prefixes have been removed and are deleted during the v2→v3 migration:

| Prefix | Former Purpose |
|--------|---------------|
| `0x1` | Validator certifications |
| `0x2` | Platform certifications |
| `0x6` | Library registrations |
| `0x7` | Certifier alias index |

## Governance And Mutation Model

### Certifier updates

Adding or removing certifiers is governance-authorized through a single execution path.

`MsgUpdateCertifier` is an authority-gated message executed only by the governance authority after proposal passage. The previous `MsgProposeCertifier` flow has been removed.

Execution flow:

1. a user submits a governance proposal
2. the proposal contains `MsgUpdateCertifier`
3. governance passes the proposal
4. the governance module executes `MsgUpdateCertifier` using module authority
5. `cert` keeper applies the certifier roster change

That path owns:

- add certifier
- remove certifier
- validation of uniqueness and minimum-roster safety

### Certificate updates

Certificate mutation is message-driven:

- `MsgIssueCertificate`
- `MsgRevokeCertificate`

Both paths maintain secondary index consistency.

## Queries

The module exposes:

- query certifier by address
- query all certifiers
- query certificate by ID
- query certificates with indexed filtering (by certifier, type, or both)

The certifier query-by-alias path has been removed together with `alias`.

The platform query has been removed together with `platform`.

### Query Behavior

- `GetCertificatesFiltered` chooses the narrowest index prefix based on filter params
- `total` returns the true count of all matching records, not page size
- pagination operates directly on the index iterator

## Messages

Public message surface:

- `MsgUpdateCertifier` (governance-authorized certifier add/remove)
- `MsgIssueCertificate`
- `MsgRevokeCertificate`

Removed messages:

- `MsgProposeCertifier` (removed; was a no-op)
- `MsgCertifyPlatform` (fully removed from proto service and generated code)

## Migration

### v2 → v3 (consensus version 3)

The module migration `Migrate2to3` performs:

1. Rebuilds all five secondary certificate indexes from existing primary certificate records
2. Deletes obsolete validator store entries (prefix `0x1`)
3. Deletes obsolete platform store entries (prefix `0x2`)
4. Deletes obsolete library store entries (prefix `0x6`)
5. Deletes obsolete certifier alias index entries (prefix `0x7`)

All deletions tolerate already-empty prefixes.

## Parameters

There are no parameters specific to the `cert` module.

`certifier` must not be moved into params. It is governance-controlled state, not module configuration.

## Deferred Items

- Refresh generated docs/swagger (requires tooling)
- Evaluate certificate content model simplification (follow-up patch)
