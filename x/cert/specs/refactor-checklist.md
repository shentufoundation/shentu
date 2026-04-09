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

- [ ] Add new certificate index key prefixes
- [ ] Add helper functions for certificate index keys
- [ ] Add helper to write all certificate indexes
- [ ] Add helper to delete all certificate indexes
- [ ] Refactor `SetCertificate` to maintain indexes
- [ ] Refactor `DeleteCertificate` to maintain indexes
- [ ] Refactor `IssueCertificate` to use indexed storage path
- [ ] Refactor `RevokeCertificate` to use indexed delete path
- [ ] Refactor genesis init/export to preserve index consistency

## Certificate Query Refactor

- [ ] Replace full-scan filtered query with prefix-based pagination
- [ ] Choose index dynamically for certifier/type/composite filters
- [ ] Fix `total` semantics in `GetCertificatesFiltered`
- [ ] Fix response assembly in gRPC query to use page result length
- [ ] Refactor `IsCertified` to use indexed lookup
- [ ] Refactor `IsContentCertified` to use indexed lookup
- [ ] Refactor `IsBountyAdmin` to use indexed lookup

## Platform Removal

- [ ] Remove `MsgCertifyPlatform`
- [ ] Remove platform gRPC query
- [ ] Remove platform CLI tx command
- [ ] Remove platform CLI query command
- [ ] Remove platform proto messages
- [ ] Remove platform keeper logic
- [ ] Remove platform genesis field handling
- [ ] Remove platform store key prefix
- [ ] Remove platform docs and swagger entries
- [ ] Remove stale sync/e2e references to platform commands

## Library Removal

- [ ] Remove library keeper logic
- [ ] Remove library genesis field handling
- [ ] Remove library proto message if no longer needed
- [ ] Remove library store key prefix
- [ ] Remove library-specific errors if no longer used
- [ ] Remove library legacy migration code if no longer needed

## Migration

- [ ] Bump cert module consensus version
- [ ] Add new migrator step for the next version
- [ ] Rebuild certificate indexes from existing primary certificate store
- [ ] Delete obsolete platform store entries during migration
- [ ] Delete obsolete library store entries during migration
- [ ] Add migration tests for pre-upgrade state

## Tests

- [x] Add unit tests for authority-controlled certifier add/remove
- [x] Add unit tests for certifier query by address after alias removal
- [x] Add unit tests for governance detection of certifier update proposal messages
- [ ] Add unit tests for certificate index writes
- [ ] Add unit tests for certificate index deletes
- [ ] Add unit tests for filtered certificate queries
- [ ] Add unit tests for `IsCertified`
- [ ] Add unit tests for `IsContentCertified`
- [ ] Add unit tests for `IsBountyAdmin`
- [ ] Add migration tests
- [ ] Run `go test ./x/cert/...`
- [ ] Run impacted module tests for `x/gov`, `x/bounty`, and `x/oracle`

## Cleanup

- [x] Remove outdated validator/platform/library/alias language from docs
- [ ] Refresh generated docs/swagger if required
- [ ] Review for dead code and unused imports after removal
- [ ] Prepare follow-up patch for deeper certificate model simplification if still desired
