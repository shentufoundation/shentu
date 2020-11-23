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
DIR_D0=$DIR/node0/certikd
DIR_CLI0=$DIR/node0/certikcli

# binary
# PROJ_ROOT=$(git rev-parse --show-toplevel)
CERTIKD=$PROJ_ROOT/tests/sync/certikd
export CERTIKCLI=$PROJ_ROOT/tests/sync/certikcli

# set up a validator node on port 20056
$CERTIKD unsafe-reset-all --home $DIR_D0
rm -rf $DIR/node0
$CERTIKD init node0 --chain-id certikchain --home $DIR_D0
$CERTIKCLI config chain-id certikchain --home $DIR_CLI0
$CERTIKCLI config keyring-backend test --home $DIR_CLI0
$CERTIKCLI keys add jack --home $DIR_CLI0
export jack=$($CERTIKCLI keys show jack -a --home $DIR_CLI0)
$CERTIKD add-genesis-account $jack 1000000000uctk --home $DIR_D0
$CERTIKD gentx --name jack --amount 2000000uctk --home-client $DIR_CLI0 --keyring-backend test --home $DIR_D0
$CERTIKD collect-gentxs --home $DIR_D0
$CERTIKD start --home $DIR_D0 >$DIR/node0/log.txt 2>&1 &
