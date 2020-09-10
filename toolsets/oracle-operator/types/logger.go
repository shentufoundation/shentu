package types

import (
	"os"

	"github.com/spf13/viper"

	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
)

var logger log.Logger

func initLogger() (err error) {
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	if viper.GetBool(cli.TraceFlag) {
		logger = log.NewTracingLogger(logger)
	}
	if viper.GetString(FlagLogLevel) != "" {
		logger, err = flags.ParseLogLevel(viper.GetString(FlagLogLevel), logger, tmconfig.DefaultLogLevel())
		if err != nil {
			return err
		}
	}
	logger = logger.With("module", "Oracle-Operator")
	return nil
}
