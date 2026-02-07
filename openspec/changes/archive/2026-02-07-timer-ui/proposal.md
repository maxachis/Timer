## Why

The Tauri app launches a window but shows only a placeholder "Timer" heading. Users need an interactive countdown timer interface with controls to start, pause, reset, and adjust the timer by +/- 5 minutes. This also requires a Tauri command bridge since the `CountdownTimer` logic lives in Rust and the UI is in Svelte.

## What Changes

- Add Tauri commands in Rust to expose timer operations to the frontend: create, start, pause, resume, reset, add time, remove time, and get remaining time/state
- Manage a `CountdownTimer` instance as Tauri app state (behind a `Mutex`)
- Replace the placeholder `App.svelte` with a full countdown timer UI: large time display (MM:SS), start/pause toggle, reset button, +5 min / -5 min buttons
- Frontend polls the backend for remaining time to keep the display updated while running
- Visual feedback when the timer finishes (display changes to indicate completion)

## Capabilities

### New Capabilities
- `timer-commands`: Tauri IPC commands that expose `CountdownTimer` operations to the frontend
- `timer-display`: Svelte UI components for countdown display and controls

### Modified Capabilities
<!-- None — the existing countdown-timer and tauri-app-shell specs are unchanged -->

## Impact

- `src-tauri/src/lib.rs`: Add Tauri commands, register state and invoke handler
- `src-tauri/src/timer.rs`: No changes (used as-is)
- `src/App.svelte`: Replace placeholder with timer UI
- `src/main.ts`: No changes expected
- `package.json`: May add `@tauri-apps/api` invoke usage (already a dependency)
