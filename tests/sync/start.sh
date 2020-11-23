#!/bin/bash

set -e
set -x

# -----------------------
# Start integration test
# -----------------------

# set up test directory
DIR=~/.synctest

# start nodes
PROJ_ROOT=$(git rev-parse --show-toplevel)
cd $PROJ_ROOT
. $PROJ_ROOT/tests/sync/node0.sh
. $PROJ_ROOT/tests/sync/node1.sh

# commence txs
echo "jack:" $jack
echo "mary:" $mary
