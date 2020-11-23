#!/bin/bash

set -e
set -x

# ------------------------------------------
#  Set up validator node running old binary
#
#  p2p port: 26656 (Cosmos default)
#  rpc port: 26657 (Cosmos default)
# ------------------------------------------

# node directory
# DIR=~/.synctest
DIR_D=$DIR/node0/certikd
DIR_CLI=$DIR/node0/certikcli

# binary
PROJ_ROOT=$(git rev-parse --show-toplevel)
CERTIKD=$PROJ_ROOT/tests/sync/certikd
export CERTIKCLI=$PROJ_ROOT/tests/sync/certikcli

# set up a validator node on port 20056
$CERTIKD unsafe-reset-all --home $DIR_D
rm -rf $DIR/node0
$CERTIKD init node0 --chain-id certikchain --home $DIR_D
$CERTIKCLI config chain-id certikchain --home $DIR_CLI
$CERTIKCLI config keyring-backend test --home $DIR_CLI
$CERTIKCLI keys add jack --home $DIR_CLI
export jack=$($CERTIKCLI keys show jack -a --home $DIR_CLI)
$CERTIKD add-genesis-account $jack 1000000000uctk --home $DIR_D
$CERTIKD gentx --name jack --amount 2000000uctk --home-client $DIR_CLI --keyring-backend test --home $DIR_D
$CERTIKD collect-gentxs --home $DIR_D
$CERTIKD start --home $DIR_D >$DIR/node0/log.txt 2>&1 &
