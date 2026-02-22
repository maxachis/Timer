## Why

When the timer completes, the current alerts (system notification + audio beep) are one-shot — if the user misses them, there's no persistent visual cue that the timer has finished. Flashing the taskbar icon provides an ongoing, attention-grabbing signal that persists until the user acknowledges the alarm by resetting the timer.

## What Changes

- Request OS-level window attention (taskbar icon flash) when the timer transitions to the "finished" state
- Continue flashing until the alarm is resolved (timer is reset)
- Stop flashing immediately when the user resets the timer

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `completion-alert`: Add a requirement for taskbar/window attention request that flashes the taskbar icon on timer completion and persists until the timer is reset

## Impact

- **Frontend**: `src/App.svelte` — call `getCurrentWindow().requestUserAttention()` on finish transition, call it again with `null` on reset
- **Tauri config**: May need `core:window:allow-request-user-attention` permission in `src-tauri/capabilities/default.json`
- **No Rust changes** — the Tauri window attention API is accessible directly from the frontend
- **Platform behavior**: On Windows this flashes the taskbar icon; on macOS it bounces the dock icon; on Linux behavior varies by desktop environment
