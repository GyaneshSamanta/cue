package main

import (
	"fmt"
	"os"

	"github.com/GyaneshSamanta/gyanesh-help/cmd"
)

var (
	Version   = "1.0.0"
	BuildDate = "dev"
)

func main() {
	cmd.SetVersionInfo(Version, BuildDate)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
