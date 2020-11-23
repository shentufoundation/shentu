#!/bin/bash

set -e
set -x

# -----------------------
# Start integration test
# -----------------------

# set up test directory
DIR=~/.synctest

# start the nodes
PROJ_ROOT=$(git rev-parse --show-toplevel)
cd $PROJ_ROOT
. $PROJ_ROOT/tests/sync/node0.sh
sleep 3
. $PROJ_ROOT/tests/sync/node1.sh
sleep 3

# commence txs
. $PROJ_ROOT/tests/sync/txs.sh
