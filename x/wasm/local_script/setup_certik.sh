#!/bin/bash
set -x

curdir=$(pwd)
binary=certik
app=certik
denom=uctk
mdenom=ctk

killall $app
rm -rf ~/.$app

$binary init --chain-id=testing testid
certik config chain-id testing
certik config keyring-backend test

echo "tribe concert jungle next slab odor mixed doll struggle crouch flush post rack pen taxi pitch first poem anxiety sea dilemma blanket virus february" | $binary keys add alice --keyring-backend test --recover
echo "aisle text grocery review hello sort ski winner foil cradle keep sight success toss garment tunnel toilet under glue plate farm century mule inmate" | $binary keys add bob --keyring-backend test --recover
sed -i 's/"voting_period": "172800s"/"voting_period": "120s"/' ~/.$app/config/genesis.json
$binary add-genesis-account $($binary keys show alice -a --keyring-backend test) 1000000000000$denom
#$binary add-genesis-certifier $($binary keys show alice -a --keyring-backend test)
#$binary add-genesis-shield-admin $($binary keys show alice -a --keyring-backend test)
$binary gentx alice 100000000000$denom --keyring-backend test --chain-id testing
$binary collect-gentxs

echo "$binary start &" | bash - > /tmp/$app.log 2>&1
sleep 5
