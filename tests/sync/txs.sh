#!/bin/bash

set -e
set -x

# ---------------------
# Transaction sequence
# ---------------------

# Add tokens to mary
$CERTIKCLI tx send $jack $mary 100000000uctk --from $jack -y --home ~/.synctest/node0/certikcli
sleep 3
$CERTIKCLI query account $jack --home ~/.synctest/node0/certikcli
$CERTIKCLI query account $mary --home ~/.synctest/node0/certikcli

# auth
# bank
# cert
# cvm
# gov
# oracle
# shield
