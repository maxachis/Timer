## Why

A countdown timer is most useful when visible while working in other apps. Currently the timer window disappears behind other windows, defeating the purpose of a visual countdown. An always-on-top toggle lets users pin the timer above everything else.

## What Changes

- Add a pin/unpin toggle button to the UI that calls Tauri's `setAlwaysOnTop()` API
- Add `core:window:allow-set-always-on-top` and `core:window:allow-is-always-on-top` permissions to the capabilities file
- No new Rust code or plugins — this uses the built-in Tauri 2 window API from the frontend
- Add a keyboard shortcut (`T`) to toggle always-on-top

## Capabilities

### New Capabilities
- `always-on-top`: Toggle to pin the timer window above all other windows

### Modified Capabilities
- `keyboard-shortcuts`: Add `T` key binding for toggling always-on-top

## Impact

- `src-tauri/capabilities/default.json`: Add window always-on-top permissions
- `src/App.svelte`: Add toggle button, state tracking, and keyboard shortcut
