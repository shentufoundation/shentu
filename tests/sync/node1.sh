#!/bin/bash

set -e
set -x

# -------------------------------------------------
#  Set up non-validator node running latest binary
#
#  p2p port: 27756
#  rpc port: 26657 (Cosmos default)
#
#  mary: a normal account
# -------------------------------------------------

# directory

# DIR=~/.synctest
DIR_D1=$DIR/node1/certikd
DIR_CLI1=$DIR/node1/certikcli

GENESIS=$DIR/node0/certikd/config/genesis.json
FILE=$(ls $DIR/node0/certikd/config/gentx/)
PEER=${FILE:6:40}"@127.0.0.1:26656"

# binary

cd $PROJ_ROOT
make install
CERTIKD1=certikd" --home $DIR_D1"
export CERTIKCLI1=certikcli" --home $DIR_CLI1"

# set up a non-validator node

$CERTIKD1 unsafe-reset-all
rm -rf $DIR/node1
$CERTIKD1 init node1 --chain-id shentuchain
sed -i "" 's/26656/27756/g' $DIR_D1/config/config.toml                                        # p2p port
sed -i "" 's/persistent_peers = ""/persistent_peers = "'$PEER'"/g' $DIR_D1/config/config.toml # persistent peers
cp $GENESIS $DIR_D1/config/genesis.json
$CERTIKCLI1 config chain-id shentuchain
$CERTIKCLI1 config keyring-backend test

$CERTIKCLI1 keys add mary
export mary=$($CERTIKCLI1 keys show mary -a)

$CERTIKD1 start >$DIR/node1/log.txt 2>&1 &
