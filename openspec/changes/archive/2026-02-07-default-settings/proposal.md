## Why

The timer's default start time (5 minutes) and time adjustment increment (5 minutes) are hardcoded across both the Rust backend and Svelte frontend. Users need to be able to configure these defaults to match their typical usage patterns (e.g., a 1-hour focus session with 10-minute increments) without manually creating a new timer each time.

## What Changes

- Add a settings view accessible via a cog icon button, positioned alongside the existing pin button
- The settings view replaces the timer view (not a modal/overlay) and provides fields for configuring default start time and default increment
- Persist settings across app restarts using Tauri's store plugin
- Apply saved settings on startup: the timer initializes with the user's configured default duration instead of the hardcoded 300 seconds
- The time adjustment buttons and keyboard shortcuts use the configured increment instead of the hardcoded 5 minutes
- The adjust button labels dynamically reflect the configured increment (e.g., "10 min" instead of "5 min")

## Capabilities

### New Capabilities
- `settings-view`: Settings UI with cog icon navigation, form fields for default start time and default increment, and a back button to return to the timer
- `settings-persistence`: Backend storage and retrieval of user settings using Tauri's store plugin, with Tauri commands for getting/saving settings

### Modified Capabilities
- `countdown-timer`: The default increment/decrement constant becomes configurable rather than hardcoded at 300 seconds
- `timer-commands`: Timer initialization uses the persisted default duration instead of hardcoded 300 seconds; add_time and remove_time accept a configurable increment from settings

## Impact

- **Rust backend** (`timer.rs`): `DEFAULT_INCREMENT` becomes a parameter rather than a constant; `add_time`/`remove_time` will take the increment from settings
- **Rust backend** (`lib.rs`): New Tauri commands for get/save settings; startup reads persisted settings to initialize timer with the correct default duration; add_time/remove_time pass the configured increment
- **Frontend** (`App.svelte`): Add cog icon button, settings view with form fields, dynamic adjust button labels
- **Dependencies**: `tauri-plugin-store` for persisting settings to disk
- **Config** (`capabilities/*.json`): New permissions for the store plugin
