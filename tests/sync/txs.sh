#!/bin/bash

set -e
set -x

# ---------------------
# Transaction sequence
# ---------------------

# Add tokens to mary
$CERTIKCLI tx send $jack $mary 100000000uctk --from $jack -y --home $DIR_CLI0
sleep 3
$CERTIKCLI query account $jack --home $DIR_CLI0
$CERTIKCLI query account $mary --home $DIR_CLI0
certikcli query account $jack --home $DIR_CLI1
certikcli query account $mary --home $DIR_CLI1

# auth

# bank
# cert
# cvm
# gov
# oracle
# shield
