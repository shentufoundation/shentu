#/bin/bash

certikd unsafe-reset-all
#rm -rf ~/.certikd
#rm -rf ~/.certikcli
certikd init node0 --chain-id hychain

certikcli config chain-id hychain
certikcli keys add jack
certikcli keys add alice

#certikd add-genesis-account $(certikcli keys show jack -a) 2000000000uctk
#certikd add-genesis-account $(certikcli keys show alice -a) 2000000000uctk
#certikd add-genesis-account $(certikcli keys show jack -a) 1000000000uctk --vesting-amount 500000000uctk --vesting-end-time 1596003000
#certikd add-genesis-account $(certikcli keys show jack -a) 1000000000uctk --vesting-amount 500000000uctk --vesting-start-time 1599011600 --period 500 --num-periods 2 --triggered=false 
#certikd add-genesis-account $(certikcli keys show jack -a) 1000000000uctk --vesting-amount 500000000uctk --vesting-start-time 1599011600 --period 500 --num-periods 2 --triggered=false 
certikd add-genesis-account $(certikcli keys show jack -a) 1000000000uctk --vesting-amount 500000000uctk --manual
certikd add-genesis-account $(certikcli keys show alice -a) 1000000000uctk --vesting-amount 500000000uctk --vesting-start-time 1599011600 --period 500 --num-periods 2
#certikd add-genesis-account $(certikcli keys show alice -a) 1000000000uctk --vesting-amount 500000000uctk --vesting-start-time 1599011600 --period 500 --num-periods 2

certikd gentx --name jack --amount 100000000uctk
certikd collect-gentxs
certikd start
