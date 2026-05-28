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
			if msg := ce.Error(); msg != "" {
				fmt.Fprintln(os.Stderr, "error:", msg)
			}
			os.Exit(ce.code)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
