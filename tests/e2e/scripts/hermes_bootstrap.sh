#!/bin/bash

set -ex

# initialize Hermes relayer configuration
mkdir -p /root/.hermes/
touch /root/.hermes/config.toml

# setup Hermes relayer configuration
tee /root/.hermes/config.toml <<EOF
[global]
log_level = 'info'

[mode]

[mode.clients]
enabled = true
refresh = true
misbehaviour = true

[mode.connections]
enabled = false

[mode.channels]
enabled = false

[mode.packets]
enabled = true
clear_interval = 100
clear_on_start = true
tx_confirmation = true

[rest]
enabled = true
host = '0.0.0.0'
port = 3031

[telemetry]
enabled = true
host = '127.0.0.1'
port = 3001

[[chains]]
id = '$CERTIK_A_E2E_CHAIN_ID'
rpc_addr = 'http://$CERTIK_A_E2E_VAL_HOST:26657'
grpc_addr = 'http://$CERTIK_A_E2E_VAL_HOST:9090'
websocket_addr = 'ws://$CERTIK_A_E2E_VAL_HOST:26657/websocket'
rpc_timeout = '10s'
account_prefix = 'certik'
key_name = 'val01-certik-a'
store_prefix = 'ibc'
max_gas = 6000000
gas_price = { price = 0.001, denom = 'photon' }
gas_adjustment = 1.0
clock_drift = '1m' # to accomdate docker containers
trusting_period = '14days'
trust_threshold = { numerator = '1', denominator = '3' }

[[chains]]
id = '$CERTIK_B_E2E_CHAIN_ID'
rpc_addr = 'http://$CERTIK_B_E2E_VAL_HOST:26657'
grpc_addr = 'http://$CERTIK_B_E2E_VAL_HOST:9090'
websocket_addr = 'ws://$CERTIK_B_E2E_VAL_HOST:26657/websocket'
rpc_timeout = '10s'
account_prefix = 'certik'
key_name = 'val01-certik-b'
store_prefix = 'ibc'
max_gas = 6000000
gas_price = { price = 0.001, denom = 'photon' }
gas_adjustment = 1.0
clock_drift = '1m' # to accomdate docker containers
trusting_period = '14days'
trust_threshold = { numerator = '1', denominator = '3' }
EOF

# import keys
hermes keys restore ${CERTIK_B_E2E_CHAIN_ID} -n "val01-certik-b" -m "${CERTIK_B_E2E_VAL_MNEMONIC}"
hermes keys restore ${CERTIK_A_E2E_CHAIN_ID} -n "val01-certik-a" -m "${CERTIK_A_E2E_VAL_MNEMONIC}"

# start Hermes relayer
hermes start
