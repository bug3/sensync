# sensync

Cross-platform mouse and trackpad sensitivity sync. One config, three operating systems.

## What it does

`sensync` reads a single `sensync.toml` and applies the equivalent OS settings on Linux (Hyprland), macOS, and Windows. At `sensitivity = 1.0` with `acceleration = false` all three operating systems are guaranteed to produce raw 1:1 input; other values use per-OS approximations.

## What it does not do

- Per-app sensitivity overrides.
- Polling rate (firmware concern).
- Daemon / boot-time auto-apply.
- Cloud config sync (git is the sync layer).
- Linux backends other than Hyprland.
- Windows non-Precision-Touchpad trackpads.

## Install

```sh
git clone https://github.com/bug3dev/sensync
cd sensync
go build ./cmd/sensync
```

The resulting `./sensync` binary is self-contained.

## Quick start

1. Copy `configs/sensync.example.toml` to `./sensync.toml` and edit values.
2. On Linux (Hyprland), add `source = ~/.config/hypr/sensync.conf` to your `hyprland.conf` once.
3. Run `./sensync apply` on each host. Re-run after changing `./sensync.toml`.

## Commands

```
sensync apply              Apply ./sensync.toml on this host
sensync apply --dry-run    Print planned changes without applying
sensync apply --yes        Skip confirmation prompts
sensync apply --config X   Explicit config path
sensync get                Print live system state as TOML
sensync version            Print sensync version
```

## Caveats

- macOS: `defaults` writes take effect for newly launched processes; log out for full effect.
- Hyprland: `input { sensitivity = X }` is global. Diverging mouse/trackpad sensitivity yields a warning and trackpad values win.
- Windows: natural scroll inversion is not supported in MVP.

## License

MIT.
