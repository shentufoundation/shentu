#!/bin/bash
set -x

curdir=$(pwd)
binary=certik
app=certik
denom=uctk
mdenom=ctk
ethaddr=0xF8F79be5D59c58CAf09560D94e6943d231852706

killall $app
rm -rf ~/.$app

$binary init --chain-id=test testid

# only for making stake -> ustake
# cat ~/.$app/config/genesis.json | sed -e 's:stake:ustake:g' > tmp.json
# mv tmp.json ~/.$app/config/genesis.json

echo "tribe concert jungle next slab odor mixed doll struggle crouch flush post rack pen taxi pitch first poem anxiety sea dilemma blanket virus february" | $binary keys add alice --keyring-backend test --recover
echo "aisle text grocery review hello sort ski winner foil cradle keep sight success toss garment tunnel toilet under glue plate farm century mule inmate" | $binary keys add bob --keyring-backend test --recover
sed -i 's/"voting_period": "172800s"/"voting_period": "3s"/' ~/.$app/config/genesis.json
sed -i 's/"unbonding_time": "1814400s"/"unbonding_time": "30s"/' ~/.$app/config/genesis.json
sed -i 's/timeout_propose = "3s"/timeout_propose = "600ms"/' ~/.$app/config/config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/' ~/.$app/config/config.toml
$binary add-genesis-account $($binary keys show alice -a --keyring-backend test) 1000000000000000000000000$denom
$binary add-genesis-certifier $($binary keys show alice -a --keyring-backend test)
$binary add-genesis-shield-admin $($binary keys show alice -a --keyring-backend test)
#$binary gentx alice 1000002345$denom $ethaddr $($binary keys show alice -a --keyring-backend test) --keyring-backend test --chain-id testing
$binary gentx alice 10000000000000$denom --keyring-backend test --chain-id test
$binary collect-gentxs

jq '.app_state.bank.denom_metadata = [{"base": "'$denom'", "denom_units": [{"aliases": [], "denom": "'$denom'", "exponent": 0},{"aliases": [], "denom": "'$mdenom'", "exponent": 6}],"description": "Test token of certik chain", "display": "'$mdenom'"}]' ~/.$app/config/genesis.json > tmp.json
mv tmp.json ~/.$app/config/genesis.json
chmod 777 ~/.$app/config/*
sed -i 's/enable = false/enable = true/g' ~/.$app/config/app.toml

echo "$binary start &" | bash - > "$curdir"/custom/$app.log 2>&1
sleep 2
$binary tx staking unbond $($binary keys show alice -a --bech=val --keyring-backend test) 10000000$denom --from alice --keyring-backend test -y -b block --chain-id test
$binary tx bank locked-send alice $($binary keys show bob -a --keyring-backend test) 10000000000$denom --unlocker $($binary keys show alice -a --keyring-backend test) --from alice --keyring-backend test -y -b block --chain-id test
$binary tx staking delegate $($binary keys show alice -a --bech=val --keyring-backend test) 20000000$denom --from bob --keyring-backend test -y -b block --chain-id test
$binary tx staking unbond $($binary keys show alice -a --bech=val --keyring-backend test) 10000000$denom --from bob --keyring-backend test -y -b block --chain-id test
