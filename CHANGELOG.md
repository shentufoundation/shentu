<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.

Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## [v2.6.0] - 10-21-2022

### Client Breaking Changes
### API Breaking Changes
### State Machine Breaking Changes
* (app) [\#473](https://github.com/shentufoundation/shentu/pull/473) Store migration for mixcoins removal

### Features
### Improvements
* (app) [\#501](https://github.com/shentufoundation/shentu/pull/501) Apply Dragonberry patch
* (app) [\#506](https://github.com/shentufoundation/shentu/pull/506) Disable ica module init upgradeHandler


### Tests
### Bug Fixes
* (tests) [\#489](https://github.com/shentufoundation/shentu/pull/489) Fix e2e test bug
* (app) [\#471](https://github.com/shentufoundation/shentu/pull/471) Add shield crisis invariant fix for the upgrade
* (x/auth) [\#452](https://github.com/shentufoundation/shentu/pull/452) fix auth v1->v2 migration bug




## [v2.5.0] - 06-17-2022

### Client Breaking Changes
* (app) [\#467](https://github.com/shentufoundation/shentu/pull/467) Change binary name to `shentud` from `certik`
### API Breaking Changes
* (app) [\#467](https://github.com/shentufoundation/shentu/pull/467) Change package path to "github.com/shentufoundation/shentu" from "github.com/certikfoundation/shentu"
### State Machine Breaking Changes
* (app) [\#390](https://github.com/shentufoundation/shentu/pull/390) Bump Cosmos SDK to v0.45.4

### Features
### Improvements
* (app) [\#394](https://github.com/shentufoundation/shentu/pull/394) Bump Go version to 1.18
* (app) [\#399](https://github.com/shentufoundation/shentu/pull/399) Format & fix bugs in shield module


### Tests
### Bug Fixes

## [v2.4.0] - 06-17-2022

### State Machine Breaking Changes
* (x/shield) [\#380](https://github.com/shentufoundation/shentu/pull/380) Add more fine-grained and additional shield rewards to the collateral providers.

### Features
* (app) [\#387](https://github.com/shentufoundation/shentu/pull/387) Add genesis validator key rotation command.
* (x/shield) [\#437](https://github.com/shentufoundation/shentu/pull/437) Add automatic reward withdrawals on deposit/withdraw collateral.

## [v2.3.2] - 05-02-2022


### Improvements
* (app) [\#376](https://github.com/shentufoundation/shentu/pull/376) Add Rosetta Command & config command.
* (x/gov) [\#373](https://github.com/shentufoundation/shentu/pull/373) Add a fix to correctly export gov genesis state.

# Changelog

## [v2.3.1] - 01-24-2022

### State Machine Breaking Changes
* (x/auth) [\#371](https://github.com/shentufoundation/shentu/pull/371) Add an upgrade height to fix the outdated MVA accounts.


## [v2.3.0] - 01-04-2022

### Client Breaking Changes
* (x/cert) [\#326](https://github.com/shentufoundation/shentu/pull/326) Remove `Bech32` encoding for validator pubkeys.

### API Breaking Changes
### State Machine Breaking Changes
* (app) [\#326](https://github.com/shentufoundation/shentu/pull/326) Bump Cosmos SDK to v0.44.3.
* (app) [\#334](https://github.com/shentufoundation/shentu/pull/334) Implement in-store migration from v2.2.0 to v2.3.0.
* (x/gov) [\#334](https://github.com/shentufoundation/shentu/pull/334) `TxHash` field has been removed from `Vote` and `Deposit` types.

### Features
* (app) [\#326](https://github.com/shentufoundation/shentu/pull/326) Add `authz` and `feegrant` modules.

### Improvements
* (deps) Bump Tendermint to v0.34.14.

### Tests
### Bug Fixes
* (x/shield) [\#326](https://github.com/shentufoundation/shentu/pull/326) Add checks for expired entries in shield purchase.
* (x/gov) [\#331](https://github.com/shentufoundation/shentu/pull/331) Fix gov tally logic.
* (x/oracle) [\345](https://github.com/shentufoundation/shentu/pull/345) Fix oracle operator update txs permission checks.
* (x/cvm) [\#357](https://github.com/shentufoundation/shentu/pull/357) Fix cvm contract view.

## [v2.2.0] - 09-08-2021

### Bug Fixes
* (SDK) [\#323](https://github.com/shentufoundation/shentu/pull/323) Bump SDK to 0.42.9


## [v2.1.0] - 09-03-2021

Version 2.1.0 re-enables endblocker in the staking module, and bumps SDK to 0.42.9 for necessary query route.

### Client Breaking Changes
### API Breaking Changes
### State Machine Breaking Changes
* (x/staking) [\#323](https://github.com/shentufoundation/shentu/pull/241) Re-enable staking endblockers.

### Features
### Improvements
 * (app) [\#323](https://github.com/shentufoundation/shentu/pull/241) Bump Cosmos SDK to 0.42.9.

### Tests
### Bug Fixes

## [v2.0.0] - 08-09-2021

Version 2.0.0 brings many breaking changes with SDK upgrading to Stargate version. For more information on the SDK upgrade, visit [CosmosSDK Release Notes](https://github.com/cosmos/cosmos-sdk/releases) 

### Client Breaking Changes
* (app) [\#241](https://github.com/shentufoundation/shentu/pull/241) Renamed default binary name to `certik`.
* (x/oracle) [\#303](https://github.com/shentufoundation/shentu/pull/303) Oracle client commands refactor.

### API Breaking Changes
* (x/cvm) [\#231](https://github.com/shentufoundation/shentu/pull/231) Remove direct solidity file deployment.
* (x/shield) [\#244](https://github.com/shentufoundation/shentu/pull/244) Fix shield query & state export.
* (x/cert) [\#249](https://github.com/shentufoundation/shentu/pull/249) Certification module refactor.
* (x/shield) [\#269](https://github.com/shentufoundation/shentu/pull/269) Shield gRPC query refactor.

### State Machine Breaking Changes
* (app) [\#221](https://github.com/shentufoundation/shentu/pull/221) Upgraded SDK to 0.42.9.
* (x/shield) [\#286](https://github.com/shentufoundation/shentu/pull/286) Fix shield emitted events to include sender.
* (x/cert) [\#302](https://github.com/shentufoundation/shentu/pull/302) Removed validator certificate.
* (x/cvm) [\#301](https://github.com/shentufoundation/shentu/pull/301) Removed zero-address coins recycling.

### Features
* (ibc) [\#251](https://github.com/shentufoundation/shentu/pull/251) Add IBC support.

### Improvements
* [\#230](https://github.com/shentufoundation/shentu/pull/230) Optimized shield invariant & removed crisis module from endblocker.
* [\#286](https://github.com/shentufoundation/shentu/pull/286) Fix shield emitted events to include sender.
* [\#283](https://github.com/shentufoundation/shentu/pull/283) Improve vesting account generation.

### Tests
* [\#296](https://github.com/shentufoundation/shentu/pull/296) General test improvements over the modules.
* [\#280](https://github.com/shentufoundation/shentu/pull/280) Additional test cases for cert and oracle modules.

### Bug Fixes
* (app) [\#254](https://github.com/shentufoundation/shentu/pull/254) Disable module account receiving coins.
* (x/gov) [\#259](https://github.com/shentufoundation/shentu/pull/259) Gov module bug fixes. 
* (x/gov) [\#268](https://github.com/shentufoundation/shentu/pull/268) Fix proposal migration bug.

## [v1.3.1] - 02-05-2021

### Client Breaking Changes
### API Breaking Changes
### State Machine Breaking Changes
### Features
### Improvements
* [\#230](https://github.com/shentufoundation/shentu/pull/230) Optimized shield invariant & removed crisis module from endblocker

### Tests
### Bug Fixes

## [v1.3.0] - 01-15-2021

### Client Breaking Changes
### API Breaking Changes
### State Machine Breaking Changes
### Features
### Improvements
* [\#219](https://github.com/shentufoundation/shentu/pull/219) Remove internal sub-packages.
* [\#177](https://github.com/shentufoundation/shentu/pull/177) Updated Swagger docs.
### Tests
* [\#219](https://github.com/shentufoundation/shentu/pull/180) Implement SimApp for testing.
### Bug Fixes
* [\#216](https://github.com/shentufoundation/shentu/pull/216) Fixed Shield fee distribution.

## [v1.2.0] - 11-20-2020

### Client Breaking Changes
### API Breaking Changes
### State Machine Breaking Changes
### Features
### Improvements
* (rest) [\#131](https://github.com/shentufoundation/shentu/pull/171) Set default query limit to 100.

### Bug Fixes
* (x/shield) [\#173](https://github.com/shentufoundation/shentu/pull/173) Fixed indexing problem when paying out from unbonding delegations
* (x/shield) [\#170](https://github.com/shentufoundation/shentu/pull/170) Fixed conditional check for depositing collateral


## [v1.1.0] - 11-11-2020

### Client Breaking Changes
### API Breaking Changes
### State Machine Breaking Changes
* (assets) [\#131](https://github.com/shentufoundation/shentu/pull/131) Added height checks for newly added tx routes.

### Features
* (x/shield) [\#132](https://github.com/shentufoundation/shentu/pull/132) Enabled Shield claim proposals for reimbursements.
* (x/shield) [\#131](https://github.com/shentufoundation/shentu/pull/131) Enabled Staking for Shield.

### Improvements
* (x/cvm) [\#129](https://github.com/shentufoundation/shentu/pull/129) Integrated CVM info to account query.
* (specs) [\#149](https://github.com/shentufoundation/shentu/pull/149) Add module specs.

### Bug Fixes
* (x/auth) [\#124](https://github.com/shentufoundation/shentu/pull/124) Fixed locked send event output.
* (x/gov) [\#145](https://github.com/shentufoundation/shentu/pull/145) Fixed param change proposal for simulations.

## [v1.0.0] - 10-24-2020

### Client Breaking Changes
### API Breaking Changes
* (x/oracle) [\#6](https://github.com/shentufoundation/shentu/pull/6) Updated the `aggregate_task` event.
* (x/gov) [\#9](https://github.com/shentufoundation/shentu/pull/9) Paginated query and next page field in votes query. 

### State Machine Breaking Changes
### Features
* (x/cvm) [\#15](https://github.com/shentufoundation/shentu/pull/15) Enabled EWASM supoort.
* (x/auth) [\#7](https://github.com/shentufoundation/shentu/pull/7) Added new vesting account type ManualVestingAccount.
* (x/auth) [\#13](https://github.com/shentufoundation/shentu/pull/13) New locked-send tx type to ManualVestingAccounts.
* (toolsets/oracle-operator) [\#2](https://github.com/shentufoundation/shentu/pull/2) Added toolset oracle-operator.
* (toolsets/oracle-operator) [\#5](https://github.com/shentufoundation/shentu/pull/5) Added multi-client support.

### Improvements
* (x/oracle) [\#6](https://github.com/shentufoundation/shentu/pull/6) Updated events and added useful fields in task types.
* (toolsets/oracle-operator) [\#5](https://github.com/shentufoundation/shentu/pull/5) Operator refactor.
* (circleci) [\#4](https://github.com/shentufoundation/shentu/pull/4) Circleci project setup.
* (x/oracle) [\#8](https://github.com/shentufoundation/shentu/pull/8) Added simulation package to x/oracle.

### Bug Fixes
