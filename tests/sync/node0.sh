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

# node directory
# DIR=~/.synctest
DIR_D0=$DIR/node0/certikd
DIR_CLI0=$DIR/node0/certikcli

# binary
# PROJ_ROOT=$(git rev-parse --show-toplevel)
CERTIKD=$PROJ_ROOT/tests/sync/certikd
CERTIKD0=$CERTIKD" --home $DIR_D0"
CERTIKCLI=$PROJ_ROOT/tests/sync/certikcli
export CERTIKCLI0=$CERTIKCLI" --home $DIR_CLI0"

# set up a validator node
$CERTIKD0 unsafe-reset-all
rm -rf $DIR/node0
$CERTIKD0 init node0 --chain-id certikchain
sed -i "" 's/26657/20057/g' $DIR_D0/config/config.toml # rpc port
$CERTIKCLI0 config chain-id certikchain
$CERTIKCLI0 config keyring-backend test

$CERTIKCLI0 keys add jack
export jack=$($CERTIKCLI0 keys show jack -a)
$CERTIKCLI0 keys add bob
export bob=$($CERTIKCLI0 keys show bob -a)
$CERTIKD0 add-genesis-account $jack 1000000000uctk --vesting-amount=1000000uctk --manual --unlocker $bob
$CERTIKD0 add-genesis-account $bob 1000000000uctk
$CERTIKD0 add-genesis-certifier $bob

$CERTIKD0 gentx --name jack --amount 2000000uctk --home-client $DIR_CLI0 --keyring-backend test
$CERTIKD0 collect-gentxs

$CERTIKD0 start >$DIR/node0/log.txt 2>&1 &
