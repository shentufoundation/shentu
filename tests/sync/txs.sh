#!/bin/bash

set -e
set -x

# ---------------------
# Transaction sequence
# ---------------------

# Add tokens to mary
$CERTIKCLI tx send $jack $mary 100000000uctk --from $jack -y --home $DIR_CLI0
sleep 5
$CERTIKCLI query account $jack --home $DIR_CLI0
$CERTIKCLI query account $bob --home $DIR_CLI0
$CERTIKCLI query account $mary --home $DIR_CLI0

# auth
$CERTIKCLI tx unlock $jack 500000uctk --from $bob -y --home $DIR_CLI0
sleep 5
$CERTIKCLI query account $jack --home $DIR_CLI0

# bank
certikcli tx locked-send $mary $jack 500000uctk --from $mary -y --home $DIR_CLI1
sleep 5
certikcli query account $jack --home $DIR_CLI1
certikcli query account $mary --home $DIR_CLI1

# cert

# cvm
# gov
# oracle
# shield
