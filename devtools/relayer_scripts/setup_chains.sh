# setup-chains.sh script creates a set of keys for ibc testing
# and starts shentud and regen chains.
#
# Prereq: shentud, regen, relayer (rly) installed
#

set -x

killall shentud
killall regen

rm -rf ~/.shentud
rm -rf ~/.regen

set -e

CURDIR=$(dirname "$0")
coins="10000000000000stake,10000000000000uctk"
delegate="100000000000stake"

# shentu chain setup
shentud init validator --chain-id yulei-1

shentud keys add validator --keyring-backend test
shentud keys add user --output json --keyring-backend test > $CURDIR/yulei_user_key.json
shentud keys add faucet --keyring-backend test

shentud add-genesis-account validator 10000000000000uctk --keyring-backend test
shentud add-genesis-account user $coins --keyring-backend test
shentud add-genesis-account faucet 10000000000000uctk --keyring-backend test

shentud gentx validator 10000000000uctk --chain-id yulei-1 --keyring-backend test

shentud collect-gentxs

# regen chain setup
regen init validator --chain-id aplikigo-1

regen keys add validator --keyring-backend test
regen keys add user --output json --keyring-backend test > $CURDIR/regen_user_key.json
regen keys add faucet --keyring-backend test

regen add-genesis-account validator 10000000000000stake --keyring-backend test
regen add-genesis-account user $coins --keyring-backend test
regen add-genesis-account faucet 10000000000000stake --keyring-backend test

regen gentx validator 10000000000stake --chain-id aplikigo-1 --keyring-backend test

regen collect-gentxs

$CURDIR/one_chain_start.sh shentud yulei-1 ~/.shentud 26657 26656 6060 9090
$CURDIR/one_chain_start.sh regen aplikigo-1 ~/.regen 26557 26556 6061 9091
