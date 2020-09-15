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

## [Unreleased] - TBD

### Client Breaking Changes

### API Breaking Changes
* (x/oracle) [\#6](https://github.com/certikfoundation/shentu/pull/6) Updated the `aggregate_task` event.

### State Machine Breaking Changes

### Features
* (toolsets/oracle-operator) [\#2](https://github.com/certikfoundation/shentu/pull/2) Added toolset oracle-operator.
* (toolsets/oracle-operator) [\#5](https://github.com/certikfoundation/shentu/pull/5) Added multi-client support.

### Improvements
* (x/oracle) [\#6](https://github.com/certikfoundation/shentu/pull/6) Updated events and added useful fields in task types.
* (toolsets/oracle-operator) [\#5](https://github.com/certikfoundation/shentu/pull/5) Operator refactor.
* (circleci) [\#4](https://github.com/certikfoundation/shentu/pull/4) Circleci project setup.

### Bug Fixes
