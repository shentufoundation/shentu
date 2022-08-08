#!/bin/bash
set -o errexit -o nounset -o pipefail

PASSWORD=${PASSWORD:-1234567890}
echo "$PASSWORD"

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

echo "-----------------------"
echo "## Add new CosmWasm contract"
RESP=$(certik tx wasm store "$DIR/../keeper/testdata/hackatom.wasm" \
  --from alice --gas 1500000 --gas-prices=0.025uctk -y --chain-id=testing --node=http://localhost:26657 -b block -o json)

CODE_ID=$(echo "$RESP" | jq -r '.logs[0].events[1].attributes[-1].value')
echo "* Code id: $CODE_ID"
echo "* Download code"
TMPDIR=$(mktemp -t wasmcodeXXXXXX)
certik q wasm code "$CODE_ID" "$TMPDIR"
rm -f "$TMPDIR"
echo "-----------------------"
echo "## List code"
certik query wasm list-code --node=http://localhost:26657 --chain-id=testing -o json | jq

echo "-----------------------"
echo "## Create new contract instance"
INIT="{\"verifier\":\"$(certik keys show alice -a)\", \"beneficiary\":\"$(certik keys show bob -a)\"}"
certik tx wasm instantiate "$CODE_ID" "$INIT" --admin="$(certik keys show alice -a)" \
  --from alice --amount="10000000uctk" --label "local0.1.0" \
  --gas 1000000 --gas-prices=0.025uctk -y --chain-id=testing -b block -o json | jq

CONTRACT=$(certik query wasm list-contract-by-code "$CODE_ID" -o json | jq -r '.contracts[-1]')
echo "* Contract address: $CONTRACT"
echo "### Query all"
RESP=$(certik query wasm contract-state all "$CONTRACT" -o json)
echo "$RESP" | jq
echo "### Query smart"
certik query wasm contract-state smart "$CONTRACT" '{"verifier":{}}' -o json | jq
echo "### Query raw"
KEY=$(echo "$RESP" | jq -r ".models[0].key")
certik query wasm contract-state raw "$CONTRACT" "$KEY" -o json | jq

echo "-----------------------"
echo "## Execute contract $CONTRACT"
MSG='{"release":{}}'
certik tx wasm execute "$CONTRACT" "$MSG" \
  --from alice \
  --gas 1000000 --gas-prices=0.025uctk -y --chain-id=testing -b block -o json | jq

echo "-----------------------"
echo "## Set new admin"
echo "### Query old admin: $(certik q wasm contract "$CONTRACT" -o json | jq -r '.contract_info.admin')"
echo "### Update contract"
certik tx wasm set-contract-admin "$CONTRACT" "$(certik keys show bob -a)" \
  --from alice --gas 1000000 --gas-prices=0.025uctk -y --chain-id=testing -b block -o json | jq
echo "### Query new admin: $(certik q wasm contract "$CONTRACT" -o json | jq -r '.contract_info.admin')"

echo "-----------------------"
echo "## Migrate contract"
echo "### Upload new code"
RESP=$(certik tx wasm store "$DIR/../keeper/testdata/burner.wasm" \
  --from alice --gas 1000000 --gas-prices=0.025uctk -y --chain-id=testing --node=http://localhost:26657 -b block -o json)

BURNER_CODE_ID=$(echo "$RESP" | jq -r '.logs[0].events[1].attributes[-1].value')
echo "### Migrate to code id: $BURNER_CODE_ID"

DEST_ACCOUNT=$(certik keys show bob -a)

echo "### Query destination account before migration"
certik q bank balances "$DEST_ACCOUNT" -o json | jq

certik tx wasm migrate "$CONTRACT" "$BURNER_CODE_ID" "{\"payout\": \"$DEST_ACCOUNT\"}" --from bob \
  --gas 1000000 --gas-prices=0.025uctk --chain-id=testing -b block -y -o json | jq

echo "### Query destination account: $BURNER_CODE_ID"
certik q bank balances "$DEST_ACCOUNT" -o json | jq
echo "### Query contract meta data: $CONTRACT"
certik q wasm contract "$CONTRACT" -o json | jq

echo "### Query contract meta history: $CONTRACT"
certik q wasm contract-history "$CONTRACT" -o json | jq

echo "-----------------------"
echo "## Clear contract admin"
echo "### Query old admin: $(certik q wasm contract "$CONTRACT" -o json | jq -r '.contract_info.admin')"
echo "### Update contract"
certik tx wasm clear-contract-admin "$CONTRACT" \
  --from bob --gas 1000000 --gas-prices=0.025uctk -y --chain-id=testing -b block -o json | jq
echo "### Query new admin: $(certik q wasm contract "$CONTRACT" -o json | jq -r '.contract_info.admin')"
