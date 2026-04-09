# Cert Refactor Checklist

Date: 2026-04-02

## Execution Order

Work should be executed in this order:

1. certifier architecture refactor
2. certificate storage refactor
3. certificate query refactor
4. platform removal
5. library removal
6. migration
7. tests
8. cleanup

## Status Discipline

- [x] Keep plan/checklist in sync with the current design decisions
- [x] Update checklist state immediately after each completed code task
- [x] Update checklist state immediately after each completed test run
- [x] Do not mark implementation items complete before matching tests or verification are finished where applicable

## Discovery

- [x] Audited current `x/cert` structure and keeper responsibilities
- [x] Confirmed `certifier` is a real dependency surface for other modules
- [x] Confirmed certificate query path is currently scan-based
- [x] Confirmed `platform` has tx/query/proto/CLI/docs but weak business semantics
- [x] Confirmed `library` has store/genesis logic but no complete public API
- [x] Confirmed `platform` and `library` have no meaningful downstream module dependency
- [x] Confirmed current migration implementation is effectively a no-op
- [x] Confirmed keeper/query test coverage is insufficient for safe refactor
- [x] Confirmed current chain export has `platforms = []`
- [x] Confirmed current chain export has `libraries = []`
- [x] Confirmed current chain export has empty certifier aliases
- [x] Confirmed current chain export still has non-empty certificate state

## Scope Decisions

- [x] Decide to keep `certifier`
- [x] Decide to keep `certificate`
- [x] Decide to remove `platform`
- [x] Decide to remove `library`
- [x] Decide to remove `certifier.alias`
- [x] Decide that `certifier` remains module state, not params
- [x] Decide to replace direct certifier proposal flow with a governance-authorized mutation path
- [x] Decide to refactor certificate storage with secondary indexes
- [x] Decide to defer certificate content model redesign to a later patch
- [x] Decide to keep the refactor plan and execution checklist in `x/cert/specs`

## Plan Docs

- [x] Create refactor plan document
- [x] Create execution checklist document
- [x] Update `specs.md` to reflect the new target boundary

## Certifier Architecture Refactor

- [x] Remove `alias` from certifier state model
- [x] Remove alias query path from certifier gRPC API
- [x] Remove alias uniqueness checks from certifier mutation flow
- [x] Remove alias-related store key helpers
- [x] Remove alias-related docs and CLI examples
- [x] Remove no-op `MsgProposeCertifier`
- [x] Introduce governance-authorized certifier update message
- [x] Route certifier update through governance proposal messages
- [x] Route certifier add/remove through one keeper service path
- [x] Keep certifier state as dedicated module KV state
- [x] Ensure external modules continue to depend only on address-based `IsCertifier`

## Certificate Storage Refactor

- [x] Add new certificate index key prefixes
- [x] Add helper functions for certificate index keys
- [x] Add helper to write all certificate indexes
- [x] Add helper to delete all certificate indexes
- [x] Refactor `SetCertificate` to maintain indexes
- [x] Refactor `DeleteCertificate` to maintain indexes
- [x] Refactor `IssueCertificate` to use indexed storage path
- [x] Refactor `RevokeCertificate` to use indexed delete path
- [x] Refactor genesis init/export to preserve index consistency

## Certificate Query Refactor

- [x] Replace full-scan filtered query with prefix-based pagination
- [x] Choose index dynamically for certifier/type/composite filters
- [x] Fix `total` semantics in `GetCertificatesFiltered`
- [x] Fix response assembly in gRPC query to use page result length
- [x] Refactor `IsCertified` to use indexed lookup
- [x] Refactor `IsContentCertified` to use indexed lookup
- [x] Refactor `IsBountyAdmin` to use indexed lookup

## Platform Removal

- [x] Remove `MsgCertifyPlatform` CLI tx command
- [x] Remove platform gRPC query (stub returning Unimplemented)
- [x] Remove platform CLI query command
- [x] Remove platform keeper logic (`CertifyPlatform`, `GetPlatform`, `GetAllPlatforms`)
- [x] Remove platform genesis field handling (init/export no-op)
- [x] Remove platform store key prefix from `types/keys.go`
- [x] Remove stale sync/e2e references to platform commands (CLI removed)
- [ ] Remove platform proto messages (deferred; requires protoc regeneration)
- [ ] Remove platform docs and swagger entries (deferred; requires swagger regen)

## Library Removal

- [x] Remove library keeper logic (gutted `keeper/library.go`)
- [x] Remove library genesis field handling (init/export no-op)
- [x] Remove library store key prefix from `types/keys.go`
- [x] Remove library-specific errors (`ErrLibraryNotExists`, `ErrLibraryAlreadyExists`)
- [ ] Remove library proto message (deferred; requires protoc regeneration)
- [ ] Remove library legacy migration code if no longer needed (kept in legacy/v2 as no-op reference)

## Migration

- [x] Bump cert module consensus version (2 → 3)
- [x] Add new migrator step `Migrate2to3`
- [x] Rebuild certificate indexes from existing primary certificate store
- [x] Delete obsolete validator store entries during migration (prefix 0x1)
- [x] Delete obsolete platform store entries during migration (prefix 0x2)
- [x] Delete obsolete library store entries during migration (prefix 0x6)
- [x] Delete obsolete certifier alias index entries during migration (prefix 0x7)
- [x] Add migration tests (`Test_Migrate2to3IndexRebuild`, `Test_Migrate2to3DeletesObsoletePrefixes`)

## Tests

- [x] Add unit tests for authority-controlled certifier add/remove
- [x] Add unit tests for certifier query by address after alias removal
- [x] Add unit tests for governance detection of certifier update proposal messages
- [x] Add unit tests for certificate index writes (`Test_CertificateIndexWrites`)
- [x] Add unit tests for certificate index deletes (`Test_CertificateIndexDeletes`)
- [x] Add unit tests for filtered certificate queries (`Test_CertificateQueries`)
- [x] Add unit tests for `IsCertified` (`Test_IsCertified`)
- [x] Add unit tests for `IsContentCertified` (`Test_IsContentCertified`)
- [x] Add unit tests for `IsBountyAdmin` (`Test_IsBountyAdmin`)
- [x] Add migration tests (`Test_Migrate2to3IndexRebuild`)
- [x] Run `go test ./x/cert/...`
- [x] Run impacted module tests for `x/gov`, `x/bounty`, and `x/oracle`

## Cleanup

- [x] Remove outdated validator/platform/library/alias language from docs
- [x] Remove `GetCmdCertifyPlatform` from CLI tx commands
- [x] Update `specs.md` to reflect completed state (index key layout, migration, removed prefixes, query behavior)
- [x] Update `refactor-plan.md` with completion status for all phases
- [ ] Refresh generated docs/swagger if required (deferred; requires tooling)
- [ ] Review for dead code and unused imports after removal (minor; `MsgCertifyPlatform` methods retained for proto compatibility)
- [ ] Prepare follow-up patch for deeper certificate model simplification if still desired
