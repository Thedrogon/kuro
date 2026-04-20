package main

import (
	"os"
	"github.com/thedrogon/kuro/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}