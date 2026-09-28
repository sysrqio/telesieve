package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/sysrqio/telesieve/internal/cmdroot"
)

var version = "0.1.0"

func main() {
	root := cmdroot.NewRoot(version)
	if err := root.Execute(); err != nil {
		var exitErr *cmdroot.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.Code)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
