#!/bin/bash

set -e
set -x

# ------------------------------------------------
# Set up non-validator node running latest binary
# ------------------------------------------------

# node directory
DIR=~/.synctest
NODE="node1"
DIR_D=$DIR/$NODE/certikd
DIR_CLI=$DIR/$NODE/certikcli

# binary paths
PROJ_ROOT=$(git rev-parse --show-toplevel)
cd $PROJ_ROOT
make install

# set up a non-validator node on port 20057 using current binary
certikd unsafe-reset-all --home $DIR_D
rm -rf $DIR/$NODE
certikd init node1 --chain-id certikchain --home $DIR_D
sed -i "" 's/26656/20156/g' $DIR_D/config/config.toml  # p2p port
sed -i "" 's/26657/20157/g' $DIR_D/config/config.toml  # rpc port
# add a persistent peer
cp $1 $DIR_D/config/genesis.json
certikcli config chain-id certikchain --home $DIR_CLI
certikcli config keyring-backend test --home $DIR_CLI
certikcli keys add mary --home $DIR_CLI
mary=$(certikcli keys show mary -a --home $DIR_CLI)
