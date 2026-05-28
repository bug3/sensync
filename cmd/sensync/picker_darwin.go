//go:build darwin

package main

import "github.com/bug3/sensync/internal/adapter"

func pickAdapter() (adapter.Adapter, error) { return adapter.NewMacOSAdapter() }
