## Context

The timer app currently hardcodes two values: the default start duration (300 seconds in `lib.rs`) and the time increment/decrement (300 seconds as `DEFAULT_INCREMENT` in `timer.rs`). The frontend also hardcodes the "5 min" label on adjust buttons. Users cannot change these defaults without editing source code.

The app is a single-view Svelte frontend inside a Tauri shell. There is no routing or multi-view navigation. The UI has a top-right corner button (pin/always-on-top), a central timer face, and bottom controls.

## Goals / Non-Goals

**Goals:**
- Let users configure default start time and default increment from a settings view
- Persist settings across app restarts
- Apply settings on startup (timer initializes with saved defaults)
- Keep the settings UI consistent with the existing visual design

**Non-Goals:**
- Per-timer profiles or multiple saved presets
- Syncing settings across devices
- Configuring keyboard shortcut bindings
- Configuring notification/alert preferences

## Decisions

### 1. Settings persistence: `tauri-plugin-store`

Use `tauri-plugin-store` for settings persistence. This writes a JSON file to the Tauri app data directory and provides a simple key-value API from both Rust and JS.

**Alternatives considered:**
- **Raw file I/O in Rust**: More control but requires implementing JSON serialization, file paths, error handling — all of which the store plugin handles.
- **localStorage in the webview**: Not accessible from Rust at startup for timer initialization. Would require a round-trip from frontend to backend on every app launch.

**Rationale:** The store plugin is first-party Tauri, already used in the ecosystem, and provides both Rust and JS access. Settings are small (two numbers), so the key-value model is ideal.

### 2. Settings data model

Store two values:
- `default_duration_secs: u64` — default start time in seconds (default: 300)
- `default_increment_secs: u64` — time adjustment increment in seconds (default: 300)

These are stored in a Tauri store file (e.g., `settings.json`) as a flat JSON object. A Rust struct `AppSettings` with `serde` derive handles serialization.

**Rationale:** Keeping it as raw seconds simplifies backend logic. The frontend converts to a user-friendly format (minutes) for display in the settings UI.

### 3. Settings view: inline view swap, not a modal

The settings view replaces the timer view entirely when active. Navigation is via a cog icon button (top-left, opposite the pin button) and a back button within the settings view.

**Alternatives considered:**
- **Modal/overlay**: Adds z-index complexity and feels heavy for two fields.
- **Separate Tauri window**: Overkill for simple settings; complicates state sharing.
- **Drawer/sidebar**: The app window is small and compact; a sidebar would cramp the timer.

**Rationale:** A full view swap is the simplest approach. A boolean `showSettings` state variable toggles between the timer view and settings view. This avoids routing libraries and keeps the single-component architecture.

### 4. Backend architecture: settings as managed state

Add an `AppSettings` struct to Tauri managed state alongside the existing `Mutex<CountdownTimer>`. Two new commands:
- `get_settings` — returns current settings
- `save_settings` — validates, persists to store, and updates managed state

On startup, load settings from the store (falling back to defaults if no file exists), then initialize the timer with `default_duration_secs`.

The `add_time` and `remove_time` commands will read the increment from `AppSettings` state rather than using the hardcoded constant.

### 5. Settings UI input: minute-based number inputs

The settings form uses numeric inputs in **minutes** (not seconds) since that's the natural unit for timer durations. The frontend converts minutes ↔ seconds when reading/writing settings.

Provide reasonable constraints:
- Default start time: 1–180 minutes (1 min to 3 hours)
- Default increment: 1–60 minutes

### 6. Keyboard shortcuts disabled in settings view

Keyboard shortcuts for timer control (Space, R, arrows, etc.) should not fire when the settings view is active, to avoid interfering with form input. The keydown handler checks whether the settings view is open before processing shortcuts.

## Risks / Trade-offs

- **First launch without store file** → Fallback to hardcoded defaults (300s / 300s). The store file is created on first save.
- **Invalid stored values** (manual file edit, corruption) → Validate on load; fall back to defaults if values are out of range or unparseable.
- **Timer already running when settings changed** → Settings apply on next timer creation/reset, not mid-countdown. The current running timer keeps its duration. This avoids confusing behavior.
- **Store plugin version compatibility** → Pin `tauri-plugin-store` to a specific v2 release matching the existing Tauri v2 setup.
