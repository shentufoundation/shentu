#!/bin/bash

set +e
set -x

checkConsensus() {
  if grep -q "CONSENSUS FAILURE" $DIR/node1/log.txt; then
    set -e
    killall certikd
    echo "CONSENSUS FAILURE!"
    exit 1
  fi
}

# ------------------------------------------------------------------
# Transaction sequence
#
# use `$CERTIKCLI0` to send txs from jack or bob running old binary
# use `$CERTIKCLI1` to send txs from mary running latest binary
# ------------------------------------------------------------------

# Add tokens to mary
$CERTIKCLI0 tx send $jack $mary 200000000uctk --from $jack -y
sleep 6
$CERTIKCLI1 query account $jack
$CERTIKCLI1 query account $bob
$CERTIKCLI1 query account $mary

checkConsensus

# auth

$CERTIKCLI0 tx unlock $jack 500000uctk --from $bob -y
sleep 6
$CERTIKCLI1 query account $jack

checkConsensus

# bank

$CERTIKCLI1 tx locked-send $mary $jack 500000uctk --from $mary -y
sleep 6
$CERTIKCLI1 query account $jack
$CERTIKCLI1 query account $mary

checkConsensus

# cert

$CERTIKCLI1 query cert certifiers
validator=$($CERTIKCLI1 query staking validators | grep conspubkey)
validator=${validator:14}
$CERTIKCLI0 tx cert certify-validator $validator --from $bob -y
sleep 6
$CERTIKCLI1 query cert validator $validator

$CERTIKCLI0 tx cert certify-platform $validator xxxx --from $bob -y
sleep 6
$CERTIKCLI1 query cert platform $validator

$CERTIKCLI0 tx cert issue-certificate AUDITING ADDRESS C --from $bob -y
sleep 6
$CERTIKCLI0 tx cert issue-certificate COMPILATION SOURCECODEHASH C --compiler A --bytecode-hash B --from $bob -y
sleep 6
$CERTIKCLI1 query cert certificates
id=$($CERTIKCLI1 query cert certificates | grep certificateid)
id=${id:17:60}

$CERTIKCLI0 tx cert revoke-certificate $id --from $bob -y
sleep 6
$CERTIKCLI1 query cert certificates

checkConsensus

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

checkConsensus

# oracle

$CERTIKCLI1 tx oracle create-operator $mary 100000uctk --from $mary -y
sleep 6
$CERTIKCLI1 query oracle operators

$CERTIKCLI0 tx oracle create-task --contract A --function B --bounty 10000uctk --wait 4 --from $bob -y
sleep 6
$CERTIKCLI1 query oracle task --contract A --function B

$CERTIKCLI1 tx oracle deposit-collateral $mary 30000uctk --from $mary -y
sleep 6

$CERTIKCLI1 tx oracle withdraw-collateral $mary 10000uctk --from $mary -y
sleep 6
$CERTIKCLI1 query oracle operators

$CERTIKCLI1 tx oracle respond-to-task --contract A --function B --score 99 --from $mary -y
sleep 6
$CERTIKCLI1 query oracle response --contract A --function B --operator $mary
$CERTIKCLI1 query oracle operator $mary

$CERTIKCLI1 tx oracle claim-reward $mary --from $mary -y
sleep 6
$CERTIKCLI1 query oracle operator $mary

$CERTIKCLI0 tx oracle delete-task --contract A --function B --force=true --from $bob -y
sleep 6

$CERTIKCLI1 tx oracle remove-operator $mary --from $mary -y
sleep 6
$CERTIKCLI1 query oracle operators
$CERTIKCLI1 query oracle withdraws

checkConsensus

# shield

val=$($CERTIKCLI1 query staking validators | grep operatoraddress)
val=${val:19}
$CERTIKCLI0 tx staking delegate $val 100000000uctk --from $jack -y
$CERTIKCLI0 tx staking delegate $val 100000000uctk --from $bob -y
$CERTIKCLI1 tx staking delegate $val 50000000uctk --from $mary -y
sleep 6
$CERTIKCLI1 query account $jack
$CERTIKCLI1 query account $bob
$CERTIKCLI1 query account $mary

$CERTIKCLI0 tx shield deposit-collateral 100000000uctk --from $jack -y
$CERTIKCLI0 tx shield deposit-collateral 100000000uctk --from $bob -y
$CERTIKCLI1 tx shield deposit-collateral 50000000uctk --from $mary -y
sleep 6
$CERTIKCLI1 query shield provider $jack
$CERTIKCLI1 query shield provider $bob
$CERTIKCLI1 query shield provider $mary

$CERTIKCLI0 tx shield withdraw-collateral 1000000uctk --from $bob -y
sleep 6
$CERTIKCLI1 query shield provider $bob

$CERTIKCLI0 tx shield create-pool 1000000uctk bob $bob --native-deposit 110000uctk --shield-limit 100000000 --from $bob -y
sleep 6
$CERTIKCLI1 query shield pool 1

$CERTIKCLI0 tx shield update-pool 1 --shield 4000000uctk --native-deposit 120000uctk --shield-limit 150000000 --from $bob -y
sleep 6
$CERTIKCLI1 query shield pool 1

$CERTIKCLI0 tx shield pause-pool 1 --from $bob -y
sleep 6
$CERTIKCLI1 query shield pool 1

$CERTIKCLI0 tx shield resume-pool 1 --from $bob -y
sleep 6
$CERTIKCLI1 query shield pool 1

$CERTIKCLI1 tx shield purchase 1 50000000uctk haha --from $mary -y
sleep 6
$CERTIKCLI1 query shield pool-purchaser 1 $mary

$CERTIKCLI0 tx shield update-sponsor 1 mary $mary --from $bob -y
sleep 6
$CERTIKCLI1 query shield pool 1

$CERTIKCLI0 tx shield stake-for-shield 1 50000000uctk haha --from $jack -y
sleep 6
$CERTIKCLI1 query shield staked-for-shield 1 $jack

$CERTIKCLI0 tx shield unstake-from-shield 1 30000000uctk --from $jack -y
sleep 6
$CERTIKCLI1 query shield staked-for-shield 1 $jack

checkConsensus

# gov

$CERTIKCLI1 tx gov submit-proposal certifier-update $PROJ_ROOT/tests/sync/certifier_update.json --from $mary -y
sleep 6
$CERTIKCLI1 query gov proposal 1

$CERTIKCLI0 tx gov deposit 1 520000000uctk --from $bob -y
sleep 6
$CERTIKCLI1 query gov proposal 1

$CERTIKCLI0 tx gov vote 1 yes --from $bob -y
sleep 6
$CERTIKCLI1 query gov proposal 1
$CERTIKCLI1 query cert certifiers

$CERTIKCLI1 tx gov submit-proposal shield-claim $PROJ_ROOT/tests/sync/shield_claim.json --from $mary -y
sleep 6
$CERTIKCLI1 query gov proposal 2

$CERTIKCLI0 tx gov vote 2 yes --from $bob -y
sleep 6
$CERTIKCLI1 query gov proposal 2

$CERTIKCLI0 tx cert issue-certificate IDENTITY ADDRESS $jack --from $bob -y
sleep 6
$CERTIKCLI1 query cert certificates

$CERTIKCLI0 tx cert issue-certificate IDENTITY ADDRESS $bob --from $bob -y
sleep 6
$CERTIKCLI1 query cert certificates

$CERTIKCLI0 tx cert issue-certificate IDENTITY ADDRESS $mary --from $bob -y
sleep 6
$CERTIKCLI1 query cert certificates

$CERTIKCLI0 tx gov vote 2 yes --from $jack -y
$CERTIKCLI0 tx gov vote 2 yes --from $bob -y
$CERTIKCLI1 tx gov vote 2 yes --from $mary -y
sleep 6
$CERTIKCLI1 query gov proposal 2

checkConsensus

# exiting without consensus failure

set -e
