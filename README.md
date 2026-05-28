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
git clone https://github.com/bug3/sensync
cd sensync
go build -o ~/.local/bin/sensync ./cmd/sensync
```

Make sure `~/.local/bin` is on your `PATH` (or place the binary wherever you keep CLIs).

## Quick start

1. `sensync init` writes an example config to your user config directory:
   - Linux: `~/.config/sensync/config.toml`
   - macOS: `~/Library/Application Support/sensync/config.toml`
   - Windows: `%APPDATA%\sensync\config.toml`
2. Edit the file to set your preferred sensitivity, acceleration, scroll behavior.
3. On Linux (Hyprland), add `source = ~/.config/hypr/sensync.conf` to your `hyprland.conf` once.
4. Run `sensync apply` on each host. Re-run after changing the config.

If you prefer to keep the config in a git-synced dotfiles repo, drop a `sensync.toml` in the directory you run `sensync` from, or pass `--config <path>` explicitly. The CLI looks in this order: `--config` flag, `./sensync.toml`, then the user config directory.

## Commands

```
sensync init               Write example config to the user config directory
sensync init --force       Overwrite an existing config file
sensync apply              Apply the resolved sensync config on this host
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
