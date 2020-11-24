#!/bin/bash

set -e
set -x

# ---------------------------------------------------------------
# Transaction sequence
#
# use `$CERTIKCLI0` to send txs from jack or bob
# use `$CERTIKCLI1` to send txs from mary
# ---------------------------------------------------------------

# Add tokens to mary
$CERTIKCLI0 tx send $jack $mary 100000000uctk --from $jack -y
sleep 6
$CERTIKCLI1 query account $jack
$CERTIKCLI1 query account $bob
$CERTIKCLI1 query account $mary

# auth
$CERTIKCLI0 tx unlock $jack 500000uctk --from $bob -y
sleep 6
$CERTIKCLI1 query account $jack

# bank
$CERTIKCLI1 tx locked-send $mary $jack 500000uctk --from $mary -y
sleep 6
$CERTIKCLI1 query account $jack
$CERTIKCLI1 query account $mary

# cert
$CERTIKCLI1 query cert certifiers

$CERTIKCLI0 tx cert certify-validator certikvalconspub1zcjduepqff623akv26we89w9qz6nk7yq66ms5tlhmn5p7v8rqv4z2ur9puhqmxvkpk --from $bob -y
sleep 6
$CERTIKCLI1 query cert validators

$CERTIKCLI0 tx cert certify-platform certikvalconspub1zcjduepqff623akv26we89w9qz6nk7yq66ms5tlhmn5p7v8rqv4z2ur9puhqmxvkpk xxxx --from $bob -y
sleep 6
$CERTIKCLI1 query cert platform certikvalconspub1zcjduepqff623akv26we89w9qz6nk7yq66ms5tlhmn5p7v8rqv4z2ur9puhqmxvkpk

$CERTIKCLI0 tx cert issue-certificate AUDITING ADDRESS C --from $bob -y
sleep 6
$CERTIKCLI0 tx cert issue-certificate COMPILATION SOURCECODEHASH C --compiler A --bytecode-hash B --from $bob -y
sleep 6
$CERTIKCLI1 query cert certificates
id=$($CERTIKCLI1 query cert certificates | grep certificateid)
id=${id:17:60}

$CERTIKCLI0 tx cert decertify-validator certikvalconspub1zcjduepqff623akv26we89w9qz6nk7yq66ms5tlhmn5p7v8rqv4z2ur9puhqmxvkpk --from $bob -y
sleep 6
$CERTIKCLI1 query cert validators

$CERTIKCLI0 tx cert revoke-certificate $id --from $bob -y
sleep 6
$CERTIKCLI1 query cert certificates

# cvm
txhash=$($CERTIKCLI1 tx cvm deploy $PROJ_ROOT/tests/simple.sol --from $mary -y | grep txhash)
txhash=${txhash:8}
sleep 6
addr=$($CERTIKCLI1 query tx $txhash | grep value | sed -n '2p')
addr=${addr:13}

$CERTIKCLI1 tx cvm call $addr set 123 --from $mary -y
sleep 6

$CERTIKCLI1 tx cvm call $addr get --from $mary -y
sleep 6

# oracle
# shield
