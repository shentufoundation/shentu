# Module Specs

CertiK Chain is built with Cosmos SDK, which organizes the application's functionality into several components, or modules. Here we document the modules that are unique to CertiK Chain.

This documentation dives into the technical specifics of the CertiK Chain modules, and is geared towards a developer audience. See the [CertiK Chain Whitepaper](https://www.certik.foundation/whitepaper) for a conceptual overview of the CertiK Chain ecosystem.

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

These modules are unique to CertiK Chain.

- [Cert](cert.md) - Validator certification and certificate issuance
- [CVM](cvm.md) - CertiK Virtual Machine: smart contract deployment and execution
- [Oracle](oracle.md) - CertiK Security Oracle: on-chain scores for real-time security checks
- [Shield](shield.md) - CertiKShield: collateral pool to reimburse stolen assets
