# Cert Refactor Plan

Date: 2026-04-02
Status: Phases 1–6 complete. Phase 7 (optional follow-up) deferred.

## Goal

Refactor `x/cert` into a smaller and clearer module focused on:

- certifier roster management
- certificate issuance, revocation, and query

The refactor will remove `platform` and `library`, which are either weakly modeled or not productized, and redesign certificate storage/query to avoid full-store scans.
This phase does not redesign the certificate content proto model yet. It first fixes module boundaries, certifier governance flow, keeper structure, and certificate storage/query architecture.

## Current Assessment

### What stays

- `certifier` is the module's permission source.
- `certificate` is the module's core business record.
- external modules depend on `IsCertifier`, `IsCertified`, and `IsBountyAdmin`.
- `certifier` remains module state and must not be moved into params.

### What will be removed

- `platform`
  - has tx/query/proto/CLI/docs, but the model is too weak
  - stores only `validator_pubkey + description`
  - does not persist the certifier who certified it
  - has no meaningful downstream dependency
- `library`
  - exists only as keeper/genesis/store logic
  - has no tx/query/proto/CLI surface
  - has no downstream dependency
  - appears to be unfinished registry logic
- `certifier.alias`
  - is not required for identity or authorization
  - adds an extra uniqueness rule and index with weak product value
  - should not remain a primary query path
- `MsgProposeCertifier`
  - is currently a no-op message path
  - should be removed instead of repaired

### Main problems to solve

- certificate queries do full scans over the entire certificate store
- `GetCertificatesFiltered` paginates after materializing all matches
- `IsCertified`, `IsContentCertified`, and `IsBountyAdmin` also rely on scans
- query response assembly has incorrect assumptions around `total`
- migration scaffolding exists, but useful store migration is missing
- tests do not protect keeper/query behavior well enough
- specs and code are out of sync
- certifier mutation is split between a legacy proposal path and a broken no-op msg path
- certifier alias adds state/index complexity without being permission-critical

### Current chain state observations

Based on current chain export reviewed during refactor planning:

- `platforms` is empty
- `libraries` is empty
- certifier aliases are empty
- certificate state is non-empty

Implications:

- removing `platform` and `library` is primarily an API/schema cleanup, not a live-state migration burden
- removing `alias` is low-risk from a live-state perspective
- certificate migration and query refactor remain the main state-sensitive part of the work

## Target State

The module should expose only two durable domain concepts:

- `Certifier`
- `Certificate`

`Certifier` should become a small governance-controlled roster keyed only by address.

`Certificate` remains the canonical record for:

- issuer certifier
- certificate type
- content
- optional compilation metadata
- description
- certificate id

The refactor should keep the external certificate query API shape stable where possible, but replace the internal storage/query path with indexed lookups.

The refactor should also replace certifier mutation with one governance-authorized execution path. The preferred end state is an authority-gated certifier update message executed by governance proposal messages after proposal passage.

## Storage Design

### Primary store

- keep the current primary certificate record by ID
- continue to use certificate ID as the canonical lookup key

### Certifier store

Keep certifiers in dedicated module state:

- primary key by certifier address
- no alias index

Certifier state remains outside params because it is permission state rather than runtime configuration.

### Secondary indexes

Add secondary indexes so common reads no longer scan the full certificate store.

Required indexes:

- by certifier
- by certificate type
- by certifier + certificate type
- by content hash
- by certificate type + content hash

These indexes support:

- `query certificates --certifier`
- `query certificates --certificate-type`
- `query certificates --certifier --certificate-type`
- `IsCertified`
- `IsContentCertified`
- `IsBountyAdmin`

### Index write rules

Every certificate write path must keep indexes consistent:

- `SetCertificate`
- `IssueCertificate`
- `DeleteCertificate`
- `RevokeCertificate`
- genesis init/export
- migration rebuild

## Query Design

### Query path

`grpc_query.Certificates` should choose the narrowest prefix store first:

- no filter: primary certificate store
- certifier only: certifier index
- type only: type index
- certifier + type: composite index

Then paginate directly on that prefix iterator instead of loading all results first.

### Query behavior

The refactor should normalize query behavior:

- `total` must mean total matched rows, not page length
- result array length must equal the current page size
- pagination should be stable and deterministic

## Platform And Library Removal Plan

### Remove platform

Remove:

- `MsgCertifyPlatform`
- `Query Platform`
- platform CLI commands
- platform docs and swagger references
- platform proto messages
- platform keeper logic
- platform genesis field
- platform store prefix and legacy references

### Remove library

Remove:

- library keeper logic
- library genesis field
- library store prefix and errors if no longer needed
- legacy references if they are only for removed state

### Compatibility note

Removing `platform` is an API change.

Removing `library` is mainly a state/schema cleanup.

Both removals require a store migration and a module consensus version bump.

## Certifier Governance Refactor

### Direction

Refactor certifier mutation so it no longer depends on a user-facing no-op message and no longer uses alias-based identity.

Target behavior:

- certifier identity is address-based only
- add/remove certifier is governance-authorized
- one execution path owns all certifier state mutation

### API direction

Planned changes:

- remove `MsgProposeCertifier`
- remove alias from certifier data model and related query paths
- introduce an authority-controlled certifier update message, for example `MsgUpdateCertifier`

Target governance flow:

1. submit governance proposal carrying `MsgUpdateCertifier`
2. governance passes the proposal
3. governance authority executes `MsgUpdateCertifier`
4. `cert` keeper performs add or remove

That message should:

- accept governance authority as signer
- support add and remove operations
- validate duplicate address, minimum roster size, and any remaining metadata rules

### Compatibility note

This is a public API change for certifier tx/query surfaces:

- alias lookup is removed
- the direct user proposal message is removed
- certifier mutation is routed through governance authority

## Migration Plan

Add a real module migration for the next consensus version.

Migration responsibilities:

- rebuild certificate secondary indexes from existing certificate records
- remove obsolete certifier alias index entries
- delete obsolete `platform` store entries
- delete obsolete `library` store entries
- tolerate already-empty legacy prefixes

If old genesis import/export compatibility must be preserved for one release, keep a temporary compatibility layer at the genesis boundary only. Otherwise, remove the fields directly and update the module spec.

## Testing Plan

### Unit tests

Add or update tests for:

- authority-controlled certifier add/remove flow
- certifier queries by address after alias removal
- certificate issue populates all indexes
- certificate delete removes all indexes
- certificate revoke removes all indexes
- filtered certificate queries by certifier, type, and certifier+type
- `IsCertified`
- `IsContentCertified`
- `IsBountyAdmin`
- migration rebuilds indexes correctly
- migration deletes old platform/library state safely

### Integration tests

Update or add tests for:

- gRPC certificate query pagination
- e2e certificate queries after migration
- absence of platform commands and queries after removal

### Documentation tests

Update generated CLI/docs/swagger outputs after API removal.

## Work Phases

### Phase 0: Baseline — Complete

- confirmed final scope
- froze `platform` and `library` as removal targets
- documented target state

### Phase 1: Architecture Refactor — Complete

- removed certifier alias from docs, API, and state design
- removed the no-op `MsgProposeCertifier` path
- introduced governance-authorized certifier mutation path (`MsgUpdateCertifier`)
- split keeper responsibilities by domain

### Phase 2: Storage Refactor — Complete

- added index key prefixes (0x10–0x14) and helpers in `types/keys.go`
- implemented `writeCertificateIndexes` / `deleteCertificateIndexes` in `keeper/certificate.go`
- refactored `SetCertificate` and `DeleteCertificate` to maintain indexes

### Phase 3: Query Refactor — Complete

- replaced scan-based filtering with prefix-based pagination (`paginateIndex`)
- fixed `total` semantics (returns true match count, not page size)
- fixed gRPC response assembly (array sized by page length, not total)

### Phase 4: Remove Platform And Library — Complete

- removed platform keeper logic, CLI commands, and gRPC query (stubbed with Unimplemented)
- removed library keeper logic and genesis handling
- removed platform/library store key prefixes and related errors
- proto messages retained for backward compatibility (deferred removal requires protoc)

### Phase 5: Migration — Complete

- implemented `Migrate2to3`: rebuilds certificate indexes, deletes prefixes 0x1/0x2/0x6/0x7
- bumped module consensus version to 3
- registered migration step in `module.go`

### Phase 6: Test And Docs — Complete

- added keeper/query/migration tests in `certificate_index_test.go`
- updated specs to reflect new module boundary and storage layout
- removed stale references to validator/platform/library/alias

### Phase 7: Optional Follow-Up — Deferred

- evaluate simplification of the `Any`-based certificate content model

## Risks

- stale indexes if delete/revoke/genesis paths are missed
- certifier governance flow breakage if old and new mutation paths coexist during transition
- API breakage from removing platform query/tx
- API breakage from removing alias lookup and direct certifier proposal msg
- migration bugs on existing chain state
- e2e/docs drift if generated artifacts are not refreshed

## Acceptance Criteria

- certificate queries no longer scan the whole certificate store for common filters
- `IsCertified` and `IsBountyAdmin` no longer depend on full scans
- `platform` and `library` are fully removed from runtime logic and public interfaces
- certifier mutation is performed only through the governance-authorized path
- certifier state no longer contains alias-based identity or alias query paths
- migration succeeds from existing cert store state
- keeper/query tests cover the new storage/query paths
- specs reflect the new module boundary
