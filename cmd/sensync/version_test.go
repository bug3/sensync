package main

import (
	"bytes"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"version"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	const want = "sensync 0.0.1-dev\n"
	if got := out.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
