module github.com/certikfoundation/shentu

go 1.14

require (
	github.com/cosmos/cosmos-sdk v0.39.1
	github.com/gorilla/mux v1.7.4
	github.com/hyperledger/burrow v0.30.5
	github.com/magiconair/properties v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.6
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.3 // new version doesnt work
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.7
	github.com/tendermint/tm-db v0.5.1
	github.com/tmthrgd/go-hex v0.0.0-20190303111820-0bdcb15db631
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/go-interpreter/wagon v0.0.0 => github.com/perlin-network/wagon v0.3.1-0.20180825141017-f8cb99b55a39

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4

replace github.com/spf13/viper v1.6.1 => github.com/spf13/viper v1.4.0
