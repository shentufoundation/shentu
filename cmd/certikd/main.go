package main

import (
	"os"

	"github.com/certikfoundation/shentu/cmd/certikd/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := cmd.Execute(rootCmd); err != nil {
		os.Exit(1)
	}
}
