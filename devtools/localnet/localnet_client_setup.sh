#!/bin/bash

rm -rf ~/.shentud
ln -s ~/node0/shentud  ~/.shentud

### public keys
#echo 'export NODE0_KEY=$(shentud keys show node0 -a --home /root/node0/shentud)' >> ~/.bashrc
#echo 'export NODE1_KEY=$(shentud keys show node1 -a --home /root/node1/shentud)' >> ~/.bashrc
#echo 'export NODE2_KEY=$(shentud keys show node2 -a --home /root/node2/shentud)' >> ~/.bashrc
#echo 'export NODE3_KEY=$(shentud keys show node3 -a --home /root/node3/shentud)' >> ~/.bashrc

### private keys added to keychain
#echo -e "$(cat /root/node0/shentud/key_seed.json | sed -r 's/^([^\"]+\"+){3}((\"*[^\"]+)).*/\2/')" "\n" | shentud keys add --recover node0
#echo -e "$(cat /root/node1/shentud/key_seed.json | sed -r 's/^([^\"]+\"+){3}((\"*[^\"]+)).*/\2/')" "\n" | shentud keys add --recover node1
#echo -e "$(cat /root/node2/shentud/key_seed.json | sed -r 's/^([^\"]+\"+){3}((\"*[^\"]+)).*/\2/')" "\n" | shentud keys add --recover node2
#echo -e "$(cat /root/node3/shentud/key_seed.json | sed -r 's/^([^\"]+\"+){3}((\"*[^\"]+)).*/\2/')" "\n" | shentud keys add --recover node3