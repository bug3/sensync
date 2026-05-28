//go:build windows

package main

import "github.com/bug3dev/sensync/internal/adapter"

func pickAdapter() (adapter.Adapter, error) { return adapter.NewWindowsAdapter() }
