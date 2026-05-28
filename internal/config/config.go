package config

import (
	"fmt"
	"math"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Version  int    `toml:"version"`
	Mouse    Device `toml:"mouse"`
	Trackpad Device `toml:"trackpad"`
}

type Device struct {
	Sensitivity   float64 `toml:"sensitivity"`
	Acceleration  bool    `toml:"acceleration"`
	NaturalScroll bool    `toml:"natural_scroll"`
	ScrollSpeed   float64 `toml:"scroll_speed"`
}

func Default() Config {
	return Config{
		Version: 1,
		Mouse: Device{
			Sensitivity:   1.0,
			Acceleration:  false,
			NaturalScroll: false,
			ScrollSpeed:   1.0,
		},
		Trackpad: Device{
			Sensitivity:   1.0,
			Acceleration:  false,
			NaturalScroll: true,
			ScrollSpeed:   1.0,
		},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	meta, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("load config %s: %w", path, err)
	}
	if undecoded := meta.Undecoded(); len(undecoded) > 0 {
		return Config{}, fmt.Errorf("unknown keys in %s: %v", path, undecoded)
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	if c.Version != 1 {
		return fmt.Errorf("unsupported version %d (only 1 is supported)", c.Version)
	}
	if err := validateDevice("mouse", c.Mouse); err != nil {
		return err
	}
	if err := validateDevice("trackpad", c.Trackpad); err != nil {
		return err
	}
	return nil
}

func validateDevice(name string, d Device) error {
	if math.IsNaN(d.Sensitivity) || math.IsInf(d.Sensitivity, 0) {
		return fmt.Errorf("%s.sensitivity must be a finite number", name)
	}
	if d.Sensitivity < 0.1 || d.Sensitivity > 3.0 {
		return fmt.Errorf("%s.sensitivity must be in [0.1, 3.0], got %v", name, d.Sensitivity)
	}
	if math.IsNaN(d.ScrollSpeed) || math.IsInf(d.ScrollSpeed, 0) {
		return fmt.Errorf("%s.scroll_speed must be a finite number", name)
	}
	if d.ScrollSpeed < 0.1 || d.ScrollSpeed > 5.0 {
		return fmt.Errorf("%s.scroll_speed must be in [0.1, 5.0], got %v", name, d.ScrollSpeed)
	}
	return nil
}
