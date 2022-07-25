#!/bin/bash
set -o errexit -o nounset -o pipefail
command -v shellcheck > /dev/null && shellcheck "$0"

tag=v1.0.0
for contract in hackatom reflect
do
  echo $contract
  url="https://github.com/CosmWasm/cosmwasm/releases/download/$tag/${contract}.wasm"
  echo "Downloading $url ..."
  wget -O "${contract}.wasm" "$url"
done
