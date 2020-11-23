#!/bin/bash

set -e
set -x

# -------------------------------------------------
#  Set up non-validator node running latest binary
#
#  p2p port: 27756
#  rpc port: 27757 (never used)
# -------------------------------------------------

# node directory
# DIR=~/.synctest
DIR_D1=$DIR/node1/certikd
DIR_CLI1=$DIR/node1/certikcli

GENESIS=$DIR/node0/certikd/config/genesis.json
FILE=$(ls $DIR/node0/certikd/config/gentx/)
PEER=${FILE:6:40}"@127.0.0.1:26656"

# binary
make install

# set up a non-validator node on port 20156 using current binary
certikd unsafe-reset-all --home $DIR_D1
rm -rf $DIR/node1
certikd init node1 --chain-id certikchain --home $DIR_D1
sed -i "" 's/26656/27756/g' $DIR_D1/config/config.toml                                        # p2p port
sed -i "" 's/26657/27757/g' $DIR_D1/config/config.toml                                        # rpc port
sed -i "" 's/persistent_peers = ""/persistent_peers = "'$PEER'"/g' $DIR_D1/config/config.toml # peer
cp $GENESIS $DIR_D1/config/genesis.json
certikcli config chain-id certikchain --home $DIR_CLI1
certikcli config keyring-backend test --home $DIR_CLI1

certikcli keys add mary --home $DIR_CLI1
export mary=$(certikcli keys show mary -a --home $DIR_CLI1)

certikd start --home $DIR_D1 >$DIR/node1/log.txt 2>&1 &
