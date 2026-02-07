## 1. Tauri commands (Rust backend)

- [x] 1.1 Add `TimerStatus` serializable struct to `lib.rs` (remaining_secs, state, is_finished)
- [x] 1.2 Add `get_timer_status` command — locks mutex, returns `TimerStatus`
- [x] 1.3 Add `start_timer`, `pause_timer`, `resume_timer`, `reset_timer` commands
- [x] 1.4 Add `add_time` and `remove_time` commands (use default 5-min increment)
- [x] 1.5 Add `create_timer` command — reinitializes timer with new duration
- [x] 1.6 Register managed state `Mutex<CountdownTimer>` (default 300s) and all commands in `run()`

## 2. Frontend timer display (Svelte)

- [x] 2.1 Create `src/App.svelte` with timer state variables using Svelte 5 `$state` runes
- [x] 2.2 Implement `fetchStatus()` function that invokes `get_timer_status` and updates state
- [x] 2.3 Implement polling: start 100ms interval when running, clear when not running
- [x] 2.4 Add MM:SS countdown display (large, centered)
- [x] 2.5 Add start/pause toggle button (label changes based on state)
- [x] 2.6 Add reset button (disabled when idle)
- [x] 2.7 Add +5:00 and -5:00 buttons (disabled when idle or finished)
- [x] 2.8 Add finished state visual feedback (color change + "Time's up!" message)

## 3. Styling

- [x] 3.1 Style the timer display — large monospace font for MM:SS, centered layout
- [x] 3.2 Style control buttons — clear visual hierarchy, appropriate spacing
- [x] 3.3 Style finished state — distinct color or visual treatment

## 4. Verify

- [x] 4.1 Run `cargo build` in src-tauri to verify commands compile
- [x] 4.2 Run `cargo test` in src-tauri to verify timer tests still pass
- [x] 4.3 Run `cargo tauri dev` to verify the full app launches and timer works end-to-end
