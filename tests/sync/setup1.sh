#!/bin/bash

set -e
set -x

# ------------------------------------------------
# Set up non-validator node running latest binary
# ------------------------------------------------

# node directory
DIR=~/.synctest
DIR_D=$DIR/node1/certikd
DIR_CLI=$DIR/node1/certikcli
GENESIS=$DIR/node0/certikd/config/genesis.json
PEERID=$(ls $DIR/node0/certikd/config/gentx/)
PEER=${PEERID:6:40}"@127.0.0.1:20056"

# binary paths
PROJ_ROOT=$(git rev-parse --show-toplevel)
cd $PROJ_ROOT
make install

# set up a non-validator node on port 20156 using current binary
certikd unsafe-reset-all --home $DIR_D
rm -rf $DIR/node1
certikd init node1 --chain-id certikchain --home $DIR_D
sed -i "" 's/26656/20156/g' $DIR_D/config/config.toml  # p2p port
sed -i "" 's/26657/20157/g' $DIR_D/config/config.toml  # rpc port
sed -i "" 's/addr_book_strict = true/addr_book_strict = false/g' $DIR_D/config/config.toml
sed -i "" 's/persistent_peers = ""/persistent_peers = "'$PEER'"/g' $DIR_D/config/config.toml
cp $GENESIS $DIR_D/config/genesis.json
certikcli config chain-id certikchain --home $DIR_CLI
certikcli config keyring-backend test --home $DIR_CLI
certikcli keys add mary --home $DIR_CLI
mary=$(certikcli keys show mary -a --home $DIR_CLI)
certikd start --home $DIR_D
