# sensync

Sync mouse sensitivity across Linux, macOS, and Windows from a single config.

## How it works

`sensync` reads one `sensync.toml` and applies the equivalent OS settings on Linux (Hyprland), macOS, and Windows. At `sensitivity = 1.0` with `acceleration = false` all three operating systems produce raw 1:1 input; other values use per-OS approximations.

Trackpad settings live in the same config under `[trackpad]` and follow the same rules.

## Install

```sh
git clone https://github.com/bug3/sensync
cd sensync
go build -o ~/.local/bin/sensync ./cmd/sensync
```

Make sure `~/.local/bin` is on your `PATH`, or drop the binary anywhere you keep CLIs.

## Quick start

1. `sensync init` writes an example config under your user config directory:
   - Linux: `~/.config/sensync/config.toml`
   - macOS: `~/Library/Application Support/sensync/config.toml`
   - Windows: `%APPDATA%\sensync\config.toml`
2. Edit the file to set the values you want.
3. On Linux (Hyprland), add `source = ~/.config/hypr/sensync.conf` to your `hyprland.conf` once.
4. Run `sensync apply` on each host. Re-run after changing the config.

For a git-synced dotfiles flow, drop a `sensync.toml` in the directory you run `sensync` from, or pass `--config <path>` explicitly. Lookup order: `--config` flag, `./sensync.toml`, then the user config directory.

## Commands

| Command | Description |
| --- | --- |
| `sensync init` | Write the example config to the user config directory |
| `sensync init --force` | Overwrite an existing config file |
| `sensync apply` | Apply the resolved config to this host |
| `sensync apply --dry-run` | Print planned changes without applying |
| `sensync apply --yes` | Skip confirmation prompts |
| `sensync apply --config <path>` | Use an explicit config path |
| `sensync get` | Print the live system state as TOML |
| `sensync version` | Print the sensync version |

## Notes per OS

- **macOS**: `defaults` writes apply to new processes; log out for full effect.
- **Hyprland**: `input { sensitivity = X }` is global. Diverging mouse and trackpad sensitivity emits a warning, and trackpad values win.
- **Windows**: natural scroll inversion is not supported in the MVP.

## What is out of scope

- Per-app sensitivity overrides.
- Polling rate (firmware concern).
- Daemon or boot-time auto-apply.
- Cloud config sync (git is the sync layer).
- Linux compositors other than Hyprland.
- Windows non-Precision-Touchpad trackpads.

## License

[MIT](./LICENSE)
