// Package types includes types used by querier.
package types

// Config defines the data type of configurations for Querier.
type Config struct {
	Endpoint   string `mapstructure:"endpoint"`
	Method     string `mapstructure:"method"` // HTTP method, `GET` or `POST`
	Timeout    int    `mapstructure:"timeout"`
	RetryTimes int    `mapstructure:"retry_times"`
}

// DefaultConfig is the default configuration for querier.
func DefaultConfig() Config {
	return Config{
		Endpoint:   "",
		Method:     "GET",
		Timeout:    300,
		RetryTimes: 3,
	}
}
