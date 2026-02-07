## 1. Dependencies and Permissions

- [x] 1.1 Add `tauri-plugin-store` dependency to `src-tauri/Cargo.toml` and register the plugin in `lib.rs`
- [x] 1.2 Add store plugin permissions to the Tauri capabilities configuration

## 2. Settings Backend (Rust)

- [x] 2.1 Create `AppSettings` struct with `default_duration_secs` and `default_increment_secs` fields, with serde derives and defaults of 300 each
- [x] 2.2 Add settings load function that reads from the Tauri store on startup, falling back to defaults if the file is missing or values are invalid
- [x] 2.3 Add `AppSettings` as Tauri managed state (wrapped in `Mutex`) and initialize it at startup alongside the timer
- [x] 2.4 Implement `get_settings` Tauri command that returns the current `AppSettings`
- [x] 2.5 Implement `save_settings` Tauri command that validates input ranges (duration 60–10800s, increment 60–3600s), persists to the store, and updates managed state
- [x] 2.6 Update timer initialization at startup to use the loaded `default_duration_secs` instead of hardcoded 300
- [x] 2.7 Update `add_time` and `remove_time` commands to read the increment from `AppSettings` state instead of using the hardcoded constant

## 3. Settings Frontend (Svelte)

- [x] 3.1 Add `showSettings` state variable and cog icon button in the top-left corner to toggle the settings view
- [x] 3.2 Build the settings view with numeric inputs for default start time (minutes) and time increment (minutes), populated from `get_settings` on open
- [x] 3.3 Implement save button that converts minutes to seconds, calls `save_settings`, reinitializes the timer if idle, and returns to timer view
- [x] 3.4 Implement back button that returns to timer view without saving
- [x] 3.5 Add input validation in the UI (clamp to 1–180 min for duration, 1–60 min for increment)
- [x] 3.6 Update the adjust button labels to display the configured increment dynamically (e.g., "10 min") by reading settings on mount
- [x] 3.7 Guard the `handleKeydown` function to skip timer shortcuts when `showSettings` is true
