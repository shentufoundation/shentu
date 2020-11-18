#!/bin/bash

set -e
set -x

# -----------------------------------------
# Set up validator node running old binary
# -----------------------------------------

# node directory
DIR=~/.synctest
NODE="node0"
DIR_D=$DIR/$NODE/certikd
DIR_CLI=$DIR/$NODE/certikcli

# binary paths
PROJ_ROOT=$(git rev-parse --show-toplevel)
CERTIKD=$PROJ_ROOT/tests/sync/certikd
CERTIKCLI=$PROJ_ROOT/tests/sync/certikcli

# set up a validator node on port 20056
$CERTIKD unsafe-reset-all --home $DIR_D
rm -rf $DIR/$NODE
$CERTIKD init node0 --chain-id certikchain --home $DIR_D
sed -i "" 's/26656/20056/g' $DIR_D/config/config.toml  # p2p port
sed -i "" 's/26657/20057/g' $DIR_D/config/config.toml  # rpc port
$CERTIKCLI config chain-id certikchain --home $DIR_CLI
$CERTIKCLI config keyring-backend test --home $DIR_CLI
$CERTIKCLI keys add jack --home $DIR_CLI
jack=$($CERTIKCLI keys show jack -a --home $DIR_CLI)
$CERTIKD add-genesis-account $jack 1000000000uctk --home $DIR_D
$CERTIKD gentx --name jack --amount 2000000uctk --home-client $DIR_CLI --keyring-backend test --home $DIR_D
$CERTIKD collect-gentxs --home $DIR_D
$CERTIKD start --home $DIR_D
