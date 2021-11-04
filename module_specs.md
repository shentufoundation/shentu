# Module Specs

Shentu Chain is built with Cosmos SDK, which organizes the application's functionality into several components, or modules. Here we document the modules that are unique to Shentu Chain.

This documentation dives into the technical specifics of the Shentu Chain modules, and is geared towards a developer audience. See the [Shentu Chain Whitepaper](https://www.certik.foundation/whitepaper) for a conceptual overview of the Shentu Chain ecosystem.

## Cosmos SDK Modules

The following modules largely inherit from Cosmos SDK. Please see [Cosmos's module documentation](https://docs.cosmos.network/master/modules/) for more info.

- [Auth](https://docs.cosmos.network/master/modules/auth/)
- [Bank](https://docs.cosmos.network/master/modules/bank/)
- [Crisis](https://docs.cosmos.network/master/modules/crisis/)
- [Distribution](https://docs.cosmos.network/master/modules/distribution/)
- [Governance](https://docs.cosmos.network/master/modules/gov/)
- [Mint](https://docs.cosmos.network/master/modules/mint/)
- [Slashing](https://docs.cosmos.network/master/modules/slashing/)
- [Staking](https://docs.cosmos.network/master/modules/staking/)
- [Upgrade](https://docs.cosmos.network/master/modules/upgrade/)

## Custom Modules

These modules are unique to Shentu Chain.

- [Cert](x/cert/specs/specs.md) - Validator certification and certificate issuance
- [CVM](x/cvm/specs/specs.md) - CertiK Virtual Machine: smart contract deployment and execution
- [Oracle](x/oracle/specs/specs.md) - CertiK Security Oracle: on-chain scores for real-time security checks
- [Shield](x/shield/specs/specs.md) - CertiKShield: collateral pool to reimburse stolen assets
