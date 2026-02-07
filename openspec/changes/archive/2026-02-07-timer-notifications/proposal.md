## Why

When the timer finishes, there's only a visual color change in the app window. If the user has minimized the app or is focused on another task, they'll miss it entirely. A system notification and an audible alert ensure the user knows their timer is done regardless of what they're doing.

## What Changes

- Add `tauri-plugin-notification` to send an OS-level toast notification when the timer finishes ("Timer Complete — Your 5:00 timer has finished.")
- Play a short alert sound from the frontend using the Web Audio API when the timer reaches zero
- Add `notification:default` permission to the Tauri capabilities file
- Frontend detects the idle → finished transition during polling and triggers both the notification and the sound

## Capabilities

### New Capabilities
- `completion-alert`: System notification and audio alert when the countdown timer reaches zero

### Modified Capabilities
<!-- None — existing countdown-timer, tauri-app-shell, timer-commands, and timer-display specs are unchanged -->

## Impact

- `src-tauri/Cargo.toml`: Add `tauri-plugin-notification` dependency
- `src-tauri/src/lib.rs`: Register the notification plugin
- `src-tauri/capabilities/default.json`: Add `notification:default` permission
- `package.json`: Add `@tauri-apps/plugin-notification` dependency
- `src/App.svelte`: Add notification + sound logic on timer completion
- Sound asset: Include a short alert tone (generated or bundled)
