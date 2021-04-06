set -e
set -x

BINARY=$1
CHAINID=$2
CHAINDIR=$3
RPCPORT=$4
P2PPORT=$5
PROFPORT=$6
GRPCPORT=$7

CURDIR=$(dirname "$0")

# Check platform
platform='unknown'
unamestr=`uname`
if [ "$unamestr" = 'Linux' ]; then
	platform='linux'
fi
if [ "$unamestr" = 'Darwin' ]; then
	platform='darwin'
fi

# Set proper defaults and change ports (use a different sed for Mac or Linux)
echo "Change settings in config.toml file..."
if [ $platform = 'linux' ] || [ $platform = 'darwin' ]; then
	sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPCPORT"'"#g' $CHAINDIR/config/config.toml
	sed -i 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"$P2PPORT"'"#g' $CHAINDIR/config/config.toml
	sed -i 's#"localhost:6060"#"localhost:'"$P2PPORT"'"#g' $CHAINDIR/config/config.toml
	sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAINDIR/config/config.toml
	sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAINDIR/config/config.toml
	sed -i 's/index_all_keys = false/index_all_keys = true/g' $CHAINDIR/config/config.toml
else
	sed -i '' 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPCPORT"'"#g' $CHAINDIR/$CHAINID/config/config.toml
	sed -i '' 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"$P2PPORT"'"#g' $CHAINDIR/$CHAINID/config/config.toml
	sed -i '' 's#"localhost:6060"#"localhost:'"$P2PPORT"'"#g' $CHAINDIR/$CHAINID/config/config.toml
	sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAINDIR/$CHAINID/config/config.toml
	sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAINDIR/$CHAINID/config/config.toml
	sed -i '' 's/index_all_keys = false/index_all_keys = true/g' $CHAINDIR/$CHAINID/config/config.toml
fi

# Start the chain
$BINARY --home $CHAINDIR start --pruning=nothing --grpc.address="0.0.0.0:$GRPCPORT" --rpc.laddr="tcp://0.0.0.0:$RPCPORT" > $CURDIR/$CHAINID.log 2>&1 &
