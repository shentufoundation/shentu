#!/bin/bash

set -e
set -x

# ---------------------
# Transaction sequence
# ---------------------

# Add tokens to mary
$CERTIKCLI tx send $jack $mary 100000000uctk --from $jack -y --home $DIR_CLI0
sleep 6
certikcli query account $jack --home $DIR_CLI1
certikcli query account $bob --home $DIR_CLI1
certikcli query account $mary --home $DIR_CLI1

# auth
$CERTIKCLI tx unlock $jack 500000uctk --from $bob -y --home $DIR_CLI0
sleep 6
certikcli query account $jack --home $DIR_CLI1

# bank
certikcli tx locked-send $mary $jack 500000uctk --from $mary -y --home $DIR_CLI1
sleep 6
certikcli query account $jack --home $DIR_CLI1
certikcli query account $mary --home $DIR_CLI1

# cert
certikcli query cert certifiers --home $DIR_CLI1

$CERTIKCLI tx cert certify-validator certikvalconspub1zcjduepqff623akv26we89w9qz6nk7yq66ms5tlhmn5p7v8rqv4z2ur9puhqmxvkpk --from $bob -y --home $DIR_CLI0
sleep 6
certikcli query cert validators --home $DIR_CLI1

$CERTIKCLI tx cert certify-platform certikvalconspub1zcjduepqff623akv26we89w9qz6nk7yq66ms5tlhmn5p7v8rqv4z2ur9puhqmxvkpk xxxx --from $bob -y --home $DIR_CLI0
sleep 6
certikcli query cert platform certikvalconspub1zcjduepqff623akv26we89w9qz6nk7yq66ms5tlhmn5p7v8rqv4z2ur9puhqmxvkpk --home $DIR_CLI1

$CERTIKCLI tx cert issue-certificate AUDITING ADDRESS C --from $bob -y --home $DIR_CLI0
sleep 6
$CERTIKCLI tx cert issue-certificate COMPILATION SOURCECODEHASH C --compiler A --bytecode-hash B --from $bob -y --home $DIR_CLI0
sleep 6
certikcli query cert certificates --home $DIR_CLI1
id=$(certikcli query cert certificates --home $DIR_CLI1 | grep certificateid)
id=${id:17:60}

$CERTIKCLI tx cert decertify-validator certikvalconspub1zcjduepqff623akv26we89w9qz6nk7yq66ms5tlhmn5p7v8rqv4z2ur9puhqmxvkpk --from $bob -y --home $DIR_CLI0
sleep 6
certikcli query cert validators --home $DIR_CLI1

$CERTIKCLI tx cert revoke-certificate $id --from $bob -y --home $DIR_CLI0
sleep 6
certikcli query cert certificates --home $DIR_CLI1

# cvm
# gov
# oracle
# shield
