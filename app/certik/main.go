package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"

	"github.com/certikfoundation/shentu/v2/app"
	"github.com/certikfoundation/shentu/v2/app/certik/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()

	if err := cmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
