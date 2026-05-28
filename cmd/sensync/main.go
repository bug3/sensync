package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		var ce cliError
		if errors.As(err, &ce) {
			fmt.Fprintln(os.Stderr, "error:", ce.Error())
			os.Exit(ce.code)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
