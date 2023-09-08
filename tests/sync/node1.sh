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
DIR_D1=$DIR/node1/shentud

GENESIS=$DIR/node0/shentud/config/genesis.json
FILE=$(ls $DIR/node0/shentud/config/gentx/)
PEER=${FILE:6:40}"@127.0.0.1:26656"

# binary

cd $PROJ_ROOT
make install
SHENTUD1=shentud" --home $DIR_D1"
export $SHENTUD1

# set up a non-validator node

$SHENTUD1 tendermint unsafe-reset-all
rm -rf $DIR/node1
$SHENTUD1 init node1 --chain-id shentuchain
sed -i "" 's/26656/27756/g' $DIR_D1/config/config.toml                                        # p2p port
sed -i "" 's/persistent_peers = ""/persistent_peers = "'$PEER'"/g' $DIR_D1/config/config.toml # persistent peers
cp $GENESIS $DIR_D1/config/genesis.json
$SHENTUD1 config chain-id shentuchain
$SHENTUD1 config keyring-backend test

$SHENTUD1 keys add mary
export mary=$($SHENTUCLI1 keys show mary -a)

$SHENTUD1 start >$DIR/node1/log.txt 2>&1 &
