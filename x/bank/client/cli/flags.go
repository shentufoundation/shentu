package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagUnlocker = "unlocker"
)

func FlagAddUnlocker() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(FlagUnlocker, "", "FlagUnlocking")
	return fs
}
