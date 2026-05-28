//go:build windows

package adapter

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// writeRegistryValue parses a Step of Kind StepRegSet with
// Target = `HKCU\...\...` and Args = [name, type, value], then performs the write.
// Only HKCU is supported in MVP (no elevation required).
func writeRegistryValue(s Step) error {
	if !strings.HasPrefix(s.Target, `HKCU\`) {
		return fmt.Errorf("registry root must be HKCU, got %q", s.Target)
	}
	subkey := strings.TrimPrefix(s.Target, `HKCU\`)
	if len(s.Args) != 3 {
		return fmt.Errorf("reg set needs Args=[name, type, value], got %d args", len(s.Args))
	}
	name, typ, value := s.Args[0], s.Args[1], s.Args[2]
	k, _, err := registry.CreateKey(registry.CURRENT_USER, subkey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("open %s: %w", s.Target, err)
	}
	defer k.Close()
	switch typ {
	case "REG_SZ":
		return k.SetStringValue(name, value)
	case "REG_DWORD":
		n, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return fmt.Errorf("parse DWORD value %q: %w", value, err)
		}
		return k.SetDWordValue(name, uint32(n))
	default:
		return fmt.Errorf("unsupported registry type %q", typ)
	}
}
