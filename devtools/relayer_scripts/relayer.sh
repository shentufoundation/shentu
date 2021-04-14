set -e
set -x

CURDIR=$(dirname "$0")

rm -rf ~/.relayer

rly config init

rly chains add -f $CURDIR/shentu.json
rly chains add -f $CURDIR/regen.json

SHENTU_MNEMONIC=$(jq -r '.mnemonic' $CURDIR/yulei_user_key.json)
REGEN_MNEMONIC=$(jq -r '.mnemonic' $CURDIR/regen_user_key.json)

rly keys restore yulei-1 testkey "$SHENTU_MNEMONIC"
rly keys restore aplikigo-1 testkey "$REGEN_MNEMONIC"

rly light init yulei-1 -f
rly light init aplikigo-1 -f

rly paths add yulei-1 aplikigo-1 demo -f $CURDIR/paths/demo.json

rly tx link demo -d -o 3s
# rly tx link demo -d -o 100s -r 10

# rly chains list
# rly paths list

# transfer from certik to regen
# rly tx transfer yulei-1 aplikigo-1 1000000000000uctk $(rly chains address aplikigo-1)

# rly query unrelayed-packets demo
# rly tx relay-packets demo -d
# rly query unrelayed-acknowledgements demo
# rly tx relay-acknowledgements demo -d

# rly query bal yulei-1 testkey
# rly query bal aplikigo-1 testkey
# certik query bank balances $(rly chains address yulei-1)
# regen query bank balances $(rly chains address aplikigo-1) --node tcp://0.0.0.0:26557

# clean up
# rm $CURDIR/yulei_user_key.json $CURDIR/regen_user_key.json
# rm $CURDIR/yulei-1.log $CURDIR/aplikigo-1.log 
# killall certik && killall regen
