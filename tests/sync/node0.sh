#!/bin/bash

set -e
set -x

# ------------------------------------------
#  Set up validator node running old binary
#
#  p2p port: 26656 (Cosmos default)
#  rpc port: 20057 (never used)
#
#  jack: validator, manual-vesting account
#  bob: jack's unlocker, certifier
# ------------------------------------------

# directory

# DIR=~/.synctest
DIR_D0=$DIR/node0/shentud

# binary

# PROJ_ROOT=$(git rev-parse --show-toplevel)
SHENTUD=$PROJ_ROOT/tests/sync/shentud
SHENTUD0=$SHENTUD" --home $DIR_D0"
export $SHENTUD0

# set up a validator node

$SHENTUD0 tendermint unsafe-reset-all
rm -rf $DIR/node0
$SHENTUD0 init node0 --chain-id shentuchain
sed -i "" 's/26657/20057/g' $DIR_D0/config/config.toml # rpc port
$SHENTUD0 config chain-id shentuchain
$SHENTUD0 config keyring-backend test

$SHENTUD0 keys add jack
export jack=$($SHENTUD0 keys show jack -a)
$SHENTUD0 keys add bob
export bob=$($SHENTUD0 keys show bob -a)
$SHENTUD0 add-genesis-account $jack 1000000000uctk --vesting-amount=1000000uctk --manual --unlocker $bob
$SHENTUD0 add-genesis-account $bob 1000000000uctk
$SHENTUD0 add-genesis-certifier $bob
$SHENTUD0 add-genesis-shield-admin $bob

$SHENTUD0 gentx --name jack --amount 2000000uctk --home-client $DIR_D0 --keyring-backend test
$SHENTUD0 collect-gentxs

$SHENTUD0 start >$DIR/node0/log.txt 2>&1 &
