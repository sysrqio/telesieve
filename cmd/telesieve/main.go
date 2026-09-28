package main

import (
	"fmt"
	"os"

	"github.com/sysrqio/telesieve/internal/cmdroot"
)

var version = "0.1.0"

func main() {
	root := cmdroot.NewRoot(version)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
