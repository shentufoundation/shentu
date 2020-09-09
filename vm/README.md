# CVM

CertiK virtual machine (CVM) handles the underlying gas calculation and corresponding error handling for contract deployments and calls.

## Behavior
These are some differences in behavior between CVM and EVM.

1. GASPRICE operation
    1. In EVM, `tx.gasprice` (triggers GASPRICE opcode) returns the gasprice for the transaction.
    2. Instead, CVM returns 0. It has to do with the opcode `GASPRICE_DEPRECATED`. Relevant handling can be found in `vm.go`.
    3. They consume the same amount of gas.
    4. This might change in the future.
2. DIFFICULTY operation
    1. In EVM, `DIFFICULTY` returns the current mining difficulty of the chain.
    2. Instead, CVM returns 0. This is due to the consensus model difference.
    3. Will not be changed in the future (might need update when Ethereum goes POS).
3. COINBASE operation
    1. In EVM, `COINBASE` pushes the miner address to the stack.
    2. Instaed, CVM pushes 0 to the stack, as there is no concept of mining.
    3. Will not change in the future.


## Gas Cost  
These are some differences in gas cost between CVM and EVM.
    
1. SSTORE opcode gas calculation
    1. CVM does not implement EIP-1283 or EIP-2200.
    2. Instead, it adds a simple NOOP gas (200) case to the original Petersburg gas calculation logic, for the case where the new value is equal to the original value.
2. SELFDESTRUCT opcode gas calculation
    1. as CVM doesn't have the identical access structure as EVM, EIP-158 is ignored.
    2. However, EIP 150 is taken account, so if `selfdestruct()` is called to an unexisting address,
    it will consume `CreateBySelfDestruct` amount of gas.

## Minor difference
These are some differences that do not affect behavior or gas cost.

1. BLOCKHEIGHT operation
    1. BLOCKHEIGHT is named `NUMBER` in EVM.
    2. They do the same thing, consume the same gas.
2. There are some missing opcodes in Burrow compared to Geth
    1. COINBASE
    2. NUMBER (Block height)
    3. CHAINID