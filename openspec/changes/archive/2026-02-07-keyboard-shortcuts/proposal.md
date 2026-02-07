## Why

The timer app requires clicking buttons for every action — start, pause, resume, reset, and time adjustments. For a utility app that users switch to briefly, keyboard shortcuts let them control the timer without hunting for the right button. This is especially useful when returning to the app from another window: press Space to start/pause instead of aiming for a 64px circle.

## What Changes

- Add a global `keydown` listener in `App.svelte` that maps keys to existing timer actions
- Space bar toggles start/pause/resume (context-aware based on current timer state)
- `R` / `ArrowDown` resets the timer (when not idle)
- `+` / `=` / `ArrowRight` adds 5 minutes, `-` / `ArrowLeft` removes 5 minutes (matching the existing adjust buttons)
- Escape resets a finished timer back to idle
- No new backend commands — all shortcuts invoke the existing `handle*` functions

## Capabilities

### New Capabilities
- `keyboard-shortcuts`: Keyboard shortcut bindings for timer control actions

### Modified Capabilities
<!-- None — existing countdown-timer, tauri-app-shell, timer-commands, timer-display, and completion-alert specs are unchanged -->

## Impact

- `src/App.svelte`: Add `keydown` event listener with key-to-action mapping, cleanup on component destroy
