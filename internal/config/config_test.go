package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "sensync.toml")
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	return p
}

func TestLoadValidConfig(t *testing.T) {
	body := `
version = 1

[mouse]
sensitivity = 1.0
acceleration = false
natural_scroll = false
scroll_speed = 1.0

[trackpad]
sensitivity = 1.0
acceleration = false
natural_scroll = true
scroll_speed = 1.0
`
	cfg, err := Load(writeTempConfig(t, body))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Version != 1 {
		t.Errorf("version: got %d, want 1", cfg.Version)
	}
	if cfg.Mouse.Sensitivity != 1.0 {
		t.Errorf("mouse.sensitivity: got %v, want 1.0", cfg.Mouse.Sensitivity)
	}
	if cfg.Mouse.Acceleration {
		t.Errorf("mouse.acceleration: got true, want false")
	}
	if !cfg.Trackpad.NaturalScroll {
		t.Errorf("trackpad.natural_scroll: got false, want true")
	}
}

func TestLoadRejectsUnknownKey(t *testing.T) {
	body := `
version = 1
mystery_field = "boo"

[mouse]
sensitivity = 1.0
acceleration = false
natural_scroll = false
scroll_speed = 1.0

[trackpad]
sensitivity = 1.0
acceleration = false
natural_scroll = true
scroll_speed = 1.0
`
	_, err := Load(writeTempConfig(t, body))
	if err == nil {
		t.Fatal("expected unknown-key error, got nil")
	}
}

func TestLoadRejectsBadVersion(t *testing.T) {
	body := `
version = 99

[mouse]
sensitivity = 1.0
acceleration = false
natural_scroll = false
scroll_speed = 1.0

[trackpad]
sensitivity = 1.0
acceleration = false
natural_scroll = true
scroll_speed = 1.0
`
	_, err := Load(writeTempConfig(t, body))
	if err == nil {
		t.Fatal("expected version error, got nil")
	}
}

func TestLoadRejectsOutOfRangeSensitivity(t *testing.T) {
	body := `
version = 1

[mouse]
sensitivity = 10.0
acceleration = false
natural_scroll = false
scroll_speed = 1.0

[trackpad]
sensitivity = 1.0
acceleration = false
natural_scroll = true
scroll_speed = 1.0
`
	_, err := Load(writeTempConfig(t, body))
	if err == nil {
		t.Fatal("expected range error, got nil")
	}
}

func TestLoadRejectsOutOfRangeScrollSpeed(t *testing.T) {
	body := `
version = 1

[mouse]
sensitivity = 1.0
acceleration = false
natural_scroll = false
scroll_speed = 10.0

[trackpad]
sensitivity = 1.0
acceleration = false
natural_scroll = true
scroll_speed = 1.0
`
	_, err := Load(writeTempConfig(t, body))
	if err == nil {
		t.Fatal("expected scroll_speed range error, got nil")
	}
}

func TestLoadAcceptsMissingFieldsViaDefaults(t *testing.T) {
	body := `
version = 1

[mouse]
sensitivity = 1.0

[trackpad]
sensitivity = 1.0
`
	cfg, err := Load(writeTempConfig(t, body))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	// BurntSushi/toml only overwrites fields present in the file. Since
	// Default() populates the struct first, fields omitted from the section
	// retain their Default() values.
	if cfg.Mouse.ScrollSpeed != 1.0 {
		t.Errorf("missing scroll_speed should fall back to default 1.0, got %v", cfg.Mouse.ScrollSpeed)
	}
	if cfg.Mouse.Acceleration != false {
		t.Errorf("missing acceleration should fall back to default false, got %v", cfg.Mouse.Acceleration)
	}
}
