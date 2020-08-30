package main

import (
	"fmt"
	"os"

	"github.com/dstdfx/solid-broccoli/cmd/solid-broccoli/app"
)

func main() {
	if err := app.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
