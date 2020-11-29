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

. $PROJ_ROOT/tests/sync/node0.sh
sleep 6

. $PROJ_ROOT/tests/sync/node1.sh
sleep 6

# commence tx sequence

. $PROJ_ROOT/tests/sync/txs.sh

# exiting

# killall certikd
echo "Compatibility check passed!"
exit 0
