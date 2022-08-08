#!/bin/bash
set -o errexit -o nounset -o pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
BASE_ACCOUNT=$(certik keys show alice -a)

echo "-----------------------"
echo "## Genesis CosmWasm contract"
certik add-wasm-genesis-message store "$DIR/../keeper/testdata/hackatom.wasm" --instantiate-everybody true --run-as $BASE_ACCOUNT

echo "-----------------------"
echo "## Genesis CosmWasm instance"
INIT="{\"verifier\":\"$(certik keys show alice -a)\", \"beneficiary\":\"$(certik keys show bob -a)\"}"
certik add-wasm-genesis-message instantiate-contract 1 "$INIT" --run-as $BASE_ACCOUNT --label=foobar --amount=100uctk --admin "$BASE_ACCOUNT"

echo "-----------------------"
echo "## Genesis CosmWasm execute"
FIRST_CONTRACT_ADDR=certik14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s2anp5v
MSG='{"release":{}}'
certik add-wasm-genesis-message execute $FIRST_CONTRACT_ADDR "$MSG" --run-as $BASE_ACCOUNT --amount=1uctk

echo "-----------------------"
echo "## List Genesis CosmWasm codes"
certik add-wasm-genesis-message list-codes

echo "-----------------------"
echo "## List Genesis CosmWasm contracts"
certik add-wasm-genesis-message list-contracts
