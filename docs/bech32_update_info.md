# Shentu Bech32 Address Prefix Update

## Overview

The bech32 address prefix for Shentu Chain is being updated to "shentu" for consistency with the chain's name and to foster the development of Shentu Chain ecosystem.


## New bech32 address prefix "shentu" - FAQ

1. What is the Bech32 Address Prefix?

- In the Cosmos ecosystem, each blockchain uses a unique address format with a chain-specific prefix. For instance, "cosmos" for the Cosmos Hub, "osmo" for Osmosis, and "juno" for Juno.

2. Why is the Shentu Chain switching to use the new address prefix "shentu"?

- This update aims to align the address prefix with the Shentu Chain's name and further support the Shentu Chain ecosystem growth.

3. What are the impacts on Centralized Exchanges (CEX)?

- Depositing addresses should be updated to the "shentu" prefix format, if applicable. 
- While the old address prefix format will still be compatible with the CLI, it will be converted to the new prefix format for persistent state and event generation purposes.
- If event handling involves parsing bech32 address strings, one should switch to use the "shentu" prefix for filtering instead of the old prefix.

4. What are the impacts on Validators?

- There are no expected impacts on validators.

5. What are the impacts on Existing Users?

- On some DApps, like Keplr and Cosmostation, users may need to manually update their settings to use the new prefix format.

## Tool for mapping old/new address

- API to get the new address

    - URL: http://44.214.37.77:8081/shentu/v1/bech32/converts
    - Method: POST
    - Params: List of addresses
    - Swagger UI: http://44.214.37.77:8081/swagger/index.html#/default/post_shentu_v1_bech32_converts
