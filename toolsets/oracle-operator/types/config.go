// Package types includes global types for syncer.
package types

import (
	"path"
	"strings"

	"github.com/spf13/viper"

	"github.com/tendermint/tmlibs/cli"
)

// DefaultConfigFileName defines the default config file name.
const DefaultConfigFileName = "oracle-operator.toml"

// Config for Relayer
type Config struct {
	Strategy   map[Client]Strategy `mapstructure:"strategy"`
	Method     string              `mapstructure:"method"` // HTTP method, `GET` or `POST`
	Timeout    int                 `mapstructure:"timeout"`
	RetryTimes int                 `mapstructure:"retry_times"`
}

// Default config values.
var config = Config{
	Strategy:   make(map[Client]Strategy),
	Method:     "GET",
	Timeout:    300,
	RetryTimes: 3,
}

func initConfig() error {
	v := viper.New()
	configFileName := viper.GetString(FlagConfigFile)
	if configFileName == "" {
		configFileName = DefaultConfigFileName
	}
	configName, configType := getBasenameAndExtension(configFileName)
	v.SetConfigName(configName)
	v.SetConfigType(configType)
	if home := viper.GetString(cli.HomeFlag); home == "" {
		v.AddConfigPath(".")
		v.AddConfigPath("..")
		v.AddConfigPath("../..")
	} else {
		v.AddConfigPath(home)
		v.AddConfigPath(path.Join(home, "config"))
		v.AddConfigPath(path.Join(home, "oracle"))
		v.AddConfigPath(path.Join(home, "operator"))
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

func getBasenameAndExtension(fileName string) (string, string) {
	fileName = path.Base(fileName)
	extension := path.Ext(fileName)
	basename := strings.TrimSuffix(fileName, extension)
	return basename, strings.TrimPrefix(extension, ".")
}
