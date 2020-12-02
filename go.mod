module github.com/certikfoundation/shentu

go 1.15

require (
	github.com/armon/go-metrics v0.3.4
	github.com/cosmos/cosmos-sdk v0.40.0-rc3
	github.com/gogo/protobuf v1.3.1
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.15.2
	github.com/hyperledger/burrow v0.30.5
	github.com/magiconair/properties v1.8.4
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/afero v1.3.4 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/tendermint v0.34.0-rc6
	github.com/tendermint/tm-db v0.6.2
	github.com/tendermint/tmlibs v0.9.0
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc
	golang.org/x/crypto v0.0.0-20201124201722-c8d3bf9c5392
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20201119123407-9b1e624d6bc4 // indirect
	google.golang.org/grpc v1.33.2
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
