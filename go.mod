module github.com/certikfoundation/shentu

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.39.1
	github.com/gorilla/mux v1.8.0
	github.com/hyperledger/burrow v0.30.5
	github.com/magiconair/properties v1.8.4
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/tendermint v0.33.8
	github.com/tendermint/tm-db v0.5.1
	github.com/tendermint/tmlibs v0.9.0
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/hyperledger/burrow v0.30.5 => github.com/certikfoundation/burrow v0.1.1
