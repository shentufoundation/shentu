// Package types includes global types for syncer.
package types

import (
	"path"

	"github.com/spf13/viper"

	"github.com/tendermint/tmlibs/cli"

	querierTypes "github.com/certikfoundation/shentu/toolsets/oracle-operator/querier/types"
	runnerTypes "github.com/certikfoundation/shentu/toolsets/oracle-operator/runner/types"
)

// Config for Relayer
type Config struct {
	Runner  runnerTypes.Config  `mapstructure:"runner"`
	Querier querierTypes.Config `mapstructure:"querier"`
}

var config = Config{
	Runner:  runnerTypes.DefaultConfig(),
	Querier: querierTypes.DefaultConfig(),
}

func initConfig() error {
	v := viper.New()
	v.SetConfigName("oracle-operator")
	v.SetConfigType("toml")
	if home := viper.GetString(cli.HomeFlag); home == "" {
		v.AddConfigPath(".")
		v.AddConfigPath("..")
		v.AddConfigPath("../..")
		v.AddConfigPath("../../..")
		v.AddConfigPath("../../../..")
	} else {
		v.AddConfigPath(home)
		v.AddConfigPath(path.Join(viper.GetString(cli.HomeFlag), "config"))
	}
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	if err := v.Unmarshal(&config); err != nil {
		return err
	}

	return nil
}
