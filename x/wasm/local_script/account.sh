#!/bin/bash
set -o errexit -o nounset -o pipefail

BASE_ACCOUNT=$(certik keys show alice -a)
certik q account "$BASE_ACCOUNT" -o json | jq

echo "## Check balance"
NEW_ACCOUNT=$(certik keys show bob -a)
certik q bank balances "$NEW_ACCOUNT" -o json || true

echo "## Transfer tokens"
certik tx bank send alice "$NEW_ACCOUNT" 1uctk --gas 1000000 --gas-prices=0.025uctk -y --chain-id=testing --node=http://localhost:26657 -b block -o json | jq

echo "## Check balance again"
certik q bank balances "$NEW_ACCOUNT" -o json | jq
