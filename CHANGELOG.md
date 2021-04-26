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

## [2.0.0] - TBD

Version 2.0.0 brings many changes with SDK upgrading to Stargate version. For more information on the SDK upgrade, visit [CosmosSDK Release Notes](https://github.com/cosmos/cosmos-sdk/releases) 

### Client Breaking Changes
* (app) [\#241](https://github.com/certikfoundation/shentu/pull/241) Renamed default binary name to `certik`.

### API Breaking Changes
* (cvm) [\#231](https://github.com/certikfoundation/shentu/pull/231) Remove direct solidity file deployment.
* (shield) [\#244](https://github.com/certikfoundation/shentu/pull/244) Fix shield query & state export.
* (cert) [\#249](https://github.com/certikfoundation/shentu/pull/249) Certification module refactor.
* (shield) [\#269](https://github.com/certikfoundation/shentu/pull/269) Shield gRPC query refactor.

### State Machine Breaking Changes
* (app) [\#221](https://github.com/certikfoundation/shentu/pull/221) Upgraded SDK to 0.42.x.

### Features
* (ibc) [\#251](https://github.com/certikfoundation/shentu/pull/251) Add IBC support.

### Improvements
* [\#230](https://github.com/certikfoundation/shentu/pull/230) Optimized shield invariant & removed crisis module from endblocker.

### Tests

### Bug Fixes
* (app) [\#254](https://github.com/certikfoundation/shentu/pull/254) Disable module account receiving coins.
* (gov) [\#259](https://github.com/certikfoundation/shentu/pull/259) Gov module bug fixes. 
* (gov) [\#268](https://github.com/certikfoundation/shentu/pull/268) Fix proposal migration bug.

## [v1.3.1] - 02-05-2021

### Client Breaking Changes
### API Breaking Changes
### State Machine Breaking Changes
### Features
### Improvements
* [\#230](https://github.com/certikfoundation/shentu/pull/230) Optimized shield invariant & removed crisis module from endblocker

### Tests
### Bug Fixes

## [v1.3.0] - 01-15-2021

### Client Breaking Changes
### API Breaking Changes
### State Machine Breaking Changes
### Features
### Improvements
* [\#219](https://github.com/certikfoundation/shentu/pull/219) Remove internal sub-packages.
* [\#177](https://github.com/certikfoundation/shentu/pull/177) Updated Swagger docs.
### Tests
* [\#219](https://github.com/certikfoundation/shentu/pull/180) Implement SimApp for testing.
### Bug Fixes
* [\#216](https://github.com/certikfoundation/shentu/pull/216) Fixed Shield fee distribution.

## [v1.2.0] - 11-20-2020

### Client Breaking Changes
### API Breaking Changes
### State Machine Breaking Changes
### Features
### Improvements
* (rest) [\#131](https://github.com/certikfoundation/shentu/pull/171) Set default query limit to 100.

### Bug Fixes
* (x/shield) [\#173](https://github.com/certikfoundation/shentu/pull/173) Fixed indexing problem when paying out from unbonding delegations
* (x/shield) [\#170](https://github.com/certikfoundation/shentu/pull/170) Fixed conditional check for depositing collateral


## v1.1.0 - 11-11-2020

### Client Breaking Changes
### API Breaking Changes
### State Machine Breaking Changes
* (assets) [\#131](https://github.com/certikfoundation/shentu/pull/131) Added height checks for newly added tx routes.

### Features
* (x/shield) [\#132](https://github.com/certikfoundation/shentu/pull/132) Enabled Shield claim proposals for reimbursements.
* (x/shield) [\#131](https://github.com/certikfoundation/shentu/pull/131) Enabled Staking for Shield.

### Improvements
* (x/cvm) [\#129](https://github.com/certikfoundation/shentu/pull/129) Integrated CVM info to account query.
* (specs) [\#149](https://github.com/certikfoundation/shentu/pull/149) Add module specs.

### Bug Fixes
* (x/auth) [\#124](https://github.com/certikfoundation/shentu/pull/124) Fixed locked send event output.
* (x/gov) [\#145](https://github.com/certikfoundation/shentu/pull/145) Fixed param change proposal for simulations.

## v1.0.0 - 10-24-2020

### Client Breaking Changes

### API Breaking Changes
* (x/oracle) [\#6](https://github.com/certikfoundation/shentu/pull/6) Updated the `aggregate_task` event.
* (x/gov) [\#9](https://github.com/certikfoundation/shentu/pull/9) Paginated query and next page field in votes query. 

### State Machine Breaking Changes

### Features
* (x/cvm) [\#15](https://github.com/certikfoundation/shentu/pull/15) Enabled EWASM supoort.
* (x/auth) [\#7](https://github.com/certikfoundation/shentu/pull/7) Added new vesting account type ManualVestingAccount.
* (x/auth) [\#13](https://github.com/certikfoundation/shentu/pull/13) New locked-send tx type to ManualVestingAccounts.
* (toolsets/oracle-operator) [\#2](https://github.com/certikfoundation/shentu/pull/2) Added toolset oracle-operator.
* (toolsets/oracle-operator) [\#5](https://github.com/certikfoundation/shentu/pull/5) Added multi-client support.

### Improvements
* (x/oracle) [\#6](https://github.com/certikfoundation/shentu/pull/6) Updated events and added useful fields in task types.
* (toolsets/oracle-operator) [\#5](https://github.com/certikfoundation/shentu/pull/5) Operator refactor.
* (circleci) [\#4](https://github.com/certikfoundation/shentu/pull/4) Circleci project setup.
* (x/oracle) [\#8](https://github.com/certikfoundation/shentu/pull/8) Added simulation package to x/oracle.

### Bug Fixes
