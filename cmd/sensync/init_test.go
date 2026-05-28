package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitWritesExampleToUserConfigDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("HOME", dir)
	t.Setenv("APPDATA", dir)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"init"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init: %v", err)
	}

	path, err := userConfigPath()
	if err != nil {
		t.Fatalf("userConfigPath: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !strings.Contains(string(got), "version = 1") {
		t.Errorf("config missing version line; contents:\n%s", string(got))
	}
}

func TestInitRefusesToOverwriteWithoutForce(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("HOME", dir)
	t.Setenv("APPDATA", dir)

	path, _ := userConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("seed dir: %v", err)
	}
	if err := os.WriteFile(path, []byte("existing"), 0o600); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	cmd := newRootCmd()
	cmd.SetArgs([]string{"init"})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected init to refuse overwriting existing file")
	}
	var ce cliError
	if !errors.As(err, &ce) || ce.code != 1 {
		t.Errorf("expected cliError with code 1, got %v", err)
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' in error, got %q", err.Error())
	}

	// Original contents preserved.
	got, _ := os.ReadFile(path)
	if string(got) != "existing" {
		t.Errorf("file was overwritten: %q", string(got))
	}
}

func TestInitForceOverwrites(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("HOME", dir)
	t.Setenv("APPDATA", dir)

	path, _ := userConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("seed dir: %v", err)
	}
	if err := os.WriteFile(path, []byte("stale"), 0o600); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	cmd := newRootCmd()
	cmd.SetArgs([]string{"init", "--force"})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init --force: %v", err)
	}
	got, _ := os.ReadFile(path)
	if !strings.Contains(string(got), "version = 1") {
		t.Errorf("force did not write template; got: %q", string(got))
	}
}
